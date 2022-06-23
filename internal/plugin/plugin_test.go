package plugin

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestNewClientError(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("WORKING DIR: " + path)
	client, err := NewPluginClient("blah/notexisting", "", hclog.L())
	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestNameRegEx(t *testing.T) {
	assert.True(t, nameRegexp.Match([]byte("plugin")))
	assert.True(t, nameRegexp.Match([]byte("p0n")))
	assert.True(t, nameRegexp.Match([]byte("p0")))
	assert.True(t, nameRegexp.Match([]byte("pn")))
	assert.True(t, nameRegexp.Match([]byte("plu-gin")))
	assert.True(t, nameRegexp.Match([]byte("p-lugin")))
	assert.True(t, nameRegexp.Match([]byte("plugin0")))

	assert.False(t, nameRegexp.Match([]byte("0plugin")), "should not start with a number")
	assert.False(t, nameRegexp.Match([]byte("p")), "should be more than 1 character")
	assert.False(t, nameRegexp.Match([]byte("-plugin")), "should not start with a dash")
	assert.False(t, nameRegexp.Match([]byte("plugin-")), "should not end with a dash")
}

func TestVersionRegEx(t *testing.T) {
	assert.True(t, versionRegexp.Match([]byte("0.0.0")))
	assert.True(t, versionRegexp.Match([]byte("0.1.2")))
	assert.True(t, versionRegexp.Match([]byte("11.0.0")))
	assert.True(t, versionRegexp.Match([]byte("111.222.33333")))

	assert.False(t, versionRegexp.Match([]byte(".0.0")))
	assert.False(t, versionRegexp.Match([]byte("jos")))
	assert.False(t, versionRegexp.Match([]byte("1")))
	assert.False(t, versionRegexp.Match([]byte("1.2")))
	assert.False(t, versionRegexp.Match([]byte("0.0.")))
	assert.False(t, versionRegexp.Match([]byte("1.2.3.")))
	assert.False(t, versionRegexp.Match([]byte("1.x.3")))
	assert.False(t, versionRegexp.Match([]byte("1.2.3x")))
}

func TestGetLatestVersion(t *testing.T) {
	assert.Equal(t, "1.1.2", getLatestVersion([]string{"1.1.0", "1.1.2", "1.1.1"}))
	assert.Equal(t, "2.1.0", getLatestVersion([]string{"2.1.0", "1.1.2", "1.1.1"}))
	assert.Equal(t, "2.1.0", getLatestVersion([]string{"2.1.0", "1.1.2", "2.1.0"}))
	assert.Equal(t, "2.1.1", getLatestVersion([]string{"2.1.0", "1.1.2", "2.1.1"}))
	assert.Equal(t, "2.2.2", getLatestVersion([]string{"2.1.0", "2.2.2", "2.1.1"}))
}

func TestGetLatestVersionFromFiles(t *testing.T) {
	assert.Equal(t, "path/group/my-file-1.1.2", getLatestVersionFromFiles([]string{"path/group/my-file-1.1.0", "path/group/my-file-1.1.2", "path/group/my-file-1.1.1"}))
	assert.Equal(t, "./path/gro-up/my_file-2.1.0", getLatestVersionFromFiles([]string{"./path/gro-up/my_file-2.1.0", "./path/gro-up/my_file-1.1.2", "./path/gro-up/my_file-1.1.1"}))
	assert.Equal(t, "/pa-th/group/myfile-2.1.0", getLatestVersionFromFiles([]string{"/pa-th/group/myfile-2.1.0", "/pa-th/group/myfile-1.1.2", "/pa-th/group/myfile-2.1.0"}))
}

// Commented as we should not be downloading stuff during unit tests. Can we fake this?
/*func TestGetPluginFromPublicRegistry(t *testing.T) {
	tmpDir := os.TempDir()
	if !strings.HasSuffix(tmpDir, "/") {
		tmpDir = tmpDir + "/"
	}
	tmpDir = tmpDir + "raito-plugin-test"
	file, err := getPluginFromPublicRegistry(&pluginRequest{
		Group: "raito",
		Name: "okta",
		Version: "latest",
	}, tmpDir)

	defer os.RemoveAll(tmpDir+"/")

	fmt.Println("Full file: "+file)

	assert.Nil(t, err)
	assert.Equal(t, tmpDir + "/raito/okta-0.1.0", file)

	fileInfo, err := os.Stat(file)
	assert.Nil(t, err)
	assert.True(t, fileInfo.Size() > 1000000)
	assert.Equal(t, fs.FileMode(0750), fileInfo.Mode())
}*/

// TODO fix plugin testing (should have a better way than having a binary in the repo).
/*func TestNewClientOK(t *testing.T) {
	err := os.Chmod("./testdata/okta", 0777)
	assert.Nil(t, err)
	client, err := NewPluginClient("./testdata/okta", hclog.L())
	assert.Nil(t, err)
	assert.NotNil(t, client)
	if client != nil {
		defer client.Close()
	}
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestNewClientIdentityStore(t *testing.T) {
	err := os.Chmod("./testdata/okta", 0777)
	assert.Nil(t, err)
	client, err := NewPluginClient("./testdata/okta", hclog.L())
	assert.Nil(t, err)
	assert.NotNil(t, client)
	if client != nil {
		defer client.Close()
	} else {
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, client)

	sync, err := client.GetIdentityStoreSyncer()
	assert.Nil(t, err)
	assert.NotNil(t, sync)
}

func TestNewClientNotImplemented(t *testing.T) {
	err := os.Chmod("./testdata/okta", 0777)
	assert.Nil(t, err)
	client, err := NewPluginClient("./testdata/okta", hclog.L())
	assert.Nil(t, err)
	assert.NotNil(t, client)
	if client != nil {
		defer client.Close()
	} else {
		return
	}
	assert.Nil(t, err)
	assert.NotNil(t, client)

	sync, err := client.GetDataAccessSyncer()
	assert.NotNil(t, err)
	assert.Nil(t, sync)

	sync2, err := client.GetDataSourceSyncer()
	assert.NotNil(t, err)
	assert.Nil(t, sync2)
}*/
