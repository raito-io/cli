// Code generated by "enumer -json -type=ActionType -transform=lower"; DO NOT EDIT.

package data_usage

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _ActionTypeName = "unknownactionreadwriteadmin"

var _ActionTypeIndex = [...]uint8{0, 13, 17, 22, 27}

const _ActionTypeLowerName = "unknownactionreadwriteadmin"

func (i ActionType) String() string {
	if i < 0 || i >= ActionType(len(_ActionTypeIndex)-1) {
		return fmt.Sprintf("ActionType(%d)", i)
	}
	return _ActionTypeName[_ActionTypeIndex[i]:_ActionTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ActionTypeNoOp() {
	var x [1]struct{}
	_ = x[UnknownAction-(0)]
	_ = x[Read-(1)]
	_ = x[Write-(2)]
	_ = x[Admin-(3)]
}

var _ActionTypeValues = []ActionType{UnknownAction, Read, Write, Admin}

var _ActionTypeNameToValueMap = map[string]ActionType{
	_ActionTypeName[0:13]:       UnknownAction,
	_ActionTypeLowerName[0:13]:  UnknownAction,
	_ActionTypeName[13:17]:      Read,
	_ActionTypeLowerName[13:17]: Read,
	_ActionTypeName[17:22]:      Write,
	_ActionTypeLowerName[17:22]: Write,
	_ActionTypeName[22:27]:      Admin,
	_ActionTypeLowerName[22:27]: Admin,
}

var _ActionTypeNames = []string{
	_ActionTypeName[0:13],
	_ActionTypeName[13:17],
	_ActionTypeName[17:22],
	_ActionTypeName[22:27],
}

// ActionTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ActionTypeString(s string) (ActionType, error) {
	if val, ok := _ActionTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ActionTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ActionType values", s)
}

// ActionTypeValues returns all values of the enum
func ActionTypeValues() []ActionType {
	return _ActionTypeValues
}

// ActionTypeStrings returns a slice of all String values of the enum
func ActionTypeStrings() []string {
	strs := make([]string, len(_ActionTypeNames))
	copy(strs, _ActionTypeNames)
	return strs
}

// IsAActionType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ActionType) IsAActionType() bool {
	for _, v := range _ActionTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for ActionType
func (i ActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for ActionType
func (i *ActionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("ActionType should be a string, got %s", data)
	}

	var err error
	*i, err = ActionTypeString(s)
	return err
}