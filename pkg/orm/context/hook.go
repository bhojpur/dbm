package context

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
	"time"
)

// ContextHook represents a hook context
type ContextHook struct {
	start       time.Time
	Ctx         context.Context
	SQL         string        // log content or SQL
	Args        []interface{} // if it's a SQL, it's the arguments
	Result      sql.Result
	ExecuteTime time.Duration
	Err         error // SQL executed error
}

// NewContextHook return context for hook
func NewContextHook(ctx context.Context, sql string, args []interface{}) *ContextHook {
	return &ContextHook{
		start: time.Now(),
		Ctx:   ctx,
		SQL:   sql,
		Args:  args,
	}
}

// End finish the hook invokation
func (c *ContextHook) End(ctx context.Context, result sql.Result, err error) {
	c.Ctx = ctx
	c.Result = result
	c.Err = err
	c.ExecuteTime = time.Since(c.start)
}

// Hook represents a hook behaviour
type Hook interface {
	BeforeProcess(c *ContextHook) (context.Context, error)
	AfterProcess(c *ContextHook) error
}

// Hooks implements Hook interface but contains multiple Hook
type Hooks struct {
	hooks []Hook
}

// AddHook adds a Hook
func (h *Hooks) AddHook(hooks ...Hook) {
	h.hooks = append(h.hooks, hooks...)
}

// BeforeProcess invoked before execute the process
func (h *Hooks) BeforeProcess(c *ContextHook) (context.Context, error) {
	ctx := c.Ctx
	for _, h := range h.hooks {
		var err error
		ctx, err = h.BeforeProcess(c)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

// AfterProcess invoked after exetue the process
func (h *Hooks) AfterProcess(c *ContextHook) error {
	firstErr := c.Err
	for _, h := range h.hooks {
		err := h.AfterProcess(c)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
