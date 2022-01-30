//go:build jsoniter
// +build jsoniter

package json

import (
	jsoniter "github.com/json-iterator/go"
)

func init() {
	DefaultJSONHandler = JSONiter{}
}

// JSONiter implements JSONInterface via jsoniter
type JSONiter struct{}

// Marshal implements JSONInterface
func (JSONiter) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

// Unmarshal implements JSONInterface
func (JSONiter) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
