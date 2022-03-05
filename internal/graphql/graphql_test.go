package graphql

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

func TestGraphQL(t *testing.T) {
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

	buf, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.Nil(t, err)
	assert.NotNil(t, buf)
	assert.Equal(t, "{ \"operationName\": \"nastyOperation\" }", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "TestRaito", domainHeader)
	assert.Equal(t, "Userke", user)
	assert.Equal(t, "SecretStuff", secret)
	assert.Equal(t, "/query", url)
	assert.Equal(t, "POST", method)
	assert.Equal(t, "body", string(buf))
}

func TestGraphQLError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte("bad stuff"))
	}))
	defer testServer.Close()

	url2.TestURL = testServer.URL

	config := target.BaseTargetConfig{
		Logger:    hclog.Default(),
	}

	buf, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.NotNil(t, err)
	assert.Nil(t, buf)
}

func TestGraphQLIllegalURL(t *testing.T) {
	url2.TestURL = "//\nbadbadbad"

	config := target.BaseTargetConfig{
		Logger:    hclog.Default(),
	}

	buf, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.NotNil(t, err)
	assert.Nil(t, buf)
}
