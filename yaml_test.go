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

var _ = Describe("YAML printer", func() {

	It("prints YAML", func() {
		p := GoodPrinter(NewYAMLPrinter())
		PrinterPass(p, struct{ Foo string }{Foo: "bar"}, `Foo: bar
`)
	})

	It("handles YAML failures", func() {
		// For those 100% code coverage aficionados, let's test that we
		// correctly handle making YAML marshalling fail and that we correctly
		// handle it in the printer ... luckily,
		// https://stackoverflow.com/a/33964549 has the answer as to how make
		// JSON marshalling fail, which in turn makes the YAML marshaller fail
		// which we use, kind of a chain reaction...
		p := GoodPrinter(NewYAMLPrinter())
		PrinterFail(p, make(chan struct{}))
	})

})
