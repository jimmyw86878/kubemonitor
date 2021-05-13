package kubeutil

import (
	"fmt"
	log "kubemonitor/internal/logger"
	"kubemonitor/internal/models"
	"kubemonitor/internal/util"
	"strings"
	"time"
)

//Store is to store pod status
type Store struct {
	List      []*models.Pod
	CachePath string
	Config    *models.Config
}

//InitStore function is to generate store
func InitStore() *Store {
	config, err := models.ReadConfig(util.LoadStrEnv("configPath", "config.json"))
	if err != nil {
		log.Error.Println(err)
		return nil
	}
	store := &Store{
		List:   make([]*models.Pod, 0),
		Config: config,
	}
	//make sure all target pod does not restart recently (checking start time of pod)
	after := time.After(60 * time.Second)
check:
	for {
		select {
		case <-after:
			break check
		default:
			err := store.InitCheckTargetPod()
			if err == nil {
				break check
			}
			time.Sleep(10 * time.Second)
		}
	}
	log.Info.Println("End for checking start time of pod. Start to get all pod status...")
	//break the loop until get all pod status
	for {
		err := store.GetCurrentAllPodStatus()
		if err != nil {
			log.Error.Println(err)
		} else {
			break
		}
		time.Sleep(10 * time.Second)
	}
	return store
}

//GetCurrentAllPodStatus define
func (s *Store) GetCurrentAllPodStatus() error {
	res := make([]*models.Pod, 0)
	for _, ns := range s.Config.WatchList {
		out, err := util.Exec(fmt.Sprintf("kubectl get pods -n %s | grep %s | grep Running", ns.Namespace, ns.MonitorTarget))
		if err != nil {
			return err
		}
		pod := models.TransPodResp(out)
		pod.NameSpace = ns.Namespace
		pod.DeploymentName = ns.MonitorTarget
		pod.RestartList = ns.RestartList
		// err = models.Writeintofile(pod, s.CachePath+fmt.Sprintf("%s#%s.json", ns.Namespace, ns.MonitorTarget))
		// if err != nil {
		// 	return err
		// }
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
				log.Warning.Printf("Deployment %s pod name has changed in %s namespace", dep, ns)
				RestartTargetDeployment(ns, pod.RestartList)
			}
			//update information of pod
			pod.UpdatePodStatus(curpodname, curpodrestart, s.CachePath)
		} else {
			if pod.Restarts != curpodrestart {
				log.Warning.Printf("Deployment %s pod has restarted in %s namespace", dep, ns)
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

//InitCheckTargetPod is to check start time of target monitor pod(vCU). if start time <= 1 min represents that target has restarted
//,and should restart the target list
func (s *Store) InitCheckTargetPod() error {
	var errSum error
	for _, ns := range s.Config.WatchList {
		if ns.Checked {
			//represent that this target pod already checked, skip it
			log.Info.Printf("Deployment %s pod in %s namespace checked and skip it", ns.MonitorTarget, ns.Namespace)
			continue
		}
		out, err := util.Exec(fmt.Sprintf("kubectl get pods -n %s | grep %s | grep Running", ns.Namespace, ns.MonitorTarget))
		if err != nil {
			log.Error.Printf(err.Error())
			errSum = err
			continue
		}
		pod := models.TransPodResp(out)
		state, err := util.Exec(fmt.Sprintf("kubectl describe pod %s -n %s | grep 'State:          Running' -A1 | grep -v 'State:          Running'", pod.Name, ns.Namespace))
		if err != nil {
			log.Error.Printf(err.Error())
			errSum = err
			continue
		}
		if len(strings.Split(state, ", ")) == 2 {
			starttime := strings.Split(state, ", ")[1]
			//true if the target has restarted in past 60 seconds
			if util.CompareCurrentTime(strings.Replace(starttime, "+0800\n", "CST", 1), 60.0) {
				log.Warning.Printf("Deployment %s pod in %s namespace has restarted before in past 60 seconds", ns.MonitorTarget, ns.Namespace)
				RestartTargetDeployment(ns.Namespace, ns.RestartList)
			}
		}
		ns.Checked = true
	}
	return errSum
}
