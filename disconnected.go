package disconnected

import (
	"github.com/Meduzz/disconnected/pkg/event"
	"github.com/Meduzz/disconnected/pkg/web"
	"github.com/Meduzz/rpc/encoding"
)

func HttpServer(contextPath string, fn func(*web.Server)) error {
	return web.HttpServer(contextPath, fn)
}

func RpcServer(codec encoding.Codec, fn func(*event.Server) error) error {
	return event.EventServer(codec, fn)
}
