package jsonSanitizer

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSanitizeSliceJSON(t *testing.T) {

	// generate deep slice BIG JSON
	raw := GenerateDeepJSON(20)
	var slicesJSON []map[string]interface{}
	slicesJSON = append(slicesJSON, raw)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)
	slicesJSON = append(slicesJSON, slicesJSON...)

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(slicesJSON)

	out, _ := json.Marshal(sanitized)
	fmt.Println(string(out))

}

func TestSanitizeJSON(t *testing.T) {

	// generate deep  JSON
	raw := GenerateDeepJSON(20)

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(raw)

	out, _ := json.Marshal(sanitized)
	fmt.Println(string(out))

}
