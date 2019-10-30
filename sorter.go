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
	"errors"
	"io"
	"reflect"
	"sort"

	"k8s.io/client-go/util/jsonpath"
	"vbom.ml/util/sortorder"
)

// SortingPrinter sorts slice values first, before it writes them to the
// next printer in the chain.
type SortingPrinter struct {
	ChainedPrinter ValuePrinter       // Next ValuePrinter we chain to.
	SortExpr       *jsonpath.JSONPath // Compiled JSONPath expression.
	raw            string             // Original JSONPath expression, to ease debugging.
}

// NewSortingPrinter returns a printer for outputting values in YAML format.
func NewSortingPrinter(expr string, p ValuePrinter) (ValuePrinter, error) {
	jp := jsonpath.New("sort")
	if err := jp.Parse(expr); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("nil ValuePrint to chain (hint: we cannot)")
	}
	return &SortingPrinter{
		ChainedPrinter: p,
		SortExpr:       jp,
		raw:            expr,
	}, nil
}

// Fprint first sorts values according to a JSONPath expression used for
// sorting, then chains to the next ValuePrinter for printing.
func (sp *SortingPrinter) Fprint(w io.Writer, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return sp.ChainedPrinter.Fprint(w, val.Interface())
	}
	//
	slicelen := val.Len()
	index := keyedItems{
		keys:  make([]reflect.Value, slicelen),
		items: make([]reflect.Value, slicelen),
	}
	for idx := 0; idx < slicelen; idx++ {
		index.items[idx] = val.Index(idx)
		key, err := sp.SortExpr.FindResults(index.items[idx].Interface())
		if err != nil {
			return err
		}
		// Depending on the JSONPath expression, the key for this item (column)
		// might consist of multiple values, or even none at all.
		if len(key) == 0 || len(key[0]) == 0 {
			index.keys[idx] = reflect.ValueOf("<none>")
		} else if len(key) == 1 && len(key[0]) == 1 {
			index.keys[idx] = key[0][0]
		} else {
			index.keys[idx] = reflect.ValueOf(key)
		}
	}
	sort.Sort(index)
	// That's it: hand over the sorted items to the chained printer so it can
	// carry out its part of the job.
	return sp.ChainedPrinter.Fprint(w, index.items)
}

//
type keyedItems struct {
	keys  []reflect.Value // results of evaluating JSONPath expressions.
	items []reflect.Value // references the items slice to be sorted.
}

// Returns number of items (interface sort.Interface).
func (ki keyedItems) Len() int { return len(ki.keys) }

// Compares two items by their keys (interface sort.Interface).
func (ki keyedItems) Less(i, j int) bool { return reflectedLess(ki.keys[i], ki.keys[j]) }

// Swaps two child nodes (interface sort.Interface). Swapping the item
// references that are in the form of reflect.Values is unfortunately slightly
// involved, as a simple swap without an intermediate temporary would fail.
func (ki keyedItems) Swap(i, j int) {
	ki.keys[i], ki.keys[j] = ki.keys[j], ki.keys[i]
	// We cannot simply swap two elements in a slice if they are more
	// intricate types, such as strings, interfaces, maps, et cetera, as
	// opposed to ints. See also: https://github.com/golang/go/issues/3126
	// Instead, we need to dance around and sacrifice the Go(ds) of
	// reflection.
	temp := reflect.New(ki.items[i].Type()).Elem()
	temp.Set(ki.items[i])
	ki.items[i].Set(ki.items[j])
	ki.items[j].Set(temp)
}

// reflectedLess compares two values and returns true if i<j. In case of
// incompatible values, it resorts to comparing the string representations of
// the values instead. Since we have no idea as how to define a sorting order on
// arrays, slices, structs, et cetera, we just consider them incompatible too,
// and resort to sorting their string representation, whatever sense that might
// make.
//
// Oh, and in contrast to kubectl's isLess() version, we don't panic, because
// that's really not nice in the face of CLI users.
func reflectedLess(i, j reflect.Value) bool {
	// First, follow pointers, so we don't need to care about them later...
	for i.Kind() == reflect.Ptr {
		i = i.Elem()
	}
	for j.Kind() == reflect.Ptr {
		j = j.Elem()
	}
	// Now let's compare...
	switch i.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch j.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return i.Int() < j.Int()
		case reflect.Float32, reflect.Float64:
			return float64(i.Int()) < j.Float()
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch j.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return i.Uint() < j.Uint()
		case reflect.Float32, reflect.Float64:
			return float64(i.Uint()) < j.Float()
		}
	case reflect.Float32, reflect.Float64:
		switch j.Kind() {
		case reflect.Float32, reflect.Float64:
			return i.Float() < j.Float()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return i.Float() < float64(j.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return i.Float() < float64(j.Uint())
		}
	case reflect.String:
		if j.Kind() == reflect.String {
			return sortorder.NaturalLess(i.String(), j.String())
		}
	}
	// We've fallen through because the two values to be compared are of
	// incompatible types, so let's compare their stringified values instead.
	return sortorder.NaturalLess(reflect.ValueOf(i).String(), reflect.ValueOf(j).String())
}
