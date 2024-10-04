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
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/jsonpath"
)

var _ = Describe("-o output options", func() {

	It("doesn't accept botched JSONPath expressions for sorting", func() {
		BadPrinter(NewSortingPrinter("{.A", nil))
		BadPrinter(NewSortingPrinter("{.A}", nil))
	})

	It("compares i<j reflection values", func() {
		type test struct {
			i, j    reflect.Value
			outcome bool
		}

		v1 := "666"
		v2 := 42

		tests := []test{
			{reflect.ValueOf(int8(42)), reflect.ValueOf(int8(43)), true},
			{reflect.ValueOf(int8(127)), reflect.ValueOf(int8(43)), false},
			{reflect.ValueOf(int16(42)), reflect.ValueOf(int16(43)), true},
			{reflect.ValueOf(int32(42)), reflect.ValueOf(int32(43)), true},
			{reflect.ValueOf(int64(42)), reflect.ValueOf(int64(43)), true},
			{reflect.ValueOf(int8(42)), reflect.ValueOf(int64(666)), true},
			{reflect.ValueOf(int8(42)), reflect.ValueOf(float32(42.5)), true},
			{reflect.ValueOf(int8(42)), reflect.ValueOf(float64(42.5)), true},

			{reflect.ValueOf(uint8(42)), reflect.ValueOf(uint8(43)), true},
			{reflect.ValueOf(uint8(255)), reflect.ValueOf(uint8(43)), false},
			{reflect.ValueOf(uint16(42)), reflect.ValueOf(uint16(43)), true},
			{reflect.ValueOf(uint32(42)), reflect.ValueOf(uint32(43)), true},
			{reflect.ValueOf(uint64(42)), reflect.ValueOf(uint64(43)), true},
			{reflect.ValueOf(uint8(42)), reflect.ValueOf(uint64(666)), true},
			{reflect.ValueOf(uint8(42)), reflect.ValueOf(float32(42.5)), true},
			{reflect.ValueOf(uint8(42)), reflect.ValueOf(float64(42.5)), true},

			{reflect.ValueOf(float32(42.9)), reflect.ValueOf(float32(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(float64(43)), true},
			{reflect.ValueOf(float32(42.9)), reflect.ValueOf(int8(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int16(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int32(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int64(43)), true},
			{reflect.ValueOf(float32(42.9)), reflect.ValueOf(uint8(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint16(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint32(43)), true},
			{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint64(43)), true},

			{reflect.ValueOf("bar"), reflect.ValueOf("foo"), true},
			{reflect.ValueOf("foo"), reflect.ValueOf("bar"), false},
			{reflect.ValueOf("bar666"), reflect.ValueOf("bar42"), false},

			{reflect.ValueOf("666"), reflect.ValueOf(int16(42)), false},

			{reflect.ValueOf(&v1), reflect.ValueOf(&v2), false},
		}

		for _, t := range tests {
			Expect(reflectedLess(t.i, t.j)).Should(
				Equal(t.outcome),
				fmt.Sprintf("comparing %s(%+v)<%s(%+v) failed",
					reflect.TypeOf(t.i.Interface()).String(), t.i,
					reflect.TypeOf(t.j.Interface()).String(), t.j))
		}
	})

	It("sorts before printing", func() {
		type row struct {
			A string
			B int
		}
		table := []row{
			{A: "foo", B: 666},
			{A: "bar", B: 42},
			{A: "aaa", B: 420},
		}

		ccp := GoodPrinter(NewCustomColumnsPrinterFromSpec("A:{.A},B:{.B}"))
		PrinterPass(GoodPrinter(NewSortingPrinter("{.A}", ccp)), table,
			`A    B
aaa  420
bar  42
foo  666
`)
		PrinterPass(GoodPrinter(NewSortingPrinter("{.B}", ccp)), table,
			`A    B
bar  42
aaa  420
foo  666
`)
		PrinterPass(GoodPrinter(NewSortingPrinter("{.B}", ccp)), &table,
			`A    B
bar  42
aaa  420
foo  666
`)
		table = []row{
			{A: "foo", B: 42},
			{A: "bar", B: 42},
			{A: "aaa", B: 420},
		}
		PrinterPass(GoodPrinter(NewSortingPrinter("{.B}{'/'}{.A}", ccp)), table,
			`A    B
bar  42
foo  42
aaa  420
`)
		type anotherrow struct {
			A string
		}
		othertable := []anotherrow{
			{A: "foo"},
		}
		PrinterPass(GoodPrinter(NewSortingPrinter("{.A}", ccp)), &othertable,
			`A    B
foo  <none>
`)
	})

	It("simply passes on non-sliced things", func() {
		r := struct {
			A string
			B int
		}{A: "foo", B: 42}
		ccp := GoodPrinter(NewCustomColumnsPrinterFromSpec("A:{.A},B:{.B}"))
		PrinterPass(GoodPrinter(NewSortingPrinter("{.A}", ccp)), reflect.ValueOf(r),
			`A    B
foo  42
`)
	})

	It("handles deeper errors in sort expression evaluation", func() {
		type row struct {
			A string
			B int
		}
		table := []row{
			{A: "foo", B: 666},
			{A: "bar", B: 42},
			{A: "aaa", B: 420},
		}
		ccp := GoodPrinter(NewCustomColumnsPrinterFromSpec("A:{.A},B:{.B}"))
		p := GoodPrinter(NewSortingPrinter("{.A}", ccp))
		sp := p.(*SortingPrinter)
		sp.SortExpr = jsonpath.New("zero")
		PrinterFail(sp, table)
	})

	It("prints <none> for JSONPaths to non-existing fields ", func() {
		type row struct {
			A []string
		}
		table := []row{
			{},
		}
		ccp := GoodPrinter(NewCustomColumnsPrinterFromSpec("A:{.A[*]}"))
		PrinterPass(GoodPrinter(NewSortingPrinter("", ccp)), &table,
			`A
<none>
`)
	})

})
