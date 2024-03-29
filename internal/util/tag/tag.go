package tag

import (
	"context"
	"fmt"
	"strings"

	"github.com/raito-io/cli/internal/plugin"
)

func FetchTagSourceFromPlugin(ctx context.Context, client plugin.PluginClient, tagSourcesScope []string) ([]string, error) {
	infoClient, err := client.GetInfo()
	if err != nil {
		return tagSourcesScope, fmt.Errorf("fetching info interface from plugin: %w", err)
	}

	pluginInfo, err := infoClient.GetInfo(ctx)
	if err != nil {
		return tagSourcesScope, fmt.Errorf("calling info from plugin: %w", err)
	}

	if pluginInfo.TagSource != "" {
		tagSourcesScope = append(tagSourcesScope, pluginInfo.TagSource)
	}

	return tagSourcesScope, nil
}

func SerializeTagList(tagList []string) string {
	sb := strings.Builder{}
	sb.WriteString("[")

	for i, tag := range tagList {
		sb.WriteString("\"")
		sb.WriteString(tag)
		sb.WriteString("\"")

		if i < len(tagList)-1 {
			sb.WriteString(",")
		}
	}

	sb.WriteString("]")

	return sb.String()
}
