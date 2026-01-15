# 🛡️ Ward

`ward` is a shared utility library for the Mbizmarket ecosystem.
It provides a **standardized foundation** for:

- 📜 **Structured Logging** (`slog`) with environment-based levels
- 🧹 **JSON Sanitizer** for safe logging & payload storage

---

## 🚀 Features

### 1. Logging (Go `slog`)
- Standard log levels based on **RFC 5424** (`DEBUG < INFO < WARN < ERROR`).
- Configurable via `LOG_LEVEL` environment variable.
- Supports `TextHandler` (local dev) and `JSONHandler` (production/observability).

NOTE :
 -  `ward` uses `slog` internally, so you can use it as well.
 - args must be in pairs k/v, e.g. `key1, val1, key2, val2` don't set single key args


### 2. Usage
```go
import (
	"github.com/yudhiana/ward/logging"
)


func main() {
	logging.NewLogger().Debug("debug message")
	logging.NewLogger().Info("service started", "version", "1.0.0")
	logging.NewLogger().Warn("retrying connection", "attempt", 3)
	logging.NewLogger().Error("failed to insert", "err", "db timeout")
}
