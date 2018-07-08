package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/followboard/api/mapping"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

const (
	eventHeader = "x-github-event"
	prEvent     = "pull_request"
	prType      = "pr"
)

// HandleHookRequest to POST /hook/event
type HandleHookRequest struct {
	PR     PR   `json:"pull_request"`
	Sender User `json:"sender"`
}

// PR body
type PR struct {
	ID           int       `json:"id"`
	Number       int       `json:"number"`
	URL          string    `json:"url"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	State        string    `json:"state"`
	Commits      int       `json:"commits"`
	Additions    int       `json:"additions"`
	Deletions    int       `json:"deletions"`
	ChangedFiles int       `json:"changed_files"`
	User         User      `json:"user"`
	Head         Reference `json:"head"`
	Base         Reference `json:"base"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// User body
type User struct {
	Login string `json:"login"`
}

// Repo body
type Repo struct {
	Name  string `json:"name"`
	Owner User   `json:"owner"`
}

// Reference body
type Reference struct {
	Repo Repo   `json:"repo"`
	Ref  string `json:"ref"`
}

// Handle PR hook
func (s *Server) handleHook(c echo.Context) error {
	eventHeader := c.Request().Header[eventHeader]
	if eventHeader == nil || len(eventHeader) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	event := eventHeader[0]
	if event != prEvent {
		return c.NoContent(http.StatusNoContent)
	}

	req := new(HandleHookRequest)
	if err := c.Bind(req); err != nil {
		glog.Errorf("Failed parsing PR: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err := s.Elastic.Index(
		s.Config.PR.Index,
		prType,
		strconv.Itoa(req.PR.ID),
		s.parsePR(req.PR, req.Sender.Login),
	)

	if err != nil {
		glog.Errorf("Failed indexing PR: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// Parse GitHub PR into Elastic document
func (s *Server) parsePR(pr PR, senderLogin string) mapping.PR {
	return mapping.PR{
		FollowboardUserLogin: senderLogin,
		Number:               pr.Number,
		URL:                  pr.URL,
		Title:                pr.Title,
		Body:                 pr.Body,
		State:                pr.State,
		Commits:              pr.Commits,
		Additions:            pr.Additions,
		Deletions:            pr.Deletions,
		ChangedFiles:         pr.ChangedFiles,
		UserLogin:            pr.User.Login,
		HeadRepoOwnerLogin:   pr.Head.Repo.Owner.Login,
		HeadRepoName:         pr.Head.Repo.Name,
		HeadRef:              pr.Head.Ref,
		BaseRepoOwnerLogin:   pr.Base.Repo.Owner.Login,
		BaseRepoName:         pr.Base.Repo.Name,
		BaseRef:              pr.Base.Ref,
		CreatedAt:            pr.CreatedAt,
		UpdatedAt:            pr.UpdatedAt,
	}
}
