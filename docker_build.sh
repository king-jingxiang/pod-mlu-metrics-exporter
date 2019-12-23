#!/usr/bin/env bash

ROOT_PATH=$(pwd)

cd ${ROOT_PATH}/src
go build -o pod-mlu-metrics-exporter -v github.com/ruanxingbaozi/pod-mlu-metrics-exporter/src

cd ${ROOT_PATH}
sudo docker build -t ruanxingbaozi/pod-mlu-metrics-exporter:mlu .

kubectl apply -f pod-mlu-metrics-exporter-daemonset.yaml

kubectl get po -nkube-system | grep pod-mlu-metrics-exporter | awk '{print $1}' | xargs kubectl delete po -nkube-system

sleep 3
kubectl get po -nkube-system | grep pod-mlu-metrics-exporter | awk '{print $1}' | xargs kubectl logs -f -nkube-system -c pod-cambricon-mlu-metrics-exporter


# kubectl exec -it $(kubectl get po -nkube-system | grep pod-mlu-metrics-exporter | awk '{print $1}') -nkube-system -c pod-cambricon-mlu-metrics-exporter bash
# kubectl exec -it $(kubectl get po -nkube-system | grep pod-mlu-metrics-exporter | awk '{print $1}') -nkube-system -c cambricon-cnmon-exporter bash