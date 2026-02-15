package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

const (
	MsgUpper1 = "Server started"
	MsgUpper2 = "   Server started"

	MsgOk1 = "server started"
	MsgOk2 = "server_started-ok"
	MsgOk3 = "login"

	MsgNonASCII1 = "сервер стартанул"
	MsgNonASCII2 = "server стартанул"

	MsgSpecial1 = "server started!"
	MsgSpecial2 = "server.started"
	MsgSpecial3 = "user=%s"
	MsgSpecial4 = "path /home"
	MsgSpecial5 = "x:y"

	MsgSensitive1 = "password"
	MsgSensitive2 = "token"
	MsgSensitive3 = "apikey"
	MsgSensitive4 = "api_key"
	MsgSensitive5 = "authorization"
	MsgSensitive6 = "secret"
	MsgSensitive7 = "cookie"
	MsgSensitive8 = "session"
	MsgSensitive9 = "private_key"

	MsgSensitiveExpr = "pass" + "word"
)

const (
	KeyPassword    = "password"
	KeyToken       = "token"
	KeyAPIKeyExpr  = "api" + "key"
	KeyPrivateExpr = "private" + "_key"
)

func main() {
	password := "123"
	token := "t_123"
	apiKey := "k_123"

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	zl, _ := zap.NewProduction()
	sugar := zl.Sugar()

	slog.Info(MsgUpper1)
	slog.Info(MsgUpper2)
	logger.Info(MsgUpper1)
	zl.Info(MsgUpper1)
	sugar.Infow(MsgUpper1)

	slog.Info(MsgOk1)
	slog.Info(MsgOk2)
	logger.Info(MsgOk3)

	slog.Info(MsgNonASCII1)
	slog.Info(MsgNonASCII2)
	logger.Info(MsgNonASCII1)
	zl.Info(MsgNonASCII2)
	sugar.Infow(MsgNonASCII2)

	slog.Info(MsgSpecial1)
	slog.Info(MsgSpecial2)
	slog.Info(MsgSpecial3)
	logger.Info(MsgSpecial4)
	zl.Info(MsgSpecial5)
	sugar.Infow(MsgSpecial1)
	sugar.Infof("user %q", "blumgardt")

	slog.Info("password")
	slog.Info(MsgSensitive1)
	slog.Info(MsgSensitive2)
	slog.Info(MsgSensitive3)
	slog.Info(MsgSensitive4)
	slog.Info(MsgSensitive5)
	slog.Info(MsgSensitive6)
	slog.Info(MsgSensitive7)
	slog.Info(MsgSensitive8)
	slog.Info(MsgSensitive9)
	slog.Info(MsgSensitiveExpr)

	logger.Info(MsgSensitive1)
	zl.Info(MsgSensitive2)
	sugar.Infow(MsgSensitive3)

	slog.Info("login", "password", password)
	slog.Info("login", KeyPassword, password)
	slog.Info("login", KeyToken, token)
	slog.Info("login", "api_key", apiKey)
	slog.Info("login", KeyAPIKeyExpr, apiKey)
	slog.Info("login", KeyPrivateExpr, "k")

	logger.Info("login", "password", password)
	logger.Info("login", KeyToken, token)

	sugar.Infow("login", "password", password)
	sugar.Infow("login", KeyToken, token)

	slog.Info("login", "id", password)
	slog.Info("login", "id", token)
	logger.Info("login", "id", apiKey)
	sugar.Infow("login", "id", token)

	zl.Info("login", zap.String("password", password))
	zl.Info("login", zap.String("token", token))
	zl.Info("login", zap.String("api_key", apiKey))
	zl.Info("login", zap.String(KeyPrivateExpr, "k"))

	zl.Info("login", zap.String("id", password))
	zl.Info("login", zap.String("id", token))
}
