package server

import (
	v1 "github.com/deepaksing/FampayAssignment/server/api/v1"
	"github.com/deepaksing/FampayAssignment/store"
	"github.com/labstack/echo"
)

type Server struct {
	e     *echo.Echo
	store *store.Store
}

func NewServer(store *store.Store) *Server {
	e := echo.New()

	//Register api routes
	rootGroup := e.Group("")
	apiv1Service := v1.NewApiv1Service(store)
	apiv1Service.Register(rootGroup)

	return &Server{
		e:     e,
		store: store,
	}
}

func (s *Server) StartServer() {
	s.e.Start(":8080")
}
