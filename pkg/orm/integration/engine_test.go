package integration

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	_ "gitee.com/travelliu/dm"
	"github.com/bhojpur/dbm/pkg/orm"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	_ "github.com/ziutek/mymysql/godrv"
	_ "modernc.org/sqlite"
)

func TestPing(t *testing.T) {
	if err := testEngine.Ping(); err != nil {
		t.Fatal(err)
	}
}
func TestPingContext(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	ctx, canceled := context.WithTimeout(context.Background(), time.Nanosecond)
	defer canceled()
	time.Sleep(time.Nanosecond)
	err := testEngine.(*orm.Engine).PingContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}
func TestAutoTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type TestTx struct {
		Id      int64     `orm:"autoincr pk"`
		Msg     string    `orm:"varchar(255)"`
		Created time.Time `orm:"created"`
	}
	assert.NoError(t, testEngine.Sync(new(TestTx)))
	engine := testEngine.(*orm.Engine)
	// will success
	_, err := engine.Transaction(func(session *orm.Session) (interface{}, error) {
		_, err := session.Insert(TestTx{Msg: "hi"})
		assert.NoError(t, err)
		return nil, nil
	})
	assert.NoError(t, err)
	has, err := engine.Exist(&TestTx{Msg: "hi"})
	assert.NoError(t, err)
	assert.EqualValues(t, true, has)
	// will rollback
	_, err = engine.Transaction(func(session *orm.Session) (interface{}, error) {
		_, err := session.Insert(TestTx{Msg: "hello"})
		assert.NoError(t, err)
		return nil, fmt.Errorf("rollback")
	})
	assert.Error(t, err)
	has, err = engine.Exist(&TestTx{Msg: "hello"})
	assert.NoError(t, err)
	assert.EqualValues(t, false, has)
}
func assertSync(t *testing.T, beans ...interface{}) {
	for _, bean := range beans {
		t.Run(testEngine.TableName(bean, true), func(t *testing.T) {
			assert.NoError(t, testEngine.DropTables(bean))
			assert.NoError(t, testEngine.Sync(bean))
		})
	}
}
func TestDump(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type TestDumpStruct struct {
		Id      int64
		Name    string
		IsMan   bool
		Created time.Time `orm:"created"`
	}
	assertSync(t, new(TestDumpStruct))
	cnt, err := testEngine.Insert([]TestDumpStruct{
		{Name: "1", IsMan: true},
		{Name: "2\n"},
		{Name: "3;"},
		{Name: "4\n;\n''"},
		{Name: "5'\n"},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 5, cnt)
	fp := fmt.Sprintf("%v.sql", testEngine.Dialect().URI().DBType)
	os.Remove(fp)
	assert.NoError(t, testEngine.DumpAllToFile(fp))
	assert.NoError(t, PrepareEngine())
	sess := testEngine.NewSession()
	defer sess.Close()
	assert.NoError(t, sess.Begin())
	_, err = sess.ImportFile(fp)
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())
	for _, tp := range []schemasvr.DBType{schemasvr.SQLITE, schemasvr.MYSQL, schemasvr.POSTGRES, schemasvr.MSSQL} {
		name := fmt.Sprintf("dump_%v.sql", tp)
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, testEngine.DumpAllToFile(name, tp))
		})
	}
}

var dbtypes = []schemasvr.DBType{schemasvr.SQLITE, schemasvr.MYSQL, schemasvr.POSTGRES, schemasvr.MSSQL}

func TestDumpTables(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type TestDumpTableStruct struct {
		Id      int64
		Data    []byte `orm:"BLOB"`
		Name    string
		IsMan   bool
		Created time.Time `orm:"created"`
	}
	assertSync(t, new(TestDumpTableStruct))
	_, err := testEngine.Insert([]TestDumpTableStruct{
		{Name: "1", IsMan: true},
		{Name: "2\n", Data: []byte{'\000', '\001', '\002'}},
		{Name: "3;", Data: []byte("0x000102")},
		{Name: "4\n;\n''", Data: []byte("Help")},
		{Name: "5'\n", Data: []byte("0x48656c70")},
		{Name: "6\\n'\n", Data: []byte("48656c70")},
		{Name: "7\\n'\r\n", Data: []byte("7\\n'\r\n")},
		{Name: "x0809ee"},
		{Name: "090a10"},
	})
	assert.NoError(t, err)
	fp := fmt.Sprintf("%v-table.sql", testEngine.Dialect().URI().DBType)
	os.Remove(fp)
	tb, err := testEngine.TableInfo(new(TestDumpTableStruct))
	assert.NoError(t, err)
	assert.NoError(t, testEngine.(*orm.Engine).DumpTablesToFile([]*schemasvr.Table{tb}, fp))
	assert.NoError(t, PrepareEngine())
	sess := testEngine.NewSession()
	defer sess.Close()
	assert.NoError(t, sess.Begin())
	_, err = sess.ImportFile(fp)
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())
	for _, tp := range dbtypes {
		name := fmt.Sprintf("dump_%v-table.sql", tp)
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, testEngine.(*orm.Engine).DumpTablesToFile([]*schemasvr.Table{tb}, name, tp))
		})
	}
	assert.NoError(t, testEngine.DropTables(new(TestDumpTableStruct)))
	importPath := fmt.Sprintf("dump_%v-table.sql", testEngine.Dialect().URI().DBType)
	t.Run("import_"+importPath, func(t *testing.T) {
		sess := testEngine.NewSession()
		defer sess.Close()
		assert.NoError(t, sess.Begin())
		_, err = sess.ImportFile(importPath)
		assert.NoError(t, err)
		assert.NoError(t, sess.Commit())
	})
}
func TestDumpTables2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type TestDumpTableStruct2 struct {
		Id      int64
		Created time.Time `orm:"Default CURRENT_TIMESTAMP"`
	}
	assertSync(t, new(TestDumpTableStruct2))
	fp := fmt.Sprintf("./dump2-%v-table.sql", testEngine.Dialect().URI().DBType)
	os.Remove(fp)
	tb, err := testEngine.TableInfo(new(TestDumpTableStruct2))
	assert.NoError(t, err)
	assert.NoError(t, testEngine.(*orm.Engine).DumpTablesToFile([]*schemasvr.Table{tb}, fp))
}
func TestSetSchema(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	if testEngine.Dialect().URI().DBType == schemasvr.POSTGRES {
		oldSchema := testEngine.Dialect().URI().Schema
		testEngine.SetSchema("my_schema")
		assert.EqualValues(t, "my_schema", testEngine.Dialect().URI().Schema)
		testEngine.SetSchema(oldSchema)
		assert.EqualValues(t, oldSchema, testEngine.Dialect().URI().Schema)
	}
}
func TestImport(t *testing.T) {
	if testEngine.Dialect().URI().DBType != schemasvr.MYSQL {
		t.Skip()
		return
	}
	sess := testEngine.NewSession()
	defer sess.Close()
	assert.NoError(t, sess.Begin())
	_, err := sess.ImportFile("./testdata/import1.sql")
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())
	assert.NoError(t, sess.Begin())
	_, err = sess.ImportFile("./testdata/import2.sql")
	assert.NoError(t, err)
	assert.NoError(t, sess.Commit())
}
func TestDBVersion(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	version, err := testEngine.DBVersion()
	assert.NoError(t, err)
	fmt.Println(testEngine.Dialect().URI().DBType, "version is", version)
}
func TestGetColumns(t *testing.T) {
	if testEngine.Dialect().URI().DBType != schemasvr.POSTGRES {
		t.Skip()
		return
	}
	type TestCommentStruct struct {
		HasComment int
		NoComment  int
	}
	assertSync(t, new(TestCommentStruct))
	comment := "this is a comment"
	sql := fmt.Sprintf("comment on column %s.%s is '%s'", testEngine.TableName(new(TestCommentStruct), true), "has_comment", comment)
	_, err := testEngine.Exec(sql)
	assert.NoError(t, err)
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	tableName := testEngine.GetColumnMapper().Obj2Table("TestCommentStruct")
	var hasComment, noComment string
	for _, table := range tables {
		if table.Name == tableName {
			col := table.GetColumn("has_comment")
			assert.NotNil(t, col)
			hasComment = col.Comment
			col2 := table.GetColumn("no_comment")
			assert.NotNil(t, col2)
			noComment = col2.Comment
			break
		}
	}
	assert.Equal(t, comment, hasComment)
	assert.Zero(t, noComment)
}
