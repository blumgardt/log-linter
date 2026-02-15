package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

func main() {
	slog.Info("server started", "port", 8080)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Error("db failed", "code", 500)

	zl, _ := zap.NewProduction()
	zl.Info("zap started", zap.Int("port", 8080))

	sugar := zl.Sugar()
	sugar.Infow("sugar started", "env", "prod")
	sugar.Infof("user %s logged in", "blumgardt")
}
