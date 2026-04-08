package disconnected

import (
	"github.com/Meduzz/disconnected/pkg/event"
	"github.com/Meduzz/disconnected/pkg/filters"
	"github.com/Meduzz/disconnected/pkg/web"
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	serviceref "github.com/Meduzz/dsl/serviceRef"
	"github.com/Meduzz/rpc/encoding"
)

func HttpServer(fn func(*web.Server) error) error {
	return web.HttpServer(fn)
}

func RpcServer(codec encoding.Codec, fn func(*event.Server) error) error {
	return event.EventServer(codec, fn, true)
}

func App(it *app.App, only ...serviceref.ServiceRef) error {
	hasEvents := len(filters.FetchEndpointsByType(it, endpoint.RpcKind, only...)) > 0
	hasWeb := len(filters.FetchEndpointsByType(it, endpoint.HttpKind, only...)) > 0

	if hasEvents {
		err := event.EventServer(encoding.Json(), func(s *event.Server) error {
			return s.ForApp(it, only...)
		}, !hasWeb)

		if err != nil {
			return err
		}
	}

	if hasWeb {
		// dig out the context path for later
		return web.HttpServer(func(s *web.Server) error {
			return s.WithApp(it, only...)
		})
	}

	return nil
}
