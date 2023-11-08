package clitrigger

import (
	"context"
	"reflect"
	"sync"

	"github.com/spf13/viper"

	"github.com/raito-io/cli/internal/constants"
)

type SyncTriggerHandler struct {
	syncQueue []SyncTrigger

	m sync.Mutex

	outputChan chan struct{}
}

func NewSyncTriggerHandler() *SyncTriggerHandler {
	targets := len(viper.Get(constants.Targets).([]interface{}))

	h := &SyncTriggerHandler{
		syncQueue:  make([]SyncTrigger, 0, targets),
		outputChan: make(chan struct{}, 1),
	}

	return h
}

func (h *SyncTriggerHandler) HandleTriggerEvent(_ context.Context, triggerEvent *TriggerEvent) {
	if triggerEvent.SyncTrigger == nil {
		return
	}

	h.m.Lock()
	defer h.m.Unlock()

	// Checking if there is an existing trigger in the queue that is equal to the requested one
	for _, queued := range h.syncQueue {
		if reflect.DeepEqual(queued, *triggerEvent.SyncTrigger) {
			// Already queued, ignoring
			return
		}
	}

	h.syncQueue = append(h.syncQueue, *triggerEvent.SyncTrigger)
	h.notifyChannel()
}

func (h *SyncTriggerHandler) Close() {
	close(h.outputChan)
}

func (h *SyncTriggerHandler) Pop() *SyncTrigger {
	h.m.Lock()
	defer h.m.Unlock()

	if len(h.syncQueue) == 0 {
		return nil
	}

	syncTrigger := h.syncQueue[0]
	h.syncQueue = h.syncQueue[1:]

	if len(h.syncQueue) > 0 {
		h.notifyChannel()
	}

	return &syncTrigger
}

func (h *SyncTriggerHandler) TriggerChannel() <-chan struct{} {
	return h.outputChan
}

func (h *SyncTriggerHandler) notifyChannel() {
	select {
	case h.outputChan <- struct{}{}:
		return
	default:
		return
	}
}
