package statement

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
	"strings"

	"github.com/bhojpur/dbm/pkg/orm/schema"
	"github.com/bhojpur/sql/pkg/builder"
)

// ErrUnsupportedExprType represents an error with unsupported express type
type ErrUnsupportedExprType struct {
	tp string
}

func (err ErrUnsupportedExprType) Error() string {
	return fmt.Sprintf("Unsupported expression type: %v", err.tp)
}

// Expr represents an SQL express
type Expr struct {
	ColName string
	Arg     interface{}
}

// WriteArgs writes args to the writer
func (expr *Expr) WriteArgs(w *builder.BytesWriter) error {
	switch arg := expr.Arg.(type) {
	case *builder.Builder:
		if _, err := w.WriteString("("); err != nil {
			return err
		}
		if err := arg.WriteTo(w); err != nil {
			return err
		}
		if _, err := w.WriteString(")"); err != nil {
			return err
		}
	case string:
		if arg == "" {
			arg = "''"
		}
		if _, err := w.WriteString(fmt.Sprintf("%v", arg)); err != nil {
			return err
		}
	default:
		if _, err := w.WriteString("?"); err != nil {
			return err
		}
		w.Append(arg)
	}
	return nil
}

type exprParams []Expr

func (exprs exprParams) ColNames() []string {
	var cols = make([]string, 0, len(exprs))
	for _, expr := range exprs {
		cols = append(cols, expr.ColName)
	}
	return cols
}
func (exprs *exprParams) Add(name string, arg interface{}) {
	*exprs = append(*exprs, Expr{name, arg})
}
func (exprs exprParams) IsColExist(colName string) bool {
	for _, expr := range exprs {
		if strings.EqualFold(schema.CommonQuoter.Trim(expr.ColName), schema.CommonQuoter.Trim(colName)) {
			return true
		}
	}
	return false
}
func (exprs exprParams) WriteArgs(w *builder.BytesWriter) error {
	for i, expr := range exprs {
		if err := expr.WriteArgs(w); err != nil {
			return err
		}
		if i != len(exprs)-1 {
			if _, err := w.WriteString(","); err != nil {
				return err
			}
		}
	}
	return nil
}
