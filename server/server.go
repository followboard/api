package server

import (
	"net/http"

	"github.com/followboard/api/github"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Server serves HTTP requests
type Server struct {
	Echo   *echo.Echo
	GitHub *github.GitHub
}

// New creates the server
func New() *Server {
	s := &Server{
		Echo:   echo.New(),
		GitHub: github.New(),
	}

	s.Echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	s.Echo.Use(s.tokenMiddleware)

	s.Echo.GET("/pr", s.getPRs)
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
