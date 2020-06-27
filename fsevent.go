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
