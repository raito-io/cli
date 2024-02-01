package data_source

import (
	"errors"

	"github.com/bcicen/jstream"

	"github.com/raito-io/cli/base/tag"
)

// DataObject represents a data object in the format that is suitable to be imported into a Raito data source.
type DataObject struct {
	ExternalId       string       `json:"externalId"`
	Name             string       `json:"name"`
	FullName         string       `json:"fullName"`
	Type             string       `json:"type"`
	Description      string       `json:"description"`
	ParentExternalId string       `json:"parentExternalId"`
	Tags             []*tag.Tag   `json:"tags"`
	DataType         *string      `json:"dataType,omitempty"`
	Owners           *OwnersInput `json:"owners"`
}

type OwnersInput struct {
	Users []string `yaml:"users" json:"users"`
}

func CreateDataObjectFromRow(row *jstream.MetaValue) (*DataObject, error) {
	if row.ValueType != jstream.Object {
		return nil, errors.New("illegal format for data object definition in source file")
	}

	var values = row.Value.(map[string]interface{})

	do := DataObject{
		ExternalId:       getStringValue(values, "externalId"),
		Name:             getStringValue(values, "name"),
		FullName:         getStringValue(values, "fullName"),
		Type:             getStringValue(values, "type"),
		Description:      getStringValue(values, "description"),
		ParentExternalId: getStringValue(values, "parentExternalId"),
	}

	if t, found := values["tags"]; found && t != nil {
		if tags, ok := t.([]interface{}); ok {
			do.Tags = make([]*tag.Tag, 0, len(tags))

			for _, tagInput := range tags {
				if tagObj, tok := tagInput.(map[string]interface{}); tok {
					do.Tags = append(do.Tags, &tag.Tag{
						Key:    getStringValue(tagObj, "key"),
						Value:  getStringValue(tagObj, "value"),
						Source: getStringValue(tagObj, "source"),
					})
				}
			}
		}
	}

	if o, found := values["owners"]; found && o != nil {
		if owner, ok := o.(map[string]interface{}); ok {
			do.Owners = &OwnersInput{
				Users: getStringSliceValue(owner, "users"),
			}
		}
	}

	return &do, nil
}

func getStringValue(row map[string]interface{}, key string) string {
	if v, found := row[key]; found {
		if vs, ok := v.(string); ok {
			return vs
		}
	}

	return ""
}
func getStringSliceValue(row map[string]interface{}, key string) []string {
	out := []string{}

	v, f := row[key]

	if !f || v == nil {
		return out
	}

	tmp := v.([]interface{})

	for _, i := range tmp {
		out = append(out, i.(string))
	}

	return out
}
