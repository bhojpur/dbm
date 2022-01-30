package dialect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bhojpur/dbm/pkg/orm/internal/utils"
	"github.com/bhojpur/dbm/pkg/orm/name"
)

// TableNameWithSchema will add schema prefix on table name if possible
func TableNameWithSchema(dialect Dialect, tableName string) string {
	// Add schema name as prefix of table name.
	// Only for postgres database.
	if dialect.URI().Schema != "" && !strings.Contains(tableName, ".") {
		return fmt.Sprintf("%s.%s", dialect.URI().Schema, tableName)
	}
	return tableName
}

// TableNameNoSchema returns table name with given tableName
func TableNameNoSchema(dialect Dialect, mapper name.Mapper, tableName interface{}) string {
	quote := dialect.Quoter().Quote
	switch tt := tableName.(type) {
	case []string:
		if len(tt) > 1 {
			return fmt.Sprintf("%v AS %v", quote(tt[0]), quote(tt[1]))
		} else if len(tt) == 1 {
			return quote(tt[0])
		}
	case []interface{}:
		l := len(tt)
		var table string
		if l > 0 {
			f := tt[0]
			switch f.(type) {
			case string:
				table = f.(string)
			case name.TableName:
				table = f.(name.TableName).TableName()
			default:
				v := utils.ReflectValue(f)
				t := v.Type()
				if t.Kind() == reflect.Struct {
					table = name.GetTableName(mapper, v)
				} else {
					table = quote(fmt.Sprintf("%v", f))
				}
			}
		}
		if l > 1 {
			return fmt.Sprintf("%v AS %v", quote(table), quote(fmt.Sprintf("%v", tt[1])))
		} else if l == 1 {
			return quote(table)
		}
	case name.TableName:
		return tableName.(name.TableName).TableName()
	case string:
		return tableName.(string)
	case reflect.Value:
		v := tableName.(reflect.Value)
		return name.GetTableName(mapper, v)
	default:
		v := utils.ReflectValue(tableName)
		t := v.Type()
		if t.Kind() == reflect.Struct {
			return name.GetTableName(mapper, v)
		}
		return quote(fmt.Sprintf("%v", tableName))
	}
	return ""
}

// FullTableName returns table name with quote and schema according parameter
func FullTableName(dialect Dialect, mapper name.Mapper, bean interface{}, includeSchema ...bool) string {
	tbName := TableNameNoSchema(dialect, mapper, bean)
	if len(includeSchema) > 0 && includeSchema[0] && !utils.IsSubQuery(tbName) {
		tbName = TableNameWithSchema(dialect, tbName)
	}
	return tbName
}
