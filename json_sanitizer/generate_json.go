package jsonSanitizer

import (
	"fmt"
	"time"
)

func GenerateDeepJSON() map[string]interface{} {
	return map[string]any{
		"user_id":  time.Now().Unix(),
		"username": fmt.Sprintf("user_%d", time.Now().UnixNano()),
		"password": fmt.Sprintf("my_secret_%d", time.Now().UnixNano()),
		"profile": map[string]any{
			"email":  fmt.Sprintf("user_%d@example.com", time.Now().UnixNano()),
			"token":  fmt.Sprintf("my_secret_%d", time.Now().UnixNano()),
			"secret": fmt.Sprintf("my_secret_%d", time.Now().UnixNano()),
			"otp":    fmt.Sprintf("my_secret_%d", time.Now().UnixNano()),
			"auth":   fmt.Sprintf("my_secret_%d", time.Now().UnixNano()),
			"pin":    fmt.Sprintf("%04d", time.Now().UnixNano()),
			"details": map[string]any{
				"device_id": fmt.Sprintf("device_%d", time.Now().UnixNano()),
				"ip":        "127.0.0.1",
			},
		},
		"history": []map[string]any{
			{
				"timestamp": time.Now().String(),
				"event":     "login",
				"auth":      "my_secret",
				"token":     "my_secret",
			},
		},
	}
}

func GenerateDeepSliceString(count int) []interface{} {
	if count == 0 {
		return []interface{}{GenerateDeepJSON()}
	}

	out := make([]interface{}, count)
	for i := 0; i < count; i++ {
		out[i] = GenerateDeepJSON()
	}

	return out
}
