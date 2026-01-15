# 🧾 LOGOS

- 📜 **Structured Logging** (`slog`) with environment-based levels

---

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
	"github.com/yudhiana/logos"
)


func main() {
	logos.NewLogger().Debug("debug message")
	logos.NewLogger().Info("service started", "version", "1.0.0")
	logos.NewLogger().Warn("retrying connection", "attempt", 3)
	logos.NewLogger().Error("failed to insert", "err", "db timeout")
}