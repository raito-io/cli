package clitrigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"nhooyr.io/websocket"

	plugin2 "github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/plugin"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/target/types"
)

const heartbeatTimeout = time.Minute * 5
const websocketReset = time.Minute * 90

type WebsocketClient struct {
	wg           sync.WaitGroup
	config       *types.BaseConfig
	websocketUrl string
}

type WebsocketMessageError struct {
	err error
}

func (e *WebsocketMessageError) Error() string {
	return fmt.Sprintf("websocket message error: %s", e.err.Error())
}

func NewWebsocketClient(config *types.BaseConfig, websocketUrl string) *WebsocketClient {
	return &WebsocketClient{
		wg:           sync.WaitGroup{},
		config:       config,
		websocketUrl: websocketUrl,
	}
}

func (s *WebsocketClient) Start(ctx context.Context) (<-chan interface{}, error) {
	options := websocket.DialOptions{
		HTTPHeader: map[string][]string{},
	}

	err := auth.AddTokenToHeader(&options.HTTPHeader, s.config)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.Dial(ctx, s.websocketUrl, &options)
	if err != nil {
		return nil, err
	}

	err = s.config.HealthChecker.MarkLiveness()
	if err != nil {
		s.config.BaseLogger.Warn(fmt.Sprintf("Unable to set liveness marker: %s", err.Error()))
	}

	err = s.heartbeat(ctx, conn)
	if err != nil {
		return nil, err
	}

	return s.readMessageFromWebsocket(ctx, conn), nil
}

func (s *WebsocketClient) Wait() {
	s.wg.Wait()
}

func (s *WebsocketClient) readMessageFromWebsocket(ctx context.Context, conn *websocket.Conn) <-chan interface{} {
	ch := make(chan interface{}, 256) // Small buffer to avoid dropped events

	pushToChannel := func(i interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- i:
			return true
		}
	}

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		defer close(ch)

		for {
			_, msg, err := conn.Read(ctx)

			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					// Do not send any error if the deadline is exceeded
					return
				}

				if !pushToChannel(err) {
					return
				}
			}

			triggerEvent := TriggerEvent{}
			err = json.Unmarshal(msg, &triggerEvent)

			if err != nil {
				if !pushToChannel(err) {
					return
				}
			}

			if !pushToChannel(triggerEvent) {
				return
			}
		}
	}()

	return ch
}

func (s *WebsocketClient) heartbeat(ctx context.Context, conn *websocket.Conn) error {
	s.wg.Add(1)

	hbTargetSync := heartBeatTargetSync{}

	// Get the list of full ds sync targets
	err := target.RunTargets(ctx, s.config, &hbTargetSync)

	if err != nil {
		return fmt.Errorf("target info: %w", err)
	}

	heartbeatMsgObject := struct {
		Message     string   `json:"message"`
		DataSources []string `json:"datasources"`
	}{
		Message:     "heartbeat",
		DataSources: hbTargetSync.DataSources,
	}

	heartbeatMsg, err := json.Marshal(heartbeatMsgObject)
	if err != nil {
		return err
	}

	go func() {
		defer s.wg.Done()

		defer conn.Close(websocket.StatusNormalClosure, "Closing websocket")

		timer := time.NewTimer(0)

		failed := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				s.config.BaseLogger.Debug("Send websocket heartbeat")

				err := conn.Write(ctx, websocket.MessageText, heartbeatMsg)
				if err != nil {
					failed += 1
					s.config.BaseLogger.Warn(fmt.Sprintf("Failed to connect with websocket for %d times", failed))

					timer.Reset(time.Duration(2) * time.Second)

					if failed >= 5 {
						healthErr := s.config.HealthChecker.RemoveLivenessMark()
						if healthErr != nil {
							s.config.BaseLogger.Warn(fmt.Sprintf("Unable to set liveness marker: %s", err.Error()))
						}

						return
					}

					continue
				} else {
					failed = 0

					timer.Reset(heartbeatTimeout)
				}
			}
		}
	}()

	return nil
}

type TriggerHandler interface {
	HandleTriggerEvent(ctx context.Context, triggerEvent *TriggerEvent)
}

type WebsocketCliTrigger struct {
	client *WebsocketClient
	logger hclog.Logger

	subscriberMutex sync.Mutex
	subscribers     []TriggerHandler

	m sync.Mutex

	cancelFn func()
}

func NewWebsocketCliTrigger(config *types.BaseConfig, websocketUrl string) *WebsocketCliTrigger {
	return &WebsocketCliTrigger{
		client: NewWebsocketClient(config, websocketUrl),
		logger: config.BaseLogger,
	}
}

func (s *WebsocketCliTrigger) Subscribe(handler TriggerHandler) {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()

	s.subscribers = append(s.subscribers, handler)
}

func (s *WebsocketCliTrigger) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := s.readChannel(ctx)

				if err != nil {
					wserr := &WebsocketMessageError{}
					if errors.As(err, &wserr) {
						s.logger.Warn(fmt.Sprintf("Received error: %s, Will try to restart websocket.", err.Error()))

						continue
					} else if websocket.CloseStatus(err) > 0 {
						s.logger.Warn(fmt.Sprintf("Failed to create websocket: %s. Will try to restart websocket.", err.Error()))

						continue
					} else {
						s.logger.Error(fmt.Sprintf("Failed to create websocket: %s", err.Error()))

						return
					}
				}

				s.logger.Info("Restart websocket.")
			}
		}
	}()
}

func (s *WebsocketCliTrigger) Reset() {
	s.m.Lock()
	defer s.m.Unlock()

	if s.cancelFn != nil {
		s.cancelFn()
	}
}

func (s *WebsocketCliTrigger) Wait() {
	s.client.Wait()
}

func (s *WebsocketCliTrigger) readChannel(ctx context.Context) error {
	internalCtx, cancelFn := context.WithTimeout(ctx, websocketReset)
	defer func() {
		cancelFn()

		s.m.Lock()
		s.cancelFn = nil
		s.m.Unlock()
	}()

	s.m.Lock()
	s.cancelFn = cancelFn
	s.m.Unlock()

	s.logger.Debug("Creating websocket connection")

	internalChannel, err := s.client.Start(internalCtx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-internalChannel:
			if !ok {
				s.logger.Debug("Websocket message channel closed")
				return nil
			}

			switch m := msg.(type) {
			case error:
				select {
				case <-ctx.Done():
					s.logger.Debug("Websocket closed. Will try to reconnect")

					return nil
				default:
					return &WebsocketMessageError{err: m}
				}
			case TriggerEvent:
				s.subscriberMutex.Lock()
				wg := sync.WaitGroup{}

				for i := range s.subscribers {
					wg.Add(1)

					go func(subscriber TriggerHandler) {
						defer wg.Done()

						subscriber.HandleTriggerEvent(ctx, &m)
					}(s.subscribers[i])
				}

				wg.Wait()
				s.subscriberMutex.Unlock()
			}
		}
	}
}

type heartBeatTargetSync struct {
	DataSources []string
}

func (s *heartBeatTargetSync) TargetSync(ctx context.Context, tConfig *types.BaseTargetConfig) error {
	client, err := plugin.NewPluginClient(tConfig.ConnectorName, tConfig.ConnectorVersion, tConfig.TargetLogger)
	if err != nil {
		return fmt.Errorf("new plugin: %w", err)
	}

	defer client.Close()

	infoClient, err := client.GetInfo()
	if err != nil {
		return fmt.Errorf("get info: %w", err)
	}

	pluginInfo, err := infoClient.GetInfo(ctx)
	if err != nil {
		return fmt.Errorf("get info: %w", err)
	}

	if len(pluginInfo.Type) == 0 {
		// Fallback
		s.DataSources = append(s.DataSources, tConfig.Name)

		return nil
	}

	for _, pluginType := range pluginInfo.Type {
		if pluginType == plugin2.PluginType_PLUGIN_TYPE_FULL_DS_SYNC {
			s.DataSources = append(s.DataSources, tConfig.Name)

			return nil
		}
	}

	return nil
}

func (s *heartBeatTargetSync) Finalize(_ context.Context, _ *types.BaseConfig, _ *target.Options) error {
	return nil
}
