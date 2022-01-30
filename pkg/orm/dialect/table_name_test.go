package dialect

import (
	"testing"

	"github.com/bhojpur/dbm/pkg/orm/name"
	"github.com/stretchr/testify/assert"
)

type MCC struct {
	ID          int64  `orm:"pk 'id'"`
	Code        string `orm:"'code'"`
	Description string `orm:"'description'"`
}

func (mcc *MCC) TableName() string {
	return "mcc"
}
func TestFullTableName(t *testing.T) {
	dialect := QueryDialect("mysql")
	assert.EqualValues(t, "mcc", FullTableName(dialect, name.SnakeMapper{}, &MCC{}))
	assert.EqualValues(t, "mcc", FullTableName(dialect, name.SnakeMapper{}, "mcc"))
}
