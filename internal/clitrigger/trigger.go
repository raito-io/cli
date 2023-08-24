package clitrigger

import (
	"context"
)

type CliTrigger interface {
	Start(ctx context.Context)
	Subscribe(handler TriggerHandler)
	Reset()
	Wait()
}

var _ CliTrigger = (*DummyCliTrigger)(nil)

type DummyCliTrigger struct {
}

func (d DummyCliTrigger) Start(_ context.Context) {
}

func (d DummyCliTrigger) Subscribe(_ TriggerHandler) {

}

func (d DummyCliTrigger) Reset() {
}

func (d DummyCliTrigger) Wait() {
}
