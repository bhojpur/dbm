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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bhojpur/dbm/pkg/orm/core"
	schemasvr "github.com/bhojpur/dbm/pkg/orm/schema"
)

// from http://www.postgresql.org/docs/current/static/sql-keywords-appendix.html
var (
	postgresReservedWords = map[string]bool{
		"A":                                true,
		"ABORT":                            true,
		"ABS":                              true,
		"ABSENT":                           true,
		"ABSOLUTE":                         true,
		"ACCESS":                           true,
		"ACCORDING":                        true,
		"ACTION":                           true,
		"ADA":                              true,
		"ADD":                              true,
		"ADMIN":                            true,
		"AFTER":                            true,
		"AGGREGATE":                        true,
		"ALL":                              true,
		"ALLOCATE":                         true,
		"ALSO":                             true,
		"ALTER":                            true,
		"ALWAYS":                           true,
		"ANALYSE":                          true,
		"ANALYZE":                          true,
		"AND":                              true,
		"ANY":                              true,
		"ARE":                              true,
		"ARRAY":                            true,
		"ARRAY_AGG":                        true,
		"ARRAY_MAX_CARDINALITY":            true,
		"AS":                               true,
		"ASC":                              true,
		"ASENSITIVE":                       true,
		"ASSERTION":                        true,
		"ASSIGNMENT":                       true,
		"ASYMMETRIC":                       true,
		"AT":                               true,
		"ATOMIC":                           true,
		"ATTRIBUTE":                        true,
		"ATTRIBUTES":                       true,
		"AUTHORIZATION":                    true,
		"AVG":                              true,
		"BACKWARD":                         true,
		"BASE64":                           true,
		"BEFORE":                           true,
		"BEGIN":                            true,
		"BEGIN_FRAME":                      true,
		"BEGIN_PARTITION":                  true,
		"BERNOULLI":                        true,
		"BETWEEN":                          true,
		"BIGINT":                           true,
		"BINARY":                           true,
		"BIT":                              true,
		"BIT_LENGTH":                       true,
		"BLOB":                             true,
		"BLOCKED":                          true,
		"BOM":                              true,
		"BOOLEAN":                          true,
		"BOTH":                             true,
		"BREADTH":                          true,
		"BY":                               true,
		"C":                                true,
		"CACHE":                            true,
		"CALL":                             true,
		"CALLED":                           true,
		"CARDINALITY":                      true,
		"CASCADE":                          true,
		"CASCADED":                         true,
		"CASE":                             true,
		"CAST":                             true,
		"CATALOG":                          true,
		"CATALOG_NAME":                     true,
		"CEIL":                             true,
		"CEILING":                          true,
		"CHAIN":                            true,
		"CHAR":                             true,
		"CHARACTER":                        true,
		"CHARACTERISTICS":                  true,
		"CHARACTERS":                       true,
		"CHARACTER_LENGTH":                 true,
		"CHARACTER_SET_CATALOG":            true,
		"CHARACTER_SET_NAME":               true,
		"CHARACTER_SET_SCHEMA":             true,
		"CHAR_LENGTH":                      true,
		"CHECK":                            true,
		"CHECKPOINT":                       true,
		"CLASS":                            true,
		"CLASS_ORIGIN":                     true,
		"CLOB":                             true,
		"CLOSE":                            true,
		"CLUSTER":                          true,
		"COALESCE":                         true,
		"COBOL":                            true,
		"COLLATE":                          true,
		"COLLATION":                        true,
		"COLLATION_CATALOG":                true,
		"COLLATION_NAME":                   true,
		"COLLATION_SCHEMA":                 true,
		"COLLECT":                          true,
		"COLUMN":                           true,
		"COLUMNS":                          true,
		"COLUMN_NAME":                      true,
		"COMMAND_FUNCTION":                 true,
		"COMMAND_FUNCTION_CODE":            true,
		"COMMENT":                          true,
		"COMMENTS":                         true,
		"COMMIT":                           true,
		"COMMITTED":                        true,
		"CONCURRENTLY":                     true,
		"CONDITION":                        true,
		"CONDITION_NUMBER":                 true,
		"CONFIGURATION":                    true,
		"CONNECT":                          true,
		"CONNECTION":                       true,
		"CONNECTION_NAME":                  true,
		"CONSTRAINT":                       true,
		"CONSTRAINTS":                      true,
		"CONSTRAINT_CATALOG":               true,
		"CONSTRAINT_NAME":                  true,
		"CONSTRAINT_SCHEMA":                true,
		"CONSTRUCTOR":                      true,
		"CONTAINS":                         true,
		"CONTENT":                          true,
		"CONTINUE":                         true,
		"CONTROL":                          true,
		"CONVERSION":                       true,
		"CONVERT":                          true,
		"COPY":                             true,
		"CORR":                             true,
		"CORRESPONDING":                    true,
		"COST":                             true,
		"COUNT":                            true,
		"COVAR_POP":                        true,
		"COVAR_SAMP":                       true,
		"CREATE":                           true,
		"CROSS":                            true,
		"CSV":                              true,
		"CUBE":                             true,
		"CUME_DIST":                        true,
		"CURRENT":                          true,
		"CURRENT_CATALOG":                  true,
		"CURRENT_DATE":                     true,
		"CURRENT_DEFAULT_TRANSFORM_GROUP":  true,
		"CURRENT_PATH":                     true,
		"CURRENT_ROLE":                     true,
		"CURRENT_ROW":                      true,
		"CURRENT_SCHEMA":                   true,
		"CURRENT_TIME":                     true,
		"CURRENT_TIMESTAMP":                true,
		"CURRENT_TRANSFORM_GROUP_FOR_TYPE": true,
		"CURRENT_USER":                     true,
		"CURSOR":                           true,
		"CURSOR_NAME":                      true,
		"CYCLE":                            true,
		"DATA":                             true,
		"DATABASE":                         true,
		"DATALINK":                         true,
		"DATE":                             true,
		"DATETIME_INTERVAL_CODE":           true,
		"DATETIME_INTERVAL_PRECISION":      true,
		"DAY":                              true,
		"DB":                               true,
		"DEALLOCATE":                       true,
		"DEC":                              true,
		"DECIMAL":                          true,
		"DECLARE":                          true,
		"DEFAULT":                          true,
		"DEFAULTS":                         true,
		"DEFERRABLE":                       true,
		"DEFERRED":                         true,
		"DEFINED":                          true,
		"DEFINER":                          true,
		"DEGREE":                           true,
		"DELETE":                           true,
		"DELIMITER":                        true,
		"DELIMITERS":                       true,
		"DENSE_RANK":                       true,
		"DEPTH":                            true,
		"DEREF":                            true,
		"DERIVED":                          true,
		"DESC":                             true,
		"DESCRIBE":                         true,
		"DESCRIPTOR":                       true,
		"DETERMINISTIC":                    true,
		"DIAGNOSTICS":                      true,
		"DICTIONARY":                       true,
		"DISABLE":                          true,
		"DISCARD":                          true,
		"DISCONNECT":                       true,
		"DISPATCH":                         true,
		"DISTINCT":                         true,
		"DLNEWCOPY":                        true,
		"DLPREVIOUSCOPY":                   true,
		"DLURLCOMPLETE":                    true,
		"DLURLCOMPLETEONLY":                true,
		"DLURLCOMPLETEWRITE":               true,
		"DLURLPATH":                        true,
		"DLURLPATHONLY":                    true,
		"DLURLPATHWRITE":                   true,
		"DLURLSCHEME":                      true,
		"DLURLSERVER":                      true,
		"DLVALUE":                          true,
		"DO":                               true,
		"DOCUMENT":                         true,
		"DOMAIN":                           true,
		"DOUBLE":                           true,
		"DROP":                             true,
		"DYNAMIC":                          true,
		"DYNAMIC_FUNCTION":                 true,
		"DYNAMIC_FUNCTION_CODE":            true,
		"EACH":                             true,
		"ELEMENT":                          true,
		"ELSE":                             true,
		"EMPTY":                            true,
		"ENABLE":                           true,
		"ENCODING":                         true,
		"ENCRYPTED":                        true,
		"END":                              true,
		"END-EXEC":                         true,
		"END_FRAME":                        true,
		"END_PARTITION":                    true,
		"ENFORCED":                         true,
		"ENUM":                             true,
		"EQUALS":                           true,
		"ESCAPE":                           true,
		"EVENT":                            true,
		"EVERY":                            true,
		"EXCEPT":                           true,
		"EXCEPTION":                        true,
		"EXCLUDE":                          true,
		"EXCLUDING":                        true,
		"EXCLUSIVE":                        true,
		"EXEC":                             true,
		"EXECUTE":                          true,
		"EXISTS":                           true,
		"EXP":                              true,
		"EXPLAIN":                          true,
		"EXPRESSION":                       true,
		"EXTENSION":                        true,
		"EXTERNAL":                         true,
		"EXTRACT":                          true,
		"FALSE":                            true,
		"FAMILY":                           true,
		"FETCH":                            true,
		"FILE":                             true,
		"FILTER":                           true,
		"FINAL":                            true,
		"FIRST":                            true,
		"FIRST_VALUE":                      true,
		"FLAG":                             true,
		"FLOAT":                            true,
		"FLOOR":                            true,
		"FOLLOWING":                        true,
		"FOR":                              true,
		"FORCE":                            true,
		"FOREIGN":                          true,
		"FORTRAN":                          true,
		"FORWARD":                          true,
		"FOUND":                            true,
		"FRAME_ROW":                        true,
		"FREE":                             true,
		"FREEZE":                           true,
		"FROM":                             true,
		"FS":                               true,
		"FULL":                             true,
		"FUNCTION":                         true,
		"FUNCTIONS":                        true,
		"FUSION":                           true,
		"G":                                true,
		"GENERAL":                          true,
		"GENERATED":                        true,
		"GET":                              true,
		"GLOBAL":                           true,
		"GO":                               true,
		"GOTO":                             true,
		"GRANT":                            true,
		"GRANTED":                          true,
		"GREATEST":                         true,
		"GROUP":                            true,
		"GROUPING":                         true,
		"GROUPS":                           true,
		"HANDLER":                          true,
		"HAVING":                           true,
		"HEADER":                           true,
		"HEX":                              true,
		"HIERARCHY":                        true,
		"HOLD":                             true,
		"HOUR":                             true,
		"ID":                               true,
		"IDENTITY":                         true,
		"IF":                               true,
		"IGNORE":                           true,
		"ILIKE":                            true,
		"IMMEDIATE":                        true,
		"IMMEDIATELY":                      true,
		"IMMUTABLE":                        true,
		"IMPLEMENTATION":                   true,
		"IMPLICIT":                         true,
		"IMPORT":                           true,
		"IN":                               true,
		"INCLUDING":                        true,
		"INCREMENT":                        true,
		"INDENT":                           true,
		"INDEX":                            true,
		"INDEXES":                          true,
		"INDICATOR":                        true,
		"INHERIT":                          true,
		"INHERITS":                         true,
		"INITIALLY":                        true,
		"INLINE":                           true,
		"INNER":                            true,
		"INOUT":                            true,
		"INPUT":                            true,
		"INSENSITIVE":                      true,
		"INSERT":                           true,
		"INSTANCE":                         true,
		"INSTANTIABLE":                     true,
		"INSTEAD":                          true,
		"INT":                              true,
		"INTEGER":                          true,
		"INTEGRITY":                        true,
		"INTERSECT":                        true,
		"INTERSECTION":                     true,
		"INTERVAL":                         true,
		"INTO":                             true,
		"INVOKER":                          true,
		"IS":                               true,
		"ISNULL":                           true,
		"ISOLATION":                        true,
		"JOIN":                             true,
		"K":                                true,
		"KEY":                              true,
		"KEY_MEMBER":                       true,
		"KEY_TYPE":                         true,
		"LABEL":                            true,
		"LAG":                              true,
		"LANGUAGE":                         true,
		"LARGE":                            true,
		"LAST":                             true,
		"LAST_VALUE":                       true,
		"LATERAL":                          true,
		"LC_COLLATE":                       true,
		"LC_CTYPE":                         true,
		"LEAD":                             true,
		"LEADING":                          true,
		"LEAKPROOF":                        true,
		"LEAST":                            true,
		"LEFT":                             true,
		"LENGTH":                           true,
		"LEVEL":                            true,
		"LIBRARY":                          true,
		"LIKE":                             true,
		"LIKE_REGEX":                       true,
		"LIMIT":                            true,
		"LINK":                             true,
		"LISTEN":                           true,
		"LN":                               true,
		"LOAD":                             true,
		"LOCAL":                            true,
		"LOCALTIME":                        true,
		"LOCALTIMESTAMP":                   true,
		"LOCATION":                         true,
		"LOCATOR":                          true,
		"LOCK":                             true,
		"LOWER":                            true,
		"M":                                true,
		"MAP":                              true,
		"MAPPING":                          true,
		"MATCH":                            true,
		"MATCHED":                          true,
		"MATERIALIZED":                     true,
		"MAX":                              true,
		"MAXVALUE":                         true,
		"MAX_CARDINALITY":                  true,
		"MEMBER":                           true,
		"MERGE":                            true,
		"MESSAGE_LENGTH":                   true,
		"MESSAGE_OCTET_LENGTH":             true,
		"MESSAGE_TEXT":                     true,
		"METHOD":                           true,
		"MIN":                              true,
		"MINUTE":                           true,
		"MINVALUE":                         true,
		"MOD":                              true,
		"MODE":                             true,
		"MODIFIES":                         true,
		"MODULE":                           true,
		"MONTH":                            true,
		"MORE":                             true,
		"MOVE":                             true,
		"MULTISET":                         true,
		"MUMPS":                            true,
		"NAME":                             true,
		"NAMES":                            true,
		"NAMESPACE":                        true,
		"NATIONAL":                         true,
		"NATURAL":                          true,
		"NCHAR":                            true,
		"NCLOB":                            true,
		"NESTING":                          true,
		"NEW":                              true,
		"NEXT":                             true,
		"NFC":                              true,
		"NFD":                              true,
		"NFKC":                             true,
		"NFKD":                             true,
		"NIL":                              true,
		"NO":                               true,
		"NONE":                             true,
		"NORMALIZE":                        true,
		"NORMALIZED":                       true,
		"NOT":                              true,
		"NOTHING":                          true,
		"NOTIFY":                           true,
		"NOTNULL":                          true,
		"NOWAIT":                           true,
		"NTH_VALUE":                        true,
		"NTILE":                            true,
		"NULL":                             true,
		"NULLABLE":                         true,
		"NULLIF":                           true,
		"NULLS":                            true,
		"NUMBER":                           true,
		"NUMERIC":                          true,
		"OBJECT":                           true,
		"OCCURRENCES_REGEX":                true,
		"OCTETS":                           true,
		"OCTET_LENGTH":                     true,
		"OF":                               true,
		"OFF":                              true,
		"OFFSET":                           true,
		"OIDS":                             true,
		"OLD":                              true,
		"ON":                               true,
		"ONLY":                             true,
		"OPEN":                             true,
		"OPERATOR":                         true,
		"OPTION":                           true,
		"OPTIONS":                          true,
		"OR":                               true,
		"ORDER":                            true,
		"ORDERING":                         true,
		"ORDINALITY":                       true,
		"OTHERS":                           true,
		"OUT":                              true,
		"OUTER":                            true,
		"OUTPUT":                           true,
		"OVER":                             true,
		"OVERLAPS":                         true,
		"OVERLAY":                          true,
		"OVERRIDING":                       true,
		"OWNED":                            true,
		"OWNER":                            true,
		"P":                                true,
		"PAD":                              true,
		"PARAMETER":                        true,
		"PARAMETER_MODE":                   true,
		"PARAMETER_NAME":                   true,
		"PARAMETER_ORDINAL_POSITION":       true,
		"PARAMETER_SPECIFIC_CATALOG":       true,
		"PARAMETER_SPECIFIC_NAME":          true,
		"PARAMETER_SPECIFIC_SCHEMA":        true,
		"PARSER":                           true,
		"PARTIAL":                          true,
		"PARTITION":                        true,
		"PASCAL":                           true,
		"PASSING":                          true,
		"PASSTHROUGH":                      true,
		"PASSWORD":                         true,
		"PATH":                             true,
		"PERCENT":                          true,
		"PERCENTILE_CONT":                  true,
		"PERCENTILE_DISC":                  true,
		"PERCENT_RANK":                     true,
		"PERIOD":                           true,
		"PERMISSION":                       true,
		"PLACING":                          true,
		"PLANS":                            true,
		"PLI":                              true,
		"PORTION":                          true,
		"POSITION":                         true,
		"POSITION_REGEX":                   true,
		"POWER":                            true,
		"PRECEDES":                         true,
		"PRECEDING":                        true,
		"PRECISION":                        true,
		"PREPARE":                          true,
		"PREPARED":                         true,
		"PRESERVE":                         true,
		"PRIMARY":                          true,
		"PRIOR":                            true,
		"PRIVILEGES":                       true,
		"PROCEDURAL":                       true,
		"PROCEDURE":                        true,
		"PROGRAM":                          true,
		"PUBLIC":                           true,
		"QUOTE":                            true,
		"RANGE":                            true,
		"RANK":                             true,
		"READ":                             true,
		"READS":                            true,
		"REAL":                             true,
		"REASSIGN":                         true,
		"RECHECK":                          true,
		"RECOVERY":                         true,
		"RECURSIVE":                        true,
		"REF":                              true,
		"REFERENCES":                       true,
		"REFERENCING":                      true,
		"REFRESH":                          true,
		"REGR_AVGX":                        true,
		"REGR_AVGY":                        true,
		"REGR_COUNT":                       true,
		"REGR_INTERCEPT":                   true,
		"REGR_R2":                          true,
		"REGR_SLOPE":                       true,
		"REGR_SXX":                         true,
		"REGR_SXY":                         true,
		"REGR_SYY":                         true,
		"REINDEX":                          true,
		"RELATIVE":                         true,
		"RELEASE":                          true,
		"RENAME":                           true,
		"REPEATABLE":                       true,
		"REPLACE":                          true,
		"REPLICA":                          true,
		"REQUIRING":                        true,
		"RESET":                            true,
		"RESPECT":                          true,
		"RESTART":                          true,
		"RESTORE":                          true,
		"RESTRICT":                         true,
		"RESULT":                           true,
		"RETURN":                           true,
		"RETURNED_CARDINALITY":             true,
		"RETURNED_LENGTH":                  true,
		"RETURNED_OCTET_LENGTH":            true,
		"RETURNED_SQLSTATE":                true,
		"RETURNING":                        true,
		"RETURNS":                          true,
		"REVOKE":                           true,
		"RIGHT":                            true,
		"ROLE":                             true,
		"ROLLBACK":                         true,
		"ROLLUP":                           true,
		"ROUTINE":                          true,
		"ROUTINE_CATALOG":                  true,
		"ROUTINE_NAME":                     true,
		"ROUTINE_SCHEMA":                   true,
		"ROW":                              true,
		"ROWS":                             true,
		"ROW_COUNT":                        true,
		"ROW_NUMBER":                       true,
		"RULE":                             true,
		"SAVEPOINT":                        true,
		"SCALE":                            true,
		"SCHEMA":                           true,
		"SCHEMA_NAME":                      true,
		"SCOPE":                            true,
		"SCOPE_CATALOG":                    true,
		"SCOPE_NAME":                       true,
		"SCOPE_SCHEMA":                     true,
		"SCROLL":                           true,
		"SEARCH":                           true,
		"SECOND":                           true,
		"SECTION":                          true,
		"SECURITY":                         true,
		"SELECT":                           true,
		"SELECTIVE":                        true,
		"SELF":                             true,
		"SENSITIVE":                        true,
		"SEQUENCE":                         true,
		"SEQUENCES":                        true,
		"SERIALIZABLE":                     true,
		"SERVER":                           true,
		"SERVER_NAME":                      true,
		"SESSION":                          true,
		"SESSION_USER":                     true,
		"SET":                              true,
		"SETOF":                            true,
		"SETS":                             true,
		"SHARE":                            true,
		"SHOW":                             true,
		"SIMILAR":                          true,
		"SIMPLE":                           true,
		"SIZE":                             true,
		"SMALLINT":                         true,
		"SNAPSHOT":                         true,
		"SOME":                             true,
		"SOURCE":                           true,
		"SPACE":                            true,
		"SPECIFIC":                         true,
		"SPECIFICTYPE":                     true,
		"SPECIFIC_NAME":                    true,
		"SQL":                              true,
		"SQLCODE":                          true,
		"SQLERROR":                         true,
		"SQLEXCEPTION":                     true,
		"SQLSTATE":                         true,
		"SQLWARNING":                       true,
		"SQRT":                             true,
		"STABLE":                           true,
		"STANDALONE":                       true,
		"START":                            true,
		"STATE":                            true,
		"STATEMENT":                        true,
		"STATIC":                           true,
		"STATISTICS":                       true,
		"STDDEV_POP":                       true,
		"STDDEV_SAMP":                      true,
		"STDIN":                            true,
		"STDOUT":                           true,
		"STORAGE":                          true,
		"STRICT":                           true,
		"STRIP":                            true,
		"STRUCTURE":                        true,
		"STYLE":                            true,
		"SUBCLASS_ORIGIN":                  true,
		"SUBMULTISET":                      true,
		"SUBSTRING":                        true,
		"SUBSTRING_REGEX":                  true,
		"SUCCEEDS":                         true,
		"SUM":                              true,
		"SYMMETRIC":                        true,
		"SYSID":                            true,
		"SYSTEM":                           true,
		"SYSTEM_TIME":                      true,
		"SYSTEM_USER":                      true,
		"T":                                true,
		"TABLE":                            true,
		"TABLES":                           true,
		"TABLESAMPLE":                      true,
		"TABLESPACE":                       true,
		"TABLE_NAME":                       true,
		"TEMP":                             true,
		"TEMPLATE":                         true,
		"TEMPORARY":                        true,
		"TEXT":                             true,
		"THEN":                             true,
		"TIES":                             true,
		"TIME":                             true,
		"TIMESTAMP":                        true,
		"TIMEZONE_HOUR":                    true,
		"TIMEZONE_MINUTE":                  true,
		"TO":                               true,
		"TOKEN":                            true,
		"TOP_LEVEL_COUNT":                  true,
		"TRAILING":                         true,
		"TRANSACTION":                      true,
		"TRANSACTIONS_COMMITTED":           true,
		"TRANSACTIONS_ROLLED_BACK":         true,
		"TRANSACTION_ACTIVE":               true,
		"TRANSFORM":                        true,
		"TRANSFORMS":                       true,
		"TRANSLATE":                        true,
		"TRANSLATE_REGEX":                  true,
		"TRANSLATION":                      true,
		"TREAT":                            true,
		"TRIGGER":                          true,
		"TRIGGER_CATALOG":                  true,
		"TRIGGER_NAME":                     true,
		"TRIGGER_SCHEMA":                   true,
		"TRIM":                             true,
		"TRIM_ARRAY":                       true,
		"TRUE":                             true,
		"TRUNCATE":                         true,
		"TRUSTED":                          true,
		"TYPE":                             true,
		"TYPES":                            true,
		"UESCAPE":                          true,
		"UNBOUNDED":                        true,
		"UNCOMMITTED":                      true,
		"UNDER":                            true,
		"UNENCRYPTED":                      true,
		"UNION":                            true,
		"UNIQUE":                           true,
		"UNKNOWN":                          true,
		"UNLINK":                           true,
		"UNLISTEN":                         true,
		"UNLOGGED":                         true,
		"UNNAMED":                          true,
		"UNNEST":                           true,
		"UNTIL":                            true,
		"UNTYPED":                          true,
		"UPDATE":                           true,
		"UPPER":                            true,
		"URI":                              true,
		"USAGE":                            true,
		"USER":                             true,
		"USER_DEFINED_TYPE_CATALOG":        true,
		"USER_DEFINED_TYPE_CODE":           true,
		"USER_DEFINED_TYPE_NAME":           true,
		"USER_DEFINED_TYPE_SCHEMA":         true,
		"USING":                            true,
		"VACUUM":                           true,
		"VALID":                            true,
		"VALIDATE":                         true,
		"VALIDATOR":                        true,
		"VALUE":                            true,
		"VALUES":                           true,
		"VALUE_OF":                         true,
		"VARBINARY":                        true,
		"VARCHAR":                          true,
		"VARIADIC":                         true,
		"VARYING":                          true,
		"VAR_POP":                          true,
		"VAR_SAMP":                         true,
		"VERBOSE":                          true,
		"VERSION":                          true,
		"VERSIONING":                       true,
		"VIEW":                             true,
		"VOLATILE":                         true,
		"WHEN":                             true,
		"WHENEVER":                         true,
		"WHERE":                            true,
		"WHITESPACE":                       true,
		"WIDTH_BUCKET":                     true,
		"WINDOW":                           true,
		"WITH":                             true,
		"WITHIN":                           true,
		"WITHOUT":                          true,
		"WORK":                             true,
		"WRAPPER":                          true,
		"WRITE":                            true,
		"XML":                              true,
		"XMLAGG":                           true,
		"XMLATTRIBUTES":                    true,
		"XMLBINARY":                        true,
		"XMLCAST":                          true,
		"XMLCOMMENT":                       true,
		"XMLCONCAT":                        true,
		"XMLDECLARATION":                   true,
		"XMLDOCUMENT":                      true,
		"XMLELEMENT":                       true,
		"XMLEXISTS":                        true,
		"XMLFOREST":                        true,
		"XMLITERATE":                       true,
		"XMLNAMESPACES":                    true,
		"XMLPARSE":                         true,
		"XMLPI":                            true,
		"XMLQUERY":                         true,
		"XMLROOT":                          true,
		"XMLSCHEMA":                        true,
		"XMLSERIALIZE":                     true,
		"XMLTABLE":                         true,
		"XMLTEXT":                          true,
		"XMLVALIDATE":                      true,
		"YEAR":                             true,
		"YES":                              true,
		"ZONE":                             true,
	}
	postgresQuoter = schemasvr.Quoter{
		Prefix:     '"',
		Suffix:     '"',
		IsReserved: schemasvr.AlwaysReserve,
	}
)
var (
	// DefaultPostgresSchema default postgres schema
	DefaultPostgresSchema = "public"
	postgresColAliases    = map[string]string{
		"numeric": "decimal",
	}
)

type postgres struct {
	Base
}

// Alias returns a alias of column
func (db *postgres) Alias(col string) string {
	v, ok := postgresColAliases[strings.ToLower(col)]
	if ok {
		return v
	}
	return col
}
func (db *postgres) Init(uri *URI) error {
	db.quoter = postgresQuoter
	return db.Base.Init(db, uri)
}
func (db *postgres) Version(ctx context.Context, queryer core.Queryer) (*schemasvr.Version, error) {
	rows, err := queryer.QueryContext(ctx, "SELECT version()")
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
	// Postgres: 9.5.22 on x86_64-pc-linux-gnu (Debian 9.5.22-1.pgdg90+1), compiled by gcc (Debian 6.3.0-18+deb9u1) 6.3.0 20170516, 64-bit
	// CockroachDB CCL v19.2.4 (x86_64-unknown-linux-gnu, built
	if strings.HasPrefix(version, "CockroachDB") {
		versions := strings.Split(strings.TrimPrefix(version, "CockroachDB CCL "), " ")
		return &schemasvr.Version{
			Number:  strings.TrimPrefix(versions[0], "v"),
			Edition: "CockroachDB",
		}, nil
	} else if strings.HasPrefix(version, "PostgreSQL") {
		versions := strings.Split(strings.TrimPrefix(version, "PostgreSQL "), " on ")
		return &schemasvr.Version{
			Number:  versions[0],
			Level:   versions[1],
			Edition: "PostgreSQL",
		}, nil
	}
	return nil, errors.New("unknow database version")
}
func (db *postgres) getSchema() string {
	if db.uri.Schema != "" {
		return db.uri.Schema
	}
	return DefaultPostgresSchema
}
func (db *postgres) needQuote(name string) bool {
	if db.IsReserved(name) {
		return true
	}
	for _, c := range name {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}
func (db *postgres) SetQuotePolicy(quotePolicy QuotePolicy) {
	switch quotePolicy {
	case QuotePolicyNone:
		var q = postgresQuoter
		q.IsReserved = schemasvr.AlwaysNoReserve
		db.quoter = q
	case QuotePolicyReserved:
		var q = postgresQuoter
		q.IsReserved = db.needQuote
		db.quoter = q
	case QuotePolicyAlways:
		fallthrough
	default:
		db.quoter = postgresQuoter
	}
}
func (db *postgres) SQLType(c *schemasvr.Column) string {
	var res string
	switch t := c.SQLType.Name; t {
	case schemasvr.TinyInt, schemasvr.UnsignedTinyInt:
		res = schemasvr.SmallInt
		return res
	case schemasvr.Bit:
		res = schemasvr.Boolean
		return res
	case schemasvr.MediumInt, schemasvr.Int, schemasvr.Integer, schemasvr.UnsignedMediumInt, schemasvr.UnsignedSmallInt:
		if c.IsAutoIncrement {
			return schemasvr.Serial
		}
		return schemasvr.Integer
	case schemasvr.BigInt, schemasvr.UnsignedBigInt, schemasvr.UnsignedInt:
		if c.IsAutoIncrement {
			return schemasvr.BigSerial
		}
		return schemasvr.BigInt
	case schemasvr.Serial, schemasvr.BigSerial:
		c.IsAutoIncrement = true
		c.Nullable = false
		res = t
	case schemasvr.Binary, schemasvr.VarBinary:
		return schemasvr.Bytea
	case schemasvr.DateTime:
		res = schemasvr.TimeStamp
	case schemasvr.TimeStampz:
		return "timestamp with time zone"
	case schemasvr.Float:
		res = schemasvr.Real
	case schemasvr.TinyText, schemasvr.MediumText, schemasvr.LongText:
		res = schemasvr.Text
	case schemasvr.NChar:
		res = schemasvr.Char
	case schemasvr.NVarchar:
		res = schemasvr.Varchar
	case schemasvr.Uuid:
		return schemasvr.Uuid
	case schemasvr.Blob, schemasvr.TinyBlob, schemasvr.MediumBlob, schemasvr.LongBlob:
		return schemasvr.Bytea
	case schemasvr.Double:
		return "DOUBLE PRECISION"
	default:
		if c.IsAutoIncrement {
			return schemasvr.Serial
		}
		res = t
	}
	if strings.EqualFold(res, "bool") {
		// for bool, we don't need length information
		return res
	}
	hasLen1 := (c.Length > 0)
	hasLen2 := (c.Length2 > 0)
	if hasLen2 {
		res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
	} else if hasLen1 {
		res += "(" + strconv.Itoa(c.Length) + ")"
	}
	return res
}
func (db *postgres) Features() *DialectFeatures {
	return &DialectFeatures{
		AutoincrMode: IncrAutoincrMode,
	}
}
func (db *postgres) ColumnTypeKind(t string) int {
	switch strings.ToUpper(t) {
	case "DATETIME", "TIMESTAMP":
		return schemasvr.TIME_TYPE
	case "VARCHAR", "TEXT":
		return schemasvr.TEXT_TYPE
	case "BIGINT", "BIGSERIAL", "SMALLINT", "INT", "INT8", "INT4", "INTEGER", "SERIAL", "FLOAT", "FLOAT4", "REAL", "DOUBLE PRECISION":
		return schemasvr.NUMERIC_TYPE
	case "BOOL":
		return schemasvr.BOOL_TYPE
	default:
		return schemasvr.UNKNOW_TYPE
	}
}
func (db *postgres) IsReserved(name string) bool {
	_, ok := postgresReservedWords[strings.ToUpper(name)]
	return ok
}
func (db *postgres) AutoIncrStr() string {
	return ""
}
func (db *postgres) IndexCheckSQL(tableName, idxName string) (string, []interface{}) {
	if len(db.getSchema()) == 0 {
		args := []interface{}{tableName, idxName}
		return `SELECT indexname FROM pg_indexes WHERE tablename = ? AND indexname = ?`, args
	}
	args := []interface{}{db.getSchema(), tableName, idxName}
	return `SELECT indexname FROM pg_indexes ` +
		`WHERE schemaname = ? AND tablename = ? AND indexname = ?`, args
}
func (db *postgres) IsTableExist(queryer core.Queryer, ctx context.Context, tableName string) (bool, error) {
	if len(db.getSchema()) == 0 {
		return db.HasRecords(queryer, ctx, `SELECT tablename FROM pg_tables WHERE tablename = $1`, tableName)
	}
	return db.HasRecords(queryer, ctx, `SELECT tablename FROM pg_tables WHERE schemaname = $1 AND tablename = $2`,
		db.getSchema(), tableName)
}
func (db *postgres) AddColumnSQL(tableName string, col *schemasvr.Column) string {
	s, _ := ColumnString(db.dialect, col, true)
	quoter := db.dialect.Quoter()
	addColumnSQL := ""
	commentSQL := "; "
	if len(db.getSchema()) == 0 || strings.Contains(tableName, ".") {
		addColumnSQL = fmt.Sprintf("ALTER TABLE %s ADD %s", quoter.Quote(tableName), s)
		commentSQL += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s'", quoter.Quote(tableName), quoter.Quote(col.Name), col.Comment)
		return addColumnSQL + commentSQL
	}
	addColumnSQL = fmt.Sprintf("ALTER TABLE %s.%s ADD %s", quoter.Quote(db.getSchema()), quoter.Quote(tableName), s)
	commentSQL += fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS '%s'", quoter.Quote(db.getSchema()), quoter.Quote(tableName), quoter.Quote(col.Name), col.Comment)
	return addColumnSQL + commentSQL
}
func (db *postgres) ModifyColumnSQL(tableName string, col *schemasvr.Column) string {
	quoter := db.dialect.Quoter()
	modifyColumnSQL := ""
	commentSQL := "; "
	if len(db.getSchema()) == 0 || strings.Contains(tableName, ".") {
		modifyColumnSQL = fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s", quoter.Quote(tableName), quoter.Quote(col.Name), db.SQLType(col))
		commentSQL += fmt.Sprintf("COMMENT ON COLUMN %s.%s IS '%s'", quoter.Quote(tableName), quoter.Quote(col.Name), col.Comment)
		return modifyColumnSQL + commentSQL
	}
	modifyColumnSQL = fmt.Sprintf("ALTER TABLE %s.%s ALTER COLUMN %s TYPE %s", quoter.Quote(db.getSchema()), quoter.Quote(tableName), quoter.Quote(col.Name), db.SQLType(col))
	commentSQL += fmt.Sprintf("COMMENT ON COLUMN %s.%s.%s IS '%s'", quoter.Quote(db.getSchema()), quoter.Quote(tableName), quoter.Quote(col.Name), col.Comment)
	return modifyColumnSQL + commentSQL
}
func (db *postgres) DropIndexSQL(tableName string, index *schemasvr.Index) string {
	idxName := index.Name
	tableParts := strings.Split(strings.Replace(tableName, `"`, "", -1), ".")
	tableName = tableParts[len(tableParts)-1]
	if !strings.HasPrefix(idxName, "UQE_") &&
		!strings.HasPrefix(idxName, "IDX_") {
		if index.Type == schemasvr.UniqueType {
			idxName = fmt.Sprintf("UQE_%v_%v", tableName, index.Name)
		} else {
			idxName = fmt.Sprintf("IDX_%v_%v", tableName, index.Name)
		}
	}
	if db.getSchema() != "" {
		idxName = db.getSchema() + "." + idxName
	}
	return fmt.Sprintf("DROP INDEX %v", db.Quoter().Quote(idxName))
}
func (db *postgres) IsColumnExist(queryer core.Queryer, ctx context.Context, tableName, colName string) (bool, error) {
	args := []interface{}{db.getSchema(), tableName, colName}
	query := "SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = $1 AND table_name = $2" +
		" AND column_name = $3"
	if len(db.getSchema()) == 0 {
		args = []interface{}{tableName, colName}
		query = "SELECT column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE table_name = $1" +
			" AND column_name = $2"
	}
	rows, err := queryer.QueryContext(ctx, query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, rows.Err()
}
func (db *postgres) GetColumns(queryer core.Queryer, ctx context.Context, tableName string) ([]string, map[string]*schemasvr.Column, error) {
	args := []interface{}{tableName}
	s := `SELECT column_name, column_default, is_nullable, data_type, character_maximum_length, description,
    CASE WHEN p.contype = 'p' THEN true ELSE false END AS primarykey,
    CASE WHEN p.contype = 'u' THEN true ELSE false END AS uniquekey
FROM pg_attribute f
    JOIN pg_class c ON c.oid = f.attrelid JOIN pg_type t ON t.oid = f.atttypid
    LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = f.attnum
    LEFT JOIN pg_description de ON f.attrelid=de.objoid AND f.attnum=de.objsubid
    LEFT JOIN pg_namespace n ON n.oid = c.relnamespace
    LEFT JOIN pg_constraint p ON p.conrelid = c.oid AND f.attnum = ANY (p.conkey)
    LEFT JOIN pg_class AS g ON p.confrelid = g.oid
    LEFT JOIN INFORMATION_SCHEMA.COLUMNS s ON s.column_name=f.attname AND c.relname=s.table_name
WHERE n.nspname= s.table_schema AND c.relkind = 'r'::char AND c.relname = $1%s AND f.attnum > 0 ORDER BY f.attnum;`
	schema := db.getSchema()
	if schema != "" {
		s = fmt.Sprintf(s, " AND s.table_schema = $2")
		args = append(args, schema)
	} else {
		s = fmt.Sprintf(s, "")
	}
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	cols := make(map[string]*schemasvr.Column)
	colSeq := make([]string, 0)
	for rows.Next() {
		col := new(schemasvr.Column)
		col.Indexes = make(map[string]int)
		var colName, isNullable, dataType string
		var maxLenStr, colDefault, description *string
		var isPK, isUnique bool
		err = rows.Scan(&colName, &colDefault, &isNullable, &dataType, &maxLenStr, &description, &isPK, &isUnique)
		if err != nil {
			return nil, nil, err
		}
		var maxLen int
		if maxLenStr != nil {
			maxLen, err = strconv.Atoi(*maxLenStr)
			if err != nil {
				return nil, nil, err
			}
		}
		if colDefault != nil && *colDefault == "unique_rowid()" { // ignore the system column added by cockroach
			continue
		}
		col.Name = strings.Trim(colName, `" `)
		if colDefault != nil {
			var theDefault = *colDefault
			// cockroach has type with the default value with :::
			// and postgres with ::, we should remove them before store them
			idx := strings.Index(theDefault, ":::")
			if idx == -1 {
				idx = strings.Index(theDefault, "::")
			}
			if idx > -1 {
				theDefault = theDefault[:idx]
			}
			if strings.HasSuffix(theDefault, "+00:00'") {
				theDefault = theDefault[:len(theDefault)-7] + "'"
			}
			col.Default = theDefault
			col.DefaultIsEmpty = false
			if strings.HasPrefix(col.Default, "nextval(") {
				col.IsAutoIncrement = true
				col.Default = ""
				col.DefaultIsEmpty = true
			}
		} else {
			col.DefaultIsEmpty = true
		}
		if description != nil {
			col.Comment = *description
		}
		if isPK {
			col.IsPrimaryKey = true
		}
		col.Nullable = (isNullable == "YES")
		switch strings.ToLower(dataType) {
		case "character varying", "string":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Varchar, DefaultLength: 0, DefaultLength2: 0}
		case "character":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Char, DefaultLength: 0, DefaultLength2: 0}
		case "timestamp without time zone":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.DateTime, DefaultLength: 0, DefaultLength2: 0}
		case "timestamp with time zone":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.TimeStampz, DefaultLength: 0, DefaultLength2: 0}
		case "double precision":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Double, DefaultLength: 0, DefaultLength2: 0}
		case "boolean":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Bool, DefaultLength: 0, DefaultLength2: 0}
		case "time without time zone":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Time, DefaultLength: 0, DefaultLength2: 0}
		case "bytes":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Binary, DefaultLength: 0, DefaultLength2: 0}
		case "oid":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.BigInt, DefaultLength: 0, DefaultLength2: 0}
		case "array":
			col.SQLType = schemasvr.SQLType{Name: schemasvr.Array, DefaultLength: 0, DefaultLength2: 0}
		default:
			startIdx := strings.Index(strings.ToLower(dataType), "string(")
			if startIdx != -1 && strings.HasSuffix(dataType, ")") {
				length := dataType[startIdx+8 : len(dataType)-1]
				l, _ := strconv.Atoi(length)
				col.SQLType = schemasvr.SQLType{Name: "STRING", DefaultLength: l, DefaultLength2: 0}
			} else {
				col.SQLType = schemasvr.SQLType{Name: strings.ToUpper(dataType), DefaultLength: 0, DefaultLength2: 0}
			}
		}
		if _, ok := schemasvr.SqlTypes[col.SQLType.Name]; !ok {
			return nil, nil, fmt.Errorf("unknown colType: %s - %s", dataType, col.SQLType.Name)
		}
		col.Length = maxLen
		if !col.DefaultIsEmpty {
			if col.SQLType.IsText() {
				if strings.HasSuffix(col.Default, "::character varying") {
					col.Default = strings.TrimSuffix(col.Default, "::character varying")
				} else if !strings.HasPrefix(col.Default, "'") {
					col.Default = "'" + col.Default + "'"
				}
			} else if col.SQLType.IsTime() {
				col.Default = strings.TrimSuffix(col.Default, "::timestamp without time zone")
			}
		}
		cols[col.Name] = col
		colSeq = append(colSeq, col.Name)
	}
	if rows.Err() != nil {
		return nil, nil, rows.Err()
	}
	return colSeq, cols, nil
}
func (db *postgres) GetTables(queryer core.Queryer, ctx context.Context) ([]*schemasvr.Table, error) {
	args := []interface{}{}
	s := "SELECT tablename FROM pg_tables"
	schema := db.getSchema()
	if schema != "" {
		args = append(args, schema)
		s = s + " WHERE schemaname = $1"
	}
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tables := make([]*schemasvr.Table, 0)
	for rows.Next() {
		table := schemasvr.NewEmptyTable()
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		table.Name = name
		tables = append(tables, table)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tables, nil
}
func getIndexColName(indexdef string) []string {
	var colNames []string
	cs := strings.Split(indexdef, "(")
	for _, v := range strings.Split(strings.Split(cs[1], ")")[0], ",") {
		colNames = append(colNames, strings.Split(strings.TrimLeft(v, " "), " ")[0])
	}
	return colNames
}
func (db *postgres) GetIndexes(queryer core.Queryer, ctx context.Context, tableName string) (map[string]*schemasvr.Index, error) {
	args := []interface{}{tableName}
	s := "SELECT indexname, indexdef FROM pg_indexes WHERE tablename=$1"
	if len(db.getSchema()) != 0 {
		args = append(args, db.getSchema())
		s += " AND schemaname=$2"
	}
	rows, err := queryer.QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	indexes := make(map[string]*schemasvr.Index)
	for rows.Next() {
		var indexType int
		var indexName, indexdef string
		var colNames []string
		err = rows.Scan(&indexName, &indexdef)
		if err != nil {
			return nil, err
		}
		if indexName == "primary" {
			continue
		}
		indexName = strings.Trim(indexName, `" `)
		// ignore primary index
		if strings.HasSuffix(indexName, "_pkey") || strings.EqualFold(indexName, "primary") {
			continue
		}
		if strings.HasPrefix(indexdef, "CREATE UNIQUE INDEX") {
			indexType = schemasvr.UniqueType
		} else {
			indexType = schemasvr.IndexType
		}
		colNames = getIndexColName(indexdef)
		isSkip := false
		//Oid It's a special index. You can't put it in
		for _, element := range colNames {
			if "oid" == element {
				isSkip = true
				break
			}
		}
		if isSkip {
			continue
		}
		var isRegular bool
		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			newIdxName := indexName[5+len(tableName):]
			isRegular = true
			if newIdxName != "" {
				indexName = newIdxName
			}
		}
		index := &schemasvr.Index{Name: indexName, Type: indexType, Cols: make([]string, 0)}
		for _, colName := range colNames {
			col := strings.TrimSpace(strings.Replace(colName, `"`, "", -1))
			fields := strings.Split(col, " ")
			index.Cols = append(index.Cols, fields[0])
		}
		index.IsRegular = isRegular
		indexes[index.Name] = index
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return indexes, nil
}
func (db *postgres) CreateTableSQL(ctx context.Context, queryer core.Queryer, table *schemasvr.Table, tableName string) (string, bool, error) {
	quoter := db.dialect.Quoter()
	if len(db.getSchema()) != 0 && !strings.Contains(tableName, ".") {
		tableName = fmt.Sprintf("%s.%s", db.getSchema(), tableName)
	}
	createTableSQL, ok, err := db.Base.CreateTableSQL(ctx, queryer, table, tableName)
	if err != nil {
		return "", ok, err
	}
	commentSQL := "; "
	if table.Comment != "" {
		// support schema.table -> "schema"."table"
		commentSQL += fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", quoter.Quote(tableName), table.Comment)
	}
	return createTableSQL + commentSQL, true, nil
}
func (db *postgres) Filters() []Filter {
	return []Filter{&SeqFilter{Prefix: "$", Start: 1}}
}

type pqDriver struct {
	baseDriver
}
type values map[string]string

func parseURL(connstr string) (string, error) {
	u, err := url.Parse(connstr)
	if err != nil {
		return "", err
	}
	if u.Scheme != "postgresql" && u.Scheme != "postgres" {
		return "", fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}
	escaper := strings.NewReplacer(` `, `\ `, `'`, `\'`, `\`, `\\`)
	if u.Path != "" {
		return escaper.Replace(u.Path[1:]), nil
	}
	return "", nil
}
func parseOpts(urlStr string, o values) error {
	if len(urlStr) == 0 {
		return fmt.Errorf("invalid options: %s", urlStr)
	}
	urlStr = strings.TrimSpace(urlStr)
	var (
		inQuote bool
		state   int // 0 key, 1 space, 2 value, 3 equal
		start   int
		key     string
	)
	for i, c := range urlStr {
		switch c {
		case ' ':
			if !inQuote {
				if state == 2 {
					state = 1
					v := urlStr[start:i]
					if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
						v = v[1 : len(v)-1]
					} else if strings.HasPrefix(v, "'") || strings.HasSuffix(v, "'") {
						return fmt.Errorf("wrong single quote in %d of %s", i, urlStr)
					}
					o[key] = v
				} else if state != 1 {
					return fmt.Errorf("wrong format: %v", urlStr)
				}
			}
		case '\'':
			if state == 3 {
				state = 2
				start = i
			} else if state != 2 {
				return fmt.Errorf("wrong format: %v", urlStr)
			}
			inQuote = !inQuote
		case '=':
			if !inQuote {
				if state != 0 {
					return fmt.Errorf("wrong format: %v", urlStr)
				}
				key = urlStr[start:i]
				state = 3
			}
		default:
			if state == 3 {
				state = 2
				start = i
			} else if state == 1 {
				state = 0
				start = i
			}
		}
		if i == len(urlStr)-1 {
			if state != 2 {
				return errors.New("no value matched key")
			}
			v := urlStr[start : i+1]
			if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
				v = v[1 : len(v)-1]
			} else if strings.HasPrefix(v, "'") || strings.HasSuffix(v, "'") {
				return fmt.Errorf("wrong single quote in %d of %s", i, urlStr)
			}
			o[key] = v
		}
	}
	return nil
}
func (p *pqDriver) Features() *DriverFeatures {
	return &DriverFeatures{
		SupportReturnInsertedID: false,
	}
}
func (p *pqDriver) Parse(driverName, dataSourceName string) (*URI, error) {
	db := &URI{DBType: schemasvr.POSTGRES}
	var err error
	if strings.Contains(dataSourceName, "://") {
		if !strings.HasPrefix(dataSourceName, "postgresql://") && !strings.HasPrefix(dataSourceName, "postgres://") {
			return nil, fmt.Errorf("unsupported protocol %v", dataSourceName)
		}
		db.DBName, err = parseURL(dataSourceName)
		if err != nil {
			return nil, err
		}
	} else {
		o := make(values)
		err = parseOpts(dataSourceName, o)
		if err != nil {
			return nil, err
		}
		db.DBName = o["dbname"]
	}
	if db.DBName == "" {
		return nil, errors.New("dbname is empty")
	}
	return db, nil
}
func (p *pqDriver) GenScanResult(colType string) (interface{}, error) {
	switch colType {
	case "VARCHAR", "TEXT":
		var s sql.NullString
		return &s, nil
	case "BIGINT", "BIGSERIAL":
		var s sql.NullInt64
		return &s, nil
	case "SMALLINT", "INT", "INT8", "INT4", "INTEGER", "SERIAL":
		var s sql.NullInt32
		return &s, nil
	case "FLOAT", "FLOAT4", "REAL", "DOUBLE PRECISION":
		var s sql.NullFloat64
		return &s, nil
	case "DATETIME", "TIMESTAMP":
		var s sql.NullTime
		return &s, nil
	case "BOOL":
		var s sql.NullBool
		return &s, nil
	default:
		var r sql.RawBytes
		return &r, nil
	}
}

type pqDriverPgx struct {
	pqDriver
}

func (pgx *pqDriverPgx) Parse(driverName, dataSourceName string) (*URI, error) {
	// Remove the leading characters for driver to work
	if len(dataSourceName) >= 9 && dataSourceName[0] == 0 {
		dataSourceName = dataSourceName[9:]
	}
	return pgx.pqDriver.Parse(driverName, dataSourceName)
}

// QueryDefaultPostgresSchema returns the default postgres schema
func QueryDefaultPostgresSchema(ctx context.Context, queryer core.Queryer) (string, error) {
	rows, err := queryer.QueryContext(ctx, "SHOW SEARCH_PATH")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	if rows.Next() {
		var defaultSchema string
		if err = rows.Scan(&defaultSchema); err != nil {
			return "", err
		}
		parts := strings.Split(defaultSchema, ",")
		return strings.TrimSpace(parts[len(parts)-1]), nil
	}
	if rows.Err() != nil {
		return "", rows.Err()
	}
	return "", errors.New("no default schema")
}
