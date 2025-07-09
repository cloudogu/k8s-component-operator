package yaml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_yamlMapper_PathToYaml(t *testing.T) {

	result := map[string]interface{}{
		"manager": map[string]interface{}{
			"env": map[string]interface{}{
				"logLevel": "debug",
			},
		},
	}

	examples := map[string]interface{}{
		"manager.env.logLevel": result,
	}

	for path, result := range examples {
		yamlMap, err := PathToYAML(path, "debug", NewSerializer())
		assert.NoError(t, err)
		assert.Equal(t, result, yamlMap)
	}
}
