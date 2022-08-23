package data_source

// DataSourceDetails is a helper to abstract this from the plugin and just expose some setters. The import file creator then automatically prepends the import file with this DataObject
type DataSourceDetails struct {
	dataSource DataObject
}

func newDataSourceDetails(dataSourceId string) DataSourceDetails {
	return DataSourceDetails{
		dataSource: DataObject{
			ExternalId:       dataSourceId,
			Name:             dataSourceId,
			FullName:         dataSourceId,
			Type:             "datasource",
			Description:      "",
			ParentExternalId: "",
			Tags:             nil,
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
