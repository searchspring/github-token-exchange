package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockClient struct {
	override func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return m.override(req)
}

func TestProcess(t *testing.T) {
	token, err := process("scope=&token_type=bearer&access_token=c20478d263c282beb19451be75dbb342a35a968c")
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "c20478d263c282beb19451be75dbb342a35a968c", token)
	_, err = process("scope=&token_type=bearer&token=c20478d263c282beb19451be75dbb342a35a968c")
	require.NotNil(t, err)
}
func TestGetUser(t *testing.T) {
	dao := &DAOImpl{Client: &mockClient{override: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/login/oauth/access_token" {
			r := ioutil.NopCloser(bytes.NewReader([]byte("access_token=mytoken&scope=repo%2Cuser%3Aemail&token_type=bearer")))
			return &http.Response{StatusCode: 200, Body: r}, nil

		}
		if req.URL.Path == "/user" {
			r := ioutil.NopCloser(bytes.NewReader([]byte(`{"login":"codeallthethingz","id":1261268,"node_id":"MDQ6VXNlcjEyNjEyNjg=","avatar_url":"https://avatars3.githubusercontent.com/u/1261268?v=4","gravatar_id":"","url":"https://api.github.com/users/codeallthethingz","html_url":"https://github.com/codeallthethingz","followers_url":"https://api.github.com/users/codeallthethingz/followers","following_url":"https://api.github.com/users/codeallthethingz/following{/other_user}","gists_url":"https://api.github.com/users/codeallthethingz/gists{/gist_id}","starred_url":"https://api.github.com/users/codeallthethingz/starred{/owner}{/repo}","subscriptions_url":"https://api.github.com/users/codeallthethingz/subscriptions","organizations_url":"https://api.github.com/users/codeallthethingz/orgs","repos_url":"https://api.github.com/users/codeallthethingz/repos","events_url":"https://api.github.com/users/codeallthethingz/events{/privacy}","received_events_url":"https://api.github.com/users/codeallthethingz/received_events","type":"User","site_admin":false,"name":"Will Warren","company":null,"blog":"","location":"Could Be Anywhere","email":null,"hireable":null,"bio":null,"twitter_username":null,"public_repos":59,"public_gists":1,"followers":11,"following":3,"created_at":"2011-12-13T20:17:30Z","updated_at":"2020-09-18T16:28:17Z"}`)))
			return &http.Response{StatusCode: 200, Body: r}, nil
		}
		return nil, fmt.Errorf("unknown path: %s", req.URL.Path)
	}}}
	user, err := dao.GetUser("clientID", "clientSecret", "code", "redirectURL")
	require.NoError(t, err)
	require.Contains(t, string(user), "codeallthethingz", string(user))
}
func TestGetUserFail(t *testing.T) {
	dao := &DAOImpl{Client: &mockClient{override: func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("failure")
	}}}
	_, err := dao.GetUser("clientID", "clientSecret", "code", "redirectURL")
	require.Error(t, err)
}
