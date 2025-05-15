package shared

import "context"

type basePathCtxKey struct{}
type appURLCtxKey struct{}
type pathCtxKey struct{}

func GetBasePath(ctx context.Context) string {
	if basePath, ok := ctx.Value(basePathCtxKey{}).(string); ok {
		return basePath
	}
	return "/"
}

func SetBasePath(ctx context.Context, basePath string) context.Context {
	return context.WithValue(ctx, basePathCtxKey{}, basePath)
}

func GetAppURL(ctx context.Context) string {
	if appURL, ok := ctx.Value(appURLCtxKey{}).(string); ok {
		return appURL
	}
	return ""
}

func SetAppURL(ctx context.Context, appURL string) context.Context {
	return context.WithValue(ctx, appURLCtxKey{}, appURL)
}

func GetURLPath(ctx context.Context) string {
	if path, ok := ctx.Value(pathCtxKey{}).(string); ok {
		return path
	}
	return "/"
}

func SetPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, pathCtxKey{}, path)
}
