package server

import (
	"net/http"

	"github.com/followboard/api/config"
	"github.com/followboard/api/elastic"
	"github.com/followboard/api/github"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server serves HTTP requests
type Server struct {
	Echo    *echo.Echo
	Elastic *elastic.Elastic
	GitHub  *github.GitHub
	Config  *config.Config
}

// New creates the server
func New(c *config.Config) *Server {
	s := &Server{
		Echo:    echo.New(),
		Elastic: elastic.New(c),
		GitHub:  github.New(),
		Config:  c,
	}

	s.Echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	s.Echo.Use(s.tokenMiddleware)

	s.Echo.POST("/hook", s.createHook)
	s.Echo.POST("/hook/event", s.handleHook)

	return s
}

func (s *Server) tokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		tokenHeader := c.Request().Header["Token"]
		if tokenHeader == nil || len(tokenHeader) == 0 {
			c.NoContent(http.StatusUnauthorized)
			return nil
		}

		token := tokenHeader[0]
		if len(token) == 0 {
			c.NoContent(http.StatusUnauthorized)
			return nil
		}

		c.Set("token", token)
		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

// Get token from context
func (s *Server) getToken(c echo.Context) string {
	return c.Get("token").(string)
}

// Start initializes the server
func (s *Server) Start() {
	glog.Fatal(s.Echo.Start(":1323"))
}
