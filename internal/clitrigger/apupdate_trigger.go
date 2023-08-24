package clitrigger

import (
	"context"
	"sync"

	"github.com/raito-io/golang-set/set"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
)

type ApUpdateTriggerHandler struct {
	queuedDataSources set.Set[string]
	apUpdateQueue     []ApUpdate

	m sync.Mutex

	outputChan chan struct{}
}

func NewApUpdateTriggerHandler() *ApUpdateTriggerHandler {
	targets := len(viper.Get(constants.Targets).([]interface{}))

	h := &ApUpdateTriggerHandler{
		queuedDataSources: set.NewSet[string](),
		apUpdateQueue:     make([]ApUpdate, 0, targets),
		outputChan:        make(chan struct{}, 1),
	}

	return h
}

func (h *ApUpdateTriggerHandler) HandleTriggerEvent(_ context.Context, triggerEvent *TriggerEvent) {
	if triggerEvent.ApUpdate == nil {
		return
	}

	h.m.Lock()
	defer h.m.Unlock()

	apUpdate := ApUpdate{
		Domain: triggerEvent.ApUpdate.Domain,
	}

	for _, dataSource := range triggerEvent.ApUpdate.DataSourceNames {
		if !h.queuedDataSources.Contains(dataSource) {
			apUpdate.DataSourceNames = append(apUpdate.DataSourceNames, dataSource)
			h.queuedDataSources.Add(dataSource)
		}
	}

	if len(apUpdate.DataSourceNames) > 0 {
		h.apUpdateQueue = append(h.apUpdateQueue, apUpdate)
		h.notifyChannel()
	}
}

func (h *ApUpdateTriggerHandler) Close() {
	close(h.outputChan)
}

func (h *ApUpdateTriggerHandler) Pop() *ApUpdate {
	h.m.Lock()
	defer h.m.Unlock()

	if len(h.apUpdateQueue) == 0 {
		return nil
	}

	apUpdate := h.apUpdateQueue[0]
	h.apUpdateQueue = h.apUpdateQueue[1:]

	for _, dataSource := range apUpdate.DataSourceNames {
		h.queuedDataSources.Remove(dataSource)
	}

	if len(h.apUpdateQueue) > 0 {
		h.notifyChannel()
	}

	return &apUpdate
}

func (h *ApUpdateTriggerHandler) TriggerChannel() <-chan struct{} {
	return h.outputChan
}

func (h *ApUpdateTriggerHandler) notifyChannel() {
	select {
	case h.outputChan <- struct{}{}:
		return
	default:
		return
	}
}
