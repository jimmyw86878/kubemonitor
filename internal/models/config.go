package models

import (
	"encoding/json"
	"io/ioutil"
)

//var LoadedTarget []*TargetDeployList

//Config define
type Config struct {
	WatchList []*TargetDeployList `json:"watch_target"`
}

//TargetDeployList define
type TargetDeployList struct {
	Namespace     string   `json:"namespace"`
	MonitorTarget string   `json:"monitor_target"`
	RestartList   []string `json:"restart_list"`
}

//ReadConfig read config json to get target
func ReadConfig(path string) (*Config, error) {
	result := &Config{
		WatchList: make([]*TargetDeployList, 0),
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
