package snippetStore

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"

	"github.com/lonelycode/montag-cli/client"
)

type SnippetStore struct {
	client *client.Client
}

func NewSnippetStore(c *client.Client) *SnippetStore {
	return &SnippetStore{
		client: c,
	}
}

func (*SnippetStore) TypeName() string {
	return "SnippetStore"
}
func (*SnippetStore) String() string {
	return "SnippetStore"
}
func (*SnippetStore) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*SnippetStore) IsFalsy() bool {
	return false
}
func (*SnippetStore) Equals(another tengo.Object) bool {
	return false
}

func (*SnippetStore) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*SnippetStore) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*SnippetStore) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*SnippetStore) CanIterate() bool {
	return false
}

func (*SnippetStore) CanCall() bool {
	return true
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (a *SnippetStore) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arguments, got %d", len(args))
	}

	slug := a.cleanString(args[0].String())
	data, err := a.client.GetSnippet(slug)
	if err != nil {
		return nil, err
	}

	out := &tengo.String{Value: data}

	return out, nil
}

func (*SnippetStore) Copy() tengo.Object {
	return &SnippetStore{}
}

func (a *SnippetStore) cleanString(str string) string {
	return strings.Trim(str, "\"")
}
