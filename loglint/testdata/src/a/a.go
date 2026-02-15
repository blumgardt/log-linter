package a

import (
	"log/slog"

	"go.uber.org/zap"
)

type MyLogger struct{}

func (MyLogger) Info(string) {}

func testSlog() {
	slog.Info("Starting server") // want "log message must start with a lower case letter"
	slog.Info("server started!") // want "log message must not contain any special characters or emoji"
	slog.Info("запуск сервера")  // want "log message must contains only English characters"
	slog.Info("token: abc")      // want "log message contains sensitive keyword \"token\""
}

func testZap() {
	l, _ := zap.NewProduction()
	l.Info("Starting server") // want "log message must start with a lower case letter"

	s := l.Sugar()
	s.Infow("token leaked", "token", "abc") // want "log message contains sensitive keyword \"token\""
}

func testFalsePositive() {
	var x MyLogger
	x.Info("Starting server")
}
