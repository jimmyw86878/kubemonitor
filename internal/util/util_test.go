package util

import (
	"kubemonitor/test"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvStr(t *testing.T) {
	notExistStr := LoadStrEnv("notexistenv", "notexist")
	assert.Equal(t, "notexist", notExistStr)
	os.Setenv("FOO", "foo")
	existStr := LoadStrEnv("FOO", "notexist")
	assert.Equal(t, "foo", existStr)
}

func TestLoadEnvInt64(t *testing.T) {
	notExistInt := LoadInt64Env("notexistenv", 30)
	assert.Equal(t, int64(30), notExistInt)
	os.Setenv("testInt64", "10")
	existInt := LoadInt64Env("testInt64", 10)
	assert.Equal(t, int64(10), existInt)
}

func TestExec(t *testing.T) {
	res, err := Exec("whoami")
	assert.NoError(t, err)
	assert.NotEqual(t, "", res)
	_, err = Exec("wrong command")
	assert.Error(t, err)
}

func TestCheckCacheExist(t *testing.T) {
	test.InitTest()
	res := CheckCacheExist(test.CachePath)
	assert.Equal(t, true, res)
	res = CheckCacheExist("../../dev/not_exist")
	assert.Equal(t, false, res)
}

func TestParseTime(t *testing.T) {
	res := CompareCurrentTime("26 Apr 2021 19:56:15 CST", 10.0)
	assert.Equal(t, false, res)
	res = CompareCurrentTime("26 Apr 2021 19:56:15 CST", math.MaxFloat64)
	assert.Equal(t, true, res)
}
