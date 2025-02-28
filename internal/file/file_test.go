package file

import (
	"strings"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"
	"github.com/raito-io/cli/internal/util/test"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestFileUpload(t *testing.T) {
	var token, urlMethod, urlPath, uploadMethod, fileBody string

	uploadTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		uploadMethod = req.Method
		buf, _ := ioutil.ReadAll(req.Body)
		fileBody = string(buf)
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))

	defer uploadTestServer.Close()

	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		token = req.Header.Get("Authorization")
		urlMethod = req.Method
		urlPath = req.RequestURI
		res.WriteHeader(200)
		ret := "{ \"URL\": \"" + uploadTestServer.URL + "\", \"Key\": \"filekey\" }"
		res.Write([]byte(ret))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	baseConfig, closer := test.CreateBaseConfig("mydomain", "api-user", "api-secret", "")
	defer closer()

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig:   *baseConfig,
	})

	assert.Nil(t, err)
	assert.True(t, len(res) > 0)
	assert.Truef(t, strings.HasPrefix(urlPath, "/file/upload/signed-url?"), "%q does not have a correct prefix", urlPath)
	assert.Equal(t, "GET", urlMethod)
	assert.Equal(t, "PUT", uploadMethod)
	assert.Equal(t, "Hellow!", fileBody)
	assert.Equal(t, "token idToken", token)
}

func TestFileUploadNotFound(t *testing.T) {
	uploadTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))

	defer uploadTestServer.Close()

	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("{ \"URL\": \"" + uploadTestServer.URL + "\", \"Key\": \"filekey\" }"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/doesntexist.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadErrorUploading(t *testing.T) {
	uploadTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte("upload failed"))
	}))

	defer uploadTestServer.Close()

	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("{ \"URL\": \"" + uploadTestServer.URL + "\", \"Key\": \"filekey\" }"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "upload failed")
	assert.Equal(t, "", res)
}

func TestFileUploadGetUrlFailed(t *testing.T) {
	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		res.Write([]byte("errorerror"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadGetUrlIllegalResult(t *testing.T) {
	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(":::"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadGetUrlIllegalUrl(t *testing.T) {
	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("{ \"URL\": \":::\", \"Key\": \"filekey\" }"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadGetUrlNonExistingUrl(t *testing.T) {
	getUrlTestServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("{ \"URL\": \"http://localhost:9999\", \"Key\": \"filekey\" }"))
	}))

	defer getUrlTestServer.Close()

	viper.Set(constants.URLOverrideFlag, getUrlTestServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadNonExistingUrl(t *testing.T) {
	viper.Set(constants.URLOverrideFlag, "http://localhost:9999")
	defer viper.Set(constants.URLOverrideFlag, "")

	res, err := UploadFile("testdata/testfile.txt", &types.BaseTargetConfig{
		TargetLogger: hclog.L(),
		BaseConfig: types.BaseConfig{
			BaseLogger: hclog.L(),
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}
