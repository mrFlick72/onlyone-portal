package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ContextFactoryConverter interface {
	CreateContextFromGin(c *gin.Context) context.Context
}

type GinContextToPlainContextFactory struct{}

func (f *GinContextToPlainContextFactory) CreateContextFromGin(c *gin.Context) context.Context {
	newCtx := c.Request.Context()
	for k, v := range c.Keys {
		newCtx = context.WithValue(newCtx, k, v)
	}
	return newCtx
}

func CopyGinKeysToRequestContext(c *gin.Context) context.Context {
	newCtx := c.Request.Context()
	for k, v := range c.Keys {
		newCtx = context.WithValue(newCtx, k, v)
	}
	return newCtx
}
