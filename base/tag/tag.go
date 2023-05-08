package tag

// Tag represents a tag, which can be unused on different entity types like data objects, groups, users etc.
type Tag struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
}
