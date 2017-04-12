package confWriters

import (
	"bytes"
	"fmt"
)

func ConfToProperties(conf map[string]string) (string, error) {
	var buffer bytes.Buffer
	for k, v := range conf {
		buffer.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}
	return buffer.String(), nil
}