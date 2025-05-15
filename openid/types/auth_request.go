package types

import (
	"errors"
	"net/url"

	"github.com/akatranlp/sentinel/openid/enums"
)

type AuthRequest struct {
	ClientID            string                    `mapstructure:"client_id"`
	RedirectURI         *url.URL                  `mapstructure:"redirect_uri"`
	Scope               enums.Scopes              `mapstructure:"scope"`
	ResponseType        enums.ResponseType        `mapstructure:"response_type"`
	ResponseMode        enums.ResponseMode        `mapstructure:"response_mode"`
	CodeChallengeMethod enums.CodeChallengeMethod `mapstructure:"code_challenge_method"`
	CodeChallenge       string                    `mapstructure:"code_challenge"`
	State               string                    `mapstructure:"state"`
	Nonce               string                    `mapstructure:"nonce"`
	// optional, check if i want to implement them
	Display enums.Display `mapstructure:"display"`
	Prompt  enums.Prompt  `mapstructure:"prompt"`
	MaxAge  *int          `mapstructure:"max_age"`
}

func (a *AuthRequest) IsValid() error {
	if !a.Scope.IsValid() {
		return enums.ErrInvalidScope
	}

	if !a.ResponseType.IsValid() {
		return enums.ErrInvalidResponseType
	}

	a.ResponseMode.WithDefault()
	if !a.ResponseMode.IsValid() {
		return enums.ErrInvalidResponseMode
	}

	a.Display.WithDefault()
	if !a.Display.IsValid() {
		return enums.ErrInvalidDisplay
	}

	if a.Prompt != "" && !a.Prompt.IsValid() {
		return enums.ErrInvalidPrompt
	}

	if a.CodeChallengeMethod != "" && !a.CodeChallengeMethod.IsValid() {
		return enums.ErrInvalidCodeChallengeMethod
	}

	if a.CodeChallengeMethod != "" && a.CodeChallenge == "" {
		return errors.New("code challenge without value")
	}

	if a.CodeChallenge != "" && a.CodeChallenge == "" {
		return errors.New("code challenge without method")
	}

	return nil
}
