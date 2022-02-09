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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isFloatEq(i, j float64, precision int) bool {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", i) == fmt.Sprintf("%."+strconv.Itoa(precision)+"f", j)
}
func TestSum(t *testing.T) {
	type SumStruct struct {
		Int   int
		Float float32
	}
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(SumStruct)))
	var (
		cases = []SumStruct{
			{1, 6.2},
			{2, 5.3},
			{92, -0.2},
		}
	)
	var i int
	var f float32
	for _, v := range cases {
		i += v.Int
		f += v.Float
	}
	cnt, err := testEngine.Insert(cases)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, cnt)
	colInt := testEngine.GetColumnMapper().Obj2Table("Int")
	colFloat := testEngine.GetColumnMapper().Obj2Table("Float")
	sumInt, err := testEngine.Sum(new(SumStruct), colInt)
	assert.NoError(t, err)
	assert.EqualValues(t, int(sumInt), i)
	sumFloat, err := testEngine.Sum(new(SumStruct), colFloat)
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return isFloatEq(sumFloat, float64(f), 2)
	})
	sums, err := testEngine.Sums(new(SumStruct), colInt, colFloat)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(sums))
	assert.EqualValues(t, i, int(sums[0]))
	assert.Condition(t, func() bool {
		return isFloatEq(sums[1], float64(f), 2)
	})
	sumsInt, err := testEngine.SumsInt(new(SumStruct), colInt)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(sumsInt))
	assert.EqualValues(t, i, int(sumsInt[0]))
}

type SumStructWithTableName struct {
	Int   int
	Float float32
}

func (s SumStructWithTableName) TableName() string {
	return "sum_struct_with_table_name_1"
}
func TestSumWithTableName(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assert.NoError(t, testEngine.Sync(new(SumStructWithTableName)))
	var (
		cases = []SumStructWithTableName{
			{1, 6.2},
			{2, 5.3},
			{92, -0.2},
		}
	)
	var i int
	var f float32
	for _, v := range cases {
		i += v.Int
		f += v.Float
	}
	cnt, err := testEngine.Insert(cases)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, cnt)
	colInt := testEngine.GetColumnMapper().Obj2Table("Int")
	colFloat := testEngine.GetColumnMapper().Obj2Table("Float")
	sumInt, err := testEngine.Sum(new(SumStructWithTableName), colInt)
	assert.NoError(t, err)
	assert.EqualValues(t, int(sumInt), i)
	sumFloat, err := testEngine.Sum(new(SumStructWithTableName), colFloat)
	assert.NoError(t, err)
	assert.Condition(t, func() bool {
		return isFloatEq(sumFloat, float64(f), 2)
	})
	sums, err := testEngine.Sums(new(SumStructWithTableName), colInt, colFloat)
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(sums))
	assert.EqualValues(t, i, int(sums[0]))
	assert.Condition(t, func() bool {
		return isFloatEq(sums[1], float64(f), 2)
	})
	sumsInt, err := testEngine.SumsInt(new(SumStructWithTableName), colInt)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(sumsInt))
	assert.EqualValues(t, i, int(sumsInt[0]))
}
func TestSumCustomColumn(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type SumStruct2 struct {
		Int   int
		Float float32
	}
	var (
		cases = []SumStruct2{
			{1, 6.2},
			{2, 5.3},
			{92, -0.2},
		}
	)
	assert.NoError(t, testEngine.Sync(new(SumStruct2)))
	cnt, err := testEngine.Insert(cases)
	assert.NoError(t, err)
	assert.EqualValues(t, 3, cnt)
	sumInt, err := testEngine.Sum(new(SumStruct2),
		"CASE WHEN `int` <= 2 THEN `int` ELSE 0 END")
	assert.NoError(t, err)
	assert.EqualValues(t, 3, int(sumInt))
}
