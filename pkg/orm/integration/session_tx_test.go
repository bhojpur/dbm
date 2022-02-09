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

	"github.com/bhojpur/dbm/pkg/orm/internal/utils"
	"github.com/bhojpur/dbm/pkg/orm/name"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))
	counter := func(t *testing.T) {
		_, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
	}
	counter(t)
	//defer counter()
	session := testEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	assert.NoError(t, err)
	user1 := Userinfo{Username: "bhojpur", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)
	user2 := Userinfo{Username: "yyy"}
	_, err = session.Where("`id` = ?", 0).Update(&user2)
	assert.NoError(t, err)
	_, err = session.Delete(&user2)
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)
}
func TestCombineTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	assertSync(t, new(Userinfo))
	counter := func() {
		total, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
		fmt.Printf("----now total %v records\n", total)
	}
	counter()
	//defer counter()
	session := testEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	assert.NoError(t, err)
	user1 := Userinfo{Username: "bhojpur2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)
	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("`id` = ?", 0).Update(&user2)
	assert.NoError(t, err)
	_, err = session.Exec("delete from "+testEngine.Quote(testEngine.TableName("userinfo", true))+" where `username` = ?", user2.Username)
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)
}
func TestCombineTransactionSameMapper(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	oldMapper := testEngine.GetColumnMapper()
	testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
	testEngine.SetMapper(name.SameMapper{})
	defer func() {
		testEngine.UnMapType(utils.ReflectValue(new(Userinfo)).Type())
		testEngine.SetMapper(oldMapper)
	}()
	assertSync(t, new(Userinfo))
	counter := func() {
		total, err := testEngine.Count(&Userinfo{})
		assert.NoError(t, err)
		fmt.Printf("----now total %v records\n", total)
	}
	counter()
	defer counter()
	session := testEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	assert.NoError(t, err)
	user1 := Userinfo{Username: "bhojpur2", Departname: "dev", Alias: "lunny", Created: time.Now()}
	_, err = session.Insert(&user1)
	assert.NoError(t, err)
	user2 := Userinfo{Username: "zzz"}
	_, err = session.Where("`id` = ?", 0).Update(&user2)
	assert.NoError(t, err)
	_, err = session.Exec("delete from  "+testEngine.Quote(testEngine.TableName("Userinfo", true))+" where `Username` = ?", user2.Username)
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)
}
func TestMultipleTransaction(t *testing.T) {
	assert.NoError(t, PrepareEngine())
	type MultipleTransaction struct {
		Id   int64
		Name string
	}
	assertSync(t, new(MultipleTransaction))
	session := testEngine.NewSession()
	defer session.Close()
	err := session.Begin()
	assert.NoError(t, err)
	m1 := MultipleTransaction{Name: "bhojpur2"}
	_, err = session.Insert(&m1)
	assert.NoError(t, err)
	user2 := MultipleTransaction{Name: "zzz"}
	_, err = session.Where("`id` = ?", 0).Update(&user2)
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)
	var ms []MultipleTransaction
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(ms))
	err = session.Begin()
	assert.NoError(t, err)
	_, err = session.Where("`id`=?", m1.Id).Delete(new(MultipleTransaction))
	assert.NoError(t, err)
	err = session.Commit()
	assert.NoError(t, err)
	ms = make([]MultipleTransaction, 0)
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, len(ms))
	err = session.Begin()
	assert.NoError(t, err)
	_, err = session.Insert(&MultipleTransaction{
		Name: "ssss",
	})
	assert.NoError(t, err)
	err = session.Rollback()
	assert.NoError(t, err)
	ms = make([]MultipleTransaction, 0)
	err = session.Find(&ms)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, len(ms))
}
