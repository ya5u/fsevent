package fsevent

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   *Value `json:"oldValue"`
	Value      *Value `json:"value"`
	UpdateMask *struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// Type is method of FirestoreEvent which returns Event Type.
func (e *FirestoreEvent) Type() string {
	if len(e.UpdateMask.FieldPaths) > 0 {
		return "Update"
	}
	if len(e.Value.Name) > 0 {
		return "Create"
	}
	return "Delete"
}

// BooleanValue is Firestore's BooleanValue.
type BooleanValue struct {
	BooleanValue bool `json:"booleanValue"`
}

// IntegerValue is Firestore's IntegerValue.
type IntegerValue struct {
	IntegerValue string `json:"integerValue"`
}

// DoubleValue is Firestore's DoubleValue.
type DoubleValue struct {
	DoubleValue float64 `json:"doubleValue"`
}

// TimestampValue is Firestore's TimestampValue.
type TimestampValue struct {
	TimestampValue string `json:"timestampValue"`
}

// StringValue is Firestore's StringValue.
type StringValue struct {
	StringValue string `json:"stringValue"`
}

// BytesValue is Firestore's BytesValue.
type BytesValue struct {
	BytesValue string `json:"bytesValue"`
}

// ReferenceValue is Firestore's ReferenceValue.
type ReferenceValue struct {
	ReferenceValue string `json:"referenceValue"`
}

// GeoPointValue is Firestore's GeoPointValue.
type GeoPointValue struct {
	GeoPointValue struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geoPointValue"`
}
