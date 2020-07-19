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
	Bool   bool           `firestore:"bool"`
	Int    int64          `firestore:"int"`
	Float  float64        `firestore:"float"`
	TimeP  *time.Time     `firestore:"timeP,serverTimestamp"`
	Time   time.Time      `firestore:"time"`
	Str    string         `firestore:"str,omitempty"`
	BytesP *[]byte        `firestore:"bytesP"`
	Bytes  []byte         `firestore:"bytes"`
	Ref    string         `firestore:"ref"`
	Geo    *latlng.LatLng `firestore:"geo"`
	unex   string         `firestore:"unex"`
}

type dataToOkTest struct {
	blob string
	exp  data
}
type dataToErrTest struct {
	blob string
	exp  error
}

var testTime = time.Date(2014, 10, 02, 15, 01, 23, 45123456, time.UTC)
var testBytes = []byte("this is bytes.")
var testLatLng = latlng.LatLng{Latitude: 35.8, Longitude: 135.7}
var dataToOkTests = []dataToOkTest{
	{
		`{"fields": {
			"bool": {
				"booleanValue": true
			},
			"int": {
				"integerValue": "3"
			},
			"float": {
				"doubleValue": 3.1415
			},
			"timeP": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"time": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"str": {
				"stringValue": "this is string"
			},
			"bytesP": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			},
			"bytes": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			},
			"ref": {
				"referenceValue": "projects/{project_id}/databases/{database_id}/documents/{document_path}"
			},
			"geo": {
				"geoPointValue": {
					"latitude": 35.8,
					"longitude": 135.7
				}
			}
		}}`,
		data{
			true,
			3,
			3.1415,
			&testTime,
			testTime,
			"this is string",
			&testBytes,
			testBytes,
			"projects/{project_id}/databases/{database_id}/documents/{document_path}",
			&testLatLng,
			"",
		},
	},
	{
		`{"fields": {
			"str": {
				"stringValue": "this is string"
			},
			"int": {
				"integerValue": "3"
			},
			"time": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"bytes": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			}
		}}`,
		data{
			false,
			3,
			0.0,
			nil,
			testTime,
			"this is string",
			nil,
			testBytes,
			"",
			nil,
			"",
		},
	},
	{
		`{"fields": {}}`,
		data{},
	},
	{
		`{"fields": {
			"notDefined": {
				"stringValue": "this field doesn't defined"
			}
		}}`,
		data{},
	},
	{
		`{"fields": {
			"str": {
				"integerValue": "3"
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
			false,
			3,
			0.0,
			&testTime,
			testTime,
			"",
			&testBytes,
			testBytes,
			"",
			&testLatLng,
			"",
		},
	},
}
var dataToErrTests = []dataToErrTest{
	{
		`{"fields": {
			"str": {
				"stringValue": "this is string"
			},
			"int": {
				"integerValue": 3
			},
			"time": {
				"timestampValue": "2014-10-02T15:01:23.045123456Z"
			},
			"bytes": {
				"bytesValue": "dGhpcyBpcyBieXRlcy4="
			}
		}}`,
		fmt.Errorf("fsevent: int is not int string"),
	},
}

func TestValue_DataTo(t *testing.T) {
	// Normal cases
	for _, test := range dataToOkTests {
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
		if !reflect.DeepEqual(p, test.exp) {
			t.Errorf("\nexpected:%+v \nactual:%+v", test.exp, p)
		}
	}
	// Error cases
	for _, test := range dataToErrTests {
		jsonBlob := []byte(test.blob)
		var v Value
		err := json.Unmarshal(jsonBlob, &v)
		if err != nil {
			t.Error(err)
			continue
		}
		var p data
		err = v.DataTo(&p)
		if err == nil {
			t.Errorf("%s DataTo %+v must be failed", test.blob, p)
			continue
		}
		if !reflect.DeepEqual(err, test.exp) {
			t.Errorf("\nexpected:%+v \nactual:%+v", test.exp, err)
		}
	}
}
