package models

import (
	"kubemonitor/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	test.InitTest()
	res, err := ReadConfig(test.ConfigPath)
	assert.NoError(t, err)
	assert.NotEqual(t, nil, res)
}
