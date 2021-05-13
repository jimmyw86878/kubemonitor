## Description

In some scenarios, we may need to synchronize the lifecycle of some pods since they may disconnect with each other when one of them has restarted or re-scheduled. `Kubemonitor` can detect the event when one pod has restarted and it will restart other pods. This purpose aims to eliminate the disconnection among your modules when you consider that `livenessProbe` or `readinessProbe` is not suitable for your modules.

## Prerequisites

* Golang 1.12 up
* Docker
* Make

For running go-test, you need to install:

* Docker
* Minikube ([install guide](https://minikube.sigs.k8s.io/docs/start/))
* Kubectl ([install guide](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/))


## How to build image

* Build

```
make name=your_image_name tag=your_image_tag all
```

## Config

There three things that need to preapre first.

1. /home/kubemonitor
2. /home/kubemonitor/config.json

For 1., just make sure the folder exists on host.

For 2., there is an example `config.json` in `test` folder:
```
{
    "watch_target":[
        {
            "namespace": "default",
            "monitor_target": "dep-a",
            "restart_list": ["dep-b"]
        },
        {
            "namespace": "test",
            "monitor_target": "dep-c",
            "restart_list": ["dep-d"]
        }
    ]
}
```

`kubemonitor` based on this `config.json` to determine what kind of pods(deployment) should be monitored. Feel free to adapt it in your scenario. As you can see, `restart_list` will be an array if you have multiple pods.


## Bring up

We based on OpenShift 3.11 platform to bring `Kubemonitor`.
All needed yaml files already exist in `yamls` folder. Assume we bring kubemonitor in `default` namespace.

1. Create serviceAccount

Default serviceAccount name will be `kubemonitoruser`. Feel free to change it inside `kubemonitor-deployment.yml` and `role.yml`.

```
oc create serviceaccount kubemonitoruser

oc adm policy add-scc-to-user privileged -z kubemonitoruser
```

2. Create cluster role and role binding

```
kubectl create -f role.yml
```

3. Prepare kube config file for your K8S cluster 

The config will be used for connecting K8S cluster. It should be mounted in kubemonitor. In `kubemonitor-deployment.yml`, hostpath is `/home/config` by default. Feel free to change it.

3. Create kubemonitor

Please prepare `kubemonitor` image first. Below section will show that how to build the image.

```
kubectl create -f kubemonitor-deployment.yml
```
