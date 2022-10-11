package router

import (
	"github.com/rs/zerolog"
)

// logging keys
const FunctionKey = "function"

// LogEvent
// logs the client response at the end of a handler call.
func LogEvent(log *zerolog.Logger, handler string, err error) *zerolog.Event {
	if err != nil {
		return log.Err(err).Stack()
	}
	return log.Info().Str(FunctionKey, handler)
}
