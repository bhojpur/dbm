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
	"fmt"
	"testing"
	"time"

	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/stretchr/testify/assert"
)

func TestStoreEngine(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.DropTables("user_store_engine"))
	type UserinfoStoreEngine struct {
		Id   int64
		Name string
	}
	assert.NoError(t, testEngine.StoreEngine("InnoDB").Table("user_store_engine").CreateTable(&UserinfoStoreEngine{}))
}
func TestCreateTable(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.DropTables("user_user"))
	type UserinfoCreateTable struct {
		Id   int64
		Name string
	}
	assert.NoError(t, testEngine.Table("user_user").CreateTable(&UserinfoCreateTable{}))
}
func TestCreateTable2(t *testing.T) {
	type BaseModelLogicalDel struct {
		Id        string     `orm:"varchar(46) pk"`
		CreatedAt time.Time  `orm:"created"`
		UpdatedAt time.Time  `orm:"updated"`
		DeletedAt *time.Time `orm:"deleted"`
	}
	type TestPerson struct {
		BaseModelLogicalDel `orm:"extends"`
		UserId              string `orm:"varchar(46) notnull"`
		PersonId            string `orm:"varchar(46) notnull"`
		Star                bool
		SortNo              int
		DispName            string `orm:"varchar(100)"`
		FirstName           string
		LastName            string
		FirstNameKana       string
		LastNameKana        string
		BirthYear           *int
		BirthMonth          *int
		BirthDay            *int
		ImageId             string `orm:"varchar(46)"`
		ImageDefaultId      string `orm:"varchar(46)"`
		UserText            string `orm:"varchar(2000)"`
		GenderId            *int
		At1                 string `orm:"varchar(10)"`
		At1Rate             int
		At2                 string `orm:"varchar(10)"`
		At2Rate             int
		At3                 string `orm:"varchar(10)"`
		At3Rate             int
		At4                 string `orm:"varchar(10)"`
		At4Rate             int
		At5                 string `orm:"varchar(10)"`
		At5Rate             int
		At6                 string `orm:"varchar(10)"`
		At6Rate             int
	}
	assert.NoError(t, PrepareEngine())
	tb1, err := testEngine.TableInfo(TestPerson{})
	assert.NoError(t, err)
	tb2, err := testEngine.TableInfo(new(TestPerson))
	assert.NoError(t, err)
	cols1, cols2 := tb1.ColumnsSeq(), tb2.ColumnsSeq()
	assert.EqualValues(t, len(cols1), len(cols2))
	for i, col := range cols1 {
		assert.EqualValues(t, col, cols2[i])
	}
	result, err := testEngine.IsTableExist(new(TestPerson))
	assert.NoError(t, err)
	if result {
		assert.NoError(t, testEngine.DropTables(new(TestPerson)))
	}
	assert.NoError(t, testEngine.CreateTables(new(TestPerson)))
	tables1, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.Len(t, tables1, 1)
	assert.EqualValues(t, len(cols1), len(tables1[0].Columns()))
	result, err = testEngine.IsTableExist(new(TestPerson))
	assert.NoError(t, err)
	if result {
		assert.NoError(t, testEngine.DropTables(new(TestPerson)))
	}
	assert.NoError(t, testEngine.CreateTables(TestPerson{}))
	tables2, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.Len(t, tables2, 1)
	assert.EqualValues(t, len(cols1), len(tables2[0].Columns()))
}
func TestCreateMultiTables(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	session := testEngine.NewSession()
	defer session.Close()
	type UserinfoMultiTable struct {
		Id   int64
		Name string
	}
	user := &UserinfoMultiTable{}
	assert.NoError(t, session.Begin())
	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("user_%v", i)
		assert.NoError(t, session.DropTable(tableName))
		assert.NoError(t, session.Table(tableName).CreateTable(user))
	}
	assert.NoError(t, session.Commit())
}

type SyncTable1 struct {
	Id   int64
	Name string
	Dev  int `orm:"index"`
}
type SyncTable2 struct {
	Id     int64
	Name   string `orm:"unique"`
	Number string `orm:"index"`
	Dev    int
	Age    int
}

func (SyncTable2) TableName() string {
	return "sync_table1"
}

type SyncTable3 struct {
	Id     int64
	Name   string `orm:"unique"`
	Number string `orm:"index"`
	Dev    int
	Age    int
}

func (s *SyncTable3) TableName() string {
	return "sync_table1"
}
func TestSyncTable(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(SyncTable1)))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, "sync_table1", tables[0].Name)
	tableInfo, err := testEngine.TableInfo(new(SyncTable1))
	assert.NoError(t, err)
	assert.EqualValues(t, testEngine.Dialect().SQLType(tables[0].GetColumn("name")), testEngine.Dialect().SQLType(tableInfo.GetColumn("name")))
	assert.NoError(t, testEngine.Sync(new(SyncTable2)))
	tables, err = testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, "sync_table1", tables[0].Name)
	tableInfo, err = testEngine.TableInfo(new(SyncTable2))
	assert.NoError(t, err)
	assert.EqualValues(t, testEngine.Dialect().SQLType(tables[0].GetColumn("name")), testEngine.Dialect().SQLType(tableInfo.GetColumn("name")))
	assert.NoError(t, testEngine.Sync(new(SyncTable3)))
	tables, err = testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, "sync_table1", tables[0].Name)
	tableInfo, err = testEngine.TableInfo(new(SyncTable3))
	assert.NoError(t, err)
	assert.EqualValues(t, testEngine.Dialect().SQLType(tables[0].GetColumn("name")), testEngine.Dialect().SQLType(tableInfo.GetColumn("name")))
}
func TestSyncTable2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Table("sync_tablex").Sync(new(SyncTable1)))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, "sync_tablex", tables[0].Name)
	assert.EqualValues(t, 3, len(tables[0].Columns()))
	type SyncTable4 struct {
		SyncTable1 `orm:"extends"`
		NewCol     string
	}
	assert.NoError(t, testEngine.Table("sync_tablex").Sync(new(SyncTable4)))
	tables, err = testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, "sync_tablex", tables[0].Name)
	assert.EqualValues(t, 4, len(tables[0].Columns()))
	assert.EqualValues(t, colMapper.Obj2Table("NewCol"), tables[0].Columns()[3].Name)
}
func TestSyncTable3(t *testing.T) {
	type SyncTable5 struct {
		Id         int64
		Name       string
		Text       string   `orm:"TEXT"`
		Char       byte     `orm:"CHAR(1)"`
		TenChar    [10]byte `orm:"CHAR(10)"`
		TenVarChar string   `orm:"VARCHAR(10)"`
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(SyncTable5)))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	tableInfo, err := testEngine.TableInfo(new(SyncTable5))
	assert.NoError(t, err)
	assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("name")), testEngine.Dialect().SQLType(tables[0].GetColumn("name")))
	assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("text")), testEngine.Dialect().SQLType(tables[0].GetColumn("text")))
	assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("char")), testEngine.Dialect().SQLType(tables[0].GetColumn("char")))
	assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("ten_char")), testEngine.Dialect().SQLType(tables[0].GetColumn("ten_char")))
	assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("ten_var_char")), testEngine.Dialect().SQLType(tables[0].GetColumn("ten_var_char")))
	if *doNVarcharTest {
		var oldDefaultVarchar string
		var oldDefaultChar string
		oldDefaultVarchar, *defaultVarchar = *defaultVarchar, "nvarchar"
		oldDefaultChar, *defaultChar = *defaultChar, "nchar"
		testEngine.Dialect().SetParams(map[string]string{
			"DEFAULT_VARCHAR": *defaultVarchar,
			"DEFAULT_CHAR":    *defaultChar,
		})
		defer func() {
			*defaultVarchar = oldDefaultVarchar
			*defaultChar = oldDefaultChar
			testEngine.Dialect().SetParams(map[string]string{
				"DEFAULT_VARCHAR": *defaultVarchar,
				"DEFAULT_CHAR":    *defaultChar,
			})
		}()
		assert.NoError(t, PrepareEngine())
		assert.NoError(t, testEngine.Sync(new(SyncTable5)))
		tables, err := testEngine.DBMetas()
		assert.NoError(t, err)
		tableInfo, err := testEngine.TableInfo(new(SyncTable5))
		assert.NoError(t, err)
		assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("name")), testEngine.Dialect().SQLType(tables[0].GetColumn("name")))
		assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("text")), testEngine.Dialect().SQLType(tables[0].GetColumn("text")))
		assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("char")), testEngine.Dialect().SQLType(tables[0].GetColumn("char")))
		assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("ten_char")), testEngine.Dialect().SQLType(tables[0].GetColumn("ten_char")))
		assert.EqualValues(t, testEngine.Dialect().SQLType(tableInfo.GetColumn("ten_var_char")), testEngine.Dialect().SQLType(tables[0].GetColumn("ten_var_char")))
	}
}
func TestSyncTable4(t *testing.T) {
	type SyncTable6 struct {
		Id  int64
		Qty float64 `orm:"numeric(36,2)"`
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(SyncTable6)))
	assert.NoError(t, testEngine.Sync(new(SyncTable6)))
}
func TestIsTableExist(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	exist, err := testEngine.IsTableExist(new(CustomTableName))
	assert.NoError(t, err)
	assert.False(t, exist)
	assert.NoError(t, testEngine.CreateTables(new(CustomTableName)))
	exist, err = testEngine.IsTableExist(new(CustomTableName))
	assert.NoError(t, err)
	assert.True(t, exist)
}
func TestIsTableEmpty(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type NumericEmpty struct {
		Numeric float64 `orm:"numeric(26,2)"`
	}
	type PictureEmpty struct {
		Id          int64
		Url         string `orm:"unique"` //image's url
		Title       string
		Description string
		Created     time.Time `orm:"created"`
		ILike       int
		PageView    int
		From_url    string // nolint
		Pre_url     string `orm:"unique"` //pre view image's url
		Uid         int64
	}
	assert.NoError(t, testEngine.DropTables(&PictureEmpty{}, &NumericEmpty{}))
	assert.NoError(t, testEngine.Sync(new(PictureEmpty), new(NumericEmpty)))
	isEmpty, err := testEngine.IsTableEmpty(&PictureEmpty{})
	assert.NoError(t, err)
	assert.True(t, isEmpty)
	tbName := testEngine.GetTableMapper().Obj2Table("PictureEmpty")
	isEmpty, err = testEngine.IsTableEmpty(tbName)
	assert.NoError(t, err)
	assert.True(t, isEmpty)
}

type CustomTableName struct {
	Id   int64
	Name string
}

func (c *CustomTableName) TableName() string {
	return "customtablename"
}
func TestCustomTableName(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	c := new(CustomTableName)
	assert.NoError(t, testEngine.DropTables(c))
	assert.NoError(t, testEngine.CreateTables(c))
}

type IndexOrUnique struct {
	Id        int64
	Index     int `orm:"index"`
	Unique    int `orm:"unique"`
	Group1    int `orm:"index(ttt)"`
	Group2    int `orm:"index(ttt)"`
	UniGroup1 int `orm:"unique(lll)"`
	UniGroup2 int `orm:"unique(lll)"`
}

func TestIndexAndUnique(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.CreateTables(&IndexOrUnique{}))
	assert.NoError(t, testEngine.DropTables(&IndexOrUnique{}))
	assert.NoError(t, testEngine.CreateTables(&IndexOrUnique{}))
	assert.NoError(t, testEngine.CreateIndexes(&IndexOrUnique{}))
	assert.NoError(t, testEngine.CreateUniques(&IndexOrUnique{}))
	assert.NoError(t, testEngine.DropIndexes(&IndexOrUnique{}))
}
func TestMetaInfo(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(CustomTableName), new(IndexOrUnique)))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(tables))
	tableNames := []string{tables[0].Name, tables[1].Name}
	assert.Contains(t, tableNames, "customtablename")
	assert.Contains(t, tableNames, "index_or_unique")
}
func TestCharst(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	err := testEngine.DropTables("user_charset")
	assert.NoError(t, err)
	err = testEngine.Charset("utf8").Table("user_charset").CreateTable(&Userinfo{})
	assert.NoError(t, err)
}
func TestSync2_1(t *testing.T) {
	type WxTest struct {
		Id                 int   `orm:"not null pk autoincr INT(64)"`
		Passport_user_type int16 `orm:"null int"`
		Id_delete          int8  `orm:"null int default 1"`
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.DropTables("wx_test"))
	assert.NoError(t, testEngine.Sync(new(WxTest)))
	assert.NoError(t, testEngine.Sync(new(WxTest)))
}
func TestUnique_1(t *testing.T) {
	type UserUnique struct {
		Id        int64
		UserName  string    `orm:"unique varchar(25) not null"`
		Password  string    `orm:"varchar(255) not null"`
		Admin     bool      `orm:"not null"`
		CreatedAt time.Time `orm:"created"`
		UpdatedAt time.Time `orm:"updated"`
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.DropTables("user_unique"))
	assert.NoError(t, testEngine.Sync(new(UserUnique)))
	assert.NoError(t, testEngine.DropTables("user_unique"))
	assert.NoError(t, testEngine.CreateTables(new(UserUnique)))
	assert.NoError(t, testEngine.CreateUniques(new(UserUnique)))
}
func TestSync2_2(t *testing.T) {
	type TestSync2Index struct {
		Id     int64
		UserId int64 `orm:"index"`
	}
	assert.NoError(t, PrepareEngine())
	var tableNames = make(map[string]bool)
	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("test_sync2_index_%d", i)
		tableNames[tableName] = true
		assert.NoError(t, testEngine.Table(tableName).Sync(new(TestSync2Index)))
		exist, err := testEngine.IsTableExist(tableName)
		assert.NoError(t, err)
		assert.True(t, exist)
	}
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	for _, table := range tables {
		assert.True(t, tableNames[table.Name])
	}
}
func TestSync2_Default(t *testing.T) {
	type TestSync2Default struct {
		Id       int64
		UserId   int64  `orm:"default(1)"`
		IsMember bool   `orm:"default(true)"`
		Name     string `orm:"default('my_name')"`
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestSync2Default))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default)))
}
func TestSync2_Default2(t *testing.T) {
	type TestSync2Default2 struct {
		Id       int64
		UserId   int64  `orm:"default(1)"`
		IsMember bool   `orm:"default(true)"`
		Name     string `orm:"default('')"`
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestSync2Default2))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
	assert.NoError(t, testEngine.Sync(new(TestSync2Default2)))
}
func TestModifyColum(t *testing.T) {
	// Since SQLITE don't support modify column SQL, currrently just ignore
	if testEngine.Dialect().URI().DBType == schemasvr.SQLITE {
		return
	}
	type TestModifyColumn struct {
		Id       int64
		UserId   int64  `orm:"default(1)"`
		IsMember bool   `orm:"default(true)"`
		Name     string `orm:"char(10)"`
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(TestModifyColumn))
	alterSQL := testEngine.Dialect().ModifyColumnSQL("test_modify_column", &schemasvr.Column{
		Name: "name",
		SQLType: schemasvr.SQLType{
			Name: "VARCHAR",
		},
		Length:         16,
		Nullable:       false,
		DefaultIsEmpty: true,
	})
	_, err := testEngine.Exec(alterSQL)
	assert.NoError(t, err)
}
