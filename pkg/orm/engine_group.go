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

import (
	"context"
	"time"

	"github.com/bhojpur/dbm/pkg/orm/cache"
	ctxsvr "github.com/bhojpur/dbm/pkg/orm/context"
	"github.com/bhojpur/dbm/pkg/orm/dialect"
	"github.com/bhojpur/dbm/pkg/orm/log"
	"github.com/bhojpur/dbm/pkg/orm/name"
)

// EngineGroup defines an engine group
type EngineGroup struct {
	*Engine
	slaves []*Engine
	policy GroupPolicy
}

// NewEngineGroup creates a new engine group
func NewEngineGroup(args1 interface{}, args2 interface{}, policies ...GroupPolicy) (*EngineGroup, error) {
	var eg EngineGroup
	if len(policies) > 0 {
		eg.policy = policies[0]
	} else {
		eg.policy = RoundRobinPolicy()
	}
	driverName, ok1 := args1.(string)
	conns, ok2 := args2.([]string)
	if ok1 && ok2 {
		engines := make([]*Engine, len(conns))
		for i, conn := range conns {
			engine, err := NewEngine(driverName, conn)
			if err != nil {
				return nil, err
			}
			engine.engineGroup = &eg
			engines[i] = engine
		}
		eg.Engine = engines[0]
		eg.slaves = engines[1:]
		return &eg, nil
	}
	master, ok3 := args1.(*Engine)
	slaves, ok4 := args2.([]*Engine)
	if ok3 && ok4 {
		master.engineGroup = &eg
		for i := 0; i < len(slaves); i++ {
			slaves[i].engineGroup = &eg
		}
		eg.Engine = master
		eg.slaves = slaves
		return &eg, nil
	}
	return nil, ErrParamsType
}

// Close the engine
func (eg *EngineGroup) Close() error {
	err := eg.Engine.Close()
	if err != nil {
		return err
	}
	for i := 0; i < len(eg.slaves); i++ {
		err := eg.slaves[i].Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Context returned a group session
func (eg *EngineGroup) Context(ctx context.Context) *Session {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.Context(ctx)
}

// NewSession returned a group session
func (eg *EngineGroup) NewSession() *Session {
	sess := eg.Engine.NewSession()
	sess.sessionType = groupSession
	return sess
}

// Master returns the master engine
func (eg *EngineGroup) Master() *Engine {
	return eg.Engine
}

// Ping tests if database is alive
func (eg *EngineGroup) Ping() error {
	if err := eg.Engine.Ping(); err != nil {
		return err
	}
	for _, slave := range eg.slaves {
		if err := slave.Ping(); err != nil {
			return err
		}
	}
	return nil
}

// SetColumnMapper set the column name mapping rule
func (eg *EngineGroup) SetColumnMapper(mapper name.Mapper) {
	eg.Engine.SetColumnMapper(mapper)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetColumnMapper(mapper)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (eg *EngineGroup) SetConnMaxLifetime(d time.Duration) {
	eg.Engine.SetConnMaxLifetime(d)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetConnMaxLifetime(d)
	}
}

// SetDefaultCacher set the default cacher
func (eg *EngineGroup) SetDefaultCacher(cacher cache.Cacher) {
	eg.Engine.SetDefaultCacher(cacher)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetDefaultCacher(cacher)
	}
}

// SetLogger set the new logger
func (eg *EngineGroup) SetLogger(logger interface{}) {
	eg.Engine.SetLogger(logger)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetLogger(logger)
	}
}

// AddHook adds Hook
func (eg *EngineGroup) AddHook(hook ctxsvr.Hook) {
	eg.Engine.AddHook(hook)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].AddHook(hook)
	}
}

// SetLogLevel sets the logger level
func (eg *EngineGroup) SetLogLevel(level log.LogLevel) {
	eg.Engine.SetLogLevel(level)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetLogLevel(level)
	}
}

// SetMapper set the name mapping rules
func (eg *EngineGroup) SetMapper(mapper name.Mapper) {
	eg.Engine.SetMapper(mapper)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetMapper(mapper)
	}
}

// SetTagIdentifier set the tag identifier
func (eg *EngineGroup) SetTagIdentifier(tagIdentifier string) {
	eg.Engine.SetTagIdentifier(tagIdentifier)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetTagIdentifier(tagIdentifier)
	}
}

// SetMaxIdleConns set the max idle connections on pool, default is 2
func (eg *EngineGroup) SetMaxIdleConns(conns int) {
	eg.Engine.DB().SetMaxIdleConns(conns)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].DB().SetMaxIdleConns(conns)
	}
}

// SetMaxOpenConns is only available for go 1.2+
func (eg *EngineGroup) SetMaxOpenConns(conns int) {
	eg.Engine.DB().SetMaxOpenConns(conns)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].DB().SetMaxOpenConns(conns)
	}
}

// SetPolicy set the group policy
func (eg *EngineGroup) SetPolicy(policy GroupPolicy) *EngineGroup {
	eg.policy = policy
	return eg
}

// SetQuotePolicy sets the special quote policy
func (eg *EngineGroup) SetQuotePolicy(quotePolicy dialect.QuotePolicy) {
	eg.Engine.SetQuotePolicy(quotePolicy)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetQuotePolicy(quotePolicy)
	}
}

// SetTableMapper set the table name mapping rule
func (eg *EngineGroup) SetTableMapper(mapper name.Mapper) {
	eg.Engine.SetTableMapper(mapper)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].SetTableMapper(mapper)
	}
}

// ShowSQL show SQL statement or not on logger if log level is great than INFO
func (eg *EngineGroup) ShowSQL(show ...bool) {
	eg.Engine.ShowSQL(show...)
	for i := 0; i < len(eg.slaves); i++ {
		eg.slaves[i].ShowSQL(show...)
	}
}

// Slave returns one of the physical databases which is a slave according the policy
func (eg *EngineGroup) Slave() *Engine {
	switch len(eg.slaves) {
	case 0:
		return eg.Engine
	case 1:
		return eg.slaves[0]
	}
	return eg.policy.Slave(eg)
}

// Slaves returns all the slaves
func (eg *EngineGroup) Slaves() []*Engine {
	return eg.slaves
}

// Query execcute a select SQL and return the result
func (eg *EngineGroup) Query(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.Query(sqlOrArgs...)
}

// QueryInterface execcute a select SQL and return the result
func (eg *EngineGroup) QueryInterface(sqlOrArgs ...interface{}) ([]map[string]interface{}, error) {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.QueryInterface(sqlOrArgs...)
}

// QueryString execcute a select SQL and return the result
func (eg *EngineGroup) QueryString(sqlOrArgs ...interface{}) ([]map[string]string, error) {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.QueryString(sqlOrArgs...)
}

// Rows execcute a select SQL and return the result
func (eg *EngineGroup) Rows(bean interface{}) (*Rows, error) {
	sess := eg.NewSession()
	sess.isAutoClose = true
	return sess.Rows(bean)
}
