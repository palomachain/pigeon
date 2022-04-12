package conductor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonSerialisationIsDeterministic(t *testing.T) {
	var v1 struct {
		A string
		B string
	}
	var v2 struct {
		B string
		A string
	}
	v3 := map[string]string{
		"A": "A",
		"B": "B",
	}

	v1.A, v1.B = "A", "B"
	v2.A, v2.B = "A", "B"
	v3["A"], v3["B"] = "A", "B"

	json1, err := jsonDeterministicEncoding(v1)
	assert.NoError(t, err)
	json2, err := jsonDeterministicEncoding(v2)
	assert.NoError(t, err)
	json3, err := jsonDeterministicEncoding(v3)
	assert.NoError(t, err)

	assert.Equal(t, json1, json2)
	assert.Equal(t, json1, json3)
	assert.Equal(t, json2, json3)
}
