package server

import (
	"net/http"

	"github.com/labstack/echo"
)

// PullRequest body
type PullRequest struct {
	URL   string
	Title string
	Body  string
}

// GetPRsRequest to GET /pr
type GetPRsRequest struct {
}

// GetPRsResponse from GET /pr
type GetPRsResponse struct {
	PullRequests []PullRequest
}

// Get pull requests
func (s *Server) getPRs(c echo.Context) error {
	return c.JSON(http.StatusOK, GetPRsResponse{})
}
