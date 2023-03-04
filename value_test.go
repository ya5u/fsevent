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
	Bool       bool                   `firestore:"bool"`
	Int        int64                  `firestore:"int"`
	IntP       *int64                 `firestore:"intP"`
	Float      float64                `firestore:"float"`
	FloatP     *float64               `firestore:"floatP"`
	TimeP      *time.Time             `firestore:"timeP,serverTimestamp"`
	Time       time.Time              `firestore:"time"`
	Str        string                 `firestore:"str,omitempty"`
	BytesP     *[]byte                `firestore:"bytesP"`
	Bytes      []byte                 `firestore:"bytes"`
	Ref        string                 `firestore:"ref"`
	Geo        *latlng.LatLng         `firestore:"geo"`
	ArrStr     []string               `firestore:"arrstr"`
	StructMap  smap                   `firestore:"structmap"`
	StructMapP *smap                  `firestore:"structmapp"`
	ShallowMap map[string]interface{} `firestore:"shallowmap"`
	// ShallowMapP *map[string]interface{} `firestore:"shallowmapp"`
	// DeepMap        map[string]map[string]interface{}  `firestore:"deepmap"`
	// DeepMapP       *map[string]map[string]interface{} `firestore:"deepmapp"`
	// DeepStructMap  map[string]smap                    `firestore:"deepstructmap"`
	// DeepStructMapP *map[string]smap                   `firestore:"deepstructmapp"`
	unex string `firestore:"unex"`
}

type smap struct {
	KeyBool  bool    `firestore:"keybool"`
	KeyInt   int64   `firestore:"keyint"`
	KeyFloat float64 `firestore:"keyfloat"`
	KeyStr   string  `firestore:"keystr"`
}

type dataToTest struct {
	name      string
	inputBlob string
	expData   data
	expErr    error
}

var testTime = time.Date(2014, 10, 02, 15, 01, 23, 45123456, time.UTC)
var testBytes = []byte("this is bytes.")
var testLatLng = latlng.LatLng{Latitude: 35.8, Longitude: 135.7}

func intPtr(i int64) *int64 {
	return &i
}
func floatPtr(f float64) *float64 {
	return &f
}

var dataToTests = []dataToTest{
	{
		"full_fields",
		`{"fields": {
			"bool": {
				"booleanValue": true
			},
			"int": {
				"integerValue": "3"
			},
			"intP": {
				"integerValue": "5"
			},
			"float": {
				"doubleValue": 3.14
			},
			"floatP": {
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
			},
			"arrstr": {
				"arrayValue": {
					"values": [
						{
							"stringValue": "string in array"
						},
						{
							"stringValue": "another string in array"
						}
					]
				}
			},
			"structmap": {
				"mapValue": {
					"fields": {
						"keystr": {"stringValue": "string value"},
						"keyint": {"integerValue": "123"}
					}
				}
			},
			"structmapp": {
				"mapValue": {
					"fields": {
						"keybool": {"booleanValue": true},
						"keyfloat": {"doubleValue": 1.23}
					}
				}
			},
			"shallowmap": {
				"mapValue": {
					"fields": {
						"keyint": {"integerValue": "123"},
						"keyfloat": {"doubleValue": 1.23}
					}
				}
			},
			"shallowmapp": {
				"mapValue": {
					"fields": {
						"keyint": {"integerValue": "123"},
						"keyfloat": {"doubleValue": 1.23}
					}
				}
			}
		}}`,
		data{
			true,
			3,
			intPtr(5),
			3.14,
			floatPtr(3.1415),
			&testTime,
			testTime,
			"this is string",
			&testBytes,
			testBytes,
			"projects/{project_id}/databases/{database_id}/documents/{document_path}",
			&testLatLng,
			[]string{"string in array", "another string in array"},
			smap{
				false,
				123,
				0.0,
				"string value",
			},
			&smap{
				true,
				0,
				1.23,
				"",
			},
			map[string]interface{}{
				"keyint":   123,
				"keyfloat": 1.23,
			},
			// &map[string]interface{}{
			// 	"keyint":   123,
			// 	"keyfloat": 1.23,
			// },
			// nil,
			// nil,
			// nil,
			// nil,
			"",
		},
		nil,
	},
	{
		"optional_fields",
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
			},
			"arrstr": {
				"arrayValue": {}
			},
			"structmap": {
				"mapValue": {}
			}
		}}`,
		data{
			false,
			3,
			nil,
			0.0,
			nil,
			nil,
			testTime,
			"this is string",
			nil,
			testBytes,
			"",
			nil,
			nil,
			smap{},
			nil,
			nil,
			// nil,
			// nil,
			// nil,
			// nil,
			// nil,
			"",
		},
		nil,
	},
	{
		"no_fields",
		`{"fields": {}}`,
		data{},
		nil,
	},
	{
		"only_not_defined_fields",
		`{"fields": {
			"notDefined": {
				"stringValue": "this field doesn't defined"
			}
		}}`,
		data{},
		nil,
	},
	{
		"optional_fields2",
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
			nil,
			0.0,
			nil,
			&testTime,
			testTime,
			"",
			&testBytes,
			testBytes,
			"",
			&testLatLng,
			nil,
			smap{},
			nil,
			nil,
			// nil,
			// nil,
			// nil,
			// nil,
			// nil,
			"",
		},
		nil,
	},
	{
		"int_error",
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
		data{},
		fmt.Errorf("fsevent: int is not int string"),
	},
}

func TestValue_DataTo(t *testing.T) {
	for _, test := range dataToTests {
		t.Run(test.name, func(t *testing.T) {
			jsonBlob := []byte(test.inputBlob)
			var v Value
			err := json.Unmarshal(jsonBlob, &v)
			if err != nil {
				t.Error(err)
				return
			}
			var p data
			err = v.DataTo(&p)
			if !reflect.DeepEqual(test.expErr, err) {
				t.Errorf("#%s\nwant:%+v\ngot:%+v", test.name, test.expErr, err)
			}
			if !reflect.DeepEqual(test.expData, p) {
				t.Errorf("#%s\nwant:%+v\ngot:%+v", test.name, test.expData, p)
			}
		})
	}
}
