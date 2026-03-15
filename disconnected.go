package disconnected

import (
	"errors"
	"fmt"

	"github.com/Meduzz/disconnected/pkg/event"
	"github.com/Meduzz/disconnected/pkg/handler"
	"github.com/Meduzz/disconnected/pkg/web"
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	"github.com/Meduzz/dsl/service"
	serviceref "github.com/Meduzz/dsl/serviceRef"
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/rpc/encoding"
	"github.com/gin-gonic/gin"
)

func HttpServer(fn func(*web.Server) error) error {
	return web.HttpServer(fn)
}

func RpcServer(codec encoding.Codec, fn func(*event.Server) error) error {
	return event.EventServer(codec, fn, true)
}

func App(it *app.App, only ...serviceref.ServiceRef) {
	matches := it.Services

	if len(only) > 0 {
		// validate our whitelist
		whiteList := slice.Map(slice.Filter(only, func(ref serviceref.ServiceRef) bool {
			a, ok := ref.App()
			return ref.Valid() && ok && a == it.Name
		}), func(ref serviceref.ServiceRef) string {
			s, _ := ref.Service()
			return s
		})

		// fetch services that matches whitelist
		matches = slice.Filter(matches, func(s *service.Service) bool {
			return slice.Contains(whiteList, s.Name)
		})
	}

	routes := slice.FlatMap(matches, func(s *service.Service) []*endpoint.Endpoint {
		return s.Endpoints
	})

	groupedRoutes := slice.Group(routes, func(r *endpoint.Endpoint) string {
		return string(r.Route.Kind)
	})

	events, hasEvents := groupedRoutes[string(endpoint.RpcKind)]
	webs, hasWeb := groupedRoutes[string(endpoint.HttpKind)]

	if hasEvents {
		event.EventServer(encoding.Json(), func(s *event.Server) error {
			return slice.Fold(events, nil, func(e *endpoint.Endpoint, agg error) error {
				// TODO even worth to continue if agg != nil?
				var err error
				if h := handler.HandlerRegistry.RpcHandler(e.Name); h != nil {
					err = s.HandleRPC(e.Route.Topic, e.Route.ConsumerGroup, h)
				} else if h := handler.HandlerRegistry.EvebtHandler(e.Name); h != nil {
					err = s.HandleEvent(e.Route.Topic, e.Route.ConsumerGroup, h)
				}

				if err != nil {
					return errors.Join(agg, err)
				}

				return agg
			})
		}, !hasWeb)
	}

	if hasWeb {
		// dig out the context path for later
		ctxPath := it.ContextPath

		web.HttpServer(func(s *web.Server) error {
			return s.WithRouter(func(r *gin.Engine) error {
				slice.ForEach(webs, func(e *endpoint.Endpoint) {
					if h := handler.HandlerRegistry.WebHandler(e.Name); h != nil {
						path := e.Route.Path
						if ctxPath != "" {
							path = fmt.Sprintf("%s%s", ctxPath, path)
						}

						r.Handle(e.Route.Method, path, h)
					}
				})

				return nil
			})
		})
	}
}
