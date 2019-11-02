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
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Specs specifies custom-column formats for the default columns in
// "-o=customcolumns" mode, and for the "-o=wide" wide columns mode.
type Specs struct {
	// default custom-columns spec in format
	// "<header>:<json-path-expr>[,<header>:json-path-expr>]..."
	DefaultColumnSpec string
	// wide custom-columns spec in format
	// "<header>:<json-path-expr>[,<header>:json-path-expr>]..."
	WideColumnSpec string
}

// PrinterFromFlag returns a suitable value printer according to the output
// format specified as the flagvalue. The "-o" flag value is passed in via the
// flagvalue paramter (without the "-o") and should denote one of the
// supported output formats, such as "json", "yaml", "custom-columns", et
// cetera. The Specs parameter specifies the default custom-columns output
// format for "-o=" and "-o=wide". If Specs is nil, then no default
// custom-column formats will apply.
func PrinterFromFlag(flagvalue string, specs *Specs) (ValuePrinter, error) {
	// Apply empty default specs, if necessary.
	if specs == nil {
		specs = &Specs{}
	}
	// If no output format is specified, default to custom columns.
	if flagvalue == "" {
		flagvalue = "custom-columns=" + specs.DefaultColumnSpec
	}
	// Do we support "-o wide"? Then map this to "-o customcolumns=..." for
	// the specified wide columns spec.
	if flagvalue == "wide" && specs.WideColumnSpec != "" {
		flagvalue = "custom-columns=" + specs.WideColumnSpec
	}
	ov := strings.Split(flagvalue, "=")
	switch ov[0] {
	case "custom-columns":
		if len(ov) != 2 {
			return nil, fmt.Errorf("missing custom columns specification")
		}
		return NewCustomColumnsPrinterFromSpec(ov[1])
	case "custom-columns-file":
		if len(ov) != 2 {
			return nil, fmt.Errorf("missing custom columns filename")
		}
		f, err := os.Open(ov[1])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return NewCustomColumnsPrinterFromTemplate(f)
	case "json":
		return NewJSONPrinter()
	case "jsonpath":
		if len(ov) != 2 {
			return nil, fmt.Errorf("missing JSONPath expression")
		}
		return NewJSONPathPrinter(ov[1])
	case "jsonpath-file":
		if len(ov) != 2 {
			return nil, fmt.Errorf("missing JSONPath filename")
		}
		f, err := os.Open(ov[1])
		if err != nil {
			return nil, err
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		if (!sc.Scan() && sc.Err() != nil) || (sc.Text() == "") {
			return nil, fmt.Errorf("missing JSONPath expression in %q", ov[1])
		}
		return NewJSONPathPrinter(sc.Text())
	case "yaml":
		return NewYAMLPrinter()
	}
	// Unsupported/unknown output format.
	wide := ""
	if specs.WideColumnSpec != "" {
		wide = " 'wide',"
	}
	return nil, fmt.Errorf("unexpected output format %q, expected "+
		"'custom-columns', 'custom-columns-file', "+
		"'json', 'jsonpath', 'jsonpath-file',%s or 'yaml'", ov[0], wide)
}
