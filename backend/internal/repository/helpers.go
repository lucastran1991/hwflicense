package repository

import (
	"encoding/json"
)

// Helper functions for JSON handling

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func fromJSONArray(jsonStr string) []string {
	if jsonStr == "" {
		return []string{}
	}
	var arr []string
	if err := json.Unmarshal([]byte(jsonStr), &arr); err != nil {
		return []string{}
	}
	return arr
}
