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
	. "github.com/onsi/ginkgo"
)

var _ = Describe("-o output options", func() {

	type Foo struct{ Foo string }
	foo := Foo{Foo: "Foo!"}

	It("unknown -o", func() {
		BadPrinter(PrinterFromFlag("unknown", "", "widespec"))
	})

	It("default -o", func() {
		BadPrinter(PrinterFromFlag("", "colspec", "widespec"))
		PrinterPass(GoodPrinter(PrinterFromFlag("", "FOO:Foo,BAR:bar", "")), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o wide", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("wide", "", "FOO:Foo,BAR:bar")), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o custom-columns", func() {
		BadPrinter(PrinterFromFlag("custom-columns", "", ""))
		PrinterPass(GoodPrinter(PrinterFromFlag("custom-columns=FOO:Foo,BAR:bar", "", "")), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o json", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("json", "", "")), foo, `{
    "Foo": "Foo!"
}
`)
	})

	It("-o jsonpath", func() {
		BadPrinter(PrinterFromFlag("jsonpath", "", ""))
		PrinterPass(GoodPrinter(PrinterFromFlag("jsonpath={[*].Foo}", "", "")),
			[]Foo{foo},
			`Foo!`)
	})

	It("-o jsonpath-file", func() {
		BadPrinter(PrinterFromFlag("jsonpath-file", "", ""))
		BadPrinter(PrinterFromFlag("jsonpath-file=./testdata/missing.jsonpath", "", ""))
		BadPrinter(PrinterFromFlag("jsonpath-file=./testdata/empty.jsonpath", "", ""))
		PrinterFail(GoodPrinter(PrinterFromFlag("jsonpath-file=./testdata/unknown.jsonpath", "", "")), []Foo{foo})
		PrinterPass(GoodPrinter(PrinterFromFlag("jsonpath-file=./testdata/valid.jsonpath", "", "")), []Foo{foo}, `Foo!`)
	})

	It("-o yaml", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("yaml", "", "")), foo,
			`Foo: Foo!
`)
	})

})
