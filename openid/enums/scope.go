package enums

import (
	"bytes"
	"github.com/akatranlp/go-pkg/its"
	"slices"
	"strings"
)

// ENUM(openid, profile, email, offline_access, api)
type Scope string

type Scopes []Scope

func (x Scopes) String() string {
	if len(x) == 0 {
		return ""
	}
	var buf strings.Builder
	buf.WriteString(string(x[0]))
	for i := range len(x) - 1 {
		buf.WriteByte(' ')
		buf.WriteString(string(x[i+1]))
	}

	return buf.String()
}

// MarshalText implements the text marshaller method.
func (x Scopes) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *Scopes) UnmarshalText(text []byte) error {
	var tmps Scopes
	for field := range bytes.FieldsSeq(text) {
		tmp, err := ParseScope(string(field))
		if err != nil {
			return err
		}
		tmps = append(tmps, tmp)
	}
	*x = tmps
	return nil
}

func (x Scopes) IsValid() bool {
	return its.All(slices.Values(x), func(scope Scope) bool { return scope.IsValid() })
}
