package clitrigger

import (
	"context"
)

type CliTrigger interface {
	TriggerChannel(ctx context.Context) <-chan TriggerEvent
	Reset()
	Wait()
}
