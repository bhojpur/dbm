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

// BeforeInsertProcessor executed before an object is initially persisted to the database
type BeforeInsertProcessor interface {
	BeforeInsert()
}

// BeforeUpdateProcessor executed before an object is updated
type BeforeUpdateProcessor interface {
	BeforeUpdate()
}

// BeforeDeleteProcessor executed before an object is deleted
type BeforeDeleteProcessor interface {
	BeforeDelete()
}

// BeforeSetProcessor executed before data set to the struct fields
type BeforeSetProcessor interface {
	BeforeSet(string, Cell)
}

// AfterSetProcessor executed after data set to the struct fields
type AfterSetProcessor interface {
	AfterSet(string, Cell)
}

// AfterInsertProcessor executed after an object is persisted to the database
type AfterInsertProcessor interface {
	AfterInsert()
}

// AfterUpdateProcessor executed after an object has been updated
type AfterUpdateProcessor interface {
	AfterUpdate()
}

// AfterDeleteProcessor executed after an object has been deleted
type AfterDeleteProcessor interface {
	AfterDelete()
}

// AfterLoadProcessor executed after an ojbect has been loaded from database
type AfterLoadProcessor interface {
	AfterLoad()
}

// AfterLoadSessionProcessor executed after an ojbect has been loaded from database with session parameter
type AfterLoadSessionProcessor interface {
	AfterLoad(*Session)
}
type executedProcessorFunc func(*Session, interface{}) error
type executedProcessor struct {
	fun     executedProcessorFunc
	session *Session
	bean    interface{}
}

func (executor *executedProcessor) execute() error {
	return executor.fun(executor.session, executor.bean)
}
func (session *Session) executeProcessors() error {
	processors := session.afterProcessors
	session.afterProcessors = make([]executedProcessor, 0)
	for _, processor := range processors {
		if err := processor.execute(); err != nil {
			return err
		}
	}
	return nil
}
func cleanupProcessorsClosures(slices *[]func(interface{})) {
	if len(*slices) > 0 {
		*slices = make([]func(interface{}), 0)
	}
}
func executeBeforeClosures(session *Session, bean interface{}) {
	// handle before delete processors
	for _, closure := range session.beforeClosures {
		closure(bean)
	}
	cleanupProcessorsClosures(&session.beforeClosures)
}
func executeBeforeSet(bean interface{}, fields []string, scanResults []interface{}) {
	if b, hasBeforeSet := bean.(BeforeSetProcessor); hasBeforeSet {
		for ii, key := range fields {
			b.BeforeSet(key, Cell(scanResults[ii].(*interface{})))
		}
	}
}
func executeAfterSet(bean interface{}, fields []string, scanResults []interface{}) {
	if b, hasAfterSet := bean.(AfterSetProcessor); hasAfterSet {
		for ii, key := range fields {
			b.AfterSet(key, Cell(scanResults[ii].(*interface{})))
		}
	}
}
func buildAfterProcessors(session *Session, bean interface{}) {
	// handle afterClosures
	for _, closure := range session.afterClosures {
		session.afterProcessors = append(session.afterProcessors, executedProcessor{
			fun: func(sess *Session, bean interface{}) error {
				closure(bean)
				return nil
			},
			session: session,
			bean:    bean,
		})
	}
	if a, has := bean.(AfterLoadProcessor); has {
		session.afterProcessors = append(session.afterProcessors, executedProcessor{
			fun: func(sess *Session, bean interface{}) error {
				a.AfterLoad()
				return nil
			},
			session: session,
			bean:    bean,
		})
	}
	if a, has := bean.(AfterLoadSessionProcessor); has {
		session.afterProcessors = append(session.afterProcessors, executedProcessor{
			fun: func(sess *Session, bean interface{}) error {
				a.AfterLoad(sess)
				return nil
			},
			session: session,
			bean:    bean,
		})
	}
}
