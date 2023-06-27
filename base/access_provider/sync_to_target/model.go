package sync_to_target

import (
	"encoding/json"
	"strings"

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

	// Deprecated: use ActualName and What fields directory
	Access []*Access `yaml:"access" json:"access"`

	Action     Action   `yaml:"action" json:"action"`
	Who        WhoItem  `yaml:"who" json:"who"`
	DeletedWho *WhoItem `yaml:"deletedWho" json:"deletedWho"`
	Delete     bool     `yaml:"delete" json:"delete"`
	WhoLocked  *bool    `yaml:"whoLocked" json:"whoLocked"`
	WhatLocked *bool    `yaml:"whatLocked" json:"whatLocked"`

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
	// Users contains all account names of users assigned to this access provider
	Users []string `yaml:"users,omitempty" json:"users,omitempty"`

	// Groups contains all account names of groups assigned to this access provider
	Groups []string `yaml:"groups,omitempty" json:"groups,omitempty"`

	// InheritFrom contains all access providers actual names in WHO part of this access provider
	InheritFrom []string `yaml:"inheritFrom,omitempty" json:"inheritFrom,omitempty"`

	// UsersInGroup contains all account names of users assigned to groups of this access provider
	UsersInGroups []string `yaml:"usersInGroups,omitempty" json:"usersInGroups,omitempty"`

	// UsersInherited contains all account names of users that are in an unpacked who part of a child access provider
	UsersInherited []string `yaml:"usersInherited,omitempty" json:"usersInherited,omitempty"`

	// NativeGroupsInherited contains all groups associated with the native IS that are in who part of this or a child access provider
	NativeGroupsInherited []string `yaml:"nativeGroupsInherited,omitempty" json:"nativeGroupsInherited,omitempty"`

	// UsersInheritedNativeGroupsExcluded contains all account names of users that are in who part of this or a child access provider but native groups are excluded
	UsersInheritedNativeGroupsExcluded []string `yaml:"usersInheritedNativeGroupsExcluded,omitempty" json:"usersInheritedNativeGroupsExcludedGroupsExcluded,omitempty"`
}

type WhatItem struct {
	DataObject  *data_source.DataObjectReference `yaml:"dataObject" json:"dataObject"`
	Permissions []string                         `yaml:"permissions" json:"permissions"`
}

type AccessSyncFeedbackInformation struct {
	AccessId   string `yaml:"accessId" json:"accessId"`
	ActualName string `yaml:"actualName" json:"actualName"`
}

type accessProviderFeedbackInformation struct {
	ExternalId            string                          `json:"externalId"`
	AccessFeedbackObjects []AccessSyncFeedbackInformation `json:"access"`
}

type Action int

const (
	Promise Action = iota
	Grant
	Deny
	Mask
	Filtered
)

var actionMap = map[string]Action{
	"promise":  Promise,
	"grant":    Grant,
	"deny":     Deny,
	"mask":     Mask,
	"filtered": Filtered,
}
var actionNames = [...]string{"promise", "grant", "deny", "mask", "filtered"}

func (a *Action) MarshalYAML() (interface{}, error) {
	s := actionNames[*a]

	return s, nil
}

func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	*a = actionMap[strings.ToLower(s)]

	return nil
}

func (a *Action) MarshalJSON() ([]byte, error) {
	s := actionNames[*a]

	return json.Marshal(s)
}

func (a *Action) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*a = actionMap[strings.ToLower(s)]

	return nil
}
