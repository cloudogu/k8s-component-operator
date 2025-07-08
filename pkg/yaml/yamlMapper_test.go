package yaml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_yamlMapper_PathToYaml(t *testing.T) {

	result1 := map[string]interface{}{
		"kube-prometheus-stack": map[string]interface{}{
			"prometheus": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"env": []interface{}{
							map[string]interface{}{
								"name":  "LOG_LEVEL",
								"value": "debug",
							},
						},
						"name": "auth",
					},
				},
			},
		},
	}

	result2 := map[string]interface{}{
		"manager": map[string]interface{}{
			"env": map[string]interface{}{
				"logLevel": "debug",
			},
		},
	}

	result3 := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name": "web",
							"ports": []interface{}{
								map[string]interface{}{
									"containerPort": "debug",
									"name":          "http",
								},
							},
						},
					},
				},
			},
		},
	}

	examples := map[string]interface{}{
		"kube-prometheus-stack.prometheus.containers[name=auth].env[name=LOG_LEVEL].value": result1,
		"manager.env.logLevel": result2,
		"spec.template.spec.containers[name=web].ports[name=http].containerPort": result3,
	}

	for path, result := range examples {
		yamlMap, err := PathToYAML(path, "debug", NewSerializer())
		assert.NoError(t, err)
		assert.Equal(t, result, yamlMap)
	}
}
