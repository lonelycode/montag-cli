package aifuncRunner

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"

	"github.com/lonelycode/montag-cli/client"
	"github.com/lonelycode/montag-cli/models"
)

type AIFuncRunner struct {
	client *client.Client
}

func NewAIFuncRunner(c *client.Client) *AIFuncRunner {
	return &AIFuncRunner{
		client: c,
	}
}

func (*AIFuncRunner) TypeName() string {
	return "AIFuncRunner"
}
func (*AIFuncRunner) String() string {
	return "AIFuncRunner"
}
func (*AIFuncRunner) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*AIFuncRunner) IsFalsy() bool {
	return false
}
func (*AIFuncRunner) Equals(another tengo.Object) bool {
	return false
}

func (*AIFuncRunner) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*AIFuncRunner) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*AIFuncRunner) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*AIFuncRunner) CanIterate() bool {
	return false
}

func (*AIFuncRunner) CanCall() bool {
	return true
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (a *AIFuncRunner) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	funcName := args[0].String()
	fmt.Println("Calling function", funcName)
	funcInput := args[1]
	fmt.Println("With input", funcInput)

	input, err := funcInput.IndexGet(&tengo.String{Value: "Input"})
	if err != nil {
		return nil, err
	}

	fmt.Println("Got input string", input)

	mtObj, err := funcInput.IndexGet(&tengo.String{Value: "Meta"})
	if err != nil {
		return nil, err
	}

	meta, ok := mtObj.(*tengo.Map)
	if !ok {
		return nil, fmt.Errorf("expected meta to be a map")
	}

	fmt.Println("Got meta", meta)

	call := &models.AIFuncCall{
		Input: a.cleanString(input.String()),
		Meta:  tengoMapToMap(meta.Value),
	}

	// this is odd
	fName := a.cleanString(funcName)

	out, err := a.CallAIFuncFromAPI(call, fName)
	if err != nil {
		return nil, err
	}

	return &tengo.String{Value: out.Response}, nil
}

func (*AIFuncRunner) Copy() tengo.Object {
	return &AIFuncRunner{}
}

func (a *AIFuncRunner) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func tengoMapToMap(m map[string]tengo.Object) map[string]interface{} {
	out := map[string]interface{}{}
	for key, value := range m {
		out[key] = value
	}
	return out
}

func (a *AIFuncRunner) CallAIFuncFromAPI(req *models.AIFuncCall, name string) (*models.AIFuncResponse, error) {
	return a.client.RunAIFunc(name, req)
}
