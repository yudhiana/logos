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
	var jsonString = `{"history":[{"auth":"secret","event":"login","timestamp":"2025-06-19 11:17:31.569680128 +0700 WIB m=+0.002512345","token":"my_secret"}],"password":"my_secret_1750306651569652039","profile":{"auth":"my_secret_1750306651569654189","details":{"device_id":"device_1750306651569679454","ip":"127.0.0.1"},"email":"user_1750306651569652599@example.com","otp":"my_secret_1750306651569653943","pin":"1750306651569654713","secret":"my_secret_1750306651569653347","token":"my_secret_1750306651569653052"},"user_id":1750306651,"username":"user_1750306651569649343"}`

	sanitizer := NewJsonSanitizer()
	sanitized := sanitizer.Sanitize(jsonString)

	out, _ := json.MarshalIndent(sanitized, "", "  ")
	// fmt.Println(string(out))
	if strings.Contains(string(out), "my_secret") {
		t.Errorf("JSON is not sanitized")
	}
}
