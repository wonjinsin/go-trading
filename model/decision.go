package model

import (
	"encoding/json"
	"strings"
)

// DecisionState ...
type DecisionState string

// DecisionStateNone ...
const (
	DecisionStateNone DecisionState = "None"
	DecisionStateBuy  DecisionState = "Buy"
	DecisionStateSell DecisionState = "Sell"
	DecisionStateHold DecisionState = "Hold"
)

var decisionStateMap = map[DecisionState]bool{
	DecisionStateNone: true,
	DecisionStateBuy:  true,
	DecisionStateSell: true,
	DecisionStateHold: true,
}

// String ...
func (ds DecisionState) String() string {
	return string(ds)
}

// MarshalJSON ...
func (ds *DecisionState) MarshalJSON() (data []byte, err error) {
	return json.Marshal(ds.String())
}

// UnmarshalJSON ...
func (ds *DecisionState) UnmarshalJSON(data []byte) (err error) {
	*ds = DecisionStateNone

	strData := strings.Trim(string(data), "\"")
	if decisionStateMap[DecisionState(strData)] {
		*ds = DecisionState(strData)
	}
	return nil
}

// Decision ...
type Decision struct {
	Decision DecisionState `json:"decision"`
	Percent  uint          `json:"percent"`
	Reason   string        `json:"reason"`
}
