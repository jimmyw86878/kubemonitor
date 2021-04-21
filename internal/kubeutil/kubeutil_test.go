package kubeutil

import (
	"kubemonitor/internal/models"
	"kubemonitor/test"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentAllPodStatus(t *testing.T) {
	test.InitTest()
	config, err := models.ReadConfig(test.ConfigPath)
	assert.NoError(t, err)
	store := &Store{
		CachePath: test.NoCachePath,
	}
	err = store.GetCurrentAllPodStatus(config.WatchList)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(store.List))
	for _, pod := range store.List {
		assert.NotEqual(t, "", pod.Name)
	}
	test.Teardown()
}

func TestCheckAndUpdatePodStat(t *testing.T) {
	test.InitTest()
	config, err := models.ReadConfig(test.ConfigPath)
	assert.NoError(t, err)
	store := &Store{
		CachePath: test.NoCachePath,
	}
	err = store.GetCurrentAllPodStatus(config.WatchList)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(store.List))
	//simulate some pod restart by changing their pod name
	RestartTargetDeployment("default", []string{"dep-a"})
	after := time.After(30 * time.Second)
	for {
		select {
		case <-after:
			assert.FailNow(t, "There are errors in CheckAnsUpdate function")
			test.Teardown()
			return
		default:
			time.Sleep(1 * time.Second)
			//check pod status, it should be detected and restart dep-b deployment
			errArr := store.CheckAndUpdatePodStat()
			if len(errArr) == 0 {
				test.Teardown()
				return
			}
		}
	}
}

func TestInitStore(t *testing.T) {
	test.InitTest()
	os.Setenv("configPath", test.ConfigPath)
	os.Setenv("cacheFilePath", test.CachePath)
	//start with cache file
	store := InitStore()
	assert.Equal(t, 2, len(store.List))
	os.Setenv("cacheFilePath", test.NoCachePath)
	//start without cache file and read from current kubernetes cluster
	store = InitStore()
	assert.Equal(t, 2, len(store.List))
	test.Teardown()
}

func TestCheckAndUpdatePodRestart(t *testing.T) {
	test.InitTest()
	config, err := models.ReadConfig(test.ConfigPath)
	assert.NoError(t, err)
	store := &Store{
		CachePath: test.NoCachePath,
	}
	err = store.GetCurrentAllPodStatus(config.WatchList)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(store.List))
	//simulate some pod restart number increases
	test.RestartMinikube()
	after := time.After(50 * time.Second)
	for {
		select {
		case <-after:
			assert.FailNow(t, "There are errors in CheckAnsUpdate function")
			test.Teardown()
			return
		default:
			time.Sleep(1 * time.Second)
			//check pod status, it should be detected and restart dep-b deployment
			errArr := store.CheckAndUpdatePodStat()
			if len(errArr) == 0 {
				test.Teardown()
				return
			}
		}
	}
}
