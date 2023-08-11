package readable

import (
	"fmt"
	"strings"
	"time"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
	"github.com/go-shiori/go-readability"
)

type Readable struct{}

func NewReadable() *Readable {
	return &Readable{}
}

func (*Readable) TypeName() string {
	return "Readable"
}
func (*Readable) String() string {
	return "Readable"
}
func (*Readable) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*Readable) IsFalsy() bool {
	return false
}
func (*Readable) Equals(another tengo.Object) bool {
	return false
}

func (*Readable) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*Readable) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*Readable) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*Readable) CanIterate() bool {
	return false
}

func (*Readable) CanCall() bool {
	return true
}

func getReadable(url string) (string, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return "", err
	}

	return article.TextContent, nil
}

func (a *Readable) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func (a *Readable) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arguments, got %d", len(args))
	}

	url := a.cleanString(args[0].String())
	out, err := getReadable(url)

	if err != nil {
		out = err.Error()
	}

	return &tengo.Map{
		Value: map[string]tengo.Object{
			"response": &tengo.String{Value: out},
		},
	}, err
}

func (*Readable) Copy() tengo.Object {
	return &Readable{}
}
