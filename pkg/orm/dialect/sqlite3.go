package dialect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bhojpur/dbm/pkg/orm/core"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
)

var (
	sqlite3ReservedWords = map[string]bool{
		"ABORT":             true,
		"ACTION":            true,
		"ADD":               true,
		"AFTER":             true,
		"ALL":               true,
		"ALTER":             true,
		"ANALYZE":           true,
		"AND":               true,
		"AS":                true,
		"ASC":               true,
		"ATTACH":            true,
		"AUTOINCREMENT":     true,
		"BEFORE":            true,
		"BEGIN":             true,
		"BETWEEN":           true,
		"BY":                true,
		"CASCADE":           true,
		"CASE":              true,
		"CAST":              true,
		"CHECK":             true,
		"COLLATE":           true,
		"COLUMN":            true,
		"COMMIT":            true,
		"CONFLICT":          true,
		"CONSTRAINT":        true,
		"CREATE":            true,
		"CROSS":             true,
		"CURRENT_DATE":      true,
		"CURRENT_TIME":      true,
		"CURRENT_TIMESTAMP": true,
		"DATABASE":          true,
		"DEFAULT":           true,
		"DEFERRABLE":        true,
		"DEFERRED":          true,
		"DELETE":            true,
		"DESC":              true,
		"DETACH":            true,
		"DISTINCT":          true,
		"DROP":              true,
		"EACH":              true,
		"ELSE":              true,
		"END":               true,
		"ESCAPE":            true,
		"EXCEPT":            true,
		"EXCLUSIVE":         true,
		"EXISTS":            true,
		"EXPLAIN":           true,
		"FAIL":              true,
		"FOR":               true,
		"FOREIGN":           true,
		"FROM":              true,
		"FULL":              true,
		"GLOB":              true,
		"GROUP":             true,
		"HAVING":            true,
		"IF":                true,
		"IGNORE":            true,
		"IMMEDIATE":         true,
		"IN":                true,
		"INDEX":             true,
		"INDEXED":           true,
		"INITIALLY":         true,
		"INNER":             true,
		"INSERT":            true,
		"INSTEAD":           true,
		"INTERSECT":         true,
		"INTO":              true,
		"IS":                true,
		"ISNULL":            true,
		"JOIN":              true,
		"KEY":               true,
		"LEFT":              true,
		"LIKE":              true,
		"LIMIT":             true,
		"MATCH":             true,
		"NATURAL":           true,
		"NO":                true,
		"NOT":               true,
		"NOTNULL":           true,
		"NULL":              true,
		"OF":                true,
		"OFFSET":            true,
		"ON":                true,
		"OR":                true,
		"ORDER":             true,
		"OUTER":             true,
		"PLAN":              true,
		"PRAGMA":            true,
		"PRIMARY":           true,
		"QUERY":             true,
		"RAISE":             true,
		"RECURSIVE":         true,
		"REFERENCES":        true,
		"REGEXP":            true,
		"REINDEX":           true,
		"RELEASE":           true,
		"RENAME":            true,
		"REPLACE":           true,
		"RESTRICT":          true,
		"RIGHT":             true,
		"ROLLBACK":          true,
		"ROW":               true,
		"SAVEPOINT":         true,
		"SELECT":            true,
		"SET":               true,
		"TABLE":             true,
		"TEMP":              true,
		"TEMPORARY":         true,
		"THEN":              true,
		"TO":                true,
		"TRANSACTI":         true,
		"TRIGGER":           true,
		"UNION":             true,
		"UNIQUE":            true,
		"UPDATE":            true,
		"USING":             true,
		"VACUUM":            true,
		"VALUES":            true,
		"VIEW":              true,
		"VIRTUAL":           true,
		"WHEN":              true,
		"WHERE":             true,
		"WITH":              true,
		"WITHOUT":           true,
	}
	sqlite3Quoter = schemasvr.Quoter{
		Prefix:     '`',
		Suffix:     '`',
		IsReserved: schemasvr.AlwaysReserve,
	}
)

type sqlite3 struct {
	Base
}

func (db *sqlite3) Init(uri *URI) error {
	db.quoter = sqlite3Quoter
	return db.Base.Init(db, uri)
}
func (db *sqlite3) Version(ctx context.Context, queryer core.Queryer) (*schemasvr.Version, error) {
	rows, err := queryer.QueryContext(ctx, "SELECT sqlite_version()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var version string
	if !rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		return nil, errors.New("unknow version")
	}
	if err := rows.Scan(&version); err != nil {
		return nil, err
	}
	return &schemasvr.Version{
		Number:  version,
		Edition: "sqlite",
	}, nil
}
func (db *sqlite3) Features() *DialectFeatures {
	return &DialectFeatures{
		AutoincrMode: IncrAutoincrMode,
	}
}
func (db *sqlite3) SetQuotePolicy(quotePolicy QuotePolicy) {
	switch quotePolicy {
	case QuotePolicyNone:
		var q = sqlite3Quoter
		q.IsReserved = schemasvr.AlwaysNoReserve
		db.quoter = q
	case QuotePolicyReserved:
		var q = sqlite3Quoter
		q.IsReserved = db.IsReserved
		db.quoter = q
	case QuotePolicyAlways:
		fallthrough
	default:
		db.quoter = sqlite3Quoter
	}
}
func (db *sqlite3) SQLType(c *schemasvr.Column) string {
	switch t := c.SQLType.Name; t {
	case schemasvr.Bool:
		if c.Default == "true" {
			c.Default = "1"
		} else if c.Default == "false" {
			c.Default = "0"
		}
		return schemasvr.Integer
	case schemasvr.Date, schemasvr.DateTime, schemasvr.TimeStamp, schemasvr.Time:
		return schemasvr.DateTime
	case schemasvr.TimeStampz:
		return schemasvr.Text
	case schemasvr.Char, schemasvr.Varchar, schemasvr.NVarchar, schemasvr.TinyText,
		schemasvr.Text, schemasvr.MediumText, schemasvr.LongText, schemasvr.Json:
		return schemasvr.Text
	case schemasvr.Bit, schemasvr.TinyInt, schemasvr.UnsignedTinyInt, schemasvr.SmallInt,
		schemasvr.UnsignedSmallInt, schemasvr.MediumInt, schemasvr.Int, schemasvr.UnsignedInt,
		schemasvr.BigInt, schemasvr.UnsignedBigInt, schemasvr.Integer:
		return schemasvr.Integer
	case schemasvr.Float, schemasvr.Double, schemasvr.Real:
		return schemasvr.Real
	case schemasvr.Decimal, schemasvr.Numeric:
		return schemasvr.Numeric
	case schemasvr.TinyBlob, schemasvr.Blob, schemasvr.MediumBlob, schemasvr.LongBlob, schemasvr.Bytea, schemasvr.Binary, schemasvr.VarBinary:
		return schemasvr.Blob
	case schemasvr.Serial, schemasvr.BigSerial:
		c.IsPrimaryKey = true
		c.IsAutoIncrement = true
		c.Nullable = false
		return schemasvr.Integer
	default:
		return t
	}
}
func (db *sqlite3) ColumnTypeKind(t string) int {
	switch strings.ToUpper(t) {
	case "DATETIME":
		return schemasvr.TIME_TYPE
	case "TEXT":
		return schemasvr.TEXT_TYPE
	case "INTEGER", "REAL", "NUMERIC", "DECIMAL":
		return schemasvr.NUMERIC_TYPE
	case "BLOB":
		return schemasvr.BLOB_TYPE
	default:
		return schemasvr.UNKNOW_TYPE
	}
}
func (db *sqlite3) IsReserved(name string) bool {
	_, ok := sqlite3ReservedWords[strings.ToUpper(name)]
	return ok
}
func (db *sqlite3) AutoIncrStr() string {
	return "AUTOINCREMENT"
}
func (db *sqlite3) IndexCheckSQL(tableName, idxName string) (string, []interface{}) {
	args := []interface{}{idxName}
	return "SELECT name FROM sqlite_master WHERE type='index' and name = ?", args
}
func (db *sqlite3) IsTableExist(queryer core.Queryer, ctx context.Context, tableName string) (bool, error) {
	return db.HasRecords(queryer, ctx, "SELECT name FROM sqlite_master WHERE type='table' and name = ?", tableName)
}
func (db *sqlite3) DropIndexSQL(tableName string, index *schemasvr.Index) string {
	// var unique string
	idxName := index.Name
	if !strings.HasPrefix(idxName, "UQE_") &&
		!strings.HasPrefix(idxName, "IDX_") {
		if index.Type == schemasvr.UniqueType {
			idxName = fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
		} else {
			idxName = fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
		}
	}
	return fmt.Sprintf("DROP INDEX %v", db.Quoter().Quote(idxName))
}
func (db *sqlite3) ForUpdateSQL(query string) string {
	return query
}
func (db *sqlite3) IsColumnExist(queryer core.Queryer, ctx context.Context, tableName, colName string) (bool, error) {
	query := "SELECT * FROM " + tableName + " LIMIT 0"
	rows, err := queryer.QueryContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}
	for _, col := range cols {
		if strings.EqualFold(col, colName) {
			return true, nil
		}
	}
	return false, nil
}

// splitColStr splits a sqlite col strings as fields
func splitColStr(colStr string) []string {
	colStr = strings.TrimSpace(colStr)
	var results = make([]string, 0, 10)
	var lastIdx int
	var hasC, hasQuote bool
	for i, c := range colStr {
		if c == ' ' && !hasQuote {
			if hasC {
				results = append(results, colStr[lastIdx:i])
				hasC = false
			}
		} else {
			if c == '\'' {
				hasQuote = !hasQuote
			}
			if !hasC {
				lastIdx = i
			}
			hasC = true
			if i == len(colStr)-1 {
				results = append(results, colStr[lastIdx:i+1])
			}
		}
	}
	return results
}
func parseString(colStr string) (*schemasvr.Column, error) {
	fields := splitColStr(colStr)
	col := new(schemasvr.Column)
	col.Indexes = make(map[string]int)
	col.Nullable = true
	col.DefaultIsEmpty = true
	for idx, field := range fields {
		if idx == 0 {
			col.Name = strings.Trim(strings.Trim(field, "`[] "), `"`)
			continue
		} else if idx == 1 {
			col.SQLType = schemasvr.SQLType{Name: field, DefaultLength: 0, DefaultLength2: 0}
			continue
		}
		switch field {
		case "PRIMARY":
			col.IsPrimaryKey = true
		case "AUTOINCREMENT":
			col.IsAutoIncrement = true
		case "NULL":
			if fields[idx-1] == "NOT" {
				col.Nullable = false
			} else {
				col.Nullable = true
			}
		case "DEFAULT":
			col.Default = fields[idx+1]
			col.DefaultIsEmpty = false
		}
	}
	return col, nil
}
func (db *sqlite3) GetColumns(queryer core.Queryer, ctx context.Context, tableName string) ([]string, map[string]*schemasvr.Column, error) {
	args := []interface{}{tableName}
	s := "SELECT sql FROM sqlite_master WHERE type='table' and name = ?"
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var name string
	if rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, nil, err
		}
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}
	if name == "" {
		return nil, nil, errors.New("no table named " + tableName)
	}
	nStart := strings.Index(name, "(")
	nEnd := strings.LastIndex(name, ")")
	reg := regexp.MustCompile(`[^\(,\)]*(\([^\(]*\))?`)
	colCreates := reg.FindAllString(name[nStart+1:nEnd], -1)
	cols := make(map[string]*schemasvr.Column)
	colSeq := make([]string, 0)
	for _, colStr := range colCreates {
		reg = regexp.MustCompile(`,\s`)
		colStr = reg.ReplaceAllString(colStr, ",")
		if strings.HasPrefix(strings.TrimSpace(colStr), "PRIMARY KEY") {
			parts := strings.Split(strings.TrimSpace(colStr), "(")
			if len(parts) == 2 {
				pkCols := strings.Split(strings.TrimRight(strings.TrimSpace(parts[1]), ")"), ",")
				for _, pk := range pkCols {
					if col, ok := cols[strings.Trim(strings.TrimSpace(pk), "`")]; ok {
						col.IsPrimaryKey = true
					}
				}
			}
			continue
		}
		col, err := parseString(colStr)
		if err != nil {
			return colSeq, cols, err
		}
		cols[col.Name] = col
		colSeq = append(colSeq, col.Name)
	}
	return colSeq, cols, nil
}
func (db *sqlite3) GetTables(queryer core.Queryer, ctx context.Context) ([]*schemasvr.Table, error) {
	args := []interface{}{}
	s := "SELECT name FROM sqlite_master WHERE type='table'"
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := make([]*schemasvr.Table, 0)
	for rows.Next() {
		table := schemasvr.NewEmptyTable()
		err = rows.Scan(&table.Name)
		if err != nil {
			return nil, err
		}
		if table.Name == "sqlite_sequence" {
			continue
		}
		tables = append(tables, table)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tables, nil
}
func (db *sqlite3) GetIndexes(queryer core.Queryer, ctx context.Context, tableName string) (map[string]*schemasvr.Index, error) {
	args := []interface{}{tableName}
	s := "SELECT sql FROM sqlite_master WHERE type='index' and tbl_name = ?"
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	indexes := make(map[string]*schemasvr.Index)
	for rows.Next() {
		var tmpSQL sql.NullString
		err = rows.Scan(&tmpSQL)
		if err != nil {
			return nil, err
		}
		if !tmpSQL.Valid {
			continue
		}
		sql := tmpSQL.String
		index := new(schemasvr.Index)
		nNStart := strings.Index(sql, "INDEX")
		nNEnd := strings.Index(sql, "ON")
		if nNStart == -1 || nNEnd == -1 {
			continue
		}
		indexName := strings.Trim(strings.TrimSpace(sql[nNStart+6:nNEnd]), "`[]'\"")
		var isRegular bool
		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			index.Name = indexName[5+len(tableName):]
			isRegular = true
		} else {
			index.Name = indexName
		}
		if strings.HasPrefix(sql, "CREATE UNIQUE INDEX") {
			index.Type = schemasvr.UniqueType
		} else {
			index.Type = schemasvr.IndexType
		}
		nStart := strings.Index(sql, "(")
		nEnd := strings.Index(sql, ")")
		colIndexes := strings.Split(sql[nStart+1:nEnd], ",")
		index.Cols = make([]string, 0)
		for _, col := range colIndexes {
			index.Cols = append(index.Cols, strings.Trim(col, "` []"))
		}
		index.IsRegular = isRegular
		indexes[index.Name] = index
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return indexes, nil
}
func (db *sqlite3) Filters() []Filter {
	return []Filter{}
}

type sqlite3Driver struct {
	baseDriver
}

func (p *sqlite3Driver) Features() *DriverFeatures {
	return &DriverFeatures{
		SupportReturnInsertedID: true,
	}
}
func (p *sqlite3Driver) Parse(driverName, dataSourceName string) (*URI, error) {
	if strings.Contains(dataSourceName, "?") {
		dataSourceName = dataSourceName[:strings.Index(dataSourceName, "?")]
	}
	return &URI{DBType: schemasvr.SQLITE, DBName: dataSourceName}, nil
}
func (p *sqlite3Driver) GenScanResult(colType string) (interface{}, error) {
	switch colType {
	case "TEXT":
		var s sql.NullString
		return &s, nil
	case "INTEGER":
		var s sql.NullInt64
		return &s, nil
	case "DATETIME":
		var s sql.NullTime
		return &s, nil
	case "REAL":
		var s sql.NullFloat64
		return &s, nil
	case "NUMERIC", "DECIMAL":
		var s sql.NullString
		return &s, nil
	case "BLOB":
		var s sql.RawBytes
		return &s, nil
	default:
		var r sql.NullString
		return &r, nil
	}
}
