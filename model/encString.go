package model

import (
	"encoding/json"
	"errors"
)

// Password ...
type Password string

// String ...
func (e *Password) String() string {
	return "**encrypted**"
}

// MarshalJSON ...
func (e *Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON ...
func (e *Password) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		*e = Password(value)
		return nil
	default:
		return errors.New("Invalid Data")
	}
}

// IsEmpty ...
func (e *Password) IsEmpty() bool {
	return *e == Password("")
}
