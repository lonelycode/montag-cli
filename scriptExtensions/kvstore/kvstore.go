package kvstore

import (
	"fmt"
	"log"
	"strings"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
	"github.com/lonelycode/montag-cli/models"
	"gorm.io/gorm"
)

type KVStoreSetter struct {
	db       *gorm.DB
	scriptID uint
}

func NewKVStore(db *gorm.DB, scriptID uint) *KVStoreSetter {
	return &KVStoreSetter{
		db:       db,
		scriptID: scriptID,
	}
}

func (*KVStoreSetter) TypeName() string {
	return "KVStoreSetter"
}
func (*KVStoreSetter) String() string {
	return "KVStoreSetter"
}
func (*KVStoreSetter) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	panic("not implemented")
}

func (*KVStoreSetter) IsFalsy() bool {
	return false
}
func (*KVStoreSetter) Equals(another tengo.Object) bool {
	return false
}

func (*KVStoreSetter) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	panic("not implemented")
}

func (*KVStoreSetter) IndexSet(index, value tengo.Object) error {
	panic("not implemented")
}

func (*KVStoreSetter) Iterate() tengo.Iterator {
	panic("not implemented")
}

func (*KVStoreSetter) CanIterate() bool {
	return false
}

func (*KVStoreSetter) CanCall() bool {
	return true
}

func (a *KVStoreSetter) setValue(key, value string) error {
	err := models.SetScriptKVStore(a.db, a.scriptID, key, value)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *KVStoreSetter) getValue(key string) (*tengo.Map, error) {
	kv, err := models.GetScriptKVStore(a.db, a.scriptID, key)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ret := &tengo.Map{
		Value: map[string]tengo.Object{
			"result": &tengo.String{Value: kv.Value},
			"error":  &tengo.String{Value: ""},
		},
	}

	return ret, nil

}

func (a *KVStoreSetter) cleanString(str string) string {
	return strings.Trim(str, "\"")
}

func (a *KVStoreSetter) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("expected at least 2 argument, got %d", len(args))
	}

	op := a.cleanString(args[0].String())
	key := a.cleanString(args[1].String())

	value := ""
	if op == "set" {
		value = a.cleanString(args[2].String())
	}

	switch strings.ToLower(op) {
	case "set":
		return nil, a.setValue(key, value)
	case "get":
		return a.getValue(key)
	default:
		return nil, fmt.Errorf("unknown operation: %s", op)
	}
}

func (*KVStoreSetter) Copy() tengo.Object {
	return &KVStoreSetter{}
}
