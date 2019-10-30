package klo_test

import (
	"os"

	"github.com/thediveo/klo"
)

func Example_sortedtable() {
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
	table.Fprint(os.Stdout, list)
	// This will output:
	// NAME        FOO  BAR
	// Another Two 123  Bar
	// One         42
	// Two         666  Bar
}
