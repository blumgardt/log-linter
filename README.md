# log-linter

AST/type-based линтер для логов `log/slog` и `go.uber.org/zap`, реализованный на `golang.org/x/tools/go/analysis`.

Линтер проверяет message у log-вызовов и дополнительно анализирует structured-аргументы для правила sensitive.

## Что проверяет (rules)

1. **Lowercase start**: сообщение должно начинаться со строчной буквы (после пробелов).
2. **ASCII only**: сообщение должно быть только ASCII (без кириллицы/emoji/и т.п.).
3. **No special chars**: разрешены только `[A-Za-z0-9 _-]` и пробелы. Доп. символы можно разрешить через конфиг.
4. **Sensitive**: запрещены чувствительные ключевые слова в message/ключах structured полей, а также подозрительные value-идентификаторы.

## Supported loggers (types-based)

- `log/slog`: `Logger.Info/Warn/Error/Debug` (+ `*Context` варианты)
- `go.uber.org/zap`: `Logger.Info/Warn/Error/...`
- `go.uber.org/zap`: `SugaredLogger.Infow/Infof/...`

Важно: есть быстрый фильтр по имени метода, но финальное решение “это slog/zap” делается по `TypesInfo`, чтобы не ловить чужие `Info()`.

---

# Быстрый старт

```bash
go test ./...
golangci-lint custom -v
./custom-gcl run -c .golangci.yml ./...