package apiversion

import "context"

type key int

const contextKey key = 0

func NewContext(ctx context.Context, ver *Version) context.Context {
	return context.WithValue(ctx, contextKey, ver)
}

func FromContext(ctx context.Context) *Version {
	return ctx.Value(contextKey).(*Version)
}
