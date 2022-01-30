package dialect

import (
	"strings"
	"time"

	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
)

// FormatColumnTime format column time
func FormatColumnTime(dialect Dialect, dbLocation *time.Location, col *schemasvr.Column, t time.Time) (interface{}, error) {
	if t.IsZero() {
		if col.Nullable {
			return nil, nil
		}
		if col.SQLType.IsNumeric() {
			return 0, nil
		}
	}
	var tmZone = dbLocation
	if col.TimeZone != nil {
		tmZone = col.TimeZone
	}
	t = t.In(tmZone)
	switch col.SQLType.Name {
	case schemasvr.Date:
		return t.Format("2006-01-02"), nil
	case schemasvr.Time:
		var layout = "15:04:05"
		if col.Length > 0 {
			layout += "." + strings.Repeat("0", col.Length)
		}
		return t.Format(layout), nil
	case schemasvr.DateTime, schemasvr.TimeStamp:
		var layout = "2006-01-02 15:04:05"
		if col.Length > 0 {
			layout += "." + strings.Repeat("0", col.Length)
		}
		return t.Format(layout), nil
	case schemasvr.Varchar:
		return t.Format("2006-01-02 15:04:05"), nil
	case schemasvr.TimeStampz:
		if dialect.URI().DBType == schemasvr.MSSQL {
			return t.Format("2006-01-02T15:04:05.9999999Z07:00"), nil
		} else {
			return t.Format(time.RFC3339Nano), nil
		}
	case schemasvr.BigInt, schemasvr.Int:
		return t.Unix(), nil
	default:
		return t, nil
	}
}
