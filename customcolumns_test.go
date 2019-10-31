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
	"strings"

	"k8s.io/client-go/util/jsonpath"

	t "github.com/thediveo/klo/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("custom columns printer", func() {

	type tfoo struct {
		Foo string
		Bar string
	}

	var foo = []tfoo{
		tfoo{Foo: "verylongfoo", Bar: "bar!"},
	}

	It("parses column spec expressions", func() {
		var c Column
		t.PassFail(t.PASSFAILS{
			t.PASS{"empty spec", c.SetExpression("")},
			t.PASS{"relaxed spec", c.SetExpression("foo")},
			t.PASS{"relaxed . spec", c.SetExpression(".foo")},
			t.PASS{"relaxed {} spec", c.SetExpression("{foo}")},
			t.PASS{"correct spec", c.SetExpression("{.foo}")},
			t.FAIL{"incomplete { spec", c.SetExpression("{foo")},
			t.FAIL{"incomplete [ spec", c.SetExpression("foo[0")},
		}) //nolint:composites
	})

	It("rejects bad column specs", func() {
		t.PassFail(t.PASSFAILS{
			t.FAIL{"empty spec", t.Err(NewCustomColumnsPrinterFromSpec(""))},
			t.FAIL{"missing column expressions",
				t.Err(NewCustomColumnsPrinterFromSpec("FOO,BAR"))},
			t.FAIL{"malformed column expression",
				t.Err(NewCustomColumnsPrinterFromSpec("FOO:foo,BAR:{bar"))},
		}) //nolint:composites
	})

	It("creates custom column printer from spec string", func() {
		p := GoodPrinter(NewCustomColumnsPrinterFromSpec("FOO:foo,BAR:.bar"))
		ccp := p.(*CustomColumnsPrinter)
		Expect(ccp.Columns).Should(HaveLen(2))
		Expect(*ccp.Columns[0]).Should(MatchFields(IgnoreExtras, Fields{
			"Header": Equal("FOO"),
			"Raw":    Equal("foo"),
		}))
		Expect(*ccp.Columns[1]).Should(MatchFields(IgnoreExtras, Fields{
			"Header": Equal("BAR"),
			"Raw":    Equal(".bar"),
		}))
	})

	It("prints neat tables using custom column specs", func() {
		p := GoodPrinter(NewCustomColumnsPrinterFromSpec("FOO:Foo,BAR:Bar,BAZ:blafasel"))
		PrinterPass(p, nil, `FOO  BAR  BAZ
`)
		PrinterPass(p, foo, `FOO         BAR  BAZ
verylongfoo bar! <none>
`)
		// For the (un)sake of code coverage...
		ccp := p.(*CustomColumnsPrinter)
		ccp.Columns[0].Template = jsonpath.New("zero")
		PrinterFail(p, foo)
	})

	It("rejects creating custom column printers from invalid template streams", func() {
		t.PassFail(t.PASSFAILS{
			t.FAIL{"empty template stream",
				t.Err(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
					"")))},
			t.FAIL{"2 empty lines",
				t.Err(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
					`   
   
`)))},
			t.FAIL{"only header line",
				t.Err(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
					`FOO BAR
`)))},
			t.FAIL{"inconsistent # of columns",
				t.Err(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
					`FOO BAR
foo bar baz
`)))},
			t.FAIL{"malformed column JSONPath expression",
				t.Err(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
					`FOO BAR BAZ
Foo Bar {Baz
`)))},
		})
	})

	It("prints neat tables using templates", func() {
		p := GoodPrinter(NewCustomColumnsPrinterFromTemplate(strings.NewReader(
			`FOO BAR BAZ
Foo Bar Baz
`)))
		PrinterPass(p, foo, `FOO         BAR  BAZ
verylongfoo bar! <none>
`)
		// For the (un)sake of code coverage...
		ccp := p.(*CustomColumnsPrinter)
		ccp.Columns[0].Template = jsonpath.New("zero")
		PrinterFail(p, foo)
	})

	It("allows different column padding", func() {
		p := GoodPrinter(NewCustomColumnsPrinterFromSpec("FOO:Foo,BAR:Bar,BAZ:blafasel"))
		p.(*CustomColumnsPrinter).Padding = 3
		PrinterPass(p, foo, `FOO           BAR    BAZ
verylongfoo   bar!   <none>
`)
		p.(*CustomColumnsPrinter).HideHeaders = true
		PrinterPass(p, foo, `verylongfoo   bar!   <none>
`)
	})

})
