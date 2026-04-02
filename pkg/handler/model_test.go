package handler

import (
	"testing"

	"github.com/Meduzz/rpc"
	"github.com/gin-gonic/gin"
)

func TestHandlerRegistry(t *testing.T) {
	storage := NewRegistry()
	SetRegistry(storage)
	subject := NewBuilder(storage)

	t.Run("WebHandlers", func(t *testing.T) {
		subject.RegisterWeb("web-test", gin.WrapF(nil))
		result := storage.WebHandler("web-test")

		if result == nil {
			t.Error("result was nil")
		}

		t.Run("Wrong types", func(t *testing.T) {
			rpcResult := storage.RpcHandler("web-test")

			if rpcResult != nil {
				t.Error("rpcResult was not nil")
			}

			eventResult := storage.EventHandler("web-test")

			if eventResult != nil {
				t.Error("eventResult was not nil")
			}
		})
	})

	t.Run("RpcHandlers", func(t *testing.T) {
		subject.RegisterRPC("rpc-test", func(rc *rpc.RpcContext) {})
		result := storage.RpcHandler("rpc-test")

		if result == nil {
			t.Error("result was nil")
		}

		t.Run("Wrong types", func(t *testing.T) {
			webResult := storage.WebHandler("rpc-test")

			if webResult != nil {
				t.Error("webResult was not nil")
			}

			eventResult := storage.EventHandler("rpc-test")

			if eventResult != nil {
				t.Error("eventResult was not nil")
			}
		})
	})

	t.Run("RpcHandlers", func(t *testing.T) {
		subject.RegisterEvent("event-test", func(ec *rpc.EventContext) {})
		result := storage.EventHandler("event-test")

		if result == nil {
			t.Error("result was nil")
		}

		t.Run("Wrong types", func(t *testing.T) {
			webResult := storage.WebHandler("event-test")

			if webResult != nil {
				t.Error("webResult was not nil")
			}

			rpcResult := storage.RpcHandler("event-test")

			if rpcResult != nil {
				t.Error("rpcResult was not nil")
			}
		})
	})
}
