package jose

import "context"

type contextKey struct{}

func SetJose(ctx context.Context, j *jose) context.Context {
	return context.WithValue(ctx, contextKey{}, j)
}

func GetJose(ctx context.Context) *jose {
	return ctx.Value(contextKey{}).(*jose)
}
