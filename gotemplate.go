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
	"html/template"
	"io"
)

// GoTemplatePrinter prints values in JSON format.
type GoTemplatePrinter struct {
	Template *template.Template // The compiled golang template.
	raw      string             // Original template text, to ease debugging.
}

// NewGoTemplatePrinter returns a printer for outputting values in JSON format.
func NewGoTemplatePrinter(tmpl string) (ValuePrinter, error) {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	return &GoTemplatePrinter{
		Template: t,
		raw:      tmpl,
	}, nil
}

// Fprint prints a value in JSON format.
func (p *GoTemplatePrinter) Fprint(w io.Writer, v interface{}) (err error) {
	defer func() {
		if tp := recover(); tp != nil {
			err = fmt.Errorf("template panicked: %+v", tp)
		}
	}()
	return p.Template.Execute(w, v)
}
