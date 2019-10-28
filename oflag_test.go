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
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("-o output options", func() {

	It("handles -o", func() {
		_, err := PrinterFromFlag("unknown", "", "widespec")
		Expect(err).Should(HaveOccurred())

		_, err = PrinterFromFlag("", "colspec", "widespec")
		Expect(err).Should(HaveOccurred())

		type Foo struct {
			Foo string
		}
		foo := Foo{Foo: "Foo!"}
		var out bytes.Buffer

		p, err := PrinterFromFlag("json", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, foo)).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`{
    "Foo": "Foo!"
}
`))

		out.Reset()
		p, err = PrinterFromFlag("yaml", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, foo)).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`Foo: Foo!
`))

		out.Reset()
		p, err = PrinterFromFlag("", "FOO:Foo,BAR:bar", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`FOO  BAR
Foo! <none>
`))

		out.Reset()
		p, err = PrinterFromFlag("wide", "", "FOO:Foo,BAR:bar")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`FOO  BAR
Foo! <none>
`))

		_, err = PrinterFromFlag("custom-columns", "", "")
		Expect(err).Should(HaveOccurred())

		out.Reset()
		p, err = PrinterFromFlag("custom-columns=FOO:Foo,BAR:bar", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`FOO  BAR
Foo! <none>
`))

		_, err = PrinterFromFlag("jsonpath", "", "")
		Expect(err).Should(HaveOccurred())

		out.Reset()
		p, err = PrinterFromFlag("jsonpath={[*].Foo}", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`Foo!`))

		_, err = PrinterFromFlag("jsonpath-file", "", "")
		Expect(err).Should(HaveOccurred())

		_, err = PrinterFromFlag("jsonpath-file=./test/missing.jsonpath", "", "")
		Expect(err).Should(HaveOccurred())

		out.Reset()
		_, err = PrinterFromFlag("jsonpath-file=./test/empty.jsonpath", "", "")
		Expect(err).Should(HaveOccurred())

		out.Reset()
		p, err = PrinterFromFlag("jsonpath-file=./test/unknown.jsonpath", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).Should(HaveOccurred())

		out.Reset()
		p, err = PrinterFromFlag("jsonpath-file=./test/valid.jsonpath", "", "")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(p.Fprint(&out, []Foo{foo})).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`Foo!`))
	})

})
