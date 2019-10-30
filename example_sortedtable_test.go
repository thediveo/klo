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

package klo_test

import (
	"os"

	"github.com/thediveo/klo"
	"github.com/thediveo/klo/testutil"
)

func Example_sortedtable() {
	/* ignore/for testing */ out := testutil.NewTestWriter(os.Stdout) /* end ignore */
	// Some data structure we want to print in tables.
	type myobj struct {
		Name string
		Foo  int
		Bar  string
	}
	// A slice of objects we want to print as a table with custom columns.
	list := []myobj{
		myobj{Name: "One", Foo: 42},
		myobj{Name: "Two", Foo: 666, Bar: "Bar"},
		myobj{Name: "Another Two", Foo: 123, Bar: "Bar"},
	}
	// Create a table printer with custom columns, to be filled from fields
	// of the objects (namely, Name, Foo, and Bar fields).
	prn, err := klo.PrinterFromFlag("", "NAME:{.Name},FOO:{.Foo},BAR:{.Bar}", "")
	if err != nil {
		panic(err)
	}
	// Use a table sorter and tell it to sort by the Name field of our column objects.
	table, err := klo.NewSortingPrinter("{.Name}", prn)
	if err != nil {
		panic(err)
	}
	table.Fprint(out, list)
	// Output:
	// NAME________FOO__BAR↵
	// Another_Two_123__Bar↵
	// One_________42___↵
	// Two_________666__Bar↵
}
