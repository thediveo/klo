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

package testutil

import (
	"bytes"
	"io"
)

// NewTestWriter returns a new io.Writer which translates spaces into
// underscores and adornes new lines with a down-leftwards pointing arrow, and
// finally writing the result to a chained io.Writer.
func NewTestWriter(w io.Writer) *TestWriter {
	return &TestWriter{w: w}
}

// TestWriter chains to another io.Writer and replaces certain characters as
// they pass through it for writing.
type TestWriter struct {
	w io.Writer
}

// Write replaces all spaces (0x20) with underscores and adorns all new lines
// \n with a "downwards arrow with corner leftwards" before writing the
// outcome to the chained writer.
func (nl *TestWriter) Write(p []byte) (n int, err error) {
	s := bytes.ReplaceAll(
		bytes.ReplaceAll(p, []byte("\n"), []byte("↵\n")),
		[]byte(" "), []byte("_"))
	n, err = nl.w.Write(s)
	// This is an ugly hack to make this working: as the string to be written
	// might grow due to the \n adornments, yet several callers won't take
	// kindly to getting reported more bytes actually written than was
	// originally specified to be sent, we simply clamp the amount of bytes
	// written to the original amount.
	if n > len(p) {
		n = len(p)
	}
	return
}
