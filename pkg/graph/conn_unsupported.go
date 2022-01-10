//go:build !linux || !amd64
// +build !linux !amd64

package graph

func NewSqliteConn(root string) (*Database, error) {
	panic("Not implemented")
}
