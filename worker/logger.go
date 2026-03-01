package worker

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

// Print centralizes logging by level
func (l *Logger) Print(level zerolog.Level, args ...interface{}) {
	msg := fmt.Sprint(args...)

	switch level {
	case zerolog.DebugLevel:
		log.Debug().Msg(msg)
	case zerolog.InfoLevel:
		log.Info().Msg(msg)
	case zerolog.WarnLevel:
		log.Warn().Msg(msg)
	case zerolog.ErrorLevel:
		log.Error().
			Str("component", "asynq-worker").
			Msg(fmt.Sprint(args...))
	case zerolog.FatalLevel:
		log.Fatal().Msg(msg)
	default:
		log.Info().
			Interface("args", args).
			Msg("asynq log")
	}
}

func (l *Logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)

}

func (l *Logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
}
