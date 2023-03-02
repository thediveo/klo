# klo

[![Go Reference](https://pkg.go.dev/badge/github.com/thediveo/klo.svg)](https://pkg.go.dev/github.com/thediveo/klo)
![GitHub](https://img.shields.io/github/license/thediveo/go-asciitree)
![build and test](https://github.com/TheDiveO/klo/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/klo)](https://goreportcard.com/report/github.com/thediveo/klo)
![Coverage](https://img.shields.io/badge/Coverage-95.2%25-brightgreen)

`klo` is a Go package for `kubectl`-like output of Go values (such as structs,
maps, et cetera) in several output formats. You might want to use this package
in your CLI tools to easily offer `kubectl`-like output formatting to your
Kubernetes-spoiled users.

[The Only True Go Way](https://golang.org/doc/effective_go.html#package-names)
mandates package names to be short, concise, evocative. Thus, the package name
`klo` was chosen to represent the essence of **`k`**`ubectl`-**l**ike
**o**utput in a manner as short, concise and evocative as ever. If you happen
to have sanitary associations, then better flush them now.

## Supported Output Formats

The following output formats are supported:

- ASCII columns, which optionally can be customized (`-o custom-columns=` and
  `-o custom-columns-file=`).
  - optional sorting by specific column(s) using JSONPath expressions.
- JSON and JSONPath-customized (`-o json`, `-o jsonpath=`, and `-o
  jsonpath-file=`).
- YAML (`-o yaml`).
- Go templates (`-o go-template=` and `-o go-template-file=`).

> **Note:** `-o name` and `-o wide` are application-specific and are basically
> customized ASCII column formats, with just varying custom column
> configurations. Thus, they can be easily implemented in your appliaction
> itself and then use the existing `klo` package features.

In addition, sorting is supported by wrapping an output-format printer into a
sorting printer. This allows to sort the rows in a custom-columns output based
on row values taken from one or even multiple columns.

## Basic Usage

The following code example prints a table with multiple columns, and the rows
sorted by the NAME column.

```go
import (
    "os"
    "github.com/thediveo/klo"
)

func main() {
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
    prn, err := klo.PrinterFromFlag("",
        &klo.Specs{DefaultColumnSpec: "NAME:{.Name},FOO:{.Foo},BAR:{.Bar}"})
    if err != nil {
        panic(err)
    }
    // Use a table sorter and tell it to sort by the Name field of our column objects.
    table, err := klo.NewSortingPrinter("{.Name}", prn)
    if err != nil {
        panic(err)
    }
    table.Fprint(os.Stdout, list)
```

This will output:

```text
NAME        FOO  BAR
Another Two 123  Bar
One         42
Two         666  Bar
```

## -o Usage

For supporting "-o" output format control via CLI args, choose any CLI arg
handling package you like, such as flag, pflag, cobra, et cetera. Then, call
`PrinterFromFlag(oflagvalue, &myspecs)`, where `oflagvalue` is the set/default
value of the CLI arg you use for controlling the output format in your own app's
CLI. Your `myspecs` should specify the default custom-columns format, and
optionally a wide custom-clumns format variant. If you support go templates for
output formatting, then you should also pass in the value of your `--template=`
CLI arg.

```go
import (
    "github.com/thediveo/klo"
)

func main() {
    // Get your -o and -template flag values depending on your CLI arg toolkit.
    templateflagvalue := ""
    oflagvalue := "wide"
    // Set up the specs and get a suitable output formatting printer according
    // to the specific output format choosen and the auxiliary information given
    // on specs and an optional Go template arg.
    myspecs := klo.Specs{
        DefaultColumnSpec: "FOO:{.Foo}",
        WideColumnSpec: "FOO:{.Foo},BAR:{.Bar}",
        GoTemplateArg: templateflagvalue,
    }
    prn, err := PrinterFromFlag(oflagvalue, &myspecs)
    //...
}
```

## Copyright and License

`klo` is Copyright 2019 Harald Albrecht, and licensed under the [Apache
License, Version
2.0](https://github.com/TheDiveO/go-mntinfo/blob/master/LICENSE).
