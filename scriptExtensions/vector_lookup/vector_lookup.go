package vectorlookup

import (
	"fmt"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"

	"github.com/lonelycode/montag-cli/client"
)

type VectorLookup struct {
	client *client.Client
}

func NewVectorLookup(c *client.Client) *VectorLookup {
	return &VectorLookup{
		client: c,
	}
}

func (*VectorLookup) TypeName() string {
	return "VectorLookup"
}
func (*VectorLookup) String() string {
	return "VectorLookup"
}
func (*VectorLookup) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*VectorLookup) IsFalsy() bool {
	return false
}
func (*VectorLookup) Equals(another tengo.Object) bool {
	return false
}

func (*VectorLookup) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*VectorLookup) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*VectorLookup) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*VectorLookup) CanIterate() bool {
	return false
}

func (*VectorLookup) CanCall() bool {
	return true
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (a *VectorLookup) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("expected 3 arguments, got %d", len(args))
	}

	botID := int(args[0].(*tengo.Int).Value)
	numResults := int(args[1].(*tengo.Int).Value)
	query := a.cleanString(args[2].String())
	out, err := a.search(botID, numResults, query)

	if err != nil {
		return &tengo.String{Value: ""}, err
	}

	return out, nil
}

func (*VectorLookup) Copy() tengo.Object {
	return &VectorLookup{}
}

func (a *VectorLookup) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func (a *VectorLookup) search(botID, numResults int, query string) (tengo.Object, error) {
	// Search the vector DB for the query
	results, err := a.client.VectorSearch(botID, query, numResults)
	if err != nil {
		return nil, err
	}

	tengoReturn := &tengo.Array{Value: []tengo.Object{}}
	for _, result := range results {
		tengoMap := &tengo.Map{Value: map[string]tengo.Object{}}
		tengoMap.Value["title"] = &tengo.String{Value: result.Metadata["title"]}
		tengoMap.Value["text"] = &tengo.String{Value: result.Metadata["text"]}

		tengoReturn.Value = append(tengoReturn.Value, tengoMap)
	}

	return tengoReturn, nil

}
