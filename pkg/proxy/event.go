package proxy

import (
	"github.com/Meduzz/rpc"
	"github.com/gin-gonic/gin"
)

// EventBridge - creates a gin handler that accepts a request and send an eventon the provided topic.
func EventBridge[T any](topic string, srv *rpc.RPC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := new(T)
		err := ctx.BindJSON(data)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		err = srv.Trigger(topic, data)

		if err != nil {
			ctx.Status(500)
			return
		}

		ctx.Status(201)
	}
}

// RpcBridge - creates a gin handler that accepts T, makes a request on topic and expects K back which is returned.
func RpcBridge[T any, K any](topic string, timeout int, srv *rpc.RPC) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data := new(T)
		err := ctx.BindJSON(data)

		if err != nil {
			ctx.AbortWithStatus(400)
			return
		}

		resp := new(K)
		err = srv.Request(topic, data, resp, timeout)

		if err != nil {
			ctx.Status(500)
			return
		}

		ctx.JSON(200, resp)
	}
}
