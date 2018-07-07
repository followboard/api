package server

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server serves HTTP requests
type Server struct {
	Echo *echo.Echo
}

// New creates the server
func New() *Server {
	s := &Server{
		Echo: echo.New(),
	}

	s.Echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	s.Echo.GET("/", s.index)

	return s
}

func (s *Server) index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

// Start initializes the server
func (s *Server) Start() {
	glog.Fatal(s.Echo.Start(":1323"))
}
