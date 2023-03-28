package clitrigger

import (
	"context"
)

type CliTrigger interface {
	TriggerChannel(ctx context.Context) <-chan TriggerEvent
	Reset()
	Wait()
}

var _ CliTrigger = (*DummyCliTrigger)(nil)

type DummyCliTrigger struct {
}

func (d DummyCliTrigger) TriggerChannel(ctx context.Context) <-chan TriggerEvent {
	ch := make(chan TriggerEvent)
	defer close(ch)

	return ch
}

func (d DummyCliTrigger) Reset() {
}

func (d DummyCliTrigger) Wait() {
}
