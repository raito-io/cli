package clitrigger

type ApUpdate struct {
	Domain          string   `json:"domain"`
	DataSourceNames []string `json:"dataSourceNames"`
}

type TriggerEvent struct {
	ApUpdate *ApUpdate `json:"apUpdate,omitempty"`
}
