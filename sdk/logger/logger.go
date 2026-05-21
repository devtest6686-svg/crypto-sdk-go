package logger

import (
	"context"
	"log"
)

type ILogger interface {
	Infof(ctx context.Context, format string, args ...any)
	Errorf(ctx context.Context, format string, args ...any)
}

type DefLogger struct{}

func (l *DefLogger) Infof(ctx context.Context, format string, args ...any) {
	log.Printf("[info]"+format, args...)
}

func (l *DefLogger) Errorf(ctx context.Context, format string, args ...any) {
	log.Printf("[error]"+format, args...)
}
