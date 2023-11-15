package clitrigger

type ApUpdate struct {
	Domain          string   `json:"domain"`
	DataSourceNames []string `json:"dataSourceNames"`
}

type SyncTrigger struct {
	// Domain represents the domain of the customer
	Domain string `json:"domain"`
	// The id of the datasource to sync. Note: this should be set if data source, usage or access needs to be synced
	DataSource *string `json:"dataSource"`
	// The id of the identity store to sync. Note: this should be set if identity store needs to be synced
	IdentityStore *string `json:"identityStore"`
	// Boolean to indicate if the identity store needs to be synced or not
	IdentityStoreSync bool `json:"identityStoreSync"`
	// Boolean to indicate if the data source needs to be synced or not
	DataSourceSync bool `json:"dataSourceSync"`
	// Boolean to indicate if access needs to be synced or not
	DataAccessSync bool `json:"accessSync"`
	// Boolean to indicate if usage needs to be synced or not
	DataUsageSync bool `json:"usageSync"`
	// Optional: the fullName of the data object to sync.
	// That means that, if this is specified, the import will be run with `DeleteUntouched=false`, so no cleanup will be done of removed data objects.
	DataObjectParent *string `json:"dataObjectParent"`
	// Optional: the list of data object names (as child of the DataObjectParent) which can be excluded during the sync
	// When DataObjectParent is provided, it is expected that this list contains all the child data objects that Raito already know about and so can be skipped during the sync.
	DataObjectExcludes []string `json:"dataObjectExcludes"`
}

type TriggerEvent struct {
	ApUpdate    *ApUpdate    `json:"apUpdate,omitempty"`
	SyncTrigger *SyncTrigger `json:"syncTrigger,omitempty"`
}
