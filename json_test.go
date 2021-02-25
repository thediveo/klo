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

var _ = Describe("JSON printer", func() {

	It("prints JSON", func() {
		p := GoodPrinter(NewJSONPrinter())
		PrinterPass(p, struct{ Foo string }{Foo: "bar"}, `{
    "Foo": "bar"
}
`)
	})

	It("handles JSON failures", func() {
		// For those 100% code coverage aficionados, let's test that we
		// correctly handle making JSON marshalling fail and that we correctly
		// handle it in the printer ... luckily,
		// https://stackoverflow.com/a/33964549 has the answer as to how make it
		// fail.
		p := GoodPrinter(NewJSONPrinter())
		Expect(p.Fprint(nil, make(chan struct{}))).ShouldNot(Succeed())
	})

})
