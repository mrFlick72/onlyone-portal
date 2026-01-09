package server

import (
	"context"

	"github.com/gin-gonic/gin"
)

func CopyGinKeysToRequestContext(c *gin.Context) *context.Context {
	newCtx := c.Request.Context()
	for k, v := range c.Keys {
		newCtx = context.WithValue(newCtx, k, v)
	}
	return &newCtx
}
