package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	//ConfigPath for test
	ConfigPath = ""
	//NoCachePath for test
	NoCachePath = ""
	//CachePath for test
	CachePath = ""
)

//InitTest for testing
func InitTest() {
	currentPath := GetPwd()
	ConfigPath = filepath.Join(currentPath, "/kubemonitor/test/config.json")
	NoCachePath = filepath.Join(currentPath, "/kubemonitor/test/test_no_cache") + "/"
	CachePath = filepath.Join(currentPath, "/kubemonitor/test/test_cache/") + "/"
}

//GetPwd is to get absolute value of path
func GetPwd() string {
	res, _ := os.Getwd()
	parentDir := res[:strings.LastIndex(res, "/kubemonitor")]
	return parentDir
}

//Teardown clean test
func Teardown() {
	files, _ := ioutil.ReadDir(NoCachePath)
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			os.Remove(NoCachePath + f.Name())
			fmt.Println("Deleted ", f.Name())
		}
	}
}

//RestartMinikube is to restart cluster to make pod restart and `restarts` number will be 1
func RestartMinikube() {
	fmt.Println("Stopping minikube...")
	_, err := exec.Command("/bin/sh", "-c", "minikube stop").Output()
	if err != nil {
		fmt.Printf("Error when stop minikube, err: %s", err.Error())
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Starting minikube...")
	_, err = exec.Command("/bin/sh", "-c", "minikube start").Output()
	if err != nil {
		fmt.Printf("Error when start minikube, err: %s", err.Error())
	}
}
