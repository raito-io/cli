package graphql

import (
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

type dataObject struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Mass   int     `json:"mass"`
}

func TestGraphQL(t *testing.T) {
	var body, contentType, token, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		token = req.Header.Get("Authorization")
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte(`
		{
			"data": {
				"name": "Luke Skywalker",
				"height": 1.72,
				"mass": 77
			},
			"errors": []
		}
		`))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		Domain:     "TestRaito",
		ApiUser:    "Userke",
		ApiSecret:  "SecretStuff",
		BaseLogger: hclog.Default(),
	}

	data := dataObject{}
	gqlResponse, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config, &data)

	assert.Nil(t, err)
	assert.NotNil(t, gqlResponse)
	assert.Equal(t, "{ \"operationName\": \"nastyOperation\" }", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "/query", url)
	assert.Equal(t, "POST", method)
	assert.Equal(t, "Luke Skywalker", data.Name)
	assert.Equal(t, 1.72, data.Height)
	assert.Equal(t, 77, data.Mass)
	assert.Empty(t, gqlResponse.Errors)
	assert.Equal(t, &data, gqlResponse.Data)
}

func TestGraphQLError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte("bad stuff"))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	data := dataObject{}
	gqlResponse, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config, &data)

	assert.NotNil(t, err)
	assert.Nil(t, gqlResponse)
}

func TestGraphQLIllegalURL(t *testing.T) {
	viper.Set(constants.URLOverrideFlag, "//\nbadbadbad")
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	data := dataObject{}
	gqlReponse, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config, &data)

	assert.NotNil(t, err)
	assert.Nil(t, gqlReponse)
}

func TestGraphQLServerError(t *testing.T) {
	var body, contentType, token, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		token = req.Header.Get("Authorization")
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte(`
		{
			"data": {
				"name": "Luke Skywalker",
				"height": 1.72,
				"mass": 77
			},
			"errors": [
				{
					"message": "Boom! This is an error message"
				},
				{
					"message": "A second error"
				}
			]
		}
		`))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		Domain:     "TestRaito",
		ApiUser:    "Userke",
		ApiSecret:  "SecretStuff",
		BaseLogger: hclog.Default(),
	}

	data := dataObject{}
	gqlResponse, err := ExecuteGraphQL("{ \"operationName\": \"nastyOperation\" }", &config, &data)

	assert.NotNil(t, err)
	assert.NotNil(t, gqlResponse)
	assert.Equal(t, "{ \"operationName\": \"nastyOperation\" }", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "/query", url)
	assert.Equal(t, "POST", method)
	assert.Equal(t, "Luke Skywalker", data.Name)
	assert.Equal(t, 1.72, data.Height)
	assert.Equal(t, 77, data.Mass)
	assert.Len(t, gqlResponse.Errors, 2)
	assert.Equal(t, &data, gqlResponse.Data)

	if merr, ok := err.(*multierror.Error); ok {
		assert.Len(t, merr.Errors, 2)
	} else {
		assert.Fail(t, "Expecting mutlierror")
	}
}

func TestGraphQLWithoutResponse(t *testing.T) {
	var body, contentType, token, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		token = req.Header.Get("Authorization")
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte(`
		{
			"data": {
				"name": "Luke Skywalker",
				"height": 1.72,
				"mass": 77
			},
			"errors": []
		}
		`))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		Domain:     "TestRaito",
		ApiUser:    "Userke",
		ApiSecret:  "SecretStuff",
		BaseLogger: hclog.Default(),
	}

	err := ExecuteGraphQLWithoutResponse("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.Nil(t, err)
	assert.Equal(t, "{ \"operationName\": \"nastyOperation\" }", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "/query", url)
	assert.Equal(t, "POST", method)
}

func TestGraphQLErrorWithoutResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte("bad stuff"))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	err := ExecuteGraphQLWithoutResponse("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.NotNil(t, err)
}

func TestGraphQLIllegalURLWithoutReponse(t *testing.T) {
	viper.Set(constants.URLOverrideFlag, "//\nbadbadbad")
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		BaseLogger: hclog.Default(),
	}

	err := ExecuteGraphQLWithoutResponse("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.NotNil(t, err)
}

func TestGraphQLServerErrorWithoutResponse(t *testing.T) {
	var body, contentType, token, url, method string

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		buf, _ := ioutil.ReadAll(req.Body)
		body = string(buf)
		contentType = req.Header.Get("Content-Type")
		token = req.Header.Get("Authorization")
		method = req.Method

		url = req.RequestURI
		res.WriteHeader(200)
		res.Write([]byte(`
		{
			"data": {
				"name": "Luke Skywalker",
				"height": 1.72,
				"mass": 77
			},
			"errors": [
				{
					"message": "Boom! This is an error message"
				},
				{
					"message": "A second error"
				}
			]
		}
		`))
	}))
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	config := types.BaseConfig{
		Domain:     "TestRaito",
		ApiUser:    "Userke",
		ApiSecret:  "SecretStuff",
		BaseLogger: hclog.Default(),
	}

	err := ExecuteGraphQLWithoutResponse("{ \"operationName\": \"nastyOperation\" }", &config)

	assert.NotNil(t, err)
	assert.Equal(t, "{ \"operationName\": \"nastyOperation\" }", body)
	assert.Equal(t, "application/json", contentType)
	assert.Equal(t, "token idToken", token)
	assert.Equal(t, "/query", url)
	assert.Equal(t, "POST", method)

	if merr, ok := err.(*multierror.Error); ok {
		assert.Len(t, merr.Errors, 2)
	} else {
		assert.Fail(t, "Expecting mutlierror")
	}
}
