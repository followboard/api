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
		GitHub:  github.New(c.Github.ClientID, c.Github.Secret),
		Config:  c,
	}

	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	authGroup := s.Echo.Group("/auth", s.tokenMiddleware)
	authGroup.POST("/hook", s.createHook)
	authGroup.POST("/hook/event", s.handleHook)

	s.Echo.GET("/login/callback", s.handleLoginCallback)

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

// Handle login callback by retrieving Access Token
func (s *Server) handleLoginCallback(c echo.Context) error {
	code := c.QueryParam("code")
	accessToken, err := s.GitHub.GetToken(code)
	if err != nil {
		glog.Error(err)
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, accessToken)
}

// Start initializes the server
func (s *Server) Start() {
	glog.Fatal(s.Echo.Start(":1323"))
}
