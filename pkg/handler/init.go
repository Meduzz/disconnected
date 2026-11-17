package handler

var (
	HandlerRegistry *Registry
)

func init() {
	HandlerRegistry = NewRegistry()
}

func SetRegistry(it *Registry) {
	HandlerRegistry = it
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make([]*handlerWrapper, 0),
	}
}
