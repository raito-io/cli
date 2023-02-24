package plugin

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
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
	assert.True(t, validateName("plugin"))
	assert.True(t, validateName("p0n"))
	assert.True(t, validateName("p0"))
	assert.True(t, validateName("pn"))
	assert.True(t, validateName("plu-gin"))
	assert.True(t, validateName("p-lugin"))
	assert.True(t, validateName("plugin0"))

	assert.False(t, validateName("0plugin"), "should not start with a number")
	assert.False(t, validateName("p"), "should be more than 1 character")
	assert.False(t, validateName("-plugin"), "should not start with a dash")
	assert.False(t, validateName("plugin-"), "should not end with a dash")
}

func TestVersionRegEx(t *testing.T) {
	assert.True(t, validateVersion("0.0.0"))
	assert.True(t, validateVersion("0.1.2"))
	assert.True(t, validateVersion("11.0.0"))
	assert.True(t, validateVersion("111.222.33333"))
	assert.True(t, validateVersion("111.222.33333-dev0"))
	assert.True(t, validateVersion("111.222.33333-dev0+build5"))

	assert.False(t, validateVersion(".0.0"))
	assert.False(t, validateVersion("jos"))
	assert.False(t, validateVersion("1"))
	assert.False(t, validateVersion("1.2"))
	assert.False(t, validateVersion("0.0."))
	assert.False(t, validateVersion("1.2.3."))
	assert.False(t, validateVersion("1.x.3"))
	assert.False(t, validateVersion("1.2.3x"))
}

func TestGetLatestVersion(t *testing.T) {
	assert.Equal(t, "1.1.2", getLatestVersion([]string{"1.1.0", "1.1.2", "1.1.1"}))
	assert.Equal(t, "2.1.0", getLatestVersion([]string{"2.1.0", "1.1.2", "1.1.1"}))
	assert.Equal(t, "2.1.0", getLatestVersion([]string{"2.1.0", "1.1.2", "2.1.0"}))
	assert.Equal(t, "2.1.1", getLatestVersion([]string{"2.1.0", "1.1.2", "2.1.1"}))
	assert.Equal(t, "2.2.2", getLatestVersion([]string{"2.1.0", "2.2.2", "2.1.1"}))

	// Non-valid versions are ignored
	assert.Equal(t, "2.2.2", getLatestVersion([]string{"2.1.0", "2.2.2", "blah"}))
	assert.Equal(t, "2.2.2", getLatestVersion([]string{"2.1.0", "2.2.2", "4.8"}))

	// Latest should always be returned as latest when present
	assert.Equal(t, "latest", getLatestVersion([]string{"2.1.0", "latest", "2.1.1"}))
}

func TestGetLatestVersionFromFiles(t *testing.T) {
	path, version := getLatestVersionFromFiles([]string{"path/group/my-file-1.1.0", "path/group/my-file-1.1.2", "path/group/my-file-1.1.1"})
	assert.Equal(t, "path/group/my-file-1.1.2", path)
	assert.Equal(t, "1.1.2", version)
	path, version = getLatestVersionFromFiles([]string{"./path/gro-up/my_file-2.1.0", "./path/gro-up/my_file-1.1.2", "./path/gro-up/my_file-1.1.1"})
	assert.Equal(t, "./path/gro-up/my_file-2.1.0", path)
	assert.Equal(t, "2.1.0", version)
	path, version = getLatestVersionFromFiles([]string{"/pa-th/group/myfile-2.1.0", "/pa-th/group/myfile-1.1.2", "/pa-th/group/myfile-2.2.0"})
	assert.Equal(t, "/pa-th/group/myfile-2.2.0", path)
	assert.Equal(t, "2.2.0", version)

	path, version = getLatestVersionFromFiles([]string{"/pa-th/group/myfile-2.1.0", "/pa-th/group/myfile-latest", "/pa-th/group/myfile-2.1.0"})
	assert.Equal(t, "/pa-th/group/myfile-latest", path)
	assert.Equal(t, "latest", version)
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
