package event

import (
	"errors"

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

func EventServer(codec encoding.Codec, fn func(*Server) error) error {
	conn, err := nuts.Connect()

	if err != nil {
		return err
	}

	srv := &Server{conn, codec}

	err = fn(srv)

	if err != nil {
		return err
	}

	return block.Block(func() error {
		return conn.Drain()
	})
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
