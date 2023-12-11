package clitrigger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"

	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/constants"
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

	var datasources []string
	var onlyTargets map[string]struct{}
	targets := viper.Get(constants.Targets).([]interface{})

	onlyTargetsS := viper.GetString(constants.OnlyTargetsFlag)
	if onlyTargetsS != "" {
		onlyTargets = make(map[string]struct{})
		for _, ot := range strings.Split(onlyTargetsS, ",") {
			onlyTargets[strings.TrimSpace(ot)] = struct{}{}
		}
	}

	for _, targetObj := range targets {
		target, ok := targetObj.(map[string]interface{})
		if !ok {
			continue
		}

		if dsId, found := target["data-source-id"]; found {
			if onlyTargets != nil {
				if _, onlyTargetsFound := onlyTargets[target["name"].(string)]; !onlyTargetsFound {
					continue
				}
			}

			datasources = append(datasources, dsId.(string))
		}
	}

	heartbeatMsgObject := struct {
		Message     string   `json:"message"`
		DataSources []string `json:"datasources"`
	}{
		Message:     "heartbeat",
		DataSources: datasources,
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
						s.logger.Warn(fmt.Sprintf("Received error: %s, Will try to retart websocket.", err.Error()))

						continue
					} else if websocket.CloseStatus(err) > 0 {
						s.logger.Warn(fmt.Sprintf("Failed to create websocket: %s. Will try to retart websocket.", err.Error()))

						continue
					} else {
						s.logger.Error(fmt.Sprintf("Failed to create websocket: %s", err.Error()))

						return
					}
				}

				s.logger.Info("Websocket connection ended. Restart websocket.")
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
				s.logger.Info("Websocket message channel closed")
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
