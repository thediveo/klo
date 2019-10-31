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
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKlo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "klo suite")
}

// PrinterPass checks that a ValuePrinter correctly renders the expected
// output.
func PrinterPass(p ValuePrinter, v interface{}, expected string) {
	var out bytes.Buffer
	ExpectWithOffset(1, p.Fprint(&out, v)).ShouldNot(HaveOccurred())
	ExpectWithOffset(1, out.String()).Should(Equal(expected))
}

// PrinterFail expects the ValuePrinter to correctly fail.
func PrinterFail(p ValuePrinter, v interface{}) {
	var out bytes.Buffer
	ExpectWithOffset(1, p.Fprint(&out, v)).Should(HaveOccurred())
}

// Asserts that creation of a ValuePrinter succeeded.
func GoodPrinter(p ValuePrinter, err error) ValuePrinter {
	ExpectWithOffset(1, err).ShouldNot(HaveOccurred(), "printer creation unexpectedly failed")
	return p
}

// Asserts that creation of a ValuePrinter failed.
func BadPrinter(p ValuePrinter, err error) {
	ExpectWithOffset(1, err).Should(HaveOccurred(), "printer creation should not have succeeded")
}
