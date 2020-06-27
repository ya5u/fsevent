package fsevent

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// FirestoreValue holds Firestore fields.
type Value struct {
	CreateTime time.Time                         `json:"createTime"`
	Fields     map[string]map[string]interface{} `json:"fields"`
	Name       string                            `json:"name"`
	UpdateTime time.Time                         `json:"updateTime"`
}

func (v *Value) DataTo(p interface{}) error {
	rv := reflect.ValueOf(p)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("fsevent: nil or not a pointer")
	}
	crv := rv.Elem()
	switch crv.Kind() {
	case reflect.Map:
		// TODO: process for map
	case reflect.Struct:
		rt := crv.Type()
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.PkgPath != "" {
				// skip unexported field
				continue
			}
			tag := f.Tag.Get("firestore")
			fmt.Println(tag, f.Type, f.Type.Kind())
			switch f.Type {
			case reflect.TypeOf(time.Time{}):
				ts, ok := v.Fields[tag]["timestampValue"].(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not timestamp string", tag)
				}
				t, err := time.Parse(time.RFC3339Nano, ts)
				if err != nil {
					return fmt.Errorf("fsevent: failed to parse time on %s. %v", tag, err)
				}
				crv.Field(i).Set(reflect.ValueOf(t))
				continue
			case reflect.TypeOf(&time.Time{}):
				ts, ok := v.Fields[tag]["timestampValue"].(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not timestamp string", tag)
				}
				t, err := time.Parse(time.RFC3339Nano, ts)
				if err != nil {
					return fmt.Errorf("fsevent: failed to parse time on %s. %v", tag, err)
				}
				crv.Field(i).Set(reflect.ValueOf(&t))
				continue
			case reflect.TypeOf([]byte{}):
				bs, ok := v.Fields[tag]["bytesValue"].(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not bytes string", tag)
				}
				b, err := base64.StdEncoding.DecodeString(bs)
				if err != nil {
					return fmt.Errorf("fsevent: failed to decode bytes string on %s. %v", tag, err)
				}
				crv.Field(i).Set(reflect.ValueOf(b))
				continue
			case reflect.TypeOf(&[]byte{}):
				bs, ok := v.Fields[tag]["bytesValue"].(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not bytes string", tag)
				}
				b, err := base64.StdEncoding.DecodeString(bs)
				if err != nil {
					return fmt.Errorf("fsevent: failed to decode bytes string on %s. %v", tag, err)
				}
				crv.Field(i).Set(reflect.ValueOf(&b))
				continue
			case reflect.TypeOf(latlng.LatLng{}):
				return fmt.Errorf("fsevent: LatLng must be pointer")
			case reflect.TypeOf(&latlng.LatLng{}):
				gm, ok := v.Fields[tag]["geoPointValue"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("fsevent: %s is not geoPoint map", tag)
				}
				lat, ok := gm["latitude"].(float64)
				if !ok {
					return fmt.Errorf("fsevent: %s.latitude is not float64", tag)
				}
				lng, ok := gm["longitude"].(float64)
				if !ok {
					return fmt.Errorf("fsevent: %s.longitude is not float64", tag)
				}
				ll := latlng.LatLng{
					Latitude:  lat,
					Longitude: lng,
				}
				crv.Field(i).Set(reflect.ValueOf(&ll))
				continue
			}

			switch f.Type.Kind() {
			case reflect.Bool:
				fv, ok := v.Fields[tag]["booleanValue"].(bool)
				if !ok {
					return fmt.Errorf("fsevent: %s is not bool", tag)
				}
				crv.Field(i).SetBool(fv)
			case reflect.Int64:
				fv, ok := v.Fields[tag]["integerValue"].(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not int64", tag)
				}
				ifv, err := strconv.ParseInt(fv, 10, 64)
				if err != nil {
					return fmt.Errorf("fsevent: failed to parse int64 on %s. %v", tag, err)
				}
				crv.Field(i).SetInt(ifv)
			case reflect.Float64:
				fv, ok := v.Fields[tag]["doubleValue"].(float64)
				if !ok {
					return fmt.Errorf("fsevent: %s is not float64", tag)
				}
				crv.Field(i).SetFloat(fv)
			case reflect.String:
				fv, ok := v.Fields[tag]["stringValue"].(string)
				if !ok {
					fv, ok = v.Fields[tag]["referenceValue"].(string)
					if !ok {
						return fmt.Errorf("fsevent: %s is not string", tag)
					}
				}
				crv.Field(i).SetString(fv)
			default:
				return fmt.Errorf("fsevent: %s is type of %v but not supported", f.Name, f.Type)
			}
		}
	}
	return nil
}
