package event

import (
	"errors"
	"fmt"

	"github.com/Meduzz/disconnected/pkg/filters"
	"github.com/Meduzz/disconnected/pkg/handler"
	"github.com/Meduzz/dsl/app"
	"github.com/Meduzz/dsl/endpoint"
	serviceref "github.com/Meduzz/dsl/serviceRef"
	"github.com/Meduzz/helper/block"
	"github.com/Meduzz/helper/fp/slice"
	"github.com/Meduzz/helper/nuts"
	quickapirpc "github.com/Meduzz/quickapi-rpc"
	"github.com/Meduzz/quickapi/model"
	"github.com/Meduzz/rpc"
	"github.com/Meduzz/rpc/encoding"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type (
	Server struct {
		conn  *nats.Conn
		codec encoding.Codec
	}
)

func EventServer(codec encoding.Codec, fn func(*Server) error, blocking bool) error {
	conn, err := nuts.Connect()

	if err != nil {
		return err
	}

	srv := &Server{conn, codec}

	err = fn(srv)

	if err != nil {
		return err
	}

	if blocking {
		return block.Block(func() error {
			return conn.Drain()
		})
	} else {
		block.EnableHooks() // enable shutdown hooks
		block.RegisterShutdownHook(func() error {
			return conn.Drain()
		})
	}

	return nil
}

func (s *Server) Quickapi(prefix string, db *gorm.DB, entities ...model.Entity) error {
	listOfErrors := slice.Map(entities, func(entity model.Entity) error {
		return quickapirpc.For(db, s.conn, s.codec, prefix, entity)
	})

	return slice.Fold(listOfErrors, nil, func(err, agg error) error {
		if err != nil {
			if agg != nil {
				return errors.Join(agg, err)
			}

			return err
		}

		return agg
	})
}

func (s *Server) WithConn(fn func(*nats.Conn)) {
	fn(s.conn)
}

func (s *Server) HandleRPC(topic, queue string, fn rpc.RpcHandler) error {
	_, err := rpc.HandleRPC(s.conn, s.codec, topic, queue, fn)

	return err
}

func (s *Server) HandleEvent(topic, queue string, fn rpc.EventHandler) error {
	_, err := rpc.HandleEvent(s.conn, s.codec, topic, queue, fn)

	return err
}

func (s *Server) Endpoint(e *endpoint.Endpoint) error {
	rpcHandler := handler.HandlerRegistry.RpcHandler(e.Name)
	if rpcHandler != nil {
		return s.HandleRPC(e.Route.Topic, e.Route.ConsumerGroup, rpcHandler)
	}

	eventHandler := handler.HandlerRegistry.EventHandler(e.Name)
	if eventHandler != nil {
		return s.HandleEvent(e.Route.Topic, e.Route.ConsumerGroup, eventHandler)
	}

	return fmt.Errorf("no rpcHandler with name: %s", e.Name)
}

func (s *Server) ForApp(it *app.App, only ...serviceref.ServiceRef) error {
	events := filters.FetchEndpointsByType(it, endpoint.RpcKind, only...)

	return slice.Fold(events, nil, func(e *endpoint.Endpoint, agg error) error {
		err := s.Endpoint(e)

		if err != nil {
			return errors.Join(agg, err)
		}

		return agg
	})
}
