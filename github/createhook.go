package github

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

const createHookURL = "/repos/%s/%s/hooks"

// CreateHookRequest to POST /repos/:org/:repo/hooks
type CreateHookRequest struct {
	Name   string           `json:"name"`
	Active bool             `json:"active"`
	Events []string         `json:"events"`
	Config CreateHookConfig `json:"config"`
}

// CreateHookConfig body
type CreateHookConfig struct {
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
}

// CreateHook with request body at org/repo
func (g *GitHub) CreateHook(org, repo, token string, body *CreateHookRequest) error {
	_, err := g.Fetch(
		http.MethodPost,
		fmt.Sprintf(createHookURL, org, repo),
		token,
		body,
	)

	if err != nil {
		glog.Errorf("Failed creating hook: %v", err)
		return err
	}

	return nil
}
