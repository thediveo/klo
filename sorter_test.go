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
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("-o output options", func() {

	It("doesn't accept botched expressions", func() {
		_, err := NewSortingPrinter("{.A", nil)
		Expect(err).Should(HaveOccurred())
		_, err = NewSortingPrinter("{.A}", nil)
		Expect(err).Should(HaveOccurred())
	})

	It("compares i<j reflection values", func() {
		type test struct {
			i, j    reflect.Value
			outcome bool
		}

		v1 := "666"
		v2 := 42

		tests := []test{
			test{reflect.ValueOf(int8(42)), reflect.ValueOf(int8(43)), true},
			test{reflect.ValueOf(int8(127)), reflect.ValueOf(int8(43)), false},
			test{reflect.ValueOf(int16(42)), reflect.ValueOf(int16(43)), true},
			test{reflect.ValueOf(int32(42)), reflect.ValueOf(int32(43)), true},
			test{reflect.ValueOf(int64(42)), reflect.ValueOf(int64(43)), true},
			test{reflect.ValueOf(int8(42)), reflect.ValueOf(int64(666)), true},
			test{reflect.ValueOf(int8(42)), reflect.ValueOf(float32(42.5)), true},
			test{reflect.ValueOf(int8(42)), reflect.ValueOf(float64(42.5)), true},

			test{reflect.ValueOf(uint8(42)), reflect.ValueOf(uint8(43)), true},
			test{reflect.ValueOf(uint8(255)), reflect.ValueOf(uint8(43)), false},
			test{reflect.ValueOf(uint16(42)), reflect.ValueOf(uint16(43)), true},
			test{reflect.ValueOf(uint32(42)), reflect.ValueOf(uint32(43)), true},
			test{reflect.ValueOf(uint64(42)), reflect.ValueOf(uint64(43)), true},
			test{reflect.ValueOf(uint8(42)), reflect.ValueOf(uint64(666)), true},
			test{reflect.ValueOf(uint8(42)), reflect.ValueOf(float32(42.5)), true},
			test{reflect.ValueOf(uint8(42)), reflect.ValueOf(float64(42.5)), true},

			test{reflect.ValueOf(float32(42.9)), reflect.ValueOf(float32(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(float64(43)), true},
			test{reflect.ValueOf(float32(42.9)), reflect.ValueOf(int8(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int16(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int32(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(int64(43)), true},
			test{reflect.ValueOf(float32(42.9)), reflect.ValueOf(uint8(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint16(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint32(43)), true},
			test{reflect.ValueOf(float64(42.9)), reflect.ValueOf(uint64(43)), true},

			test{reflect.ValueOf("bar"), reflect.ValueOf("foo"), true},
			test{reflect.ValueOf("foo"), reflect.ValueOf("bar"), false},
			test{reflect.ValueOf("bar666"), reflect.ValueOf("bar42"), false},

			test{reflect.ValueOf("666"), reflect.ValueOf(int16(42)), false},

			test{reflect.ValueOf(&v1), reflect.ValueOf(&v2), false},
		}

		for _, t := range tests {
			Expect(reflectedLess(t.i, t.j)).Should(
				Equal(t.outcome),
				fmt.Sprintf("%s(%+v)!<%s(%+v)",
					reflect.TypeOf(t.i.Interface()).String(), t.i, reflect.TypeOf(t.j.Interface()).String(), t.j))
		}
	})

	It("sorts before printing", func() {
		type row struct {
			A string
			B int
		}
		table := []row{
			row{A: "foo", B: 666},
			row{A: "bar", B: 42},
			row{A: "aaa", B: 420},
		}

		ccp, err := NewCustomColumnsPrinterFromSpec("A:{.A},B:{.B}")
		Expect(err).ShouldNot(HaveOccurred())
		sp, err := NewSortingPrinter("{.A}", ccp)
		Expect(err).ShouldNot(HaveOccurred())

		var out bytes.Buffer
		Expect(sp.Fprint(&out, table)).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`A    B
aaa  420
bar  42
foo  666
`))

		out.Reset()
		sp, err = NewSortingPrinter("{.B}", ccp)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(sp.Fprint(&out, table)).ShouldNot(HaveOccurred())
		Expect(out.String()).Should(Equal(`A    B
bar  42
aaa  420
foo  666
`))

		out.Reset()
		Expect(sp.Fprint(&out, &table)).ShouldNot(HaveOccurred())
	})

})
