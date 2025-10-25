package helpers

import "encoding/json"

// JsonStringToMap parses a JSON string and converts it to a map with string keys and any values
func JsonStringToMap(jsonString string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
