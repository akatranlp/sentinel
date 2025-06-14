// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package enums

import (
	"errors"
	"fmt"
)

const (
	// ClaimTypeNormal is a ClaimType of type normal.
	ClaimTypeNormal ClaimType = "normal"
)

var ErrInvalidClaimType = errors.New("not a valid ClaimType")

// ClaimTypeValues returns a list of the values for ClaimType
func ClaimTypeValues() []ClaimType {
	return []ClaimType{
		ClaimTypeNormal,
	}
}

// String implements the Stringer interface.
func (x ClaimType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ClaimType) IsValid() bool {
	_, err := ParseClaimType(string(x))
	return err == nil
}

var _ClaimTypeValue = map[string]ClaimType{
	"normal": ClaimTypeNormal,
}

// ParseClaimType attempts to convert a string to a ClaimType.
func ParseClaimType(name string) (ClaimType, error) {
	if x, ok := _ClaimTypeValue[name]; ok {
		return x, nil
	}
	return ClaimType(""), fmt.Errorf("%s is %w", name, ErrInvalidClaimType)
}

// MarshalText implements the text marshaller method.
func (x ClaimType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ClaimType) UnmarshalText(text []byte) error {
	tmp, err := ParseClaimType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
