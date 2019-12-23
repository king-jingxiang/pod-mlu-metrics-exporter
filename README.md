# Pod GPU Metrics Exporter

A simple go http server serving per pod GPU metrics at localhost:9410/gpu/metrics. The exporter connects to kubelet gRPC server (/var/lib/kubelet/pod-resources) to identify the GPUs running on a pod leveraging Kubernetes [device assignment feature](https://github.com/vikaschoudhary16/community/blob/060a25c441269be476ade624ea0347ebc113e659/keps/sig-node/compute-device-assignment.md) and appends the GPU device's pod information to metrics collected by [dcgm-exporter](https://github.com/NVIDIA/gpu-monitoring-tools/tree/master/exporters/prometheus-dcgm/dcgm-exporter).

The http server allows Prometheus to scrape GPU metrics directly via a separate endpoint without relying on node-exporter. But if you still want to scrape GPU metrics via node-exporter, follow [these instructions](https://github.com/NVIDIA/gpu-monitoring-tools/tree/master/exporters/prometheus-dcgm#node_exporter).

### Prerequisites
* NVIDIA Tesla drivers = R384+ (download from [NVIDIA Driver Downloads page](http://www.nvidia.com/drivers))
* nvidia-docker version > 2.0 (see how to [install](https://github.com/NVIDIA/nvidia-docker) and it's [prerequisites](https://github.com/nvidia/nvidia-docker/wiki/Installation-\(version-2.0\)#prerequisites))
* Set the [default runtime](https://github.com/NVIDIA/nvidia-container-runtime#daemon-configuration-file) to nvidia
* Kubernetes version = 1.13
* Set KubeletPodResources in /etc/default/kubelet: KUBELET_EXTRA_ARGS=--feature-gates=KubeletPodResources=true

#### Deploy on Kubernetes cluster 
```sh
# Deploy nvidia-k8s-device-plugin

# Deploy GPU Pods

# Create the monitoring namespace
$ kubectl create namespace monitoring

# Add gpu metrics endpoint to prometheus
$ kubectl create -f prometheus/prometheus-configmap.yaml

# Deploy prometheus
$ kubectl create -f prometheus/prometheus-deployment.yaml

$ kubectl create -f pod-gpu-metrics-exporter-daemonset.yaml

# Open in browser: localhost:9090
```

#### Docker Build and Run

```sh
$ docker build -t pod-gpu-metrics-exporter .

# Make sure to run dcgm-exporter
$ docker run -d --runtime=nvidia --rm --name=nvidia-dcgm-exporter nvidia/dcgm-exporter

$ docker run -d --privileged --rm -p 9410:9410 -v /var/lib/kubelet/pod-resources:/var/lib/kubelet/pod-resources --volumes-from nvidia-dcgm-exporter:ro nvidia/pod-gpu-metrics-exporter:v1.0.0-alpha

# Check GPU metrics
$ curl -s localhost:9410/gpu/metrics

# Sample output

# HELP dcgm_gpu_temp GPU temperature (in C).
# TYPE dcgm_gpu_temp gauge
dcgm_gpu_temp{container_name="pod1-ctr",gpu="0",pod_name="pod1",pod_namespace="default",uuid="GPU-2b399198-c670-a848-173b-d3400051a200"} 33
dcgm_gpu_temp{container_name="pod1-ctr",gpu="1",pod_name="pod1",pod_namespace="default",uuid="GPU-9567a9e7-341e-bb7e-fcf5-788d8caa50f9"} 34
# HELP dcgm_gpu_utilization GPU utilization (in %).
# TYPE dcgm_gpu_utilization gauge
dcgm_gpu_utilization{container_name="pod1-ctr",gpu="0",pod_name="pod1",pod_namespace="default",uuid="GPU-2b399198-c670-a848-173b-d3400051a200"} 0
dcgm_gpu_utilization{container_name="pod1-ctr",gpu="1",pod_name="pod1",pod_namespace="default",uuid="GPU-9567a9e7-341e-bb7e-fcf5-788d8caa50f9"} 0
# HELP dcgm_low_util_violation Throttling duration due to low utilization (in us).
# TYPE dcgm_low_util_violation counter
dcgm_low_util_violation{container_name="pod1-ctr",gpu="0",pod_name="pod1",pod_namespace="default",uuid="GPU-2b399198-c670-a848-173b-d3400051a200"} 0
dcgm_low_util_violation{container_name="pod1-ctr",gpu="1",pod_name="pod1",pod_namespace="default",uuid="GPU-9567a9e7-341e-bb7e-fcf5-788d8caa50f9"} 0
# HELP dcgm_mem_copy_utilization Memory utilization (in %).
# TYPE dcgm_mem_copy_utilization gauge
dcgm_mem_copy_utilization{container_name="pod1-ctr",gpu="0",pod_name="pod1",pod_namespace="default",uuid="GPU-2b399198-c670-a848-173b-d3400051a200"} 0
dcgm_mem_copy_utilization{container_name="pod1-ctr",gpu="1",pod_name="pod1",pod_namespace="default",uuid="GPU-9567a9e7-341e-bb7e-fcf5-788d8caa50f9"} 0
# HELP dcgm_memory_clock Memory clock frequency (in MHz).
# TYPE dcgm_memory_clock gauge
dcgm_memory_clock{container_name="pod1-ctr",gpu="0",pod_name="pod1",pod_namespace="default",uuid="GPU-2b399198-c670-a848-173b-d3400051a200"} 810
dcgm_memory_clock{container_name="pod1-ctr",gpu="1",pod_name="pod1",pod_namespace="default",uuid="GPU-9567a9e7-341e-bb7e-fcf5-788d8caa50f9"} 810
```

#### Build and Run locally
```sh
$ git clone
$ cd src && go build。
$ sudo ./src
```
### For GPUShare
#### Add gpu process memory used metrics
```
add /var/run/docker.sock    # used to get pod container pid
add hostPID: true           # used to check whether the parent porcess of the gpu process is pod container process 
use nvml                    # used to get the gpu process used memory
```
#### metrics sample output
```
# TYPE container_device_count gauge
# HELP container_device_count container used device count (in num).
container_device_count{device_name="nvidia.com/gpu",pod_name="test-pod-01",pod_namespace="default",container_name="test"} 2
container_device_count{device_name="xilinx.com/fpga-xilinx_u200_xdma_201820_1-1535712995",pod_name="test-pod-02",pod_namespace="default",container_name="test"} 1
container_device_count{device_name="cambricon.com/mlu",pod_name="test-pod-03",pod_namespace="default",container_name="test"} 1
```

### deploy
```bash
kubectl apply -f ./prometheus
kubectl apply -f ./grafana
```

```bash

sudo docker run -it --rm --device=/dev/cnmon_dev --device=/dev/cambricon_c10Dev0 cambricon-k8s-device-plugin:v0.2.0 cnmon
sudo docker run -it --rm --device=/dev/cnmon_dev --device=/dev/cambricon_c10Dev0 cambricon-office-ubuntu:v0.4 cnmon
sudo docker run -it --rm --device=/dev/cnmon_dev --device=/dev/cambricon_c10Dev0 cambricon-test1:v0.1 cnmon

```