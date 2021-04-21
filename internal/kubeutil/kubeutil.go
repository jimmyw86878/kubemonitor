package kubeutil

import (
	"fmt"
	log "kubemonitor/internal/logger"
	"kubemonitor/internal/models"
	"kubemonitor/internal/util"
	"strings"
)

//Store is to store pod status
type Store struct {
	List      []*models.Pod
	CachePath string
}

//InitStore function is to generate store
func InitStore() *Store {
	config, err := models.ReadConfig(util.LoadStrEnv("configPath", "config.json"))
	if err != nil {
		log.Error.Println(err)
		return nil
	}
	cachePath := util.LoadStrEnv("cacheFilePath", "cache/")
	store := &Store{
		List:      make([]*models.Pod, 0),
		CachePath: cachePath,
	}
	//pod status come from file if there are cache files
	if util.CheckCacheExist(cachePath) {
		store.List, err = models.Readfromfile(cachePath)
		if err != nil {
			log.Error.Println(err)
			return nil
		}
	} else {
		//break the loop until get pod status
		for {
			err := store.GetCurrentAllPodStatus(config.WatchList)
			if err != nil {
				log.Error.Println(err)
			} else {
				break
			}
		}
	}
	return store
}

//GetCurrentAllPodStatus define
func (s *Store) GetCurrentAllPodStatus(input []*models.TargetDeployList) error {
	res := make([]*models.Pod, 0)
	for _, ns := range input {
		out, err := util.Exec(fmt.Sprintf("kubectl get pods -n %s | grep %s | grep Running", ns.Namespace, ns.MonitorTarget))
		if err != nil {
			return err
		}
		pod := models.TransPodResp(out)
		pod.NameSpace = ns.Namespace
		pod.DeploymentName = ns.MonitorTarget
		pod.RestartList = ns.RestartList
		err = models.Writeintofile(pod, s.CachePath+fmt.Sprintf("%s#%s.json", ns.Namespace, ns.MonitorTarget))
		if err != nil {
			return err
		}
		res = append(res, pod)
	}
	s.List = res
	return nil
}

//CheckAndUpdatePodStat is to check pod status and restart other deployment if status change
func (s *Store) CheckAndUpdatePodStat() []error {
	errOutput := make([]error, 0)
	for _, pod := range s.List {
		ns := pod.NameSpace
		dep := pod.DeploymentName
		//check kubernetes is ok and namespace exists
		_, err := util.Exec(fmt.Sprintf("kubectl get pods -n %s", ns))
		if err != nil {
			errOutput = append(errOutput, err)
			log.Error.Printf("kubectl failed or namespace not exists, err: %s\n", err.Error())
			continue
		}
		//check pod status by deployment name and restart count of pod
		res, err := util.Exec(fmt.Sprintf("kubectl get pods -n %s | grep %s | grep Running", ns, dep))
		if err != nil {
			errOutput = append(errOutput, err)
			log.Error.Printf("Can not find pod for deployment %s, err: %s\n", dep, err.Error())
			//check if replicaset of deployment change to 0 or not
			replica, err := util.Exec(fmt.Sprintf("kubectl describe deploy %s -n %s | grep 'Replicas:' | awk '{print $2}'", dep, ns))
			if err != nil {
				log.Error.Printf("Can not find replicaset of deployment %s, err: %s\n", dep, err.Error())
				continue
			}
			if replica == "0\n" {
				log.Warning.Printf("Deployment %s replicaset change to 0\n", dep)
				//update pod name to `notexists` when last pod name is not `notexists`
				if pod.Name != "notexists" {
					pod.UpdatePodStatus("notexists", "0", s.CachePath)
				}
			}
			continue
		}
		curpodname := strings.Fields(res)[0]
		curpodrestart := strings.Fields(res)[3]
		if pod.Name != curpodname {
			//need to restart here when last pod name exist
			if pod.Name != "notexists" {
				log.Warning.Printf("Deployment %s pod name has changed", dep)
				RestartTargetDeployment(ns, pod.RestartList)
			}
			//update information of pod
			pod.UpdatePodStatus(curpodname, curpodrestart, s.CachePath)
		} else {
			if pod.Restarts != curpodrestart {
				log.Warning.Printf("Deployment %s pod has restarted", dep)
				//need to restart here
				RestartTargetDeployment(ns, pod.RestartList)
				//update information of pod
				pod.UpdatePodStatus(curpodname, curpodrestart, s.CachePath)
			}
		}
	}
	return errOutput
}

//RestartTargetDeployment is to restart deployment
func RestartTargetDeployment(ns string, deploymentList []string) {
	for _, deploy := range deploymentList {
		kubecmd := fmt.Sprintf("kubectl rollout restart deploy %s -n %s", deploy, ns)
		_, err := util.Exec(kubecmd)
		if err != nil {
			log.Error.Printf("Can not restart deployment %s in %s namespace, err: %s\n", deploy, ns, err.Error())
		} else {
			log.Info.Printf("Restart deployment %s in %s namespace successfully\n", deploy, ns)
		}
	}
}
