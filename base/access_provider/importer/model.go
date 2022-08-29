package importer

import "github.com/raito-io/cli/base/data_source"

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
	NamingHint string     `yaml:"namingHint" json:"namingHint"`
	Who        WhoItem    `yaml:"who" json:"who"`
	What       []WhatItem `yaml:"what" json:"what"`
}

type WhoItem struct {
	Users                  []string `yaml:"users" json:"users"`
	Groups                 []string `yaml:"groups" json:"groups"`
	AccessProviders        []string `yaml:"accessProviders" json:"accessProviders"`
	UsersInGroups          []string `yaml:"usersInGroups" json:"usersInGroups"`
	UsersInAccessProviders []string `yaml:"usersInAccessProviders" json:"usersInAccessProviders"`
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
