package util

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
)

//ClientAPIRequestContext represents context of client API request
type ClientAPIRequestContext struct {
	Authorization string
	Email         string
	Role          string
	UserID        int
}

// FromGinContext generates new Context from gin context
func FromGinContext(ginCtx *gin.Context) context.Context {
	email, _ := ginCtx.Get(`email`)
	userId, _ := ginCtx.Get(`id`)
	role, _ := ginCtx.Get(`role`)
	id, _ := strconv.Atoi(userId.(string))
	reqCtx := &ClientAPIRequestContext{
		Authorization: ginCtx.GetHeader("Authorization"),
		Email:         email.(string),
		UserID:        id,
		Role:          role.(string),
	}
	ctx := context.Background()
	return context.WithValue(ctx, 0, reqCtx)
}

// FromContext returns the api.ClientAPIRequestContext value stored in ctx, if any.
func FromContext(ctx context.Context) (*ClientAPIRequestContext, bool) {
	rc, ok := ctx.Value(0).(*ClientAPIRequestContext)
	return rc, ok
}
