package helper

import (
	"context"
	"github.com/google/uuid"
)

type traceId struct{}
type moduleName struct{}

func WithContextTrace(parent context.Context) context.Context {
	return context.WithValue(parent, traceId{}, uuid.NewString())
}

func WithContextModule(parent context.Context, name string) context.Context {
	return context.WithValue(parent, moduleName{}, name)
}

func ReadContextTrace(ctx context.Context) string {
	val := ctx.Value(traceId{})
	trace, _ := val.(string)
	return trace
}

func ReadContextModule(ctx context.Context) string {
	val := ctx.Value(moduleName{})
	name, _ := val.(string)
	return name
}
