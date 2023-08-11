package models

import (
	"github.com/d5/tengo/v2"
	"gorm.io/gorm"
)

type AIFuncCall struct {
	Input string                 `json:"input"`
	Meta  map[string]interface{} `json:"meta"`
}

type AIFuncResponse struct {
	ID       string `json:"id"`
	Response string `json:"response"`
}

type ScriptOutput struct {
	Response    string                 `json:"response"`
	Outputs     map[string]interface{} `json:"outputs"`
	ReturnQuery string                 `json:"return_query"`
}

func DummyMultiCallerResponse() *tengo.Map {
	r := &tengo.Map{
		Value: map[string]tengo.Object{},
	}

	r.Value["Name"] = &tengo.String{Value: "Func Name"}
	r.Value["Error"] = &tengo.String{Value: ""}
	r.Value["Output"] = &tengo.String{Value: "Foo"}

	return r
}

type ScriptKVStore struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	ScriptID uint
	Key      string
	Value    string
}

func GetScriptKVStore(db *gorm.DB, scriptID uint, key string) (*ScriptKVStore, error) {
	var kv ScriptKVStore
	err := db.Where("script_id = ? AND key = ?", scriptID, key).First(&kv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			nullValue := &ScriptKVStore{
				ScriptID: scriptID,
				Key:      key,
				Value:    "<not found>",
			}

			return nullValue, nil
		}
		return nil, err
	}
	return &kv, nil
}

func SetScriptKVStore(db *gorm.DB, scriptID uint, key string, value string) error {
	var kv ScriptKVStore
	err := db.Where("script_id = ? AND key = ?", scriptID, key).First(&kv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			kv = ScriptKVStore{
				ScriptID: scriptID,
				Key:      key,
				Value:    value,
			}
		} else {
			return err
		}
	} else {
		kv.Value = value
	}
	return db.Save(&kv).Error
}

func DeleteScriptKVStore(db *gorm.DB, scriptID uint, key string) error {
	return db.Where("script_id = ? AND key = ?", scriptID, key).Delete(&ScriptKVStore{}).Error
}

type QueryMatch struct {
	ID       string            `json:"id"`
	Score    float32           `json:"score"` // Use "score" instead of "distance"
	Metadata map[string]string `json:"metadata"`
}
