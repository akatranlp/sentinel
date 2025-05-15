package utils

func Bang[T any](v T, ok bool) T { return v }

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
