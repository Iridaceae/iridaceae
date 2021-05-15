package configparser

import (
	"bytes"
	"testing"

	helpers2 "github.com/Iridaceae/iridaceae/internal/helpers"

	"github.com/stretchr/testify/assert"
)

var testYamlParser = &YamlParser{}

func TestYamlParser_Marshal(t *testing.T) {
	assert.Nil(t, testYamlParser.Marshal(&bytes.Buffer{}, helpers2.GetSettings()))
}
