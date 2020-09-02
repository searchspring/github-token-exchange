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

type DAOImpl struct {
	Client *http.Client
}

func NewDAO() DAO {
	client := &http.Client{}
	return &DAOImpl{
		Client: client,
	}
}

type User struct {
	Login     string `json:"login"`
	NodeID    string `json:"node_id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Token     string `json:"token"`
}

func (d *DAOImpl) GetUser(clientID, clientSecret, code, redirectURL string) ([]byte, error) {
	jsonBuffer := bytes.NewBuffer([]byte(`{
		"client_id":"` + clientID + `",
		"client_secret":"` + clientSecret + `",
		"code":"` + code + `", 
		"redirect_uri":"` + redirectURL + `"
		}`))

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", jsonBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
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
	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	resp, err = client.Do(req)
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
