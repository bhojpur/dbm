package integration

import (
	"testing"

	"github.com/bhojpur/dbm/pkg/orm"
	"github.com/bhojpur/dbm/pkg/orm/log"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/stretchr/testify/assert"
)

func TestEngineGroup(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	master := testEngine.(*orm.Engine)
	if master.Dialect().URI().DBType == schemasvr.SQLITE {
		t.Skip()
		return
	}
	eg, err := orm.NewEngineGroup(master, []*orm.Engine{master})
	assert.NoError(t, err)
	eg.SetMaxIdleConns(10)
	eg.SetMaxOpenConns(100)
	eg.SetTableMapper(master.GetTableMapper())
	eg.SetColumnMapper(master.GetColumnMapper())
	eg.SetLogLevel(log.LOG_INFO)
	eg.ShowSQL(true)
}
