package constants

// A set representing all the known flags.
// If you add a new flag constant, add it in this set as well
var KnownFlags = map[string]struct{}{
	DebugFlag:                 {},
	EnvironmentFlag:           {},
	LogFileFlag:               {},
	DomainFlag:                {},
	ApiUserFlag:               {},
	ApiSecretFlag:             {},
	ConfigFileFlag:            {},
	FrequencyFlag:             {},
	SkipIdentityStoreSyncFlag: {},
	SkipDataSourceSyncFlag:    {},
	SkipDataAccessSyncFlag:    {},
	SkipDataUsageSyncFlag:     {},
	DataSourceIdFlag:          {},
	IdentityStoreIdFlag:       {},
	OnlyTargetsFlag:           {},
	ConnectorNameFlag:         {},
	ConnectorVersionFlag:      {},
	NameFlag:                  {},
	DeleteUntouchedFlag:       {},
	ReplaceTagsFlag:           {},
	DeleteTempFilesFlag:       {},
	ReplaceGroupsFlag:         {},
	AccessFileFlag:            {},
}

const (
	DebugFlag                 = "debug"
	EnvironmentFlag           = "environment"
	LogFileFlag               = "log-file"
	LogOutputFlag             = "log-output"
	DomainFlag                = "domain"
	ApiUserFlag               = "api-user"
	ApiSecretFlag             = "api-secret"
	ConfigFileFlag            = "config-file"
	FrequencyFlag             = "frequency"
	SkipDataSourceSyncFlag    = "skip-data-source-sync"
	SkipDataAccessSyncFlag    = "skip-data-access-sync"
	SkipIdentityStoreSyncFlag = "skip-identity-store-sync"
	SkipDataUsageSyncFlag     = "skip-data-usage-sync"
	DataSourceIdFlag          = "data-source-id"
	IdentityStoreIdFlag       = "identity-store-id"
	OnlyTargetsFlag           = "only-targets"

	ConnectorNameFlag    = "connector-name"
	ConnectorVersionFlag = "connector-version"
	NameFlag             = "name"

	// Environments
	EnvironmentProd = "prod"
	EnvironmentDev  = "dev"
	EnvironmentTest = "test"

	// Import specific flags
	DeleteUntouchedFlag = "delete-untouched"
	ReplaceTagsFlag     = "replace-tags"
	DeleteTempFilesFlag = "delete-temp-files"
	ReplaceGroupsFlag   = "replace-groups"

	// Access specific parameters
	AccessFileFlag = "access-file"

	Targets      = "targets"
	Repositories = "repositories"

	GitHubToken = "token"

	Started   = "STARTED"
	Completed = "COMPLETED"
	Failed    = "FAILED"
	Skipped   = "SKIPPED"

	IdentitySync   = "IS"
	DataSourceSync = "DS"
	DataAccessSync = "DA"
	DataUsageSync  = "DU"
	Job            = "JOB"
)
