package dialect

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
	"database/sql"
	"fmt"
	"time"

	"github.com/bhojpur/dbm/pkg/orm/core"
)

// ScanContext represents a context when Scan
type ScanContext struct {
	DBLocation   *time.Location
	UserLocation *time.Location
}

// DriverFeatures represents driver feature
type DriverFeatures struct {
	SupportReturnInsertedID bool
}

// Driver represents a database driver
type Driver interface {
	Parse(string, string) (*URI, error)
	Features() *DriverFeatures
	GenScanResult(string) (interface{}, error) // according given column type generating a suitable scan interface
	Scan(*ScanContext, *core.Rows, []*sql.ColumnType, ...interface{}) error
}

var (
	drivers = map[string]Driver{}
)

// RegisterDriver register a driver
func RegisterDriver(driverName string, driver Driver) {
	if driver == nil {
		panic("core: Register driver is nil")
	}
	if _, dup := drivers[driverName]; dup {
		panic("core: Register called twice for driver " + driverName)
	}
	drivers[driverName] = driver
}

// QueryDriver query a driver with name
func QueryDriver(driverName string) Driver {
	return drivers[driverName]
}

// RegisteredDriverSize returned all drivers's length
func RegisteredDriverSize() int {
	return len(drivers)
}

// OpenDialect opens a dialect via driver name and connection string
func OpenDialect(driverName, connstr string) (Dialect, error) {
	driver := QueryDriver(driverName)
	if driver == nil {
		return nil, fmt.Errorf("unsupported driver name: %v", driverName)
	}
	uri, err := driver.Parse(driverName, connstr)
	if err != nil {
		return nil, err
	}
	dialect := QueryDialect(uri.DBType)
	if dialect == nil {
		return nil, fmt.Errorf("unsupported dialect type: %v", uri.DBType)
	}
	dialect.Init(uri)
	return dialect, nil
}

type baseDriver struct{}

func (b *baseDriver) Scan(ctx *ScanContext, rows *core.Rows, types []*sql.ColumnType, v ...interface{}) error {
	return rows.Scan(v...)
}
