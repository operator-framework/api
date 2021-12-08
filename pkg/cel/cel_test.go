package cel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelLibrary(t *testing.T) {
	props := make([]map[string]interface{}, 1)
	props[0] = map[string]interface{}{
		"type":  "olm.test",
		"value": "1.0.0",
	}

	propertiesMap := map[string]interface{}{"properties": props}

	tests := []struct {
		name  string
		rule  string
		out   bool
		isErr bool
	}{
		{
			name:  "ValidCelExpression/True",
			rule:  "properties.exists(p, p.type == 'olm.test' && (semver_compare(p.value, '1.0.0') == 0))",
			out:   true,
			isErr: false,
		},
		{
			name:  "ValidCelExpression/False",
			rule:  "properties.exists(p, p.type == 'olm.test' && (semver_compare(p.value, '1.0.1') == 0))",
			out:   false,
			isErr: false,
		},
		{
			name:  "InvalidCelExpression",
			rule:  "properties.exists(p, p.type == 'olm.test' && (doesnt_exist(p.value, '1.0.0') == 0))",
			isErr: true,
		},
	}

	validator := NewCelEvaluatorProvider()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := validator.Evaluator(tt.rule)
			if tt.isErr {
				assert.Error(t, err)
			} else {
				result, err := prog.Evaluate(propertiesMap)
				assert.NoError(t, err)
				assert.Equal(t, result, tt.out)
			}
		})
	}
}
