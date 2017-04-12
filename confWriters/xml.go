package confWriters

import (
	"encoding/xml"
	"bytes"
	"reflect"
)

func toXml(conf map[string]interface{}) []xml.Token {
	tokens := make([]xml.Token, 0)

	for key, value := range conf {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t)
		if reflect.TypeOf(value).Kind() == reflect.String {
			tokens = append(tokens, xml.CharData(value.(string)))
		} else {
			tokens = append(tokens, toXml(value.(map[string]interface{}))...)
		}
		tokens = append(tokens, xml.EndElement{t.Name})
	}
	return tokens
}

func structConfToXML(conf map[string]interface{}, pretty bool) ([]byte, error) {
	out := &bytes.Buffer{}
	e := xml.NewEncoder(out)
	if pretty {
		e.Indent("", "  ")
	}
	start := xml.StartElement{Name:xml.Name{"", "config"}}
	tokens := []xml.Token{start}
	tokens = append(tokens, toXml(conf)...)
	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return nil, err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func ConfToStructuredXml(conf map[string]string, pretty bool) (string, error) {
	var xmlBytes []byte
	var err error
	structMap, err := mapToStructMap(conf, "")
	if err != nil {
		return "", err
	}
	if !pretty {
		xmlBytes, err = structConfToXML(structMap, false)
	} else {
		xmlBytes, err = structConfToXML(structMap, true)
	}
	if err != nil {
		return "", err
	}
	return string(xmlBytes), nil
}

func ConfToXml(conf map[string]string, pretty bool) (string, error) {
	var xmlBytes []byte
	var err error
	genConf := make(map[string]interface{})
	for k, v := range conf {
		genConf[k] = v
	}
	if !pretty {
		xmlBytes, err = structConfToXML(genConf, false)
	} else {
		xmlBytes, err = structConfToXML(genConf, true)
	}
	if err != nil {
		return "", err
	}
	return string(xmlBytes), nil
}


