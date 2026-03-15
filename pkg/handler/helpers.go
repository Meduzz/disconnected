package handler

func NewBuilder(registry *Registry) NamedHandlerBuilder {
	return &namedHandlerBuilder{
		registry: registry,
	}
}
