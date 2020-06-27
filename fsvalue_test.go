package fsevent

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"
)

type data struct {
	Str    string         `firestore:"str"`
	Int    int64          `firestore:"int"`
	TimeP  *time.Time     `firestore:"timeP"`
	Time   time.Time      `firestore:"time"`
	BytesP *[]byte        `firestore:"bytesP"`
	Bytes  []byte         `firestore:"bytes"`
	Geo    *latlng.LatLng `firestore:"geo"`
	unex   string         `firestore:"unex"`
}

type dataToTest struct {
	blob string
	exp  data
}

var testTime = time.Date(2014, 10, 02, 15, 01, 23, 45123456, time.UTC)
var testBytes = []byte("this is bytes.")
var testLatLng = latlng.LatLng{Latitude: 35.8, Longitude: 135.7}
var dataToTests = []dataToTest{
	{
		`{"fields": {
			"str": {
				"stringValue": "this is string"
			},
			"int": {
				"integerValue": "3"
			},
			"timeP": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"time": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"bytesP": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			},
			"bytes": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			},
			"geo": {
				"geoPointValue": {
					"latitude": 35.8,
					"longitude": 135.7
				}
			}
		}}`,
		data{
			"this is string",
			3,
			&testTime,
			testTime,
			&testBytes,
			testBytes,
			&testLatLng,
			"",
		},
	},
}

func TestValue_DataTo(t *testing.T) {
	for _, test := range dataToTests {
		jsonBlob := []byte(test.blob)
		var v Value
		err := json.Unmarshal(jsonBlob, &v)
		if err != nil {
			t.Error(err)
			continue
		}
		var p data
		err = v.DataTo(&p)
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Printf("p: %+v", p)
		reflect.DeepEqual(p, test.exp)
	}
}
