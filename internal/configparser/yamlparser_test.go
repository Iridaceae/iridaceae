package configparser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testYamlParser = &YamlParser{}

func TestYamlParser_Marshal(t *testing.T) {
	assert.Nil(t, testYamlParser.Marshal(&bytes.Buffer{}, GetSettings()))
}
