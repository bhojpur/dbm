package convert

import "database/sql"

var (
	_ sql.Scanner = &EmptyScanner{}
)

// EmptyScanner represents an empty scanner which will ignore the scan
type EmptyScanner struct{}

// Scan implements sql.Scanner
func (EmptyScanner) Scan(value interface{}) error {
	return nil
}
