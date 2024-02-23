package connect

import (
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/test"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestDoGet(t *testing.T) {
	var token, url, method, domain string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		token = req.Header.Get("Authorization")
		domain = req.Header.Get(constants.DomainHeader)
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer testServer.Close()

	config, closer := test.CreateBaseConfig("TestRaito", "Userke", "SecretStuff", testServer.URL)
	defer closer()

	res, err := DoGetToRaito("the/path", config)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "TestRaito", domain)

	assert.Equal(t, "/the/path", url)
	assert.Equal(t, "GET", method)
}

func TestDoPost(t *testing.T) {
	var body, contentType, token, url, method, domain string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		token = req.Header.Get("Authorization")
		domain = req.Header.Get(constants.DomainHeader)
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer testServer.Close()

	config, closer := test.CreateBaseConfig("TestRaito", "Userke", "SecretStuff", testServer.URL)
	defer closer()

	res, err := DoPostToRaito("the/path", "The body", "application/json", config)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "The body", body)
	assert.Equal(t, "TestRaito", domain)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "/the/path", url)
	assert.Equal(t, "POST", method)
}

func TestDoPostIllegalURL(t *testing.T) {
	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}
	res, err := doPost("\\we\nird", "illegal path", "The body", "application/json", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDoGetIllegalURL(t *testing.T) {
	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}
	res, err := doGet("\\we\nird", "illegal path", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDoGetClosed(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	}))
	url := testServer.URL
	testServer.Close()

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	res, err := doGet(url, "the/path", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDoPostClosed(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	}))
	url := testServer.URL
	testServer.Close()

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	res, err := doPost(url, "the/path", "", "", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
