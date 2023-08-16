package scriptBytesHttpCaller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type BytesHttpCaller struct {
}

func NewBytesHttpCaller() *BytesHttpCaller {
	return &BytesHttpCaller{}
}

func (*BytesHttpCaller) TypeName() string {
	return "BytesHttpCaller"
}
func (*BytesHttpCaller) String() string {
	return "BytesHttpCaller"
}
func (*BytesHttpCaller) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*BytesHttpCaller) IsFalsy() bool {
	return false
}
func (*BytesHttpCaller) Equals(another tengo.Object) bool {
	return false
}

func (*BytesHttpCaller) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*BytesHttpCaller) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*BytesHttpCaller) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*BytesHttpCaller) CanIterate() bool {
	return false
}

func (*BytesHttpCaller) CanCall() bool {
	return true
}

func makeHttpRequest(method string, url string, headers map[string]string, body string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Printf("client: could not create request: %s\n", err)
		return []byte{}, 500, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client: could not send request: %s\n", err)
		return []byte{}, 500, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, 500, err
	}

	return b, resp.StatusCode, nil
}

func (a *BytesHttpCaller) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func (a *BytesHttpCaller) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) != 4 {
		return nil, fmt.Errorf("expected 4 arguments, got %d", len(args))
	}

	method := a.cleanString(args[0].String())
	url := a.cleanString(args[1].String())
	headers := args[2].(*tengo.Map)
	body := a.cleanString(args[3].String())

	headersMap := make(map[string]string)
	if headers != nil {
		for key, value := range headers.Value {
			headersMap[a.cleanString(key)] = a.cleanString(value.String())
		}
	}

	out, status, err := makeHttpRequest(method, url, headersMap, body)
	return &tengo.Map{
		Value: map[string]tengo.Object{
			"status":   &tengo.Int{Value: int64(status)},
			"response": &tengo.Bytes{Value: out},
		},
	}, err
}

func (*BytesHttpCaller) Copy() tengo.Object {
	return &BytesHttpCaller{}
}
