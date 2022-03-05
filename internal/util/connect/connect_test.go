package connect

import (
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target"
	url2 "github.com/raito-io/cli/internal/util/url"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoGet(t *testing.T) {
	var domainHeader, user, secret, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		domainHeader = req.Header.Get(constants.OrgDomainHeader)
		user = req.Header.Get(constants.ApiUserHeader)
		secret = req.Header.Get(constants.ApiSecretHeader)
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer testServer.Close()

	url2.TestURL = testServer.URL

	config := target.BaseTargetConfig{
		Domain:    "TestRaito",
		ApiUser:   "Userke",
		ApiSecret: "SecretStuff",
		Logger:    hclog.Default(),
	}

	res, err := DoGetToRaito("the/path", &config)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "TestRaito", domainHeader)
	assert.Equal(t, "Userke", user)
	assert.Equal(t, "SecretStuff", secret)
	assert.Equal(t, "/the/path", url)
	assert.Equal(t, "GET", method)
}

func TestDoPost(t *testing.T) {
	var body, contentType, domainHeader, user, secret, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		domainHeader = req.Header.Get(constants.OrgDomainHeader)
		user = req.Header.Get(constants.ApiUserHeader)
		secret = req.Header.Get(constants.ApiSecretHeader)
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer testServer.Close()

	url2.TestURL = testServer.URL

	config := target.BaseTargetConfig{
		Domain:    "TestRaito",
		ApiUser:   "Userke",
		ApiSecret: "SecretStuff",
		Logger:    hclog.Default(),
	}

	res, err := DoPostToRaito("the/path", "The body", "application/json", &config)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "The body", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "TestRaito", domainHeader)
	assert.Equal(t, "Userke", user)
	assert.Equal(t, "SecretStuff", secret)
	assert.Equal(t, "/the/path", url)
	assert.Equal(t, "POST", method)
}

func TestDoPostIllegalURL(t *testing.T) {
	config := target.BaseTargetConfig{
		Logger: hclog.Default(),
	}
	res, err := doPost("\\we\nird", "illegal path", "The body", "application/json", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}

func TestDoGetIllegalURL(t *testing.T) {
	config := target.BaseTargetConfig{
		Logger: hclog.Default(),
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

	config := target.BaseTargetConfig{
		Logger: hclog.Default(),
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

	config := target.BaseTargetConfig{
		Logger: hclog.Default(),
	}

	res, err := doPost(url, "the/path", "", "", &config)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
