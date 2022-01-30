//go:build dm
// +build dm

package integration

import "github.com/bhojpur/dbm/pkg/orm/schema"

func init() {
	dbtypes = append(dbtypes, schema.DAMENG)
}
