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
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"testing"

	"github.com/bhojpur/dbm/pkg/orm"
	"github.com/bhojpur/dbm/pkg/orm/convert"
	"github.com/bhojpur/dbm/pkg/orm/internal/json"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/stretchr/testify/assert"
)

func TestArrayField(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type ArrayStruct struct {
		Id   int64
		Name [20]byte `orm:"char(80)"`
	}
	assert.NoError(t, testEngine.Sync(new(ArrayStruct)))
	var as = ArrayStruct{
		Name: [20]byte{
			96, 96, 96, 96, 96,
			96, 96, 96, 96, 96,
			96, 96, 96, 96, 96,
			96, 96, 96, 96, 96,
		},
	}
	cnt, err := testEngine.Insert(&as)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var arr ArrayStruct
	has, err := testEngine.ID(1).Get(&arr)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, as.Name, arr.Name)
	var arrs []ArrayStruct
	err = testEngine.Find(&arrs)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(arrs))
	assert.Equal(t, as.Name, arrs[0].Name)
	var newName = [20]byte{
		90, 96, 96, 96, 96,
		96, 96, 96, 96, 96,
		96, 96, 96, 96, 96,
		96, 96, 96, 96, 96,
	}
	cnt, err = testEngine.ID(1).Update(&ArrayStruct{
		Name: newName,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var newArr ArrayStruct
	has, err = testEngine.ID(1).Get(&newArr)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, newName, newArr.Name)
	cnt, err = testEngine.ID(1).Delete(new(ArrayStruct))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var cfgArr ArrayStruct
	has, err = testEngine.ID(1).Get(&cfgArr)
	assert.NoError(t, err)
	assert.Equal(t, false, has)
}
func TestGetBytes(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type Varbinary struct {
		Data []byte `orm:"VARBINARY(250)"`
	}
	err := testEngine.Sync(new(Varbinary))
	assert.NoError(t, err)
	cnt, err := testEngine.Insert(&Varbinary{
		Data: []byte("test"),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var b Varbinary
	has, err := testEngine.Get(&b)
	assert.NoError(t, err)
	assert.Equal(t, true, has)
	assert.Equal(t, "test", string(b.Data))
}

type ConvString string

func (s *ConvString) FromDB(data []byte) error {
	*s = ConvString("prefix---" + string(data))
	return nil
}
func (s *ConvString) ToDB() ([]byte, error) {
	return []byte(string(*s)), nil
}

type ConvConfig struct {
	Name string
	Id   int64
}

func (s *ConvConfig) FromDB(data []byte) error {
	if data == nil {
		s = nil
		return nil
	}
	return json.DefaultJSONHandler.Unmarshal(data, s)
}
func (s *ConvConfig) ToDB() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return json.DefaultJSONHandler.Marshal(s)
}

type SliceType []*ConvConfig

func (s *SliceType) FromDB(data []byte) error {
	return json.DefaultJSONHandler.Unmarshal(data, s)
}
func (s *SliceType) ToDB() ([]byte, error) {
	return json.DefaultJSONHandler.Marshal(s)
}

type Nullable struct {
	Data string
}

func (s *Nullable) FromDB(data []byte) error {
	if data == nil {
		return nil
	}
	*s = Nullable{
		Data: string(data),
	}
	return nil
}
func (s *Nullable) ToDB() ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	return []byte(s.Data), nil
}

type ConvStruct struct {
	Conv      ConvString
	Conv2     *ConvString
	Cfg1      ConvConfig
	Cfg2      *ConvConfig        `orm:"TEXT"`
	Cfg3      convert.Conversion `orm:"BLOB"`
	Slice     SliceType
	Nullable1 *Nullable `orm:"null"`
	Nullable2 *Nullable `orm:"null"`
}

func (c *ConvStruct) BeforeSet(name string, cell orm.Cell) {
	if name == "cfg3" || name == "Cfg3" {
		c.Cfg3 = new(ConvConfig)
	}
}
func TestConversion(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	c := new(ConvStruct)
	assert.NoError(t, testEngine.DropTables(c))
	assert.NoError(t, testEngine.Sync(c))
	var s ConvString = "sssss"
	c.Conv = "tttt"
	c.Conv2 = &s
	c.Cfg1 = ConvConfig{"mm", 1}
	c.Cfg2 = &ConvConfig{"xx", 2}
	c.Cfg3 = &ConvConfig{"zz", 3}
	c.Slice = []*ConvConfig{{"yy", 4}, {"ff", 5}}
	c.Nullable1 = &Nullable{Data: "test"}
	c.Nullable2 = nil
	_, err := testEngine.Nullable("nullable2").Insert(c)
	assert.NoError(t, err)
	c1 := new(ConvStruct)
	has, err := testEngine.Get(c1)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "prefix---tttt", string(c1.Conv))
	assert.NotNil(t, c1.Conv2)
	assert.EqualValues(t, "prefix---"+s, *c1.Conv2)
	assert.EqualValues(t, c.Cfg1, c1.Cfg1)
	assert.NotNil(t, c1.Cfg2)
	assert.EqualValues(t, *c.Cfg2, *c1.Cfg2)
	assert.NotNil(t, c1.Cfg3)
	assert.EqualValues(t, *c.Cfg3.(*ConvConfig), *c1.Cfg3.(*ConvConfig))
	assert.EqualValues(t, 2, len(c1.Slice))
	assert.EqualValues(t, *c.Slice[0], *c1.Slice[0])
	assert.EqualValues(t, *c.Slice[1], *c1.Slice[1])
	cnt, err := testEngine.Where("1=1").Delete(new(ConvStruct))
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	c.Cfg2 = nil
	_, err = testEngine.Insert(c)
	assert.NoError(t, err)
	c2 := new(ConvStruct)
	has, err = testEngine.Get(c2)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, "prefix---tttt", string(c2.Conv))
	assert.NotNil(t, c2.Conv2)
	assert.EqualValues(t, "prefix---"+s, *c2.Conv2)
	assert.EqualValues(t, c.Cfg1, c2.Cfg1)
	assert.Nil(t, c2.Cfg2)
	assert.NotNil(t, c2.Cfg3)
	assert.EqualValues(t, *c.Cfg3.(*ConvConfig), *c2.Cfg3.(*ConvConfig))
	assert.EqualValues(t, 2, len(c2.Slice))
	assert.EqualValues(t, *c.Slice[0], *c2.Slice[0])
	assert.EqualValues(t, *c.Slice[1], *c2.Slice[1])
	assert.NotNil(t, c1.Nullable1)
	assert.Equal(t, c1.Nullable1.Data, "test")
	assert.Nil(t, c1.Nullable2)
}

type MyInt int
type MyUInt uint
type MyFloat float64
type MyStruct struct {
	Type      MyInt
	U         MyUInt
	F         MyFloat
	S         MyString
	IA        []MyInt
	UA        []MyUInt
	FA        []MyFloat
	SA        []MyString
	NameArray []string
	Name      string
	UIA       []uint
	UIA8      []uint8
	UIA16     []uint16
	UIA32     []uint32
	UIA64     []uint64
	UI        uint
	//C64       complex64
	MSS map[string]string
}

func TestCustomType1(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	err := testEngine.DropTables(&MyStruct{})
	assert.NoError(t, err)
	err = testEngine.CreateTables(&MyStruct{})
	assert.NoError(t, err)
	i := MyStruct{Name: "Test", Type: MyInt(1)}
	i.U = 23
	i.F = 1.34
	i.S = "fafdsafdsaf"
	i.UI = 2
	i.IA = []MyInt{1, 3, 5}
	i.UIA = []uint{1, 3}
	i.UIA16 = []uint16{2}
	i.UIA32 = []uint32{4, 5}
	i.UIA64 = []uint64{6, 7, 9}
	i.UIA8 = []uint8{1, 2, 3, 4}
	i.NameArray = []string{"ssss", "fsdf", "lllll, ss"}
	i.MSS = map[string]string{"s": "sfds,ss", "x": "lfjljsl"}
	cnt, err := testEngine.Insert(&i)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	fmt.Println(i)
	i.NameArray = []string{}
	i.MSS = map[string]string{}
	i.F = 0
	has, err := testEngine.Get(&i)
	assert.NoError(t, err)
	assert.True(t, has)
	ss := []MyStruct{}
	err = testEngine.Find(&ss)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(ss))
	assert.EqualValues(t, i, ss[0])
	sss := MyStruct{}
	has, err = testEngine.Get(&sss)
	assert.NoError(t, err)
	assert.True(t, has)
	sss.NameArray = []string{}
	sss.MSS = map[string]string{}
	cnt, err = testEngine.Delete(&sss)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
}

type Status struct {
	Name  string
	Color string
}

var (
	_          convert.Conversion = &Status{}
	Registered                    = Status{"Registered", "white"}
	Approved                      = Status{"Approved", "green"}
	Removed                       = Status{"Removed", "red"}
	Statuses                      = map[string]Status{
		Registered.Name: Registered,
		Approved.Name:   Approved,
		Removed.Name:    Removed,
	}
)

func (s *Status) FromDB(bytes []byte) error {
	if r, ok := Statuses[string(bytes)]; ok {
		*s = r
		return nil
	}
	return errors.New("no this data")
}
func (s *Status) ToDB() ([]byte, error) {
	return []byte(s.Name), nil
}

type UserCus struct {
	Id     int64
	Name   string
	Status Status `orm:"varchar(40)"`
}

func TestCustomType2(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	var uc UserCus
	err := testEngine.CreateTables(&uc)
	assert.NoError(t, err)
	tableName := testEngine.TableName(&uc, true)
	_, err = testEngine.Exec("delete from " + testEngine.Quote(tableName))
	assert.NoError(t, err)
	session := testEngine.NewSession()
	defer session.Close()
	if testEngine.Dialect().URI().DBType == schemasvr.MSSQL {
		err = session.Begin()
		assert.NoError(t, err)
		_, err = session.Exec("set IDENTITY_INSERT " + tableName + " on")
		assert.NoError(t, err)
	}
	cnt, err := session.Insert(&UserCus{1, "xlw", Registered})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	if testEngine.Dialect().URI().DBType == schemasvr.MSSQL {
		err = session.Commit()
		assert.NoError(t, err)
	}
	user := UserCus{}
	exist, err := testEngine.ID(1).Get(&user)
	assert.NoError(t, err)
	assert.True(t, exist)
	fmt.Println(user)
	users := make([]UserCus, 0)
	err = testEngine.Where("`"+testEngine.GetColumnMapper().Obj2Table("Status")+"` = ?", "Registered").Find(&users)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(users))
	fmt.Println(users)
}
func TestUnsignedUint64(t *testing.T) {
	type MyUnsignedStruct struct {
		Id uint64
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MyUnsignedStruct))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 1, len(tables[0].Columns()))
	switch testEngine.Dialect().URI().DBType {
	case schemasvr.SQLITE:
		assert.EqualValues(t, "INTEGER", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.MYSQL:
		assert.EqualValues(t, "UNSIGNED BIGINT", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.POSTGRES, schemasvr.DAMENG:
		assert.EqualValues(t, "BIGINT", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.MSSQL:
		assert.EqualValues(t, "BIGINT", tables[0].Columns()[0].SQLType.Name)
	default:
		assert.False(t, true, "Unsigned is not implemented")
	}
	// Only MYSQL database supports unsigned bigint
	if testEngine.Dialect().URI().DBType != schemasvr.MYSQL {
		return
	}
	cnt, err := testEngine.Insert(&MyUnsignedStruct{
		Id: math.MaxUint64,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var v MyUnsignedStruct
	has, err := testEngine.Get(&v)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, uint64(math.MaxUint64), v.Id)
}
func TestUnsignedUint32(t *testing.T) {
	type MyUnsignedInt32Struct struct {
		Id uint32
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MyUnsignedInt32Struct))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 1, len(tables[0].Columns()))
	switch testEngine.Dialect().URI().DBType {
	case schemasvr.SQLITE:
		assert.EqualValues(t, "INTEGER", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.MYSQL:
		assert.EqualValues(t, "UNSIGNED INT", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.POSTGRES, schemasvr.MSSQL, schemasvr.DAMENG:
		assert.EqualValues(t, "BIGINT", tables[0].Columns()[0].SQLType.Name)
	default:
		assert.False(t, true, "Unsigned is not implemented")
	}
	cnt, err := testEngine.Insert(&MyUnsignedInt32Struct{
		Id: math.MaxUint32,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var v MyUnsignedInt32Struct
	has, err := testEngine.Get(&v)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, uint64(math.MaxUint32), v.Id)
}
func TestUnsignedTinyInt(t *testing.T) {
	type MyUnsignedTinyIntStruct struct {
		Id uint8 `orm:"unsigned tinyint"`
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MyUnsignedTinyIntStruct))
	tables, err := testEngine.DBMetas()
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(tables))
	assert.EqualValues(t, 1, len(tables[0].Columns()))
	switch testEngine.Dialect().URI().DBType {
	case schemasvr.SQLITE, schemasvr.DAMENG:
		assert.EqualValues(t, "INTEGER", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.MYSQL:
		assert.EqualValues(t, "UNSIGNED TINYINT", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.POSTGRES:
		assert.EqualValues(t, "SMALLINT", tables[0].Columns()[0].SQLType.Name)
	case schemasvr.MSSQL:
		assert.EqualValues(t, "INT", tables[0].Columns()[0].SQLType.Name)
	default:
		assert.False(t, true, fmt.Sprintf("Unsigned is not implemented, returned %s", tables[0].Columns()[0].SQLType.Name))
	}
	cnt, err := testEngine.Insert(&MyUnsignedTinyIntStruct{
		Id: math.MaxUint8,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, 1, cnt)
	var v MyUnsignedTinyIntStruct
	has, err := testEngine.Get(&v)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, uint64(math.MaxUint32), v.Id)
}

type MyDecimal big.Int

func (d *MyDecimal) FromDB(data []byte) error {
	i, _ := strconv.ParseInt(string(data), 10, 64)
	if d == nil {
		d = (*MyDecimal)(big.NewInt(i))
	} else {
		(*big.Int)(d).SetInt64(i)
	}
	return nil
}
func (d *MyDecimal) ToDB() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", (*big.Int)(d).Int64())), nil
}
func (d *MyDecimal) AsBigInt() *big.Int {
	return (*big.Int)(d)
}
func (d *MyDecimal) AsInt64() int64 {
	return d.AsBigInt().Int64()
}
func TestDecimal(t *testing.T) {
	type MyMoney struct {
		Id      int64
		Account *MyDecimal
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MyMoney))
	_, err := testEngine.Insert(&MyMoney{
		Account: (*MyDecimal)(big.NewInt(10000000000000000)),
	})
	assert.NoError(t, err)
	var m MyMoney
	has, err := testEngine.Get(&m)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.NotNil(t, m.Account)
	assert.EqualValues(t, 10000000000000000, m.Account.AsInt64())
}

type MyArray [20]byte

func (d *MyArray) FromDB(data []byte) error {
	for i, b := range data[:20] {
		(*d)[i] = b
	}
	return nil
}
func (d MyArray) ToDB() ([]byte, error) {
	return d[:], nil
}
func TestMyArray(t *testing.T) {
	type MyArrayStruct struct {
		Id      int64
		Content MyArray
	}
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(MyArrayStruct))
	var v = [20]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	_, err := testEngine.Insert(&MyArrayStruct{
		Content: v,
	})
	assert.NoError(t, err)
	var m MyArrayStruct
	has, err := testEngine.Get(&m)
	assert.NoError(t, err)
	assert.True(t, has)
	assert.EqualValues(t, v, m.Content)
}
