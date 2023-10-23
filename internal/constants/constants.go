package constants

var KnownFlags = map[string]struct{}{
	DebugFlag:                 {},
	LogFileFlag:               {},
	SkipAuthentication:        {},
	SkipFileUpload:            {},
	URLOverrideFlag:           {},
	DomainFlag:                {},
	ApiUserFlag:               {},
	ApiSecretFlag:             {},
	ConfigFileFlag:            {},
	FrequencyFlag:             {},
	CronFlag:                  {},
	SyncAtStartupFlag:         {},
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
	DeleteTempFilesFlag:       {},
	ReplaceGroupsFlag:         {},
	AccessFileFlag:            {},
}

const (
	DebugFlag                             = "debug"
	URLOverrideFlag                       = "raito-url-override"
	SkipAuthentication                    = "skip-authentication"
	SkipFileUpload                        = "skip-file-upload"
	LogFileFlag                           = "log-file"
	LogOutputFlag                         = "log-output"
	DomainFlag                            = "domain"
	ApiUserFlag                           = "api-user"
	ApiSecretFlag                         = "api-secret"
	ConfigFileFlag                        = "config-file"
	FrequencyFlag                         = "frequency"
	CronFlag                              = "cron"
	SyncAtStartupFlag                     = "sync-at-startup"
	SkipDataSourceSyncFlag                = "skip-data-source-sync"
	SkipDataAccessSyncFlag                = "skip-data-access-sync"
	SkipIdentityStoreSyncFlag             = "skip-identity-store-sync"
	SkipDataUsageSyncFlag                 = "skip-data-usage-sync"
	LockAllWhoFlag                        = "lock-all-who"
	LockAllInheritanceFlag                = "lock-all-inheritance"
	LockAllWhatFlag                       = "lock-all-what"
	LockAllNamesFlag                      = "lock-all-names"
	LockAllDeleteFlag                     = "lock-all-delete"
	DataSourceIdFlag                      = "data-source-id"
	IdentityStoreIdFlag                   = "identity-store-id"
	OnlyTargetsFlag                       = "only-targets"
	DisableWebsocketFlag                  = "disable-websocket"
	DisableLogForwarding                  = "disable-log-forwarding"
	DisableLogForwardingDataSourceSync    = "disable-log-forwarding-data-source-sync"
	DisableLogForwardingDataAccessSync    = "disable-log-forwarding-data-access-sync"
	DisableLogForwardingIdentityStoreSync = "disable-log-forwarding-identity-store-sync"
	DisableLogForwardingDataUsageSync     = "disable-log-forwarding-data-usage-sync"

	ConnectorNameFlag    = "connector-name"
	ConnectorVersionFlag = "connector-version"
	NameFlag             = "name"

	// Import specific flags
	DeleteUntouchedFlag = "delete-untouched"
	DeleteTempFilesFlag = "delete-temp-files"
	ReplaceGroupsFlag   = "replace-groups"

	// Access specific parameters
	AccessFileFlag = "access-file"

	Targets             = "targets"
	DataObjectEnrichers = "data-object-enrichers"
	Repositories        = "repositories"

	GitHubToken = "token"

	IdentitySync   = "IS"
	DataSourceSync = "DS"
	DataAccessSync = "DA"
	DataUsageSync  = "DU"

	SubtaskAccessSync = "AccessSync"

	// HTTP headers
	DomainHeader = "Raito-Domain"
)
