package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//LoadStrEnv load string environment value
func LoadStrEnv(envName string, defaultValue string) string {
	res := defaultValue
	if val, exists := os.LookupEnv(envName); exists {
		res = val
	}
	return res
}

//LoadInt64Env load int64 environment value
func LoadInt64Env(envName string, defaultValue int64) int64 {
	var (
		res = defaultValue
		err error
	)
	if val, exists := os.LookupEnv(envName); exists {
		res, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return defaultValue
		}
	}
	return res
}

//Exec to execute command in container
func Exec(cmd string) (string, error) {
	res, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("Error when executing command %s, err: %s", cmd, err.Error())
	}
	return string(res), nil
}

//CheckCacheExist is to check cache files exist or not, return true if files exist
func CheckCacheExist(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return false
	}
	if len(files) != 0 {
		for _, item := range files {
			if strings.Contains(item.Name(), ".json") {
				return true
			}
		}
	}
	return false
}
