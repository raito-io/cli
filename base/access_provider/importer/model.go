package importer

import (
	"encoding/json"
	"strings"

	"github.com/raito-io/cli/base/data_source"
)

type AccessProviderImport struct {
	LastCalculated  int64            `yaml:"lastCalculated" json:"lastCalculated"`
	AccessProviders []AccessProvider `yaml:"accessProviders" json:"accessProviders"`
}

type AccessProvider struct {
	Id          string    `yaml:"id" json:"id"`
	Name        string    `yaml:"name" json:"name"`
	Description string    `yaml:"description" json:"description"`
	NamingHint  string    `yaml:"namingHint" json:"namingHint"`
	Access      []*Access `yaml:"access" json:"access"`
	Action      Action    `yaml:"action" json:"action"`
	Delete      bool      `yaml:"delete" json:"delete"`
}

type Access struct {
	Id         string     `yaml:"id" json:"id"`
	NamingHint string     `yaml:"namingHint" json:"namingHint"`
	Who        WhoItem    `yaml:"who" json:"who"`
	What       []WhatItem `yaml:"what" json:"what"`
}

type WhoItem struct {
	Users          []string `yaml:"users" json:"users"`
	Groups         []string `yaml:"groups" json:"groups"`
	InheritFrom    []string `yaml:"inheritFrom" json:"inheritFrom"`
	UsersInGroups  []string `yaml:"usersInGroups" json:"usersInGroups"`
	UsersInherited []string `yaml:"usersInherited" json:"usersInherited"`
}

type WhatItem struct {
	DataObject  *data_source.DataObjectReference `yaml:"dataObject" json:"dataObject"`
	Permissions []string                         `yaml:"permissions" json:"permissions"`
}

type AccessProviderActualNameTranslation struct {
	AccessProviderId         string `yaml:"accessProviderId" json:"accessProviderId"`
	AccessProviderActualName string `yaml:"accessProviderActualName" json:"accessProviderActualName"`
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
