package testutil

import (
	. "github.com/onsi/gomega"
)

type PASSFAILS []PASSFAIL
type PASSFAIL interface {
	Description() string
	Actual() interface{}
}

// FAIL this testcase.
type FAIL struct {
	D string // testcase description
	A interface{}
}

func (f FAIL) Description() string { return f.D }
func (f FAIL) Actual() interface{} { return f.A }

// PASS is just another type of FAIL, after all. (No, this hasn't been ripped
// off straight from one of Nico Semsrott's satirical revues.)
type PASS FAIL

func (p PASS) Description() string { return p.D }
func (p PASS) Actual() interface{} { return p.A }

func PassFail(tests PASSFAILS) {
	for _, t := range tests {
		if _, ok := t.(PASS); ok {
			ExpectWithOffset(1, t.Actual()).
				Should(Succeed(), t.Description())
		} else {
			ExpectWithOffset(1, t.Actual()).
				ShouldNot(Succeed(), t.Description())
		}
	}
}
