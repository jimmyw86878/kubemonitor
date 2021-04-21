package models

import (
	"kubemonitor/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	examplePod = `kubemonitor-deployment-d44655b-bdcrh   1/1     Running   0          3h`
)

func TestTransPodResp(t *testing.T) {
	res := TransPodResp(examplePod)
	assert.Equal(t, "kubemonitor-deployment-d44655b-bdcrh", res.Name)
	assert.Equal(t, "0", res.Restarts)
}

func TestReadfromfile(t *testing.T) {
	test.InitTest()
	res, err := Readfromfile(test.NoCachePath)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(res))
	res, err = Readfromfile(test.CachePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
}

func TestWriteintofile(t *testing.T) {
	test.InitTest()
	example := &Pod{
		Name:           "test-dep1",
		NameSpace:      "ns1",
		RestartList:    []string{"test1", "test2"},
		Restarts:       "1",
		DeploymentName: "dep1",
	}
	err := Writeintofile(example, test.NoCachePath+"test.json")
	assert.NoError(t, err)
	test.Teardown()
}

func TestUpdatePodStatus(t *testing.T) {
	test.InitTest()
	example := &Pod{
		Name:           "test-dep1",
		NameSpace:      "ns1",
		RestartList:    []string{"test1", "test2"},
		Restarts:       "1",
		DeploymentName: "dep1",
	}
	err := example.UpdatePodStatus("update_name", "1", test.NoCachePath)
	assert.NoError(t, err)
	test.Teardown()
}
