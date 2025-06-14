// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package provider

import (
	"errors"
	"fmt"
)

const (
	// ProviderTypeGitlab is a ProviderType of type gitlab.
	ProviderTypeGitlab ProviderType = "gitlab"
	// ProviderTypeGithub is a ProviderType of type github.
	ProviderTypeGithub ProviderType = "github"
	// ProviderTypeGitea is a ProviderType of type gitea.
	ProviderTypeGitea ProviderType = "gitea"
	// ProviderTypeKeycloak is a ProviderType of type keycloak.
	ProviderTypeKeycloak ProviderType = "keycloak"
	// ProviderTypeUnknown is a ProviderType of type unknown.
	ProviderTypeUnknown ProviderType = "unknown"
)

var ErrInvalidProviderType = errors.New("not a valid ProviderType")

// String implements the Stringer interface.
func (x ProviderType) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ProviderType) IsValid() bool {
	_, err := ParseProviderType(string(x))
	return err == nil
}

var _ProviderTypeValue = map[string]ProviderType{
	"gitlab":   ProviderTypeGitlab,
	"github":   ProviderTypeGithub,
	"gitea":    ProviderTypeGitea,
	"keycloak": ProviderTypeKeycloak,
	"unknown":  ProviderTypeUnknown,
}

// ParseProviderType attempts to convert a string to a ProviderType.
func ParseProviderType(name string) (ProviderType, error) {
	if x, ok := _ProviderTypeValue[name]; ok {
		return x, nil
	}
	return ProviderType(""), fmt.Errorf("%s is %w", name, ErrInvalidProviderType)
}

// MarshalText implements the text marshaller method.
func (x ProviderType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *ProviderType) UnmarshalText(text []byte) error {
	tmp, err := ParseProviderType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
