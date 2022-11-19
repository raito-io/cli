package sync_from_target

import "github.com/raito-io/cli/base/data_source"

// AccessProvider describes data access in the format that is suitable to be imported into Raito.x
type AccessProvider struct {
	ExternalId        string    `json:"externalId"`
	NotInternalizable bool      `json:"notInternalizable"`
	Name              string    `json:"name"`
	NamingHint        string    `json:"namingHint"`
	Access            []*Access `yaml:"access" json:"access"`
	Action            Action    `json:"action"`
	Policy            string    `json:"policy"`
	Who               *WhoItem  `yaml:"who" json:"who"`
}

type Access struct {
	ActualName string `yaml:"actualName" json:"actualName"`
	// Who represents who has access to the 'what'. Nil means that the 'who' is unknown.
	What []WhatItem `yaml:"what" json:"what"`
}

type WhoItem struct {
	Users           []string `yaml:"users" json:"users"`
	Groups          []string `yaml:"groups" json:"groups"`
	AccessProviders []string `yaml:"accessProviders" json:"accessProviders"`
}

type WhatItem struct {
	DataObject  *data_source.DataObjectReference `yaml:"dataObject" json:"dataObject"`
	Permissions []string                         `yaml:"permissions" json:"permissions"`
}

type Action int

const (
	Promise Action = iota
	Grant
	Deny
	Mask
	Filtered
)
