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

	err := buildNodeFromPath(root, strings.Split(path, "."), value)
	if err != nil {
		return nil, err
	}

	var out strings.Builder
	enc := yaml.NewEncoder(&out)
	enc.SetIndent(2)
	if err := enc.Encode(root); err != nil {
		return nil, err
	}

	var result map[string]any

	err = serializer.Unmarshal([]byte(out.String()), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildNodeFromPath(parent *yaml.Node, parts []string, value *yaml.Node) error {
	current := parent

	for i := 0; i < len(parts); i++ {
		part := parts[i]
		isLast := i == len(parts)-1

		// Liste erkennen z.B. containers[name=auth]
		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			key, matchKey, matchVal := parseListSegment(part)

			listNode := getOrCreateMapEntry(current, key, &yaml.Node{
				Kind: yaml.SequenceNode,
			})

			// Suche oder erstelle Element mit name=...
			var item *yaml.Node
			for _, el := range listNode.Content {
				if getMapValue(el, matchKey) == matchVal {
					item = el
					break
				}
			}
			if item == nil {
				item = &yaml.Node{Kind: yaml.MappingNode}
				// name als erster Eintrag!
				item.Content = append(item.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: matchKey},
					&yaml.Node{Kind: yaml.ScalarNode, Value: matchVal},
				)
				listNode.Content = append(listNode.Content, item)
			}

			current = item
		} else {
			// Einfaches Mapping
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

	return nil
}

func parseListSegment(seg string) (key, matchKey, matchVal string) {
	start := strings.Index(seg, "[")
	end := strings.Index(seg, "]")
	if start == -1 || end == -1 {
		return
	}
	key = seg[:start]
	kv := strings.SplitN(seg[start+1:end], "=", 2)
	if len(kv) != 2 {
		return
	}
	matchKey, matchVal = kv[0], kv[1]
	return
}

func getOrCreateMapEntry(mapNode *yaml.Node, key string, defaultValue *yaml.Node) *yaml.Node {
	for i := 0; i < len(mapNode.Content); i += 2 {
		k := mapNode.Content[i]
		v := mapNode.Content[i+1]
		if k.Value == key {
			return v
		}
	}
	mapNode.Content = append(mapNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: key},
		defaultValue,
	)
	return defaultValue
}

func getMapValue(mapNode *yaml.Node, key string) string {
	for i := 0; i < len(mapNode.Content); i += 2 {
		if mapNode.Content[i].Value == key {
			return mapNode.Content[i+1].Value
		}
	}
	return ""
}
