package orm

// Exist returns true if the record exist otherwise return false
func (session *Session) Exist(bean ...interface{}) (bool, error) {
	if session.isAutoClose {
		defer session.Close()
	}
	if session.statement.LastError != nil {
		return false, session.statement.LastError
	}
	sqlStr, args, err := session.statement.GenExistSQL(bean...)
	if err != nil {
		return false, err
	}
	rows, err := session.queryRows(sqlStr, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	}
	return false, rows.Err()
}
