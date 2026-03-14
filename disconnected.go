package disconnected

import (
	"github.com/Meduzz/disconnected/pkg/event"
	"github.com/Meduzz/disconnected/pkg/web"
	"github.com/Meduzz/rpc/encoding"
)

func HttpServer(fn func(*web.Server) error) error {
	return web.HttpServer(fn)
}

func RpcServer(codec encoding.Codec, fn func(*event.Server) error) error {
	return event.EventServer(codec, fn)
}
