package openid

import (
	"errors"
	"iter"
	"net/url"
)

var (
	ErrInvalidParameter = errors.New("invalid request parameters")
)

func mapFormValueKeysWithoutError(input url.Values, keys iter.Seq[string]) map[string]string {
	result := make(map[string]string)

	for k := range keys {
		v, ok := input[k]
		if !ok || len(v) < 1 {
			continue
		}
		result[k] = input[k][0]
	}

	return result
}

func mapFormValueKeys(input url.Values, keys iter.Seq[string]) (map[string]string, error) {
	result := make(map[string]string)

	for k := range keys {
		v, ok := input[k]
		if !ok || len(v) < 1 {
			continue
		}
		if len(v) > 1 {
			return nil, ErrInvalidParameter
		}
		result[k] = input[k][0]
	}

	return result, nil
}
