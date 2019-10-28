// Copyright 2019 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package klo

import (
	"fmt"
	"io"

	"k8s.io/client-go/util/jsonpath"
)

// JSONPathPrinter prints values in JSON format.
type JSONPathPrinter struct {
	Expr *jsonpath.JSONPath // Compiled JSONPath expression.
	raw  string             // Original JSONPath expression
}

// NewJSONPathPrinter returns a printer for outputting the values that were
// filtered using a JSONPath expression. Please note that the JSONPath
// expression here must be strictly conforming to the syntax rules, in
// particular: it must be enclosed in "{...}" and leading "." must be present.
// In addition, the elements addresses must exist, or otherwise an error will
// be raised when printing objects using this expression. In consequence,
// JSONPath printers are much less forgiving than the custom-column printers.
func NewJSONPathPrinter(expr string) (ValuePrinter, error) {
	jp := jsonpath.New("expr")
	if err := jp.Parse(expr); err != nil {
		return nil, err
	}
	return &JSONPathPrinter{
		Expr: jp,
		raw:  expr,
	}, nil
}

// Fprint prints fields of a value in text format, where the values are selected
// using JSONPath expressions.
func (p *JSONPathPrinter) Fprint(w io.Writer, v interface{}) error {
	if err := p.Expr.Execute(w, v); err != nil {
		return fmt.Errorf(
			"JSONPath failure on expression %q for value %+v",
			p.raw, v)
	}
	return nil
}
