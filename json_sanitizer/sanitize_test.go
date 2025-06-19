package jsonSanitizer

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSanitizeSliceJSON(t *testing.T) {

	// generate deep slice BIG JSON
	slicesJSON := GenerateDeepSliceString(1000)

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(slicesJSON)

	out, _ := json.MarshalIndent(sanitized, "", "  ")
	// fmt.Println(string(out))
	if strings.Contains(string(out), "my_secret") {
		t.Errorf("JSON is not sanitized")
	}
}

func TestSanitizeJSON(t *testing.T) {

	// generate deep  JSON
	raw := GenerateDeepJSON()

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(raw)

	out, _ := json.MarshalIndent(sanitized, "", "  ")
	// fmt.Println(string(out))
	if strings.Contains(string(out), "my_secret") {
		t.Errorf("JSON is not sanitized")
	}
}

func TestIsJSON(t *testing.T) {
	var jsonString = `{"name": "John Doe", "age": 30, "city": "New York"}`
	if !IsJSON(jsonString) {
		t.Errorf("JSON is not valid")
	}
}

func TestSanitizeStringJSON(t *testing.T) {
	jsonRaw := GenerateDeepJSON()
	jsonByte, _ := json.Marshal(jsonRaw)
	jsonString := string(jsonByte)

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(jsonString)

	out, _ := json.MarshalIndent(sanitized, "", "  ")
	// fmt.Println(string(out))
	if strings.Contains(string(out), "my_secret") {
		t.Errorf("JSON is not sanitized")
	}
}
