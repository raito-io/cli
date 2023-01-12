package clitrigger

import (
	"context"

	"github.com/raito-io/cli/internal/target"
)

type CliTrigger interface {
	TriggerChannel(ctx context.Context, targetConfig *target.BaseConfig) (chan TriggerEvent, error)
}
