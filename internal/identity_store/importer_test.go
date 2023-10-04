package identity_store

import (
	"context"

	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/target/types"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/internal/job"
	"github.com/raito-io/cli/internal/job/mocks"
)

const GoodImportResult = "{ \"data\": { \"importIdentityRequest\": { \"subtask\" : {\"status\": \"QUEUED\", \"subtaskId\": \"ImportSync\" } } } }"
const FaultyImportResult = ":::"
const ImportResultWithErrors = "{ \"errors\": [ { \"message\": \"twisted error\" } ], \"data\": { \"importIdentityStore\": { \"usersAdded\": 1, \"usersUpdated\": 2, \"usersRemoved\": 3, \"groupsAdded\": 4, \"groupsUpdated\": 5, \"groupsRemoved\": 6 } } }"

func TestIdentityStoreImport(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	defer os.Remove(f1.Name())
	defer os.Remove(f2.Name())
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, subtaskId, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.Nil(t, err)
	assert.True(t, correctContent)
	assert.True(t, gotSignedURL)
	assert.True(t, didUpload)
	assert.True(t, calledImport)
	assert.Equal(t, job.Queued, status)
	assert.Equal(t, "ImportSync", subtaskId)
}

func TestIdentityStoreImportFailUploadUrl(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, true, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, _, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "uploading")
	assert.Contains(t, err.Error(), "upload URL")
	assert.Equal(t, job.Failed, status)
}

func TestIdentityStoreImportFailUpload(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(true, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, _, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "uploading")
	assert.Contains(t, err.Error(), "executing upload")
	assert.Equal(t, job.Failed, status)
}

func TestIdentityStoreImportFailImport(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, true, GoodImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, _, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "import")
	assert.Contains(t, strings.ToLower(err.Error()), "graphql")
	assert.Equal(t, job.Failed, status)
}

func TestIdentityStoreImportFaultyResponse(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, FaultyImportResult, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, _, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "invalid character")
	assert.Equal(t, job.Failed, status)
}

func TestIdentityStoreImportWithErrors(t *testing.T) {
	var didUpload, gotSignedURL, calledImport bool
	correctContent := true

	uploadTestServer := UploadServer(false, &didUpload, &correctContent)
	defer uploadTestServer.Close()

	testServer := RaitoServer(uploadTestServer.URL, false, false, ImportResultWithErrors, &calledImport, &gotSignedURL)
	defer testServer.Close()

	viper.Set(constants.URLOverrideFlag, testServer.URL)
	defer viper.Set(constants.URLOverrideFlag, "")

	f1, f2 := writeTempFiles()
	isi := newIdentityStoreImporter(t, f1.Name(), f2.Name())

	status, _, err := (*isi).TriggerImport(context.Background(), "someJobId")

	assert.NotNil(t, err)
	assert.Contains(t, strings.ToLower(err.Error()), "twisted")
	assert.Equal(t, job.Failed, status)
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

func newIdentityStoreImporter(t *testing.T, f1, f2 string) *IdentityStoreImporter {
	t.Helper()

	isi := NewIdentityStoreImporter(&IdentityStoreImportConfig{
		BaseTargetConfig: types.BaseTargetConfig{
			TargetLogger: hclog.L(),
			BaseConfig: types.BaseConfig{
				Domain:     "mydomain",
				ApiUser:    "myuser",
				ApiSecret:  "mysecret",
				BaseLogger: hclog.L(),
			},
		},
		UserFile:        f1,
		GroupFile:       f2,
		DeleteUntouched: true,
		ReplaceGroups:   true,
	}, dummyTaskEventUpdater(t))
	return &isi
}

func dummyTaskEventUpdater(t *testing.T) *mocks.TaskEventUpdater {
	t.Helper()

	m := mocks.NewTaskEventUpdater(t)
	m.EXPECT().SetStatusToStarted(mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToDataRetrieve(mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToDataUpload(mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToQueued(mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToDataProcessing(mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToCompleted(mock.Anything, mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToFailed(mock.Anything, mock.Anything).Return().Maybe()
	m.EXPECT().SetStatusToSkipped(mock.Anything).Return().Maybe()
	m.EXPECT().GetSubtaskEventUpdater(mock.Anything).Return(nil).Maybe()

	return m
}
