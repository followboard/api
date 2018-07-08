package server

import (
	"net/http"

	"github.com/followboard/api/github"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	hookName        = "web"
	hookContentType = "json"
	hookURL         = "https://followboard.now.sh/hook/event"
)

var hookEvents = []string{"pull_request"}

// CreateHookRequest to POST /hook
type CreateHookRequest struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}

// Create PR hook for org/repo
func (s *Server) createHook(c echo.Context) error {
	req := new(CreateHookRequest)

	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	err := s.GitHub.CreateHook(req.Org, req.Repo, s.getToken(c), &github.CreateHookRequest{
		Name:   hookName,
		Active: true,
		Events: hookEvents,
		Config: github.CreateHookConfig{
			URL:         hookURL,
			ContentType: hookContentType,
		},
	})

	if err != nil {
		glog.Errorf("Failed creating hook: %v", err)
		return c.JSONBlob(http.StatusInternalServerError, []byte(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
