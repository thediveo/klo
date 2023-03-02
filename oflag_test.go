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
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("-o output options", func() {

	type Foo struct{ Foo string }
	foo := Foo{Foo: "Foo!"}

	It("unknown -o", func() {
		BadPrinter(PrinterFromFlag("unknown", &Specs{WideColumnSpec: "foo"}))
	})

	It("default -o", func() {
		BadPrinter(PrinterFromFlag("", nil))
		PrinterPass(GoodPrinter(PrinterFromFlag("",
			&Specs{DefaultColumnSpec: "FOO:Foo,BAR:bar"})), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o wide", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("wide",
			&Specs{WideColumnSpec: "FOO:Foo,BAR:bar"})), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o custom-columns", func() {
		BadPrinter(PrinterFromFlag("custom-columns", nil))
		PrinterPass(GoodPrinter(PrinterFromFlag("custom-columns=FOO:Foo,BAR:bar", nil)), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o customs-columns-file", func() {
		BadPrinter(PrinterFromFlag("custom-columns-file", nil))
		BadPrinter(PrinterFromFlag("custom-columns-file=./testdata/missing.columns", nil))
		PrinterPass(GoodPrinter(PrinterFromFlag("custom-columns-file=./testdata/foobar.columns", nil)), []Foo{foo},
			`FOO  BAR
Foo! <none>
`)
	})

	It("-o json", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("json", nil)), foo, `{
    "Foo": "Foo!"
}
`)
	})

	It("-o jsonpath", func() {
		BadPrinter(PrinterFromFlag("jsonpath", nil))
		PrinterPass(GoodPrinter(PrinterFromFlag("jsonpath={[*].Foo}", nil)),
			[]Foo{foo},
			`Foo!`)
	})

	It("-o jsonpath-file", func() {
		BadPrinter(PrinterFromFlag("jsonpath-file", nil))
		BadPrinter(PrinterFromFlag("jsonpath-file=./testdata/missing.jsonpath", nil))
		BadPrinter(PrinterFromFlag("jsonpath-file=./testdata/empty.jsonpath", nil))
		PrinterFail(GoodPrinter(PrinterFromFlag("jsonpath-file=./testdata/unknown.jsonpath", nil)), []Foo{foo})
		PrinterPass(GoodPrinter(PrinterFromFlag("jsonpath-file=./testdata/valid.jsonpath", nil)), []Foo{foo}, `Foo!`)
	})

	It("-o yaml", func() {
		PrinterPass(GoodPrinter(PrinterFromFlag("yaml", nil)), foo,
			`Foo: Foo!
`)
	})

	It("-o go-template", func() {
		BadPrinter(PrinterFromFlag(`go-template={{oops}}`, nil))
		PrinterPass(GoodPrinter(PrinterFromFlag(`go-template`, nil)), nil,
			"")
		PrinterPass(GoodPrinter(PrinterFromFlag(`go-template={{"ok"}}`, nil)), nil,
			`ok`)
		PrinterPass(GoodPrinter(PrinterFromFlag(`go-template`, &Specs{GoTemplateArg: `{{"ok"}}`})), nil,
			`ok`)
	})

	It("-o go-template-file", func() {
		BadPrinter(PrinterFromFlag(`go-template-file`, nil))
		BadPrinter(PrinterFromFlag(`go-template-file=./testdata/missing.tpl`, nil))
		PrinterPass(GoodPrinter(PrinterFromFlag(`go-template-file=./testdata/ok.tpl`, nil)), nil,
			"ok")
		PrinterPass(GoodPrinter(PrinterFromFlag(`go-template-file=`, &Specs{GoTemplateArg: "./testdata/ok.tpl"})), nil,
			"ok")
	})

})
