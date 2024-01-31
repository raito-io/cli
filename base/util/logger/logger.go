package logger

import (
	"sync"

	"github.com/hashicorp/go-hclog"
)

var logger hclog.Logger
var onlyOnce sync.Once

// InitializeLogger creates a new logger that should be used as a basis for all logging in the plugin.
// So it's advised to call this method first and store the logger in a (global) variable.
func InitializeLogger() hclog.Logger {
	onlyOnce.Do(func() {
		logger = hclog.New(&hclog.LoggerOptions{
			JSONFormat: true,
		})
	})

	return logger
}
