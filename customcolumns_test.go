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

	"github.com/onsi/gomega/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("custom columns printer", func() {

	type tfoo struct {
		Foo string
		Bar string
	}

	var foo = []tfoo{
		{Foo: "verylongfoo", Bar: "bar!"},
	}

	var (
		PASS = Succeed()
		FAIL = Not(Succeed())
	)

	DescribeTable("parses column spec expressions",
		func(expr string, outcome types.GomegaMatcher) {
			err := (&Column{}).SetExpression(expr)
			Expect(err).To(outcome)
		},
		Entry("empty spec", "", PASS),
		Entry("relaxed spec", "foo", PASS),
		Entry("relaxed . spec", ".foo", PASS),
		Entry("relaxed {} spec", "{foo}", PASS),
		Entry("correct spec", "{.foo}", PASS),
		Entry("incomplete { spec", "{foo", FAIL),
		Entry("incomplete [ spec", "foo[0", FAIL),
	)

	DescribeTable("rejects bad column specs",
		func(spec string) {
			Expect(NewCustomColumnsPrinterFromSpec(spec)).Error().To(HaveOccurred())
		},
		Entry("empty spec", ""),
		Entry("missing column expression", "FOO,BAR"),
		Entry("malformed column expression", "FOO:foo,BAR:{bar"),
	)

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

	DescribeTable("rejects creating custom column printers from invalid template streams",
		func(stream string) {
			Expect(NewCustomColumnsPrinterFromTemplate(strings.NewReader(stream))).Error().To(HaveOccurred())
		},
		Entry("empty template stream", ""),
		Entry("2 empty lines", "\n\n"),
		Entry("only header line", "FOO BAR\n"),
		Entry("inconsistent # of columns", "FOO BAR\foo bar baz\n"),
		Entry("malformed column JSONPath expression", "FOO BAR BAZ\nFoo Bar {Baz"),
	)

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
