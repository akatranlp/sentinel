package utils

import (
	"encoding"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func BinaryUnmarshallerHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(encoding.BinaryUnmarshaler)
		if !ok {
			return data, nil
		}
		str, ok := data.(string)
		if !ok {
			str = reflect.Indirect(reflect.ValueOf(&data)).Elem().String()
		}
		if err := unmarshaller.UnmarshalBinary([]byte(str)); err != nil {
			return nil, err
		}
		return result, nil
	}
}

func CreateDecoderConfig() mapstructure.DecoderConfig {
	return mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			BinaryUnmarshallerHookFunc(),
		),
		ErrorUnused: true,
		ZeroFields:  true,
	}
}
