package jsonSanitizer

import "fmt"

func GenerateDeepJSON(depth int) map[string]interface{} {
	if depth == 0 {
		return map[string]interface{}{"password": "mysecret", "note": "<b>hello</b>"}
	}
	return map[string]interface{}{fmt.Sprintf("level-%d", depth): GenerateDeepJSON(depth - 1)}
}

func GenerateDeepSliceString(count int) []interface{} {
	if count == 0 {
		return []interface{}{"password", "note"}
	}
	return []interface{}{fmt.Sprintf("level-%d", count), GenerateDeepSliceString(count - 1)}

}

func GenerateDeepSliceJSON(count int) []map[string]interface{} {
	if count == 0 {
		return []map[string]interface{}{{"password": "mysecret", "note": "<b>hello</b>"}}
	}
	return []map[string]interface{}{{fmt.Sprintf("level-%d", count): GenerateDeepSliceJSON(count - 1)}}

}
