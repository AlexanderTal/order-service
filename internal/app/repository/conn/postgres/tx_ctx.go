package rcpostgres

import (
	"context"

	"gorm.io/gorm"
)

type txCtxKey struct{}

func ctxWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txCtxKey{}, tx)
}

func getTxFromCtx(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(txCtxKey{}).(*gorm.DB)
	return tx
}
