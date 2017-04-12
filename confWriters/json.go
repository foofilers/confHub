package confWriters

import (
	"encoding/json"
)

func ConfToStructuredJson(conf map[string]string, pretty bool) (string, error) {
	var jsonBytes []byte
	var err error
	structMap,err := mapToStructMap(conf, "")
	if err != nil {
		return "", err
	}
	if !pretty {
		jsonBytes, err = json.Marshal(structMap)
	} else {
		jsonBytes, err = json.MarshalIndent(structMap, "", "  ")
	}
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}




