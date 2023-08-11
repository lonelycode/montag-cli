package secretGetter

import (
	"fmt"
	"os"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type SecretGetter struct {
}

func NewSecretGetter() *SecretGetter {
	return &SecretGetter{}
}

func (d *SecretGetter) TypeName() string {
	return "SecretGetter"
}
func (d *SecretGetter) String() string {
	return "SecretGetter"
}
func (*SecretGetter) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*SecretGetter) IsFalsy() bool {
	return false
}
func (*SecretGetter) Equals(another tengo.Object) bool {
	return false
}

func (*SecretGetter) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*SecretGetter) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*SecretGetter) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*SecretGetter) CanIterate() bool {
	return false
}

func (*SecretGetter) CanCall() bool {
	return true
}

func (a *SecretGetter) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (a *SecretGetter) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SecretGetter requires 1 argument, got %d", len(args))
	}

	key := a.cleanString(args[0].(*tengo.String).Value)
	asEnvKey := fmt.Sprintf("SECRET_%s", strings.ToUpper(key))
	fmt.Printf("SecretGetter looking for ENV var: %s\n", asEnvKey)

	t := os.Getenv(asEnvKey)
	return &tengo.String{Value: t}, nil
}

func (*SecretGetter) Copy() tengo.Object {
	return &SecretGetter{}
}
