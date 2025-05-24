package web

import (
	"fmt"
	"strings"

	"github.com/Meduzz/helper/utilz"
	"github.com/Meduzz/quickapi/http"
	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Server struct {
		contextPath string
		srv         *gin.Engine
	}
)

func HttpServer(contextPath string, fn func(*Server)) error {
	srv := gin.Default()

	s := &Server{contextPath: contextPath, srv: srv}
	fn(s)

	port := utilz.Env("PORT", "8080")

	return srv.Run(fmt.Sprintf(":%s", port))
}

// Quickapi - mount entities at mount prefixed by contextPath if not present.
func (s *Server) Quickapi(mount string, db *gorm.DB, entities ...model.Entity) {
	path := mount

	if s.contextPath != "" && !strings.HasPrefix(mount, s.contextPath) {
		path = fmt.Sprintf("%s%s", s.contextPath, mount)
	}

	api := s.srv.Group(path)

	http.For(db, api, entities...)
}

// Static - mount root at mount, including contextPath if not already present.
func (s *Server) Static(mount string, root string) {
	path := mount

	if s.contextPath != "" && !strings.HasPrefix(mount, s.contextPath) {
		path = fmt.Sprintf("%s%s", s.contextPath, mount)
	}

	s.srv.Static(path, root)
}

// SPA - send file for all unknown routes
func (s *Server) SPA(file string) {
	s.srv.NoRoute(func(ctx *gin.Context) {
		ctx.File(file)
	})
}

// WithRouter - acts as the escape valve for all other cases.
func (s *Server) WithRouter(fn func(*gin.Engine)) {
	fn(s.srv)
}
