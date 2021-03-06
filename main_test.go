package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockDAO struct {
	Override func(clientID, clientSecret, code, redirectURL string) ([]byte, error)
}

func (m *mockDAO) GetUser(clientID, clientSecret, code, redirectURL string) ([]byte, error) {
	return m.Override(clientID, clientSecret, code, redirectURL)
}

func simple(clientID, clientSecret, code, redirectURL string) ([]byte, error) {
	return []byte("userstring"), nil
}

func returnError(clientID, clientSecret, code, redirectURL string) ([]byte, error) {
	return nil, fmt.Errorf("explosion")
}
func TestCallout(t *testing.T) {
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:1231/?code=blah", nil)
	githubDAO = &mockDAO{Override: simple}
	handler(res, req)
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, 200, res.Result().StatusCode)
	require.True(t, strings.Contains(string(body), "let user = userstring"))
}
func TestEdges(t *testing.T) {
	githubDAO = &mockDAO{Override: simple}
	res := httptest.NewRecorder()
	handler(res, httptest.NewRequest("GET", "http://localhost:1231/?code=", nil))
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	require.True(t, strings.Contains(string(body), ""))
}

func TestBadPath(t *testing.T) {
	githubDAO = &mockDAO{Override: simple}
	res := httptest.NewRecorder()
	handler(res, httptest.NewRequest("GET", "http://localhost:1231/somethingelse/?code=", nil))
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.Result().StatusCode)
	require.Equal(t, "", string(body))
}

func TestErrorFromGithub(t *testing.T) {
	githubDAO = &mockDAO{Override: returnError}
	res := httptest.NewRecorder()
	handler(res, httptest.NewRequest("GET", "http://localhost:1231/?code=aoeu", nil))
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.Equal(t, http.StatusForbidden, res.Result().StatusCode)
	require.Equal(t, "explosion\n", string(body))
}

func TestChecks(t *testing.T) {
	os.Setenv("PORT", "")
	os.Setenv("GITHUB_CLIENT_SECRET", "")
	os.Setenv("GITHUB_CLIENT_ID", "")
	os.Setenv("GITHUB_REDIRECT_URL", "")
	os.Setenv("ALLOWLIST_REDIRECT_URLS", "")
	testFail(t)

	os.Setenv("PORT", "8888")
	testFail(t)
	os.Setenv("GITHUB_REDIRECT_URL", "aoeu")
	testFail(t)
	os.Setenv("GITHUB_CLIENT_ID", "aoeu")
	testFail(t)
	os.Setenv("ALLOWLIST_REDIRECT_URLS", "http://localhost,https://localhost,https://searchspring.github.io/snapp-explorer")
	testFail(t)
	os.Setenv("ALLOWLIST_REDIRECT_URLS", "http://localhost,https://localhost ,  https://searchspring.github.io/snapp-explorer, ")
	testFail(t)
}

func TestAllowlistPass(t *testing.T) {
	githubDAO = &mockDAO{Override: simple}
	res := httptest.NewRecorder()
	handler(res, httptest.NewRequest("GET", "http://localhost:1231/?code=blah&redirect=https://searchspring.github.io/snapp-explorer", nil))
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.True(t, strings.Contains(string(body), "https://searchspring.github.io/snapp-explorer"))
	require.False(t, strings.Contains(string(body), "http://localhost:3827"))
}

func TestAllowlistFail(t *testing.T) {
	githubDAO = &mockDAO{Override: simple}
	res := httptest.NewRecorder()
	handler(res, httptest.NewRequest("GET", "http://localhost:1231/?code=blah&redirect=https://dne.searchspring.io", nil))
	body, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	require.True(t, strings.Contains(string(body), "http://localhost:3827"))
	require.False(t, strings.Contains(string(body), "https://dne.searchspring.io"))
}

func testFail(t *testing.T) {
	defer func() {
		//nolint
		if r := recover(); r != nil {
			// do nothing
		}
	}()
	main()
	t.Fail()
}
