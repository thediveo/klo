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
