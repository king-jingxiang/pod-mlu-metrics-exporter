package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
)

const (
	mluResourceName = "cambricon.com/mlu"
)

// 存储map[uuid]index
var mluUUIDMap = make(map[string]uint)

type devicePodInfo struct {
	name      string
	namespace string
	container string
}

// Helper function that creates a map of pod info for each device
func createDevicePodMap(devicePods podresourcesapi.ListPodResourcesResponse) map[string]devicePodInfo {
	deviceToPodMap := make(map[string]devicePodInfo)

	for _, pod := range devicePods.GetPodResources() {
		for _, container := range pod.GetContainers() {
			for _, device := range container.GetDevices() {
				if device.GetResourceName() == mluResourceName {
					podInfo := devicePodInfo{
						name:      pod.GetName(),
						namespace: pod.GetNamespace(),
						container: container.GetName(),
					}
					for _, uuid := range device.GetDeviceIds() {
						deviceToPodMap[uuid] = podInfo
					}
				}
			}
		}
	}
	return deviceToPodMap
}

type deviceInfo struct {
	count      int
	namespace  string
	container  string
	podName    string
	deviceName string
}

// Helper function that creates a map of pod info for each device
func getPodDeviceCount(devicePods podresourcesapi.ListPodResourcesResponse) map[string]deviceInfo {
	devicePodCount := make(map[string]deviceInfo)

	for _, pod := range devicePods.GetPodResources() {
		for _, container := range pod.GetContainers() {
			for _, device := range container.GetDevices() {
				deviceInfo := deviceInfo{
					count:      len(device.DeviceIds),
					podName:    pod.GetName(),
					deviceName: device.GetResourceName(),
					namespace:  pod.GetNamespace(),
					container:  container.GetName(),
				}
				devicePodCount[container.GetName()] = deviceInfo
			}
		}
	}
	return devicePodCount
}
func getDevicePodInfo(socket string) (map[string]devicePodInfo, error) {
	devicePods, err := getListOfPods(socket)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices Pod information: %v", err)
	}
	return createDevicePodMap(*devicePods), nil

}

func getDevicePodCount(socket string) (map[string]deviceInfo, error) {
	devicePods, err := getListOfPods(socket)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices Pod information: %v", err)
	}
	return getPodDeviceCount(*devicePods), nil

}

func addPodInfoToMetrics(dir string, srcFile string, destFile string, deviceToPodMap map[string]devicePodInfo) error {
	readFI, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", srcFile, err)
	}
	defer readFI.Close()
	reader := bufio.NewReader(readFI)

	tmpPrefix := "pod"
	tmpF, err := ioutil.TempFile(dir, tmpPrefix)
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}

	tmpFname := tmpF.Name()
	defer func() {
		tmpF.Close()
		os.Remove(tmpFname)
	}()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(line) == 0 {
				return writeDestFile(tmpFname, destFile)
			}
			return fmt.Errorf("error reading %s: %v", srcFile, err)
		}
		// Skip comments and add pod info
		if string(line[0]) != "#" {
			uuid := strings.Split(strings.Split(line, ",")[1], "\"")[1]
			if pod, exists := deviceToPodMap[uuid]; exists {
				splitLine := strings.Split(line, "}")
				line = fmt.Sprintf("%s,pod_name=\"%s\",pod_namespace=\"%s\",container_name=\"%s\"}%s", splitLine[0], pod.name, pod.namespace, pod.container, splitLine[1])
			}
			//glog.Infof("%v:%v",uuid,deviceToPodMap)
		}

		_, err = tmpF.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing to %s: %v", tmpFname, err)
		}
	}
}

func addDeviceInfoToMetrics(dir string, destFile string, deviceInfo map[string]deviceInfo) error {

	tmpPrefix := "device"
	tmpF, err := ioutil.TempFile(dir, tmpPrefix)
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	tmpFname := tmpF.Name()
	defer func() {
		tmpF.Close()
		os.Remove(tmpFname)
	}()
	_, err = tmpF.WriteString("# TYPE container_device_count gauge\n")
	_, err = tmpF.WriteString("# HELP container_device_count container used device count (in num).\n")

	if err != nil {
		return fmt.Errorf("error writing to %s: %v", tmpFname, err)
	}
	//# TYPE container_device_count gauge
	//container_device_count{device_name="nvidia.com/gpu",pod_name="test-pod-01",pod_namespace="default",container_name="test"} 2
	//container_device_count{device_name="xilinx.com/fpga-xilinx_u200_xdma_201820_1-1535712995",pod_name="test-pod-02",pod_namespace="default",container_name="test"} 1
	//container_device_count{device_name="cambricon.com/mlu",pod_name="test-pod-03",pod_namespace="default",container_name="test"} 1
	for _, dev := range deviceInfo {
		line := fmt.Sprintf("container_device_count{device_name=\"%v\",pod_name=\"%s\",pod_namespace=\"%s\",container_name=\"%s\"} %v\n",
			dev.deviceName, dev.podName, dev.namespace, dev.container, dev.count)
		_, err = tmpF.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing to %s: %v", tmpFname, err)
		}
	}
	return writeDestFile(tmpFname, destFile)
}
