package serializer

import (
	"encoding/json"
	"database/sql"
)

// Json null float 64
type JsonNullFloat64 struct {
	sql.NullFloat64
}

// Marshal json function
func (v JsonNullFloat64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Float64)
	} else {
		return json.Marshal(nil)
	}
}

// Unmarshal json function
func (v *JsonNullFloat64) UnmarshalJSON(data []byte) error {
	var x *float64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Float64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

