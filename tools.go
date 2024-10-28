package ward

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func GetEnv(key string, fallback string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		env = fallback
	}

	return env
}

// ParsePayloadData parse payload data to out struct
func ParsePayloadData(payloadData map[string]interface{}, out interface{}) error {

	// if payload hava 'data' key
	if val, ok := payloadData["data"]; ok {
		payloadData = val.(map[string]interface{})
	}

	jsonRaw, err := json.Marshal(payloadData)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonRaw, out)
}

func GetStackTrace() (stacktrace string) {
	for i := 1; ; i++ {
		pc, f, l, got := runtime.Caller(i)
		if !got {
			break
		}

		pcf := runtime.FuncForPC(pc)
		fnl := strings.Split(pcf.Name(), ".")
		ff, fl := pcf.FileLine(pcf.Entry())
		stacktrace += fmt.Sprintf("%s\n%s:%d => %s\n", fmt.Sprintf("%s:%d => runtime-caller", f, l), ff, fl, fnl[len(fnl)-1])

	}

	return fmt.Sprintf("\n:stacktrace:\n%s", stacktrace)
}
