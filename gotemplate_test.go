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
	"strconv"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Go template printer", func() {

	It("gracefully fails on invalid templates", func() {
		BadPrinter(NewGoTemplatePrinter(`{{oops}}`))
	})

	It("catches template execution panics", func() {
		p := GoodPrinter(NewGoTemplatePrinter(`{{"oops"}}`))
		ExpectWithOffset(1, p.Fprint(nil, nil)).Should(HaveOccurred())
	})

	It("templates", func() {
		p := GoodPrinter(NewGoTemplatePrinter(`{{range .}}{{.}}{{println}}{{end}}`))
		PrinterPass(p, []string{"foo", "bar"}, `foo
bar
`)
	})

})

var _ = Describe("Go template printer with funcs", func() {

	It("gracefully fails on invalid templates", func() {
		BadPrinter(NewGoTemplatePrinterWithFuncs(`{{oops}}`, nil))
	})

	It("catches template execution panics", func() {
		p := GoodPrinter(NewGoTemplatePrinterWithFuncs(`{{"oops"}}`, nil))
		ExpectWithOffset(1, p.Fprint(nil, nil)).Should(HaveOccurred())
	})

	It("templates", func() {
		p := GoodPrinter(NewGoTemplatePrinterWithFuncs(`{{range .}}{{.}}{{println}}{{end}}`, nil))
		PrinterPass(p, []string{"foo", "bar"}, `foo
bar
`)
	})

	It("gracefully fails for missing template functions", func() {
		BadPrinter(NewGoTemplatePrinterWithFuncs(`{{range .}}{{ $intValue := atoi . }}{{ add $intValue 1 }}{{end}}`, template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
		}))
	})

	It("template funcs", func() {
		p := GoodPrinter(NewGoTemplatePrinterWithFuncs(`{{range .}}{{ $intValue := atoi . }}{{ add $intValue 1 }}{{end}}`, template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"atoi": func(s string) int {
				i, _ := strconv.Atoi(s)
				return i
			},
		}))
		PrinterPass(p, []string{"10"}, `11`)
	})

})
