package filters

import (
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	"github.com/Meduzz/dsl/service"
	serviceref "github.com/Meduzz/dsl/serviceRef"
	"github.com/Meduzz/helper/fp/slice"
)

func filterService(it *app.App, only ...serviceref.ServiceRef) []*service.Service {
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

	return matches
}

func FetchEndpointsByType(it *app.App, typ endpoint.RouteKind, only ...serviceref.ServiceRef) []*endpoint.Endpoint {
	matches := filterService(it, only...)

	eps := slice.FlatMap(matches, func(svc *service.Service) []*endpoint.Endpoint {
		return svc.Endpoints
	})

	return slice.Filter(eps, func(e *endpoint.Endpoint) bool {
		return e.Route.Kind == typ
	})
}
