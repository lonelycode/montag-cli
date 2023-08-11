package dummyFunc

import (
	"fmt"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type DummyFunc struct {
	Name      string
	ReturnVal tengo.Object
}

func NewDummyFunc(name string, retVal tengo.Object) *DummyFunc {
	return &DummyFunc{
		Name:      name,
		ReturnVal: retVal,
	}
}

func (d *DummyFunc) TypeName() string {
	return d.Name
}
func (d *DummyFunc) String() string {
	return d.Name
}
func (*DummyFunc) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*DummyFunc) IsFalsy() bool {
	return false
}
func (*DummyFunc) Equals(another tengo.Object) bool {
	return false
}

func (*DummyFunc) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*DummyFunc) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*DummyFunc) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*DummyFunc) CanIterate() bool {
	return false
}

func (*DummyFunc) CanCall() bool {
	return true
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (a *DummyFunc) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	fmt.Printf("calling stub for func: '%s'\n", a.Name)
	return a.ReturnVal, nil
}

func (*DummyFunc) Copy() tengo.Object {
	return &DummyFunc{}
}
