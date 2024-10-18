package sync_from_target

import (
	"github.com/raito-io/cli/base/data_source"
	"github.com/raito-io/cli/base/tag"
)

// AccessProvider describes data access in the format that is suitable to be imported into Raito.x
type AccessProvider struct {
	ExternalId string  `json:"externalId"`
	Name       string  `json:"name"`
	NamingHint string  `json:"namingHint"`
	Type       *string `json:"type"`

	// Deprecated: use ActualName and What fields directory
	Access []*Access `yaml:"access" json:"access"`

	Action Action   `json:"action"`
	Policy string   `json:"policy"`
	Who    *WhoItem `yaml:"who" json:"who"`

	Owners *OwnersInput `json:"owners,omitempty"` // Deprecated, use RaitoOwnerTag instead

	// Locking properties

	// NotInternalizable means that the entire access provider is locked
	NotInternalizable       bool    `json:"notInternalizable"`
	WhoLocked               *bool   `json:"whoLocked"`
	WhoLockedReason         *string `json:"whoLockedReason"`
	InheritanceLocked       *bool   `json:"inheritanceLocked"`
	InheritanceLockedReason *string `json:"inheritanceLockedReason"`
	WhatLocked              *bool   `json:"whatLocked"`
	WhatLockedReason        *string `json:"whatLockedReason"`
	NameLocked              *bool   `json:"nameLocked"`
	NameLockedReason        *string `json:"nameLockedReason"`
	DeleteLocked            *bool   `json:"deleteLocked"`
	DeleteLockedReason      *string `json:"deleteLockedReason"`
	OwnersLocked            *bool   `json:"ownersLocked"`
	OwnersLockedReason      *string `json:"ownersLockedReason"`

	ActualName string `yaml:"actualName" json:"actualName"`
	// Who represents who has access to the 'what'. Nil means that the 'who' is unknown.
	What []WhatItem `yaml:"what" json:"what"`

	// Allows the plugin to indicate that the access provider is incomplete (because not all who items, what items or permissions could be handled)
	Incomplete *bool `json:"incomplete"`

	Tags []*tag.Tag `json:"tags"`
}

type OwnersInput struct {
	Users []string `yaml:"users" json:"users"`
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
