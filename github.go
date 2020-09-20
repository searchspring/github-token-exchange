package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type DAO interface {
	GetUser(clientID, clientSecret, code, redirectURL string) ([]byte, error)
}

type GithubClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GithubClientImpl struct {
	client *http.Client
}

func (g *GithubClientImpl) Do(req *http.Request) (*http.Response, error) {
	return g.client.Do(req)
}

type DAOImpl struct {
	Client GithubClient
}

func NewDAO() DAO {
	return &DAOImpl{
		Client: &GithubClientImpl{
			client: http.DefaultClient,
		},
	}
}

type User struct {
	Login     string `json:"login"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Token     string `json:"token"`
}

func mustCreateRequest(method string, url string, jsonBuffer *bytes.Buffer) *http.Request {
	if jsonBuffer == nil {
		jsonBuffer = bytes.NewBuffer([]byte(``))
	}
	req, err := http.NewRequest(method, url, jsonBuffer)
	if err != nil {
		panic(err)
	}
	return req
}
func (d *DAOImpl) GetUser(clientID, clientSecret, code, redirectURL string) ([]byte, error) {
	jsonBuffer := bytes.NewBuffer([]byte(`{
		"client_id":"` + clientID + `",
		"client_secret":"` + clientSecret + `",
		"code":"` + code + `", 
		"redirect_uri":"` + redirectURL + `"
		}`))

	req := mustCreateRequest("POST", "https://github.com/login/oauth/access_token", jsonBuffer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	token, err := process(string(body))
	if err != nil {
		return nil, err
	}
	req = mustCreateRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "token "+token)
	resp, err = d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return json.Marshal(user)
}

func process(body string) (string, error) {
	accessTokenRegex := regexp.MustCompile(`access_token=([^&]*)`)
	tokens := accessTokenRegex.FindStringSubmatch(body)
	if len(tokens) != 2 {
		return "", fmt.Errorf("could not match access_token in '%s'", body)
	}
	return tokens[1], nil
}
