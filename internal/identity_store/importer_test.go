package identity_store

import (
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/internal/target"
	"github.com/raito-io/cli/internal/util/url"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const GoodImportResult = "{ \"data\": { \"importIdentityStore\": { \"usersAdded\": 1, \"usersUpdated\": 2, \"usersRemoved\": 3, \"groupsAdded\": 4, \"groupsUpdated\": 5, \"groupsRemoved\": 6 } } }"
const FaultyImportResult = ":::"
const ImportResultWithErrors = "{ \"errors\": [ { \"message\": \"twisted error\" } ], \"data\": { \"importIdentityStore\": { \"usersAdded\": 1, \"usersUpdated\": 2, \"usersRemoved\": 3, \"groupsAdded\": 4, \"groupsUpdated\": 5, \"groupsRemoved\": 6 } } }"

func TestIdentityStoreImport(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.True(t, correctContent)
	assert.True(t, gotSignedURL)
	assert.True(t, didUpload)
	assert.True(t, calledImport)
	assert.Equal(t, 1, res.UsersAdded)
	assert.Equal(t, 2, res.UsersUpdated)
	assert.Equal(t, 3, res.UsersRemoved)
	assert.Equal(t, 4, res.GroupsAdded)
	assert.Equal(t, 5, res.GroupsUpdated)
	assert.Equal(t, 6, res.GroupsRemoved)
}

func TestIdentityStoreImportFailUploadUrl(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, true, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "uploading")
	assert.Contains(t, err.Error(), "upload URL")
	assert.Nil(t, res)
}

func TestIdentityStoreImportFailUpload(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(true, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "uploading")
	assert.Contains(t, err.Error(), "executing upload")
	assert.Nil(t, res)
}

func TestIdentityStoreImportFailImport(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, true, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "import")
	assert.Contains(t, strings.ToLower(err.Error()), "graphql")
	assert.Nil(t, res)
}

func TestIdentityStoreImportFaultyResponse(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, FaultyImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "invalid character")
	assert.Nil(t, res)
}

func TestIdentityStoreImportWithErrors(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, ImportResultWithErrors, &calledImport, &gotSignedURL)
	defer testServer.Close()

	url.TestURL = testServer.URL

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(f1.Name(), f2.Name())

	res, err := (*isi).TriggerImport()

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "twisted")
	assert.NotNil(t, res)
}

func UploadServer(fail bool, didUpload *bool, correctContent *bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if fail {
			res.WriteHeader(500)
			res.Write([]byte("failed"))
			return
		}
		buf, _ := ioutil.ReadAll(req.Body)
		body := string(buf)
		if body != "file1" && body != "file2" {
			*correctContent = false
		}
		res.WriteHeader(200)
		res.Write([]byte("body"))
		*didUpload = true
	}))
}

func RaitoServer(uploadUrl string, failUploadUrl bool, failQuery bool, importResult string, calledImport *bool, gotSignedURL *bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.RequestURI == "/file/upload/signed-url" && req.Method == "GET" {
			if failUploadUrl {
				res.WriteHeader(500)
				res.Write([]byte("failed"))
				return
			}
			res.WriteHeader(200)
			res.Write([]byte("{ \"URL\": \"" + uploadUrl + "\", \"Key\": \"filekey\" }"))
			*gotSignedURL = true
		} else if req.RequestURI == "/query" && req.Method == "POST" {
			if failQuery {
				res.WriteHeader(500)
				res.Write([]byte("failed"))
				return
			}
			res.WriteHeader(200)
			res.Write([]byte(importResult))
			*calledImport = true
		}
	}))
}

func writeTempFiles() (*os.File, *os.File) {
	f1, _ := ioutil.TempFile("", "is-import-test1.txt")
	f1.WriteString("file1")
	f1.Close()
	f2, _ := ioutil.TempFile("", "is-import-test2.txt")
	f2.WriteString("file2")
	f2.Close()

	return f1, f2
}

func newIdentityStoreImporter(f1, f2 string) *IdentityStoreImporter {
	isi := NewIdentityStoreImporter(&IdentityStoreImportConfig{
		BaseTargetConfig: target.BaseTargetConfig{
			Logger:    hclog.L(),
			Domain:    "mydomain",
			ApiUser:   "myuser",
			ApiSecret: "mysecret",
		},
		UserFile:        f1,
		GroupFile:       f2,
		DeleteUntouched: true,
		ReplaceGroups:   true,
		ReplaceTags:     true,
	})
	return &isi
}