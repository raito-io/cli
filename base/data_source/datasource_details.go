package data_source

import "strings"

// To handle availablePermissions on Data Source level correctly, we need to import a DataObject of type datasource.
// DataSourceDetails is a helper to abstract this from the plugin and just expose some setters. The import file creator then automatically prepends the import file with this DataObject
type DataSourceDetails struct {
	dataSource DataObject
}

func newDataSourceDetails(dataSourceId string) DataSourceDetails {
	return DataSourceDetails{
		dataSource: DataObject{
			ExternalId:           dataSourceId,
			Name:                 dataSourceId,
			FullName:             dataSourceId,
			Type:                 "datasource",
			Description:          "",
			ParentExternalId:     "",
			Tags:                 nil,
			AvailablePermissions: []string{},
		},
	}
}

func (d *DataSourceDetails) SetName(name string) {
	d.dataSource.Name = name
}

func (d *DataSourceDetails) SetFullname(name string) {
	d.dataSource.FullName = name
}

func (d *DataSourceDetails) SetDescription(desc string) {
	d.dataSource.Description = desc
}

func (d *DataSourceDetails) AddAvailablePermission(permission string) {
	found := false

	for _, p := range d.dataSource.AvailablePermissions {
		if strings.EqualFold(p, permission) {
			found = true
			break
		}
	}

	if !found {
		d.dataSource.AvailablePermissions = append(d.dataSource.AvailablePermissions, permission)
	}
}

func (d *DataSourceDetails) SetAvailablePermissions(permissions []string) {
	for _, p := range permissions {
		d.AddAvailablePermission(p)
	}
}
