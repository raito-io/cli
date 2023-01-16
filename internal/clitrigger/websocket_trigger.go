package clitrigger

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"nhooyr.io/websocket"

	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/target"
)

const heartbeatTimeout = time.Minute * 5

type WebsocketClient struct {
	wg           sync.WaitGroup
	config       *target.BaseConfig
	websocketUrl string
}

func NewWebsocketClient(config *target.BaseConfig, websocketUrl string) *WebsocketClient {
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

	s.heartbeat(ctx, conn)
	return s.readMessageFromWebsocket(ctx, conn), nil
}

func (s *WebsocketClient) Wait() {
	s.wg.Wait()
}

func (s *WebsocketClient) readMessageFromWebsocket(ctx context.Context, conn *websocket.Conn) <-chan interface{} {
	ch := make(chan interface{})

	pushToChannel := func(i interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		case ch <- i:
			return true
		}
	}

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()

		defer close(ch)

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

		pushToChannel(triggerEvent)
	}()

	return ch
}

func (s *WebsocketClient) heartbeat(ctx context.Context, conn *websocket.Conn) {
	go func() {
		s.wg.Add(1)
		defer s.wg.Done()

		defer conn.Close(websocket.StatusNormalClosure, "Closing websocket")

		heartbeatMsg := []byte("{\"message\": \"heartbeat\"}")

		timer := time.NewTimer(heartbeatTimeout)

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
					s.config.BaseLogger.Warn("Failed to connect with websocket for %d times", failed)

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
}

type WebsocketCliTrigger struct {
	client *WebsocketClient
	logger hclog.Logger
}

func NewWebsocketCliTrigger(config *target.BaseConfig, websocketUrl string) *WebsocketCliTrigger {
	return &WebsocketCliTrigger{
		client: NewWebsocketClient(config, websocketUrl),
		logger: config.BaseLogger,
	}
}

func (s *WebsocketCliTrigger) TriggerChannel(ctx context.Context) <-chan TriggerEvent {
	outputChannel := make(chan TriggerEvent)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				internalChannel, err := s.client.Start(ctx)
				if err != nil {
					if websocket.CloseStatus(err) > 0 {
						s.logger.Warn("Failed to create websocket. Will try again: %s", err.Error())

						continue
					} else {
						s.logger.Error("Failed to create websocket: %s", err.Error())

						return
					}
				}

				for msg := range internalChannel {
					switch m := msg.(type) {
					case error:
						s.logger.Warn("Received error on websocket: %s", m.Error())

						break
					case TriggerEvent:
						outputChannel <- m
					}
				}
			}
		}
	}()

	return outputChannel
}

func (s *WebsocketCliTrigger) Wait() {
	s.client.Wait()
}
