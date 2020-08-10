package fsevent

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

// Value holds Firestore fields.
type Value struct {
	CreateTime time.Time                         `json:"createTime"`
	Fields     map[string]map[string]interface{} `json:"fields"`
	Name       string                            `json:"name"`
	UpdateTime time.Time                         `json:"updateTime"`
}

// DataTo uses the document's fields to populate p, which can be a pointer to a
// map[string]interface{} or a pointer to a struct.
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
			parts := strings.Split(tag, ",")
			tag = parts[0]
			if v.Fields[tag] == nil {
				// skip fields that have no value
				continue
			}
			switch f.Type {
			case reflect.TypeOf(time.Time{}):
				t, err := reflectTime(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				if t == nil {
					continue
				}
				crv.Field(i).Set(reflect.ValueOf(*t))
				continue
			case reflect.TypeOf(&time.Time{}):
				t, err := reflectTime(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				crv.Field(i).Set(reflect.ValueOf(t))
				continue
			case reflect.TypeOf([]byte{}):
				b, err := reflectBytes(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				if b == nil {
					continue
				}
				crv.Field(i).Set(reflect.ValueOf(*b))
				continue
			case reflect.TypeOf(&[]byte{}):
				b, err := reflectBytes(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				crv.Field(i).Set(reflect.ValueOf(b))
				continue
			case reflect.TypeOf(latlng.LatLng{}):
				return fmt.Errorf("fsevent: LatLng must be pointer")
			case reflect.TypeOf(&latlng.LatLng{}):
				geo := v.Fields[tag]["geoPointValue"]
				if geo == nil {
					continue
				}
				gm, ok := geo.(map[string]interface{})
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
			case reflect.TypeOf(new(int64)):
				i64, err := reflectInt64(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				crv.Field(i).Set(reflect.ValueOf(i64))
				continue
			case reflect.TypeOf(new(float64)):
				f64, err := reflectFloat64(tag, v.Fields[tag])
				if err != nil {
					return err
				}
				crv.Field(i).Set(reflect.ValueOf(f64))
				continue
			}

			switch f.Type.Kind() {
			case reflect.Bool:
				bv := v.Fields[tag]["booleanValue"]
				if bv == nil {
					continue
				}
				fv, ok := bv.(bool)
				if !ok {
					return fmt.Errorf("fsevent: %s is not bool", tag)
				}
				crv.Field(i).SetBool(fv)
			case reflect.Int64:
				iv := v.Fields[tag]["integerValue"]
				if iv == nil {
					continue
				}
				fv, ok := iv.(string)
				if !ok {
					return fmt.Errorf("fsevent: %s is not int string", tag)
				}
				ifv, err := strconv.ParseInt(fv, 10, 64)
				if err != nil {
					return fmt.Errorf("fsevent: failed to parse int64 on %s. %v", tag, err)
				}
				crv.Field(i).SetInt(ifv)
			case reflect.Float64:
				dv := v.Fields[tag]["doubleValue"]
				iv := v.Fields[tag]["integerValue"]
				var fv float64
				var ok bool
				if dv != nil {
					fv, ok = dv.(float64)
				} else if iv != nil {
					fv, ok = iv.(float64)
				} else {
					continue
				}
				if !ok {
					return fmt.Errorf("fsevent: %s is not float64", tag)
				}
				crv.Field(i).SetFloat(fv)
			case reflect.String:
				sv := v.Fields[tag]["stringValue"]
				rv := v.Fields[tag]["referenceValue"]
				var fv string
				var ok bool
				if sv != nil {
					fv, ok = sv.(string)
				} else if rv != nil {
					fv, ok = rv.(string)
				} else {
					continue
				}
				if !ok {
					return fmt.Errorf("fsevent: %s is not string", tag)
				}
				crv.Field(i).SetString(fv)
			default:
				return fmt.Errorf("fsevent: %s is type of %v but not supported", f.Name, f.Type)
			}
		}
	}
	return nil
}

func reflectTime(tag string, field map[string]interface{}) (*time.Time, error) {
	tv := field["timestampValue"]
	if tv == nil {
		return nil, nil
	}
	ts, ok := tv.(string)
	if !ok {
		return nil, fmt.Errorf("fsevent: %s is not timestamp string", tag)
	}
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return nil, fmt.Errorf("fsevent: failed to parse time on %s. %v", tag, err)
	}
	return &t, nil
}

func reflectBytes(tag string, field map[string]interface{}) (*[]byte, error) {
	bv := field["bytesValue"]
	if bv == nil {
		return nil, nil
	}
	bs, ok := bv.(string)
	if !ok {
		return nil, fmt.Errorf("fsevent: %s is not bytes string", tag)
	}
	b, err := base64.StdEncoding.DecodeString(bs)
	if err != nil {
		return nil, fmt.Errorf("fsevent: failed to decode bytes string on %s. %v", tag, err)
	}
	return &b, nil
}

func reflectInt64(tag string, field map[string]interface{}) (*int64, error) {
	iv := field["integerValue"]
	if iv == nil {
		return nil, nil
	}
	is, ok := iv.(string)
	if !ok {
		return nil, fmt.Errorf("fsevent: %s is not int64 string", tag)
	}
	i, err := strconv.ParseInt(is, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fsevent: failed to parse int64 on %s. %v", tag, err)
	}
	return &i, nil
}

func reflectFloat64(tag string, field map[string]interface{}) (*float64, error) {
	fv := field["doubleValue"]
	if fv == nil {
		return nil, nil
	}
	fs, ok := fv.(string)
	if !ok {
		return nil, fmt.Errorf("fsevent: %s is not float64 string", tag)
	}
	f, err := strconv.ParseFloat(fs, 64)
	if err != nil {
		return nil, fmt.Errorf("fsevent: failed to parse float64 on %s. %v", tag, err)
	}
	return &f, nil
}
