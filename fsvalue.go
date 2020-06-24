package fsevent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// FirestoreValue holds Firestore fields.
type Value struct {
	CreateTime time.Time `json:"createTime"`
	Fields     []byte    `json:"fields"`
	Name       string    `json:"name"`
	UpdateTime time.Time `json:"updateTime"`
}

func (v *Value) DataTo(p interface{}) error {
	rv := reflect.ValueOf(p)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("fsevent: nil or not a pointer")
	}
	err := json.Unmarshal(v.Fields, p)
	if err != nil {
		return fmt.Errorf("fsevent: could not unmarshal Value.Fields %v", err)
	}
	return nil
}
