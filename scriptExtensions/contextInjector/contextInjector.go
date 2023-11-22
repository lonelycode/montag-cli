package contextInjector

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
	"github.com/lonelycode/montag-cli/interfaces"
)

type ContextInjector struct {
	prompt *interfaces.Prompt
}

func NewContextInjector(prompt *interfaces.Prompt) *ContextInjector {
	return &ContextInjector{
		prompt: prompt,
	}
}

func (*ContextInjector) TypeName() string {
	return "ContextInjector"
}
func (*ContextInjector) String() string {
	return "ContextInjector"
}
func (*ContextInjector) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*ContextInjector) IsFalsy() bool {
	return false
}
func (*ContextInjector) Equals(another tengo.Object) bool {
	return false
}

func (*ContextInjector) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*ContextInjector) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*ContextInjector) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*ContextInjector) CanIterate() bool {
	return false
}

func (*ContextInjector) CanCall() bool {
	return true
}

func (a *ContextInjector) addToContextHistory(title, content string) {
	a.prompt.ContextTitles = append(a.prompt.ContextTitles, title)
	a.prompt.ContextToRender = append(a.prompt.ContextToRender, content)
}

func (a *ContextInjector) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func (a *ContextInjector) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 argument, got %d", len(args))
	}

	title := a.cleanString(args[0].String())
	content := a.cleanString(args[1].String())

	a.addToContextHistory(title, content)
	return &tengo.Int{Value: 1}, err
}

func (*ContextInjector) Copy() tengo.Object {
	return &ContextInjector{}
}
