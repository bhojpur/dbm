package log

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

	"github.com/bhojpur/dbm/pkg/orm/context"
)

// LogContext represents a log context
type LogContext context.ContextHook

// SQLLogger represents an interface to log SQL
type SQLLogger interface {
	BeforeSQL(context LogContext) // only invoked when IsShowSQL is true
	AfterSQL(context LogContext)  // only invoked when IsShowSQL is true
}

// ContextLogger represents a logger interface with context
type ContextLogger interface {
	SQLLogger
	Debugf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Level() LogLevel
	SetLevel(l LogLevel)
	ShowSQL(show ...bool)
	IsShowSQL() bool
}

var (
	_ ContextLogger = &LoggerAdapter{}
)

// enumerate all the context keys
var (
	SessionIDKey      = "__orm_session_id"
	SessionKey        = "__orm_session_key"
	SessionShowSQLKey = "__orm_show_sql"
)

// LoggerAdapter wraps a Logger interface as LoggerContext interface
type LoggerAdapter struct {
	logger Logger
}

// NewLoggerAdapter creates an adapter for old ORM logger interface
func NewLoggerAdapter(logger Logger) ContextLogger {
	return &LoggerAdapter{
		logger: logger,
	}
}

// BeforeSQL implements ContextLogger
func (l *LoggerAdapter) BeforeSQL(ctx LogContext) {}

// AfterSQL implements ContextLogger
func (l *LoggerAdapter) AfterSQL(ctx LogContext) {
	var sessionPart string
	v := ctx.Ctx.Value(SessionIDKey)
	if key, ok := v.(string); ok {
		sessionPart = fmt.Sprintf(" [%s]", key)
	}
	if ctx.ExecuteTime > 0 {
		l.logger.Infof("[SQL]%s %s %v - %v", sessionPart, ctx.SQL, ctx.Args, ctx.ExecuteTime)
	} else {
		l.logger.Infof("[SQL]%s %s %v", sessionPart, ctx.SQL, ctx.Args)
	}
}

// Debugf implements ContextLogger
func (l *LoggerAdapter) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

// Errorf implements ContextLogger
func (l *LoggerAdapter) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

// Infof implements ContextLogger
func (l *LoggerAdapter) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Warnf implements ContextLogger
func (l *LoggerAdapter) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

// Level implements ContextLogger
func (l *LoggerAdapter) Level() LogLevel {
	return l.logger.Level()
}

// SetLevel implements ContextLogger
func (l *LoggerAdapter) SetLevel(lv LogLevel) {
	l.logger.SetLevel(lv)
}

// ShowSQL implements ContextLogger
func (l *LoggerAdapter) ShowSQL(show ...bool) {
	l.logger.ShowSQL(show...)
}

// IsShowSQL implements ContextLogger
func (l *LoggerAdapter) IsShowSQL() bool {
	return l.logger.IsShowSQL()
}
