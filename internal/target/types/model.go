package types

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"

	"github.com/raito-io/cli/base/util/config"
	iconfig "github.com/raito-io/cli/internal/config"
	"github.com/raito-io/cli/internal/constants"
	"github.com/raito-io/cli/internal/health_check"
	"github.com/raito-io/cli/internal/util/file"
)

type ConfigMap struct {
	Parameters map[string]string
}

func (c *ConfigMap) ToProtobufConfigMap() *config.ConfigMap {
	return &config.ConfigMap{
		Parameters: c.Parameters,
	}
}

type BaseConfig struct {
	ConfigMap
	ApiUser   string
	ApiSecret string
	Domain    string

	FileBackupLocation      string
	MaximumBackupsPerTarget int

	BaseLogger    hclog.Logger
	HealthChecker health_check.HealthChecker
	OtherArgs     []string
}

func (c *BaseConfig) ReloadConfig() error {
	apiUser, err := iconfig.HandleField(viper.GetString(constants.ApiUserFlag), reflect.String)
	if err != nil {
		return err
	}

	apiSecret, err := iconfig.HandleField(viper.GetString(constants.ApiSecretFlag), reflect.String)
	if err != nil {
		return err
	}

	domain, err := iconfig.HandleField(viper.GetString(constants.DomainFlag), reflect.String)
	if err != nil {
		return err
	}

	c.ApiUser = apiUser.(string)
	c.ApiSecret = apiSecret.(string)
	c.Domain = domain.(string)

	// Only read the parameters the first time as this is read from the command line + otherwise it would override the parameters as read from the
	if c.Parameters == nil {
		c.Parameters = BuildParameterMapFromArguments(c.OtherArgs)
	}

	return nil
}

func BuildParameterMapFromArguments(args []string) map[string]string {
	params := make(map[string]string)

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			arg := args[i][2:]
			if strings.Contains(arg, "=") {
				// The case where the flag is in the form of "--key=value"
				key := arg[0:strings.Index(arg, "=")]
				value := arg[strings.Index(arg, "=")+1:]
				params[key] = value
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				// The case where the flag is in the form of "--key value"
				params[arg] = args[i+1]
				i++
			} else {
				// Otherwise, we consider this a boolean flag
				params[arg] = "TRUE"
			}
		}
	}

	return params
}

type EnricherConfig struct {
	ConfigMap
	ConnectorName    string
	ConnectorVersion string
	Name             string
}

type BaseTargetConfig struct {
	BaseConfig
	ConnectorName    string
	ConnectorVersion string
	Name             string
	DataSourceId     string
	IdentityStoreId  string

	SkipIdentityStoreSync bool
	SkipDataSourceSync    bool
	SkipDataAccessSync    bool
	SkipDataUsageSync     bool
	SkipResourceProvider  bool
	SkipTagSync           bool

	LockAllWho            bool
	LockWhoByName         string
	LockWhoByTag          string
	LockWhoWhenIncomplete bool

	LockAllInheritance            bool
	LockInheritanceByName         string
	LockInheritanceByTag          string
	LockInheritanceWhenIncomplete bool

	LockAllWhat            bool
	LockWhatByName         string
	LockWhatByTag          string
	LockWhatWhenIncomplete bool

	LockAllNames            bool
	LockNamesByName         string
	LockNamesByTag          string
	LockNamesWhenIncomplete bool

	LockAllDelete            bool
	LockDeleteByName         string
	LockDeleteByTag          string
	LockDeleteWhenIncomplete bool

	LockAllOwners bool

	MakeNotInternalizable   string
	FullyLockAll            bool
	FullyLockByName         string
	FullyLockByTag          string
	FullyLockWhenIncomplete bool

	TagOverwriteKeyForAccessProviderName   string
	TagOverwriteKeyForAccessProviderOwners string
	TagOverwriteKeyForDataObjectOwners     string

	TagKeyAndValueForUserIsMachine string

	OnlyOutOfSyncData    bool
	SkipDataAccessImport bool

	DeleteUntouched bool
	DeleteTempFiles bool
	ReplaceGroups   bool

	DataObjectParent   *string
	DataObjectExcludes []string

	DataObjectEnrichers []*EnricherConfig

	TargetLogger hclog.Logger

	fileBackupLocationForRun string
}

// CalculateFileBackupLocationForRun calculated the full directory path where the backup files of this run should be stored.
// The current time is used for the directory name
func (c *BaseTargetConfig) CalculateFileBackupLocationForRun(runType string) error {
	c.fileBackupLocationForRun = ""

	if runType != "" && c.FileBackupLocation != "" {
		var err error
		c.FileBackupLocation, err = filepath.Abs(c.FileBackupLocation)

		if err != nil {
			return fmt.Errorf("cannot get absolute path for backup location %q", c.FileBackupLocation)
		}

		dir := strings.TrimSuffix(c.FileBackupLocation, string(filepath.Separator))
		targetFolder := file.GetFileNameFromName(c.Name) + "-" + runType

		dir += string(filepath.Separator) + targetFolder + string(filepath.Separator) + time.Now().Format("2006-01-02T15-04-05")

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create backup folder %q: %s", dir, err.Error())
		}

		c.TargetLogger.Debug(fmt.Sprintf("created backup directory %q", dir))

		c.fileBackupLocationForRun = dir
	}

	return nil
}

// FinalizeRun will remove old backup files if the maximum number of backups is exceeded
func (c *BaseTargetConfig) FinalizeRun() {
	if c.fileBackupLocationForRun != "" {
		empty, err := isFolderEmpty(c.fileBackupLocationForRun)
		if err != nil {
			c.TargetLogger.Error(fmt.Sprintf("unable to check if directory %q is empty or not: %s", c.FileBackupLocation, err.Error()))
		} else if empty {
			err = os.RemoveAll(c.fileBackupLocationForRun)
			if err != nil {
				c.TargetLogger.Error(fmt.Sprintf("unable to delete empty directory %q: %s", c.FileBackupLocation, err.Error()))
			} else {
				c.TargetLogger.Debug(fmt.Sprintf("deleted empty backup directory %q", c.fileBackupLocationForRun))
			}
		}

		if c.MaximumBackupsPerTarget > 0 {
			dir := strings.TrimSuffix(c.fileBackupLocationForRun, string(filepath.Separator))
			dir = c.fileBackupLocationForRun[0:strings.LastIndex(dir, string(filepath.Separator))]

			subFolders, err := fetchSubfolders(dir)

			if err != nil {
				c.TargetLogger.Error(fmt.Sprintf("unable to walk through backup folders in %q: %s", c.FileBackupLocation, err.Error()))
			}

			if len(subFolders) > c.MaximumBackupsPerTarget {
				removeOldestFolders(subFolders, c.MaximumBackupsPerTarget, c.TargetLogger)
			}
		}
	}
}

// HandleTempFile handles the temporary file by backing it up if needed and deleting it if needed
func (c *BaseTargetConfig) HandleTempFile(filePath string, neverDelete bool) {
	if c.fileBackupLocationForRun != "" {
		fileName := filePath[strings.LastIndex(filePath, string(filepath.Separator))+1:]

		input, err := os.ReadFile(filePath)
		if err != nil {
			c.TargetLogger.Error(fmt.Sprintf("unable to read file %q to backup: %s", filePath, err.Error()))
		}

		targetFile := c.fileBackupLocationForRun + string(filepath.Separator) + fileName

		err = os.WriteFile(targetFile, input, 0600)
		if err != nil {
			c.TargetLogger.Error(fmt.Sprintf("unable to write backup file %q: %s", targetFile, err.Error()))
		} else {
			c.TargetLogger.Debug(fmt.Sprintf("backed up file %q to %q", filePath, targetFile))
		}
	}

	if !neverDelete && c.DeleteTempFiles {
		err := os.RemoveAll(filePath)
		if err != nil {
			c.TargetLogger.Error(fmt.Sprintf("unable to delete temporary file %q: %s", filePath, err.Error()))
		} else {
			c.TargetLogger.Debug(fmt.Sprintf("removed temporary file %q", filePath))
		}
	}
}

func isFolderEmpty(folderPath string) (bool, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		// Handle errors (e.g., folder does not exist)
		return false, err
	}
	// If the length of the files slice is 0, the folder is empty
	return len(files) == 0, nil
}

func fetchSubfolders(folder string) ([]string, error) {
	var subfolders []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate the error upwards
		}

		if info.IsDir() && path != folder { // Check if it's a directory and not the root
			subfolders = append(subfolders, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return subfolders, nil
}

func removeOldestFolders(folders []string, maxFolders int, logger hclog.Logger) {
	if len(folders) > maxFolders {
		sort.Strings(folders)

		// Remove the oldest folders
		for i := 0; i < len(folders)-maxFolders; i++ {
			err := os.RemoveAll(folders[i])
			if err != nil {
				logger.Error(fmt.Sprintf("unable to remove old backup folder %q: %s", folders[i], err.Error()))
			} else {
				logger.Debug(fmt.Sprintf("removed old backup folder %q", folders[i]))
			}
		}
	}
}
