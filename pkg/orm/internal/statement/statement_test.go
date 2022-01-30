package statement

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bhojpur/dbm/pkg/orm/cache"
	dialectsvr "github.com/bhojpur/dbm/pkg/orm/dialect"
	"github.com/bhojpur/dbm/pkg/orm/name"
	"github.com/bhojpur/dbm/pkg/orm/schema"
	tags "github.com/bhojpur/dbm/pkg/orm/tag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var (
	dialect   dialectsvr.Dialect
	tagParser *tags.Parser
)

func TestMain(m *testing.M) {
	var err error
	dialect, err = dialectsvr.OpenDialect("sqlite3", "./test.db")
	if err != nil {
		panic("unknow dialect")
	}
	tagParser = tags.NewParser("orm", dialect, name.SnakeMapper{}, name.SnakeMapper{}, cache.NewManager())
	if tagParser == nil {
		panic("tags parser is nil")
	}
	m.Run()
	os.Exit(0)
}

var colStrTests = []struct {
	omitColumn        string
	onlyToDBColumnNdx int
	expected          string
}{
	{"", -1, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"Code2", -1, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"", 1, "`ID`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`, `Longitude`"},
	{"Code3", 1, "`ID`, `Caption`, `Code1`, `Code2`, `ParentID`, `Latitude`, `Longitude`"},
	{"Longitude", 1, "`ID`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`"},
	{"", 8, "`ID`, `IsDeleted`, `Caption`, `Code1`, `Code2`, `Code3`, `ParentID`, `Latitude`"},
}

func TestColumnsStringGeneration(t *testing.T) {
	for ndx, testCase := range colStrTests {
		statement, err := createTestStatement()
		assert.NoError(t, err)
		if testCase.omitColumn != "" {
			statement.Omit(testCase.omitColumn)
		}
		columns := statement.RefTable.Columns()
		if testCase.onlyToDBColumnNdx >= 0 {
			columns[testCase.onlyToDBColumnNdx].MapType = schema.ONLYTODB
		}
		actual := statement.genColumnStr()
		if actual != testCase.expected {
			t.Errorf("[test #%d] Unexpected columns string:\nwant:\t%s\nhave:\t%s", ndx, testCase.expected, actual)
		}
		if testCase.onlyToDBColumnNdx >= 0 {
			columns[testCase.onlyToDBColumnNdx].MapType = schema.TWOSIDES
		}
	}
}
func TestConvertSQLOrArgs(t *testing.T) {
	statement, err := createTestStatement()
	assert.NoError(t, err)
	// example orm struct
	// type Table struct {
	// 	ID  int
	// 	del *time.Time `orm:"deleted"`
	// }
	args := []interface{}{
		"INSERT `table` (`id`, `del`) VALUES (?, ?)", 1, (*time.Time)(nil),
	}
	// before fix, here will panic
	_, _, err = statement.convertSQLOrArgs(args...)
	assert.NoError(t, err)
}
func BenchmarkGetFlagForColumnWithICKey_ContainsKey(b *testing.B) {
	b.StopTimer()
	mapCols := make(map[string]bool)
	cols := []*schema.Column{
		{Name: `ID`},
		{Name: `IsDeleted`},
		{Name: `Caption`},
		{Name: `Code1`},
		{Name: `Code2`},
		{Name: `Code3`},
		{Name: `ParentID`},
		{Name: `Latitude`},
		{Name: `Longitude`},
	}
	for _, col := range cols {
		mapCols[strings.ToLower(col.Name)] = true
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, col := range cols {
			if _, ok := getFlagForColumn(mapCols, col); !ok {
				b.Fatal("Unexpected result")
			}
		}
	}
}
func BenchmarkGetFlagForColumnWithICKey_EmptyMap(b *testing.B) {
	b.StopTimer()
	mapCols := make(map[string]bool)
	cols := []*schema.Column{
		{Name: `ID`},
		{Name: `IsDeleted`},
		{Name: `Caption`},
		{Name: `Code1`},
		{Name: `Code2`},
		{Name: `Code3`},
		{Name: `ParentID`},
		{Name: `Latitude`},
		{Name: `Longitude`},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, col := range cols {
			if _, ok := getFlagForColumn(mapCols, col); ok {
				b.Fatal("Unexpected result")
			}
		}
	}
}

type TestType struct {
	ID        int64   `orm:"ID PK"`
	IsDeleted bool    `orm:"IsDeleted"`
	Caption   string  `orm:"Caption"`
	Code1     string  `orm:"Code1"`
	Code2     string  `orm:"Code2"`
	Code3     string  `orm:"Code3"`
	ParentID  int64   `orm:"ParentID"`
	Latitude  float64 `orm:"Latitude"`
	Longitude float64 `orm:"Longitude"`
}

func (TestType) TableName() string {
	return "TestTable"
}
func createTestStatement() (*Statement, error) {
	statement := NewStatement(dialect, tagParser, time.Local)
	if err := statement.SetRefValue(reflect.ValueOf(TestType{})); err != nil {
		return nil, err
	}
	return statement, nil
}
func BenchmarkColumnsStringGeneration(b *testing.B) {
	b.StopTimer()
	statement, err := createTestStatement()
	if err != nil {
		panic(err)
	}
	testCase := colStrTests[0]
	if testCase.omitColumn != "" {
		statement.Omit(testCase.omitColumn) // !nemec784! Column must be skipped
	}
	if testCase.onlyToDBColumnNdx >= 0 {
		columns := statement.RefTable.Columns()
		columns[testCase.onlyToDBColumnNdx].MapType = schema.ONLYTODB // !nemec784! Column must be skipped
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		actual := statement.genColumnStr()
		if actual != testCase.expected {
			b.Errorf("Unexpected columns string:\nwant:\t%s\nhave:\t%s", testCase.expected, actual)
		}
	}
}
