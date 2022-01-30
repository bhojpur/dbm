//go:build gojson
// +build gojson

package json

import (
	gojson "github.com/goccy/go-json"
)

func init() {
	DefaultJSONHandler = GOjson{}
}

// GOjson implements JSONInterface via gojson
type GOjson struct{}

// Marshal implements JSONInterface
func (GOjson) Marshal(v interface{}) ([]byte, error) {
	return gojson.Marshal(v)
}

// Unmarshal implements JSONInterface
func (GOjson) Unmarshal(data []byte, v interface{}) error {
	return gojson.Unmarshal(data, v)
}
