package file

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/url"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestCreateUniqueFileName(t *testing.T) {
	fileNames := make(map[string]struct{})
	for i := 0; i < 10000; i++ {
		fileName := CreateUniqueFileName("thehint", "yml")
		assert.True(t, strings.HasPrefix(fileName, "thehint"), "Filename doesn't have the right prefix")
		assert.True(t, strings.HasSuffix(fileName, ".yml"), "Filename doesn't have the right suffix")
		fmt.Println(fileName)
		_, found := fileNames[fileName]
		assert.False(t, found, "Duplicate filename found ("+strconv.Itoa(i)+")")
		fileNames[fileName] = struct{}{}
	}

}

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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger:    hclog.L(),
		Domain:    "mydomain",
		ApiUser:   "api-user",
		ApiSecret: "api-secret",
	})

	assert.Nil(t, err)
	assert.True(t, len(res) > 0)
	assert.Equal(t, "/file/upload/signed-url", urlPath)
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/doesntexist.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
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

	url.TestURL = getUrlTestServer.URL

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}

func TestFileUploadNonExistingUrl(t *testing.T) {
	url.TestURL = "http://localhost:9999"

	res, err := UploadFile("testdata/testfile.txt", &target.BaseTargetConfig{
		Logger: hclog.L(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, "", res)
}
