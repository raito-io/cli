package sync_to_target

import (
	"github.com/raito-io/bexpression"

	"github.com/raito-io/cli/base/access_provider/types"
	"github.com/raito-io/cli/base/data_source"
)

type AccessProviderImport struct {
	LastCalculated  int64             `yaml:"lastCalculated" json:"lastCalculated"`
	AccessProviders []*AccessProvider `yaml:"accessProviders" json:"accessProviders"`
}

type AccessProvider struct {
	Id          string  `yaml:"id" json:"id"`
	Name        string  `yaml:"name" json:"name"`
	Description string  `yaml:"description" json:"description"`
	NamingHint  string  `yaml:"namingHint" json:"namingHint"`
	Type        *string `yaml:"type" json:"type"`

	ExternalId *string `yaml:"externalId" json:"externalId"`

	Action            types.Action `yaml:"action" json:"action"`
	Who               WhoItem      `yaml:"who" json:"who"`
	DeletedWho        *WhoItem     `yaml:"deletedWho" json:"deletedWho"`
	Delete            bool         `yaml:"delete" json:"delete"`
	WhoLocked         *bool        `yaml:"whoLocked" json:"whoLocked"`
	InheritanceLocked *bool        `yaml:"inheritanceLocked" json:"inheritanceLocked"`
	WhatLocked        *bool        `yaml:"whatLocked" json:"whatLocked"`
	DeleteLocked      *bool        `yaml:"deleteLocked" json:"deleteLocked"`

	// Row level filter properties
	PolicyRule     *string                               `yaml:"policyRule,omitempty" json:"policyRule,omitempty"`
	FilterCriteria *bexpression.DataComparisonExpression `yaml:"filterCriteria,omitempty" json:"filterCriteria,omitempty"`

	ActualName *string    `yaml:"actualName" json:"actualName"`
	What       []WhatItem `yaml:"what" json:"what"`
	DeleteWhat []WhatItem `yaml:"deleteWhat" json:"deleteWhat"`
}

type Access struct {
	Id         string     `yaml:"id" json:"id"`
	ActualName *string    `yaml:"actualName" json:"actualName"`
	What       []WhatItem `yaml:"what" json:"what"`
}

type WhoItem struct {
	// Users contains all account names directly assigned to this access provider
	Users []string `yaml:"users,omitempty" json:"users,omitempty"`

	// Groups contains all group names assigned to this access provider
	Groups []string `yaml:"groups,omitempty" json:"groups,omitempty"`

	// InheritFrom contains all access providers actual names in WHO part of this access provider
	InheritFrom []string `yaml:"inheritFrom,omitempty" json:"inheritFrom,omitempty"`

	// Recipients contains all identifiers of data share recipients
	Recipients []string `yaml:"recipients,omitempty" json:"recipients,omitempty"`
}

type WhatItem struct {
	DataObject  *data_source.DataObjectReference `yaml:"dataObject" json:"dataObject"`
	Permissions []string                         `yaml:"permissions" json:"permissions"`
}

type AccessProviderSyncFeedback struct {
	AccessProvider string   `yaml:"accessProvider" json:"accessProvider"`
	ActualName     string   `yaml:"actualName" json:"actualName"`
	ExternalId     *string  `yaml:"externalId" json:"externalId"`
	Type           *string  `yaml:"type" json:"type"`
	Errors         []string `yaml:"errors" json:"errors"`
	Warnings       []string `yaml:"warnings" json:"warnings"`
}

// The legacy format that the appserver still supports. The CLI will convert the new format to the old for now until appserver supports the new format.
type accessSyncFeedbackInformation struct {
	AccessId   string   `yaml:"accessId" json:"accessId"`
	ActualName string   `yaml:"actualName" json:"actualName"`
	ExternalId *string  `yaml:"externalId" json:"externalId"`
	Type       *string  `yaml:"type" json:"type"`
	Errors     []string `yaml:"errors" json:"errors"`
	Warnings   []string `yaml:"warnings" json:"warnings"`
}

type accessProviderFeedbackInformation struct {
	ExternalId            string                          `json:"externalId"`
	AccessFeedbackObjects []accessSyncFeedbackInformation `json:"access"`
}
