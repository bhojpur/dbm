package migrate

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
	"log"
	"os"
	"testing"

	"github.com/bhojpur/dbm/pkg/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	ID   int64
	Name string
}
type Pet struct {
	ID       int64
	Name     string
	PersonID int
}

const (
	dbName = "testdb.sqlite3"
)

var (
	migrations = []*Migration{
		{
			ID: "201608301400",
			Migrate: func(tx *orm.Engine) error {
				return tx.Sync(&Person{})
			},
			Rollback: func(tx *orm.Engine) error {
				return tx.DropTables(&Person{})
			},
		},
		{
			ID: "201608301430",
			Migrate: func(tx *orm.Engine) error {
				return tx.Sync(&Pet{})
			},
			Rollback: func(tx *orm.Engine) error {
				return tx.DropTables(&Pet{})
			},
		},
	}
)

func TestMigration(t *testing.T) {
	_ = os.Remove(dbName)
	db, err := orm.NewEngine("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	m := New(db, DefaultOptions, migrations)
	err = m.Migrate()
	assert.NoError(t, err)
	exists, _ := db.IsTableExist(&Person{})
	assert.True(t, exists)
	exists, _ = db.IsTableExist(&Pet{})
	assert.True(t, exists)
	assert.Equal(t, 2, tableCount(db, "migrations"))
	err = m.RollbackLast()
	assert.NoError(t, err)
	exists, _ = db.IsTableExist(&Person{})
	assert.True(t, exists)
	exists, _ = db.IsTableExist(&Pet{})
	assert.False(t, exists)
	assert.Equal(t, 1, tableCount(db, "migrations"))
	err = m.RollbackLast()
	assert.NoError(t, err)
	exists, _ = db.IsTableExist(&Person{})
	assert.False(t, exists)
	exists, _ = db.IsTableExist(&Pet{})
	assert.False(t, exists)
	assert.Equal(t, 0, tableCount(db, "migrations"))
}
func TestInitSchema(t *testing.T) {
	os.Remove(dbName)
	db, err := orm.NewEngine("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	m := New(db, DefaultOptions, migrations)
	m.InitSchema(func(tx *orm.Engine) error {
		if err := tx.Sync(&Person{}); err != nil {
			return err
		}
		return tx.Sync(&Pet{})
	})
	err = m.Migrate()
	assert.NoError(t, err)
	exists, _ := db.IsTableExist(&Person{})
	assert.True(t, exists)
	exists, _ = db.IsTableExist(&Pet{})
	assert.True(t, exists)
	assert.Equal(t, 2, tableCount(db, "migrations"))
}
func TestMissingID(t *testing.T) {
	os.Remove(dbName)
	db, err := orm.NewEngine("sqlite3", dbName)
	assert.NoError(t, err)
	if db != nil {
		defer db.Close()
	}
	assert.NoError(t, db.DB().Ping())
	migrationsMissingID := []*Migration{
		{
			Migrate: func(tx *orm.Engine) error {
				return nil
			},
		},
	}
	m := New(db, DefaultOptions, migrationsMissingID)
	assert.Equal(t, ErrMissingID, m.Migrate())
}
func tableCount(db *orm.Engine, tableName string) (count int) {
	row := db.DB().QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
	_ = row.Scan(&count)
	return
}
