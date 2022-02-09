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
	"log/syslog"
)

var _ Logger = &SyslogLogger{}

// SyslogLogger will be depricated
type SyslogLogger struct {
	w       *syslog.Writer
	showSQL bool
}

// NewSyslogLogger implements Logger
func NewSyslogLogger(w *syslog.Writer) *SyslogLogger {
	return &SyslogLogger{w: w}
}

// Debug log content as Debug
func (s *SyslogLogger) Debug(v ...interface{}) {
	_ = s.w.Debug(fmt.Sprint(v...))
}

// Debugf log content as Debug and format
func (s *SyslogLogger) Debugf(format string, v ...interface{}) {
	_ = s.w.Debug(fmt.Sprintf(format, v...))
}

// Error log content as Error
func (s *SyslogLogger) Error(v ...interface{}) {
	_ = s.w.Err(fmt.Sprint(v...))
}

// Errorf log content as Errorf and format
func (s *SyslogLogger) Errorf(format string, v ...interface{}) {
	_ = s.w.Err(fmt.Sprintf(format, v...))
}

// Info log content as Info
func (s *SyslogLogger) Info(v ...interface{}) {
	_ = s.w.Info(fmt.Sprint(v...))
}

// Infof log content as Infof and format
func (s *SyslogLogger) Infof(format string, v ...interface{}) {
	_ = s.w.Info(fmt.Sprintf(format, v...))
}

// Warn log content as Warn
func (s *SyslogLogger) Warn(v ...interface{}) {
	_ = s.w.Warning(fmt.Sprint(v...))
}

// Warnf log content as Warnf and format
func (s *SyslogLogger) Warnf(format string, v ...interface{}) {
	_ = s.w.Warning(fmt.Sprintf(format, v...))
}

// Level shows log level
func (s *SyslogLogger) Level() LogLevel {
	return LOG_UNKNOWN
}

// SetLevel always return error, as current log/syslog package doesn't allow to set priority level after syslog.Writer created
func (s *SyslogLogger) SetLevel(l LogLevel) {}

// ShowSQL set if logging SQL
func (s *SyslogLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL if logging SQL
func (s *SyslogLogger) IsShowSQL() bool {
	return s.showSQL
}
