package ward

import (
	"encoding/json"
	"os"
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
