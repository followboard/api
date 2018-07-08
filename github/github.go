package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

const apiEndpoint = "https://api.github.com"

var errInternal = errors.New("{\"message\": \"Internal Server Error\"}")

// GitHub API v3 client
type GitHub struct {
	Client   *http.Client
	ClientID string
	Secret   string
}

// New creates a new GitHub client
func New(clientID string, secret string) *GitHub {
	return &GitHub{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
	}
}

// Fetch makes an HTTP request to the GitHub API v3 endpoint
func (g *GitHub) Fetch(method, url, token string, body interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(body)

	if body != nil && err != nil {
		glog.Errorf("Failed parsing request body: %v", err)
		return nil, errInternal
	}

	req, err := http.NewRequest(
		method,
		apiEndpoint+url,
		bytes.NewReader(reqBody),
	)

	if err != nil {
		glog.Errorf("Failed creating request: %v", err)
		return nil, errInternal
	}

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", token),
	)

	res, err := g.Client.Do(req)
	if err != nil {
		glog.Errorf("Failed processing request: %v", err)
		return nil, errInternal
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Errorf("Failed reading response body: %v", err)
		return nil, errInternal
	}

	if res.StatusCode >= http.StatusBadRequest {
		sBody := string(resBody)
		glog.Errorf("%s: %s [%s] %s", method, url, res.Status, sBody)
		return nil, errors.New(sBody)
	}

	return resBody, nil
}
