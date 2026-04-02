package handler

import (
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/service"
	"github.com/Meduzz/rpc"
	"github.com/gin-gonic/gin"
)

type (
	Registry struct {
		handlers []*handlerWrapper
	}

	handlerWrapper struct {
		name         string
		webHandler   gin.HandlerFunc
		rpcHandler   rpc.RpcHandler
		eventHandler rpc.EventHandler
	}

	handlerDelegate struct{}

	NamedHandlerBuilder interface {
		RegisterWeb(name string, handler gin.HandlerFunc)
		RegisterRPC(name string, handler rpc.RpcHandler)
		RegisterEvent(name string, handler rpc.EventHandler)
	}

	namedHandlerBuilder struct {
		registry *Registry
	}

	NamedHandlerProvider interface {
		Builder(NamedHandlerBuilder)
	}
)

var (
	_ service.Delegate    = &handlerDelegate{}
	_ NamedHandlerBuilder = &namedHandlerBuilder{}
)

func (h *handlerDelegate) Visit(svc service.Service) error {
	it, ok := svc.(NamedHandlerProvider)

	if ok {
		it.Builder(NewBuilder(HandlerRegistry))
	}

	return nil
}

func (h *namedHandlerBuilder) RegisterWeb(name string, handler gin.HandlerFunc) {
	h.register(name, func(h *handlerWrapper) {
		h.webHandler = handler
	})
}

func (h *namedHandlerBuilder) RegisterRPC(name string, handler rpc.RpcHandler) {
	h.register(name, func(h *handlerWrapper) {
		h.rpcHandler = handler
	})
}

func (h *namedHandlerBuilder) RegisterEvent(name string, handler rpc.EventHandler) {
	h.register(name, func(h *handlerWrapper) {
		h.eventHandler = handler
	})
}

func (h *namedHandlerBuilder) register(name string, cb func(h *handlerWrapper)) {
	wrapper := &handlerWrapper{
		name: name,
	}

	cb(wrapper)

	h.registry.handlers = append(h.registry.handlers, wrapper)
}

func (r *Registry) WebHandler(name string) gin.HandlerFunc {
	return slice.Head(slice.Filter(r.handlers, func(w *handlerWrapper) bool {
		return w.name == name
	})).webHandler
}

func (r *Registry) RpcHandler(name string) rpc.RpcHandler {
	return slice.Head(slice.Filter(r.handlers, func(w *handlerWrapper) bool {
		return w.name == name
	})).rpcHandler
}

func (r *Registry) EventHandler(name string) rpc.EventHandler {
	return slice.Head(slice.Filter(r.handlers, func(w *handlerWrapper) bool {
		return w.name == name
	})).eventHandler
}
