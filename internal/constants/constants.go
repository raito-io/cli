package constants

var KnownFlags = map[string]struct{}{
	DebugFlag:                   {},
	LogFileFlag:                 {},
	SkipAuthentication:          {},
	SkipFileUpload:              {},
	URLOverrideFlag:             {},
	DomainFlag:                  {},
	ApiUserFlag:                 {},
	ApiSecretFlag:               {},
	ConfigFileFlag:              {},
	FrequencyFlag:               {},
	CronFlag:                    {},
	SyncAtStartupFlag:           {},
	SkipIdentityStoreSyncFlag:   {},
	SkipDataSourceSyncFlag:      {},
	SkipDataAccessSyncFlag:      {},
	SkipDataUsageSyncFlag:       {},
	DataSourceIdFlag:            {},
	IdentityStoreIdFlag:         {},
	OnlyTargetsFlag:             {},
	ConnectorNameFlag:           {},
	ConnectorVersionFlag:        {},
	NameFlag:                    {},
	DeleteUntouchedFlag:         {},
	DeleteTempFilesFlag:         {},
	ReplaceGroupsFlag:           {},
	FileBackupLocationFlag:      {},
	MaximumBackupsPerTargetFlag: {},
}

const (
	DebugFlag                                = "debug"
	URLOverrideFlag                          = "raito-url-override"
	SkipAuthentication                       = "skip-authentication"
	SkipFileUpload                           = "skip-file-upload"
	LogFileFlag                              = "log-file"
	LogOutputFlag                            = "log-output"
	DomainFlag                               = "domain"
	ApiUserFlag                              = "api-user"
	ApiSecretFlag                            = "api-secret"
	ConfigFileFlag                           = "config-file"
	FrequencyFlag                            = "frequency"
	CronFlag                                 = "cron"
	SyncAtStartupFlag                        = "sync-at-startup"
	SkipDataSourceSyncFlag                   = "skip-data-source-sync"
	SkipDataAccessSyncFlag                   = "skip-data-access-sync"
	SkipIdentityStoreSyncFlag                = "skip-identity-store-sync"
	SkipDataUsageSyncFlag                    = "skip-data-usage-sync"
	SkipResourceProviderFlag                 = "skip-resource-provider-sync"
	SkipTagFlag                              = "skip-tag-sync"
	DataSourceIdFlag                         = "data-source-id"
	IdentityStoreIdFlag                      = "identity-store-id"
	OnlyTargetsFlag                          = "only-targets"
	DisableWebsocketFlag                     = "disable-websocket"
	DisableLogForwarding                     = "disable-log-forwarding"
	DisableLogForwardingDataSourceSync       = "disable-log-forwarding-data-source-sync"
	DisableLogForwardingDataAccessSync       = "disable-log-forwarding-data-access-sync"
	DisableLogForwardingIdentityStoreSync    = "disable-log-forwarding-identity-store-sync"
	DisableLogForwardingDataUsageSync        = "disable-log-forwarding-data-usage-sync"
	DisableLogForwardingResourceProviderSync = "disable-log-forwarding-resource-provider-sync"
	DisableLogForwardingTagSync              = "disable-log-forwarding-tag-sync"

	// Locking parameters
	LockAllWhoFlag            = "lock-all-who"
	LockWhoByNameFlag         = "lock-who-by-name"
	LockWhoByTagFlag          = "lock-who-by-tag"
	LockWhoWhenIncompleteFlag = "lock-who-when-incomplete"

	LockAllInheritanceFlag            = "lock-all-inheritance"
	LockInheritanceByNameFlag         = "lock-inheritance-by-name"
	LockInheritanceByTagFlag          = "lock-inheritance-by-tag"
	LockInheritanceWhenIncompleteFlag = "lock-inheritance-when-incomplete"

	LockAllWhatFlag            = "lock-all-what"
	LockWhatByNameFlag         = "lock-what-by-name"
	LockWhatByTagFlag          = "lock-what-by-tag"
	LockWhatWhenIncompleteFlag = "lock-what-when-incomplete"

	LockAllNamesFlag            = "lock-all-names"
	LockNamesByNameFlag         = "lock-names-by-name"
	LockNamesByTagFlag          = "lock-names-by-tag"
	LockNamesWhenIncompleteFlag = "lock-names-when-incomplete"

	LockAllDeleteFlag            = "lock-all-delete"
	LockDeleteByNameFlag         = "lock-delete-by-name"
	LockDeleteByTagFlag          = "lock-delete-by-tag"
	LockDeleteWhenIncompleteFlag = "lock-delete-when-incomplete"

	LockAllOwnersFlag = "lock-all-owners"

	// MakeNotInternalizableFlag is deprecated and replaced by FullyLockByNameFlag
	MakeNotInternalizableFlag   = "make-not-internalizable"
	FullyLockAllFlag            = "fully-lock-all"
	FullyLockByNameFlag         = "fully-lock-by-name"
	FullyLockByTagFlag          = "fully-lock-by-tag"
	FullyLockWhenIncompleteFlag = "fully-lock-when-incomplete"

	// File handling config tags
	FileBackupLocationFlag      = "file-backup-location"
	MaximumBackupsPerTargetFlag = "maximum-backups-per-target"
	DeleteTempFilesFlag         = "delete-temp-files"
	MaximumFileSizesFlag        = "maximum-file-size"

	TagOverwriteKeyForAccessProviderName   = "tag-overwrite-key-for-access-provider-name"
	TagOverwriteKeyForAccessProviderOwners = "tag-overwrite-key-for-access-provider-owners"
	TagOverwriteKeyForDataObjectOwners     = "tag-overwrite-key-for-data-object-owners"

	TagKeyAndValueForUserIsMachine = "tag-key-and-value-for-user-is-machine"

	ConnectorNameFlag    = "connector-name"
	ConnectorVersionFlag = "connector-version"
	NameFlag             = "name"

	// Import specific flags
	DeleteUntouchedFlag = "delete-untouched"
	ReplaceGroupsFlag   = "replace-groups"

	// For the apply-access command
	FilterAccessFlag = "filter-access"

	Targets             = "targets"
	DataObjectEnrichers = "data-object-enrichers"
	Repositories        = "repositories"

	GitHubToken = "token"

	IdentitySync         = "IS"
	DataSourceSync       = "DS"
	DataAccessSync       = "DA"
	DataUsageSync        = "DU"
	ResourceProviderSync = "RP"
	TagSync              = "TAG"

	SubtaskAccessSync = "AccessSync"

	// HTTP headers
	DomainHeader = "Raito-Domain"

	// Docker flags
	ContainerLivenessFile = "cli-container-liveness-file"
)
