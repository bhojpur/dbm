package utils

import (
	"strings"
)

// IsSubQuery returns true if it contains a sub query
func IsSubQuery(tbName string) bool {
	const selStr = "select"
	if len(tbName) <= len(selStr)+1 {
		return false
	}
	return strings.EqualFold(tbName[:len(selStr)], selStr) ||
		strings.EqualFold(tbName[:len(selStr)+1], "("+selStr)
}
