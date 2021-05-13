package models

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

//Pod define
type Pod struct {
	Name           string   `json:"name"`
	Restarts       string   `json:"restarts"`
	RestartList    []string `json:"restart_list"`
	NameSpace      string   `json:"namespace"`
	DeploymentName string   `json:"deployment_name"`
}

//TransPodResp is to transform pod fields from output of command into pod struct
func TransPodResp(input string) *Pod {
	result := &Pod{}
	detail := strings.Fields(input)
	result.Name = detail[0]
	result.Restarts = detail[3]
	return result
}

//Readfromfile is to read pod status from files if files exist
func Readfromfile(path string) ([]*Pod, error) {
	result := make([]*Pod, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		tempPod := &Pod{}
		if strings.Contains(f.Name(), ".json") {
			data, err := ioutil.ReadFile(path + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &tempPod)
			if err != nil {
				return nil, err
			}
			result = append(result, tempPod)
		}
	}
	return result, nil

}

//Writeintofile write pod status into file
func Writeintofile(pod *Pod, dest string) error {
	file, err := json.MarshalIndent(pod, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dest, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

//UpdatePodStatus is to update pod status
func (p *Pod) UpdatePodStatus(name, restarts, filepath string) error {
	p.Name = name
	p.Restarts = restarts
	// err := Writeintofile(p, filepath+fmt.Sprintf("%s#%s.json", p.NameSpace, p.DeploymentName))
	// if err != nil {
	// 	log.Error.Printf("Can not write file for %s, err: %s\n", p.NameSpace+"#"+p.DeploymentName, err.Error())
	// 	return err
	// }
	return nil
}
