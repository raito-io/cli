package constants

// A set representing all the known flags.
// If you add a new flag constant, add it in this set as well
var KnownFlags = map[string]struct{}{
	DebugFlag:                 struct{}{},
	EnvironmentFlag:           struct{}{},
	LogFileFlag:               struct{}{},
	DomainFlag:                struct{}{},
	ApiUserFlag:               struct{}{},
	ApiSecretFlag:             struct{}{},
	ConfigFileFlag:            struct{}{},
	FrequencyFlag:             struct{}{},
	SkipIdentityStoreSyncFlag: struct{}{},
	SkipDataSourceSyncFlag:    struct{}{},
	SkipDataAccessSyncFlag:    struct{}{},
	DataSourceIdFlag:          struct{}{},
	IdentityStoreIdFlag:       struct{}{},
	OnlyTargetsFlag:           struct{}{},
	ConnectorNameFlag:         struct{}{},
	ConnectorVersionFlag:      struct{}{},
	NameFlag:                  struct{}{},
	DeleteUntouchedFlag:       struct{}{},
	ReplaceTagsFlag:           struct{}{},
	DeleteTempFilesFlag:       struct{}{},
	ReplaceGroupsFlag:         struct{}{},
	AccessFileFlag:            struct{}{},
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
)
