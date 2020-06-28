package fsevent

const (
	TypeCreate = iota
	TypeUpdate
	TypeDelete
)

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   *Value `json:"oldValue"`
	Value      *Value `json:"value"`
	UpdateMask *struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// Type is method of FirestoreEvent which returns Event Type.
func (e *FirestoreEvent) Type() int {
	if len(e.UpdateMask.FieldPaths) > 0 {
		return TypeUpdate
	}
	if len(e.Value.Name) > 0 {
		return TypeCreate
	}
	return TypeDelete
}
