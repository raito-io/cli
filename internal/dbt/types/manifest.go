package types

type Manifest struct {
	Metadata Metadata        `json:"metadata"`
	Nodes    map[string]Node `json:"nodes"`
}

type Metadata struct {
	DbtSchemaVersion string `json:"dbt_schema_version"`
	DbtVersion       string `json:"dbt_version"`
	ProjectName      string `json:"project_name"`
	ProjectId        string `json:"project_id"`
}

type Node struct {
	Database         string                 `json:"database"`
	Schema           string                 `json:"schema"`
	Name             string                 `json:"name"`
	ResourceType     string                 `json:"resource_type"`
	PackageName      string                 `json:"package_name"`
	Path             string                 `json:"path"`
	OriginalFilePath string                 `json:"original_file_path"`
	UniqueId         string                 `json:"unique_id"`
	Fqn              []string               `json:"fqn"`
	Alias            string                 `json:"alias"`
	Config           NodeConfig             `json:"config"`
	Description      string                 `json:"description"`
	Columns          map[string]Column      `json:"columns"`
	Meta             Meta                   `json:"meta"`
	PatchPath        string                 `json:"patch_path"`
	BuildPath        string                 `json:"build_path"`
	Deferred         bool                   `json:"deferred"`
	UnrenderedConfig map[string]interface{} `json:"unrendered_config"`
	RelationName     string                 `json:"relation_name"`
	Language         string                 `json:"language"`
	DependsOn        NodeDependsOn          `json:"depends_on"`
	Access           string                 `json:"access"`
}

type NodeConfig struct {
	Enabled      bool                   `json:"enabled"`
	Alias        *string                `json:"alias"`
	Schema       *string                `json:"schema"`
	Database     *string                `json:"database"`
	Meta         map[string]interface{} `json:"meta"`
	Group        *string                `json:"group"`
	Materialized *string                `json:"materialized"`
}

type Column struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Meta        Meta    `json:"meta"`
	DataType    *string `json:"data_type"`
}

type NodeDependsOn struct {
	Macros []string `json:"macros"`
	Nodes  []string `json:"nodes"`
}

type Meta struct {
	Raito RaitoMeta `json:"raito"`
}

type RaitoMeta struct {
	Grant  []Grant  `json:"grant,omitempty"`
	Filter []Filter `json:"filter,omitempty"`
	Mask   *Mask    `json:"mask,omitempty"`
}

type Grant struct {
	Name              string   `json:"name"`
	Permissions       []string `json:"permissions"`
	GlobalPermissions []string `json:"global_permissions"`
}

type Filter struct {
	Name       string `json:"name"`
	PolicyRule string `json:"policy_rule"`
}

type Mask struct {
	Name string  `json:"name"`
	Type *string `json:"type,omitempty"`
}
