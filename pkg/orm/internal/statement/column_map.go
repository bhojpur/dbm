package statement

import (
	"strings"

	"github.com/bhojpur/dbm/pkg/orm/schema"
)

type columnMap []string

func (m columnMap) Contain(colName string) bool {
	if len(m) == 0 {
		return false
	}
	n := len(colName)
	for _, mk := range m {
		if len(mk) != n {
			continue
		}
		if strings.EqualFold(mk, colName) {
			return true
		}
	}
	return false
}
func (m columnMap) Len() int {
	return len(m)
}
func (m columnMap) IsEmpty() bool {
	return len(m) == 0
}
func (m *columnMap) Add(colName string) bool {
	if m.Contain(colName) {
		return false
	}
	*m = append(*m, colName)
	return true
}
func getFlagForColumn(m map[string]bool, col *schema.Column) (val bool, has bool) {
	if len(m) == 0 {
		return false, false
	}
	n := len(col.Name)
	for mk := range m {
		if len(mk) != n {
			continue
		}
		if strings.EqualFold(mk, col.Name) {
			return m[mk], true
		}
	}
	return false, false
}
