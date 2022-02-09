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
	"errors"
	"testing"
)

type testHook struct {
	before func(c *ContextHook) (context.Context, error)
	after  func(c *ContextHook) error
}

func (h *testHook) BeforeProcess(c *ContextHook) (context.Context, error) {
	if h.before != nil {
		return h.before(c)
	}
	return c.Ctx, nil
}
func (h *testHook) AfterProcess(c *ContextHook) error {
	if h.after != nil {
		return h.after(c)
	}
	return c.Err
}

var _ Hook = &testHook{}

func TestBeforeProcess(t *testing.T) {
	expectErr := errors.New("before error")
	tests := []struct {
		msg    string
		hooks  []Hook
		expect error
	}{
		{
			msg: "first hook return err",
			hooks: []Hook{
				&testHook{
					before: func(c *ContextHook) (ctx context.Context, err error) {
						return c.Ctx, expectErr
					},
				},
				&testHook{
					before: func(c *ContextHook) (ctx context.Context, err error) {
						return c.Ctx, nil
					},
				},
			},
			expect: expectErr,
		},
		{
			msg: "second hook return err",
			hooks: []Hook{
				&testHook{
					before: func(c *ContextHook) (ctx context.Context, err error) {
						return c.Ctx, nil
					},
				},
				&testHook{
					before: func(c *ContextHook) (ctx context.Context, err error) {
						return c.Ctx, expectErr
					},
				},
			},
			expect: expectErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			hooks := Hooks{}
			hooks.AddHook(tt.hooks...)
			_, err := hooks.BeforeProcess(&ContextHook{
				Ctx: context.Background(),
			})
			if err != tt.expect {
				t.Errorf("got %v, expect %v", err, tt.expect)
			}
		})
	}
}
func TestAfterProcess(t *testing.T) {
	expectErr := errors.New("expect err")
	tests := []struct {
		msg    string
		ctx    *ContextHook
		hooks  []Hook
		expect error
	}{
		{
			msg: "context has err",
			ctx: &ContextHook{
				Ctx: context.Background(),
				Err: expectErr,
			},
			hooks: []Hook{
				&testHook{
					after: func(c *ContextHook) error {
						return errors.New("hook err")
					},
				},
			},
			expect: expectErr,
		},
		{
			msg: "last hook has err",
			ctx: &ContextHook{
				Ctx: context.Background(),
				Err: nil,
			},
			hooks: []Hook{
				&testHook{
					after: func(c *ContextHook) error {
						return nil
					},
				},
				&testHook{
					after: func(c *ContextHook) error {
						return expectErr
					},
				},
			},
			expect: expectErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			hooks := Hooks{}
			hooks.AddHook(tt.hooks...)
			err := hooks.AfterProcess(tt.ctx)
			if err != tt.expect {
				t.Errorf("got %v, expect %v", err, tt.expect)
			}
		})
	}
}
