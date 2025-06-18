package jsonSanitizer

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mataharibiz/ward/logging"
	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer provides configurable JSON sanitization
type Sanitizer struct {
	MaxDepth        int
	RedactionMarker string
	SensitiveFields map[string]bool
	Policy          *bluemonday.Policy
	Ctx             context.Context
}

func getIntEnv(env string) int {
	result, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return DefaultMaxDepth
	}
	return result
}

func getSensitiveFields() (result map[string]bool) {
	defaultFields := strings.Split(DefaultSensitiveFields, "|")
	defaultFields = append(defaultFields, strings.Split(os.Getenv("JSON_SANITIZER_SENSITIVE_FIELDS"), "|")...)

	result = make(map[string]bool)
	for _, field := range defaultFields {
		result[strings.TrimSpace(strings.ToLower(field))] = true
	}
	return
}

func NewJsonSanitizer() *Sanitizer {
	return &Sanitizer{
		MaxDepth:        getIntEnv("JSON_SANITIZER_MAX_DEPTH"),
		RedactionMarker: DefaultRedactionMarker,
		SensitiveFields: getSensitiveFields(),
		Policy:          bluemonday.StrictPolicy(),
		Ctx:             context.Background(),
	}
}

func (s *Sanitizer) Sanitize(value any) any {
	ctx, cancel := context.WithTimeout(s.Ctx, DefaultSanitizeDuration)
	defer cancel()

	start := time.Now()
	result := s.sanitize(ctx, value, 0)

	if time.Since(start) > DefaultSanitizeDuration {
		logging.NewLogger().Warn("JSON sanitization took too long time", "duration_in_seconds", time.Since(start).Seconds())
	}

	if ctx.Err() != nil {
		logging.NewLogger().Warn("JSON sanitization failed", "error", ctx.Err())
		return value
	}

	return result
}

// Sanitize recursively cleans the input JSON-like structure
func (s *Sanitizer) sanitize(ctx context.Context, value any, depth int) any {
	if ctx.Err() != nil {
		return value
	}

	if depth > s.MaxDepth {
		return value
	}

	switch reflect.TypeOf(value).Kind() {
	// case reflect.String:
	// return s.RedactionMarker

	case reflect.Map:
		v := value.(map[string]any)
		for k, val := range v {
			if ctx.Err() != nil {
				break
			}

			if s.SensitiveFields[strings.ToLower(k)] {
				switch reflect.TypeOf(val).Kind() {
				case reflect.String:
					v[k] = s.RedactionMarker
				default:
					v[k] = s.sanitize(ctx, val, depth+1)
				}
			} else {
				v[k] = s.sanitize(ctx, val, depth+1)
			}
		}
		return v

	case reflect.Slice, reflect.Array:
		values := reflect.ValueOf(value).Type().Elem()
		switch values.Kind() {
		case reflect.String:
			v := value.([]string)
			for i, val := range v {
				if ctx.Err() != nil {
					break
				}

				v[i] = s.sanitize(ctx, val, depth).(string)
			}
			return v
		case reflect.Map:
			v := value.([]map[string]any)
			for i := range v {
				if ctx.Err() != nil {
					break
				}

				v[i] = s.sanitize(ctx, v[i], depth).(map[string]any)
			}
		case reflect.Interface:
			v := value.([]any)
			for i := range v {
				if ctx.Err() != nil {
					break
				}

				v[i] = s.sanitize(ctx, v[i], depth)
			}
			return v
		}
	}
	return value
}
