package clitrigger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/go-hclog"
	"nhooyr.io/websocket"

	"github.com/raito-io/cli/internal/auth"
	"github.com/raito-io/cli/internal/target"
)

type WebsocketCliTrigger struct {
	websocketUrl string
}

func NewWebsocketCliTrigger(websocketUrl string) *WebsocketCliTrigger {

	return &WebsocketCliTrigger{
		websocketUrl: websocketUrl,
	}
}

func (s *WebsocketCliTrigger) TriggerChannel(ctx context.Context, config *target.BaseConfig) (chan TriggerEvent, error) {
	options := websocket.DialOptions{
		HTTPHeader: map[string][]string{},
	}

	err := auth.AddTokenToHeader(&options.HTTPHeader, config)
	if err != nil {
		return nil, err
	}

	c, _, err := websocket.Dial(ctx, s.websocketUrl, &options)
	if err != nil {
		return nil, err
	}

	heartbeat(ctx, c, config.BaseLogger)

	return triggerChannel(ctx, c), nil
}

func heartbeat(ctx context.Context, c *websocket.Conn, logger hclog.Logger) {
	go func() {
		heartbeatMsg := []byte("{\"message\": \"heartbeat\"}")

		ticker := time.NewTicker(time.Duration(9) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := c.Write(ctx, websocket.MessageText, heartbeatMsg)
				if err != nil {
					return
				}

				logger.Debug("Send websocket heartbeat")
			}
		}
	}()
}

func triggerChannel(ctx context.Context, c *websocket.Conn) chan TriggerEvent {
	resultChannel := make(chan TriggerEvent)

	go func() {
		defer close(resultChannel)

		internalChannel := make(chan interface{}, 1)
		defer close(internalChannel)

		defer c.CloseRead(ctx)

		for {
			go readMessageFromWebsocket(ctx, c, internalChannel)

			select {
			case <-ctx.Done():
				return
			case msg := <-internalChannel:
				switch m := msg.(type) {
				case error:
					return
				case TriggerEvent:
					resultChannel <- m
				}
			}
		}

	}()

	return resultChannel
}

func readMessageFromWebsocket(ctx context.Context, c *websocket.Conn, ch chan interface{}) {
	_, msg, err := c.Read(ctx)

	if err != nil {
		ch <- err
	}

	triggerEvent := TriggerEvent{}
	err = json.Unmarshal(msg, &triggerEvent)

	if err != nil {
		ch <- err
	}

	ch <- triggerEvent
}
