package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

const (
	// rule1: uppercase start (ASCII, без спецсимволов)
	MsgUpper1 = "Server started"
	MsgUpper2 = "   Server started" // пробелы + uppercase

	// ok
	MsgOk1 = "server started"
	MsgOk2 = "server_started-ok"
	MsgOk3 = "login"

	// rule2: non-ASCII (кириллица/символы вне ASCII)
	MsgNonASCII1 = "сервер стартанул"
	MsgNonASCII2 = "server стартанул"

	// rule3: special chars (ASCII, но не в allowlist)
	MsgSpecial1 = "server started!"
	MsgSpecial2 = "server.started"
	MsgSpecial3 = "user=%s"
	MsgSpecial4 = "path /home"
	MsgSpecial5 = "x:y"

	// rule4: sensitive keywords в message (константы)
	MsgSensitive1 = "password"
	MsgSensitive2 = "token"
	MsgSensitive3 = "apikey"
	MsgSensitive4 = "api_key"
	MsgSensitive5 = "authorization"
	MsgSensitive6 = "secret"
	MsgSensitive7 = "cookie"
	MsgSensitive8 = "session"
	MsgSensitive9 = "private_key"

	// rule4: sensitive keyword как const-выражение (TypesInfo const string)
	MsgSensitiveExpr = "pass" + "word" // => "password"
)

const (
	// rule4: sensitive ключи как const/expr
	KeyPassword    = "password"
	KeyToken       = "token"
	KeyAPIKeyExpr  = "api" + "key" // => "apikey"
	KeyPrivateExpr = "private" + "_key"
)

func main() {
	// значения (неважно какие — линтер смотрит на имена/ключи/сообщения)
	password := "123"
	token := "t_123"
	apiKey := "k_123"
	secret := "s_123"
	cookie := "c_123"
	session := "s_456"

	// slog logger instance
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// zap logger + sugared
	zl, _ := zap.NewProduction()
	sugar := zl.Sugar()

	// =========================================================
	// RULE 1: message must start with lowercase
	// =========================================================
	slog.Info(MsgUpper1)
	slog.Info(MsgUpper2)
	logger.Info(MsgUpper1)
	zl.Info(MsgUpper1)
	sugar.Infow(MsgUpper1)

	// OK cases (должны пройти 1–3)
	slog.Info(MsgOk1)
	slog.Info(MsgOk2)
	logger.Info(MsgOk3)

	// =========================================================
	// RULE 2: English only (ASCII)
	// =========================================================
	slog.Info(MsgNonASCII1)
	slog.Info(MsgNonASCII2)
	logger.Info(MsgNonASCII1)
	zl.Info(MsgNonASCII2)
	sugar.Infow(MsgNonASCII2)

	// =========================================================
	// RULE 3: no special chars / emoji (ASCII, но запрещённые символы)
	// =========================================================
	slog.Info(MsgSpecial1)
	slog.Info(MsgSpecial2)
	slog.Info(MsgSpecial3)
	logger.Info(MsgSpecial4)
	zl.Info(MsgSpecial5)
	sugar.Infow(MsgSpecial1)

	// printf-style (тоже попадёт в rule3 из-за %)
	sugar.Infof("user %q", "blumgardt")

	// =========================================================
	// RULE 4: no sensitive data
	// A) sensitive keyword в message (literal/const/const-expr)
	// =========================================================
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
	slog.Info(MsgSensitiveExpr) // "pass"+"word" => "password"

	logger.Info(MsgSensitive1)
	zl.Info(MsgSensitive2)
	sugar.Infow(MsgSensitive3)

	// B) slog/sugar key-value: sensitive KEY (literal/const/const-expr)
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

	// C) slog/sugar key-value: potentially sensitive VALUE ident (ключ нейтральный, значение по имени)
	slog.Info("login", "id", password) // value ident name "password"
	slog.Info("login", "id", token)    // value ident name "token"
	logger.Info("login", "id", apiKey) // value ident name "apiKey" (если keyword list ловит "apikey" — может не поймать camelCase)
	sugar.Infow("login", "id", token)

	// D) zap fields: sensitive KEY
	zl.Info("login", zap.String("password", password))
	zl.Info("login", zap.String("token", token))
	zl.Info("login", zap.String("api_key", apiKey))
	zl.Info("login", zap.String(KeyPrivateExpr, "k"))

	// E) zap fields: potentially sensitive VALUE ident (ключ нейтральный)
	zl.Info("login", zap.String("id", password))
	zl.Info("login", zap.String("id", token))

	// чтобы не ругался компилятор на неиспользуемое
	_ = []string{secret, cookie, session}
}
