# Pod MLU Metrics Exporter

A simple go http server serving per pod MLU metrics at localhost:9410/mlu/metrics. The exporter connects to kubelet gRPC server (/var/lib/kubelet/pod-resources) to identify the MLUs running on a pod leveraging Kubernetes [device assignment feature](https://github.com/vikaschoudhary16/community/blob/060a25c441269be476ade624ea0347ebc113e659/keps/sig-node/compute-device-assignment.md) and appends the MLU device's pod information to metrics collected by [cnmon-exporter](./cnmon-exporter).

The http server allows Prometheus to scrape MLU metrics directly via a separate endpoint without relying on node-exporter. 


## MLU metrics example
```bash
curl http://localhost:9410/mlu/metrics
or
curl http://[ip]:31400/mlu/metrics
```
```
# HELP HELP cnmon_mlu_temp MLU temperature (in C).
# TYPE cnmon_mlu_temp gauge
cnmon_mlu_temp{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 28
# HELP cnmon_mlu_board_utilization MLU utilization (in %).
# TYPE cnmon_mlu_board_utilization gauge
cnmon_mlu_board_utilization{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 0
# HELP cnmon_physical_memory_used MLU utilization (in MB).
# TYPE cnmon_physical_memory_used gauge
cnmon_physical_memory_used{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 4192
# HELP cnmon_physical_memory_free MLU utilization (in MB).
# TYPE cnmon_physical_memory_free gauge
cnmon_physical_memory_free{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 4000
# HELP cnmon_virtual_memory_used MLU utilization (in MB).
# TYPE cnmon_virtual_memory_used gauge
cnmon_virtual_memory_used{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 16643
# HELP cnmon_virtual_memory_free MLU utilization (in MB).
# TYPE cnmon_virtual_memory_free gauge
cnmon_virtual_memory_free{mlu="0",uuid="MLU-11812120333",pod_name="pod1",pod_namespace="default",container_name="pod1-ctr"} 16125
cnmon_mlu_temp{mlu="1",uuid="MLU-11812120351"} 27
```
## prometheus configMap
```bash
apiVersion: v1
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: prometheus-config
  namespace: kube-system
data:
  prometheus.yaml: |
    global:
      scrape_interval: 15s
    scrape_configs:
    - job_name: prometheus
      scrape_interval: 5s
      static_configs:
      - targets:
        - localhost:9090
    - job_name: mlu_metrics
      scrape_interval: 1s
      metrics_path: /mlu/metrics
      scheme: http
      static_configs:
      - targets:
        - localhost:9410
```
## grafana metrics promql
#### Cluster MLU usage
```promql
sum((count(cnmon_mlu_temp{container_name=~".+"} or vector(0))-1) / count(cnmon_mlu_temp{uuid=~".+"}) or vector(0))*100
```

#### Used MLUs
```promql
count(cnmon_mlu_board_utilization{container_name=~".+"} or vector(0))-1
```

#### Total MLUs
```promql
count(cnmon_mlu_temp{uuid=~".+"} or vector(0)) - 1
```

#### Cluster MLU utilization (avg)
```promql
sum((sum(cnmon_mlu_board_utilization{pod_name=~".*"}) / count(cnmon_mlu_board_utilization{pod_name=~".*"}) ) or vector(0))
```

#### Cluster MLU virtual memory utilization(avg)
```promql
sum((sum(cnmon_virtual_memory_used{pod_name=~".+"}) / (sum(cnmon_virtual_memory_used) + sum(cnmon_virtual_memory_free))) or vector(0))
```

#### Cluster MLU physical memory utilization(avg)
```promql
sum((sum(cnmon_physical_memory_free{pod_name=~".+"}) / (sum(cnmon_physical_memory_free) + sum(cnmon_physical_memory_used))) or vector(0))
```

## deploy
```bash
cat <<EOF >pod-mlu-metrics-exporter-daemonset.yaml
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: pod-mlu-metrics-exporter
  namespace: kube-system
spec:
  template:
    metadata:
      labels:
        app: pod-mlu-metrics-exporter
      name: pod-mlu-metrics-exporter
    spec:
      containers:
      - image: 192.168.202.74:5000/openi/pod-mlu-metrics-exporter:v0.1
        name: pod-cambricon-mlu-metrics-exporter
        ports:
        - name: mlu-metrics
          containerPort: 9410
          hostPort: 9410
        securityContext:
          runAsNonRoot: false
          runAsUser: 0
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
        volumeMounts:
        - name: pod-mlu-resources
          readOnly: true
          mountPath: /var/lib/kubelet/pod-resources
        - name: device-metrics
          readOnly: true
          mountPath: /run/prometheus
      - image: 192.168.202.74:5000/openi/cambricon-cnmon-exporter:v0.1
        name: cambricon-mlu-exporter
        securityContext:
          privileged: true
          capabilities:
            drop: ["ALL"]
        volumeMounts:
        - name: device-metrics
          mountPath: /run/prometheus
      volumes:
      - name: pod-mlu-resources
        hostPath:
          path: /var/lib/kubelet/pod-resources
      - name: device-metrics
        emptyDir:
          medium: Memory
---
kind: Service
apiVersion: v1
metadata:
  labels:
    app: pod-mlu-metrics-svc
  annotations:
    prometheus.io/port: "9410"
    prometheus.io/path: "/mlu/metrics"
    prometheus.io/scrape: "true"
  name: pod-mlu-metrics-svc
  namespace: kube-system
spec:
  type: NodePort
  ports:
    - port: 9410
      targetPort: 9410
      nodePort: 31400
  selector:
    app: pod-mlu-metrics-exporter
EOF
kubectl appy -f pod-mlu-metrics-exporter-daemonset.yaml
```
# 192.168.202.102
```bash
# metrics
curl http://192.168.202.102:31400/mlu/metrics
# prometheus
curl http://192.168.202.102:30003
```