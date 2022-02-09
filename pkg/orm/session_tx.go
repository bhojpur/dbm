package orm

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

// Begin a transaction
func (session *Session) Begin() error {
	if session.isAutoCommit {
		tx, err := session.DB().BeginTx(session.ctx, nil)
		if err != nil {
			return err
		}
		session.isAutoCommit = false
		session.isCommitedOrRollbacked = false
		session.tx = tx
		session.saveLastSQL("BEGIN TRANSACTION")
	}
	return nil
}

// Rollback When using transaction, you can rollback if any error
func (session *Session) Rollback() error {
	if !session.isAutoCommit && !session.isCommitedOrRollbacked {
		session.saveLastSQL("ROLL BACK")
		session.isCommitedOrRollbacked = true
		session.isAutoCommit = true
		return session.tx.Rollback()
	}
	return nil
}

// Commit When using transaction, Commit will commit all operations.
func (session *Session) Commit() error {
	if !session.isAutoCommit && !session.isCommitedOrRollbacked {
		session.saveLastSQL("COMMIT")
		session.isCommitedOrRollbacked = true
		session.isAutoCommit = true
		if err := session.tx.Commit(); err != nil {
			return err
		}
		// handle processors after tx committed
		closureCallFunc := func(closuresPtr *[]func(interface{}), bean interface{}) {
			if closuresPtr != nil {
				for _, closure := range *closuresPtr {
					closure(bean)
				}
			}
		}
		for bean, closuresPtr := range session.afterInsertBeans {
			closureCallFunc(closuresPtr, bean)
			if processor, ok := interface{}(bean).(AfterInsertProcessor); ok {
				processor.AfterInsert()
			}
		}
		for bean, closuresPtr := range session.afterUpdateBeans {
			closureCallFunc(closuresPtr, bean)
			if processor, ok := interface{}(bean).(AfterUpdateProcessor); ok {
				processor.AfterUpdate()
			}
		}
		for bean, closuresPtr := range session.afterDeleteBeans {
			closureCallFunc(closuresPtr, bean)
			if processor, ok := interface{}(bean).(AfterDeleteProcessor); ok {
				processor.AfterDelete()
			}
		}
		cleanUpFunc := func(slices *map[interface{}]*[]func(interface{})) {
			if len(*slices) > 0 {
				*slices = make(map[interface{}]*[]func(interface{}))
			}
		}
		cleanUpFunc(&session.afterInsertBeans)
		cleanUpFunc(&session.afterUpdateBeans)
		cleanUpFunc(&session.afterDeleteBeans)
	}
	return nil
}

// IsInTx if current session is in a transaction
func (session *Session) IsInTx() bool {
	return !session.isAutoCommit
}
