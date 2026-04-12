package handler

import "github.com/Meduzz/helper/service"

var (
	HandlerRegistry *Registry
)

func init() {
	HandlerRegistry = NewRegistry()
	service.AddDelegate(&handlerDelegate{})
}

func SetRegistry(it *Registry) {
	HandlerRegistry = it
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make([]*handlerWrapper, 0),
	}
}
