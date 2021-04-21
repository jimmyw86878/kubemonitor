#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

function process_cmd() {
    echo "$CMD"
    $CMD
    if [[ $? -ne 0 ]]; then
        echo "Command failed $CMD"
        exit 1
    fi
}

function usage() {
    cat <<EOF
Usage: $0 COMMAND

Commands:
    start_dev               Start standalone dev environment
    stop_dev                Stop standalone dev environment
    run_test                Run go test
EOF
}

declare -a hosts=()

function start_dev() {
    echo "Start development"
    echo "Start minikube"
    minikube start --kubernetes-version=v1.14.0
    echo "Bring some deployment for testing usage"
    minikube kubectl -- create deployment dep-a --image=k8s.gcr.io/echoserver:1.4
    minikube kubectl -- create deployment dep-b --image=k8s.gcr.io/echoserver:1.4
    echo "Create test namespace"
    minikube kubectl create namespace test
    minikube kubectl -- create deployment dep-c --image=k8s.gcr.io/echoserver:1.4 -n test
    minikube kubectl -- create deployment dep-d --image=k8s.gcr.io/echoserver:1.4 -n test
}

function stop_dev() {
    echo "Stop development"
    echo "Delete minikube cluster"
    minikube delete
}

function run_test() {
    echo "run test"
    go test -p 1 --cover -coverprofile=coverage.out -v kubemonitor/internal{/kubeutil,/models,/util,/route/v1} -count=1
    go tool cover -html=coverage.out
}

OPT=
while [ "$#" -gt 0 ]; do
    case "$1" in

    (start_dev)
        OPT="start_dev"
        shift 1
        ;;
    (stop_dev)
        OPT="stop_dev"
        shift 1
        ;;
    (run_test)
        OPT="run_test"
        shift 1
        ;;
    (*)
        usage
        exit 3
        ;;
esac
done

if [[ ! -z "$OPT" ]]; then
    $OPT
else
    usage
fi
