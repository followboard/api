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
)

// CreateHookRequest to POST /hook
type CreateHookRequest struct {
	Org    string   `json:"org"`
	Repo   string   `json:"repo"`
	Events []string `json:"events"`
	URL    string   `json:"url"`
}

// CreateHookResponse from POST /hook
type CreateHookResponse struct {
}

// HandleHookRequest to POST /hook/event
type HandleHookRequest struct {
}

// HandleHookResponse from POST /hook/event
type HandleHookResponse struct {
}

// Create hook with request body at org/repo
func (s *Server) createHook(c echo.Context) error {
	req := new(CreateHookRequest)

	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	_, err := s.GitHub.CreateHook(req.Org, req.Repo, s.getToken(c), &github.CreateHookRequest{
		Name:   hookName,
		Active: true,
		Events: req.Events,
		Config: github.CreateHookConfig{
			URL:         req.URL,
			ContentType: hookContentType,
		},
	})

	if err != nil {
		glog.Errorf("Failed creating hook: %v", err.Error())
		return c.JSONBlob(http.StatusInternalServerError, []byte(err.Error()))
	}

	return c.JSON(http.StatusOK, CreateHookResponse{})
}

// Handle hook by event type
func (s *Server) handleHook(c echo.Context) error {
	return c.JSON(http.StatusOK, HandleHookResponse{})
}
