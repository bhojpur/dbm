//go:build amd64
// +build amd64

package graph

import (
	"database/sql"
	_ "github.com/bhojpur/dbm/pkg/sqlite" // registers sqlite
	"os"
)

func NewSqliteConn(root string) (*Database, error) {
	initDatabase := false
	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			initDatabase = true
		} else {
			return nil, err
		}
	}
	conn, err := sql.Open("sqlite3", root)
	if err != nil {
		return nil, err
	}
	return NewDatabase(conn, initDatabase)
}
