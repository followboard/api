package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

const authEndpoint = "https://github.com/login/oauth"

type accessTokenRequest struct {
	ClientID string `json:"client_id"`
	Secret   string `json:"client_secret"`
	Code     string `json:"code"`
}

// AccessTokenResponse handles access token retrieved by GitHub API
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// GetToken retrieves Access Token from github API using authorization code
func (g *GitHub) GetToken(code string) (*AccessTokenResponse, error) {
	if len(code) == 0 {
		glog.Errorf("Authorization Code not provided!")
		return nil, errInternal
	}

	accessTokenEndpoint := authEndpoint + "/access_token"
	accessTokenRequest := accessTokenRequest{
		ClientID: g.ClientID,
		Secret:   g.Secret,
		Code:     code,
	}

	reqBody, err := json.Marshal(accessTokenRequest)
	if err != nil {
		glog.Errorf("Failed to encode request body: %v", err)
		return nil, errInternal
	}
	glog.Error(string(reqBody))

	req, err := http.NewRequest(
		"POST",
		accessTokenEndpoint,
		bytes.NewReader(reqBody),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := g.Client.Do(req)
	if err != nil {
		glog.Errorf("Failed processing request: %v", err)
		return nil, errInternal
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	glog.Error(string(resBody))
	if res.StatusCode >= 400 {
		sBody := string(resBody)
		glog.Errorf("Error calling API: %s - %s", sBody, res)
		return nil, errors.New(sBody)
	}

	var tokenResponse AccessTokenResponse
	err = json.Unmarshal(resBody, &tokenResponse)
	if err != nil {
		glog.Errorf("Failed to parse json: %v", err)
		return nil, errInternal
	}

	return &tokenResponse, nil
}
