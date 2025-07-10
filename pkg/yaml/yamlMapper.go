package yaml

import (
	"gopkg.in/yaml.v3"
	"strings"
)

func PathToYAML(path string, val string, serializer Serializer) (map[string]any, error) {
	value := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: val,
		Tag:   "!!str",
	}

	root := &yaml.Node{
		Kind: yaml.MappingNode,
	}

	buildNodeFromPath(root, strings.Split(path, "."), value)

	var out strings.Builder
	enc := yaml.NewEncoder(&out)
	enc.SetIndent(2)
	if err := enc.Encode(root); err != nil {
		return nil, err
	}

	var result map[string]any

	err := serializer.Unmarshal([]byte(out.String()), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildNodeFromPath(parent *yaml.Node, parts []string, value *yaml.Node) {
	current := parent

	for i := 0; i < len(parts); i++ {
		part := parts[i]
		isLast := i == len(parts)-1

		if isLast {
			getOrCreateMapEntry(current, part, value)
		} else {
			next := getOrCreateMapEntry(current, part, &yaml.Node{
				Kind: yaml.MappingNode,
			})
			current = next
		}
	}
}

func getOrCreateMapEntry(mapNode *yaml.Node, key string, defaultValue *yaml.Node) *yaml.Node {
	mapNode.Content = append(mapNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		defaultValue,
	)
	return defaultValue
}
