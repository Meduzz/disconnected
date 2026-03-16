package web

import (
	"fmt"
	"io/fs"

	gohttp "net/http"

	"github.com/Meduzz/helper/utilz"
	"github.com/Meduzz/quickapi/http"
	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Server struct {
		srv *gin.Engine
	}
)

func HttpServer(fn func(*Server) error) error {
	srv := gin.Default()

	s := &Server{srv: srv}
	err := fn(s)

	if err != nil {
		return err
	}

	port := utilz.Env("PORT", "8080")

	return srv.Run(fmt.Sprintf(":%s", port))
}

// Quickapi - mount entities at mount prefixed by contextPath if not present.
func (s *Server) Quickapi(mount string, db *gorm.DB, config *http.Config, entities ...model.Entity) {
	path := mount

	api := s.srv.Group(path)

	http.For(db, api, config, entities...)
}

// Static - mount root at mount, including contextPath if not already present.
func (s *Server) Static(mount string, root string) {
	s.srv.Static(mount, root)
}

// StaticFS - like static but with a fs.FS
func (s *Server) StaticFS(mount string, fs fs.FS) {
	httpfs := gohttp.FS(fs)
	s.srv.StaticFS(mount, httpfs)
}

// SPA - send file for all unknown routes
func (s *Server) SPA(file string) {
	s.srv.NoRoute(func(ctx *gin.Context) {
		ctx.File(file)
	})
}

// SPAFs - like SPA but with a fs.FS.
func (s *Server) SPAFs(file string, fs fs.FS) {
	s.srv.NoRoute(func(ctx *gin.Context) {
		httpfs := gohttp.FS(fs)
		ctx.FileFromFS(file, httpfs)
	})
}

// WithRouter - acts as the escape valve for all other cases.
func (s *Server) WithRouter(fn func(*gin.Engine) error) error {
	return fn(s.srv)
}
