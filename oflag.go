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

// PrinterFromFlag returns a suitable value printer according to the output
// format specified as the flagvalue. The "-o" flag value is passed in via the
// flagvalue paramter (without the "-o") and should denote one of the supported
// output formats, such as "json", "yaml", "custom-columns", et cetera.
//
// If the output format flagvalue is zero (""), then the returned printer will
// default to a custom columns printer using the also specified columnspec.
//
// If an additional widecolumns custom column spec was given (non-zero), then
// "-o wide" will be supported, otherwise a user trying to use the wide output
// format will raise an error.
func PrinterFromFlag(flagvalue string, columnspec, widecolumnsspec string) (ValuePrinter, error) {
	if flagvalue == "" {
		flagvalue = "custom-columns=" + columnspec
	}
	// Do we support "-o wide"? Then map this to "-o customcolumns=..." for
	// the specified wide columns spec.
	if flagvalue == "wide" && widecolumnsspec != "" {
		flagvalue = "custom-columns=" + widecolumnsspec
	}
	ov := strings.Split(flagvalue, "=")
	switch ov[0] {
	case "custom-columns":
		if len(ov) != 2 {
			return nil, fmt.Errorf("missing custom columns specification")
		}
		return NewCustomColumnsPrinterFromSpec(ov[1])
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
	if widecolumnsspec != "" {
		wide = " 'wide',"
	}
	return nil, fmt.Errorf("unexpected output format, expected 'json',%s or 'yaml'", wide)
}
