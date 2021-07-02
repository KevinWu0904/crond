package logs

import (
	"context"
)

type key int

const (
	kvCtxKey key = iota
)

type ctxKVs struct {
	kvs []interface{}
	pre *ctxKVs
}

func getKVs(ctx context.Context) *ctxKVs {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(kvCtxKey)
	if i == nil {
		return nil
	}
	if kvs, ok := i.(*ctxKVs); ok {
		return kvs
	}
	return nil
}

// recursiveAllKVs keeps FIFO order.
func recursiveAllKVs(result *[]interface{}, kvs *ctxKVs, total int) {
	if kvs == nil {
		*result = make([]interface{}, 0, total)
		return
	}
	recursiveAllKVs(result, kvs.pre, total+len(kvs.kvs))
	*result = append(*result, kvs.kvs...)
}

// CtxAddKVs puts new kvs into predefined kvCtxKey.
func CtxAddKVs(ctx context.Context, kvs ...interface{}) context.Context {
	if len(kvs) == 0 || (len(kvs)&1 == 1) { // ignore odd kvs
		return ctx
	}

	kvList := make([]interface{}, 0, len(kvs))
	kvList = append(kvList, kvs...)

	return context.WithValue(ctx, kvCtxKey, &ctxKVs{
		kvs: kvList,
		pre: getKVs(ctx),
	})
}

// GetAllKVs retrieves all kvs from predefined kvCtxKey.
func GetAllKVs(ctx context.Context) []interface{} {
	if ctx == nil {
		return nil
	}
	kvs := getKVs(ctx)
	if kvs == nil {
		return nil
	}

	var result []interface{}
	recursiveAllKVs(&result, kvs, 0)
	return result
}
