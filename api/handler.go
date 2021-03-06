// TODO : need to change to api
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	tcp "net"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/SOUP-CE-KMITL/Thoth"
	"github.com/SOUP-CE-KMITL/Thoth/profil"
	"github.com/gorilla/mux"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// list every node
func GetNodes(w http.ResponseWriter, r *http.Request) {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get(thoth.KubeApi + "/api/v1/nodes")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	// defer for ensure that res is close.
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))
}

func GetNode(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// node name from user.
	nodesName := vars["nodeName"]
	// TODO: need to read api and port of api server from configuration file
	res, err := http.Get(thoth.KubeApi + "/api/v1/nodes/" + nodesName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var object map[string]interface{}
	err = json.Unmarshal([]byte(body), &object)
	if err == nil {
		fmt.Printf("%+v\n", object)
	} else {
		fmt.Println(err)
	}
	send_obj, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, string(send_obj))
}

func OptionCors(w http.ResponseWriter, r *http.Request) {
	// TODO: need to change origin to deployed domain name
	if origin := r.Header.Get("Origin"); origin != "http://localhost" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
}

// list specific node cpu
func NodeCpu(w http.ResponseWriter, r *http.Request) {
}

// list specifc node memory
func NodeMemory(w http.ResponseWriter, r *http.Request) {
}

// list all pods
func GetPods(w http.ResponseWriter, r *http.Request) {
	// to do need to read api and port of api server from configuration file
	res, err := http.Get(thoth.KubeApi + "/api/v1/pods")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(body))
}

// list specific pod details
func GetPod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// node name from user.
	podName := vars["podName"]
	fmt.Fprint(w, string(podName))
	// to do need to read api and port of api server from configuration file
	// TODO: change namespace to flexible.
	var dat map[string]interface{}
	res, err := http.Get(thoth.KubeApi + "/api/v1/namespaces/default/pods/" + podName)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	pretty_body, err := json.MarshalIndent(dat, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(pretty_body))
}

// list specific pod cpu
func PodCpu(w http.ResponseWriter, r *http.Request) {
}

// list specific pod memory
func PodMemory(w http.ResponseWriter, r *http.Request) {
}

// post handler for scale pod by pod name

// TODO : remove
// test mocks
func nodeTestMock(w http.ResponseWriter, r *http.Request) {
	nodes := thoth.Nodes{
		thoth.Node{Name: "node1", Ip: "192.168.1.2", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
		thoth.Node{Name: "node2", Ip: "192.169.1.4", Cpu: 5000, Memory: 3000, DiskUsage: 1000},
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		panic(err)
	}
}

// TODO : remove
// test ssh to exec command on other machine
func testExec(w http.ResponseWriter, r *http.Request) {
	commander := SSHCommander{"root", "161.246.70.75"}
	cmd := []string{
		"ls",
		".",
	}
	var (
		output []byte
		err    error
	)

	if output, err = commander.Command(cmd...).Output(); err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(output[:6]))
}

func GetApp(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// app name from user.
	appName := vars["appName"]
	fmt.Println(appName)
	// TODO: need to find new solution to get info from api like other done.
	res, err := exec.Command("kubectl", "get", "pod", "-l", "app="+appName, "-o", "json").Output()
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(res))
}

func GetApps(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// node name from user.
	namespace := vars["namespace"]

	res, err := exec.Command("kubectl", "get", "rc", "-o", "json", "--namespace="+namespace).Output()
	fmt.Println("namespace = " + namespace)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, string(res))
}

func CreatePod(w http.ResponseWriter, r *http.Request) {
	var pod thoth.Pod
	// limits json post request for prevent overflow attack.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// catch error from close reader
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// get request json information
	if err := json.Unmarshal(body, &pod); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// prepare json to send to create by kubernetes api server
	labels := map[string]interface{}{
		"app": pod.Name,
	}

	metadata := map[string]interface{}{
		"name":   pod.Name,
		"labels": labels,
	}

	ports := map[string]interface{}{
		"containerPort": 80,
	}

	containers := map[string]interface{}{
		"name":   pod.Name,
		"image":  pod.Image,
		"ports":  []map[string]interface{}{ports},
		"memory": pod.Memory,
		"cpu":    pod.Cpu,
	}

	spec := map[string]interface{}{
		"containers": []map[string]interface{}{containers},
	}

	objReq := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata":   metadata,
		"spec":       spec,
	}

	jsonReq, err := json.MarshalIndent(objReq, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println("you sent ", string(jsonReq))
	// post json to kubernete api server

	// TODO: need to change name space to user namespace
	postUrl := thoth.KubeApi + "/api/v1/namespaces/default/pods"
	req, err := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonReq))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// defer for ensure
	defer resp.Body.Close()

	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(response))
}

/**
* 	Get current resource by using cgroup
**/
func GetCurrentResourceCgroup(container_id string, metric_type int) (uint64, error) {
	// TODO : Read Latency from HA Proxy
	// file path prefix
	var path = "/sys/fs/cgroup/memory/docker/" + container_id + "/"

	if metric_type == 2 {
		// read memory usage
		current_usage, err := ioutil.ReadFile(path + "memory.usage_in_bytes")
		if err != nil {
			return binary.BigEndian.Uint64(current_usage), nil
		} else {
			return 0, err
		}
	} else if metric_type == 3 {
		// read memory usage
		current_usage, err := ioutil.ReadFile(path + "memory.usage_in_bytes")
		if err != nil {
			return 0, err
		} else {
			n := bytes.Index(current_usage, []byte{10})
			usage_str := string(current_usage[:n])
			resource_usage, _ := strconv.ParseInt(usage_str, 10, 64)
			return uint64(resource_usage), nil
		}
	} else {
		// not match any case
		return 0, errors.New("not match any case")
	}
}

/**
 *   Get List of ContainerID and pod's ip by replication name and their namespace
 **/
func GetContainerIDList(url string, rc_name string, namespace string) ([]string, []string, error) {
	// TODO : maybe user want to get container id which map with it's pod
	// initail rasult array
	container_ids := []string{}
	pod_ips := []string{}

	res, err := http.Get(url + "/api/v1/namespaces/" + namespace + "/pods")
	if err != nil {
		fmt.Println("Can't connect to cadvisor")
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, nil, err
	} else {
		// json handler type
		var res_obj map[string]interface{}
		if err := json.Unmarshal(body, &res_obj); err != nil {
			return nil, nil, err
		}
		pod_arr := res_obj["items"].([]interface{})
		// iterate to get pod of specific rc
		for _, pod := range pod_arr {
			pod_name := pod.(map[string]interface{})["metadata"].(map[string]interface{})["generateName"]
			if pod_name != nil {
				if pod_name == rc_name+"-" {
					pod_ips = append(pod_ips, pod.(map[string]interface{})["status"].(map[string]interface{})["podIP"].(string))
					containers := pod.(map[string]interface{})["status"].(map[string]interface{})["containerStatuses"].([]interface{})
					// one pod can has many container ,so iterate for get each container
					for _, container := range containers {
						container_id := container.(map[string]interface{})["containerID"].(string)[9:]
						container_ids = append(container_ids, container_id)
					}
				}
			}
		}
		return container_ids, pod_ips, nil
	}
}

/**
	read node resource usage
**/
func GetNodeResource(w http.ResponseWriter, r *http.Request) {
	// get this node memory
	memory, _ := mem.VirtualMemory()
	// get this node cpu percent usage
	cpu_percent, _ := cpu.CPUPercent(time.Duration(1)*time.Second, false)
	// Disk mount Point
	disk_partitions, _ := disk.DiskPartitions(true)
	// Disk usage
	var disk_usages []*disk.DiskUsageStat
	for _, disk_partition := range disk_partitions {
		if disk_partition.Mountpoint == "/" || disk_partition.Mountpoint == "/home" {
			disk_stat, _ := disk.DiskUsage(disk_partition.Device)
			disk_usages = append(disk_usages, disk_stat)
		}
	}
	// Network
	network, _ := net.NetIOCounters(false)

	// create new node obj with resource usage information
	node_metric := thoth.NodeMetric{
		Cpu:       cpu_percent,
		Memory:    memory,
		DiskUsage: disk_usages,
		Network:   network,
	}

	node_json, err := json.MarshalIndent(node_metric, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprint(w, string(node_json))
}

/**
 CPU Percent Calculation
**/
func DockerCPUPercent(container_id string) (float32, error) {

	res, err := profil.GetCpu(container_id)
	if err != nil {
		return 0.0, nil
	}
	fmt.Println("res : ", res)
	return float32(res), nil

}

/**
 	get resource usage of application (pods) on node
**/
func GetAppResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// get app Name
	appName := vars["appName"]
	namespace := vars["namespace"]

	var summary_cpu float32
	var memory_bundle []*docker.CgroupMemStat

	container_ids, pod_ips, err := GetContainerIDList(thoth.KubeApi, appName, namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, container_id := range container_ids {
		fmt.Println(container_id, pod_ips)
		// calculation percentage of cpu usage
		container_cpu, _ := DockerCPUPercent(container_id)
		summary_cpu += container_cpu
		// memory usage
		container_memory, _ := docker.CgroupMemDocker(container_id)
		memory_bundle = append(memory_bundle, container_memory)
	}

	podNum := len(pod_ips)

	// find the request per sec from haproxy-frontend
	res_front, err := http.Get(thoth.VampApi + "/v1/stats/frontends")
	if err != nil {
		panic(err)
	}
	body_front, err := ioutil.ReadAll(res_front.Body)
	res_front.Body.Close()
	if err != nil {
		panic(err)
	}
	//var rps uint64
	var object_front []map[string]interface{}
	err = json.Unmarshal([]byte(body_front), &object_front)
	rps := object_front[0]["req_rate"].(string)
	rps_int, _ := strconv.ParseInt(rps, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}

	//find resonse time from haproxy-backends
	//var rtime uint64
	res_back, err := http.Get(thoth.VampApi + "/v1/stats/backends")
	if err != nil {
		panic(err)
	}
	body_back, err := ioutil.ReadAll(res_back.Body)
	res_back.Body.Close()
	if err != nil {
		panic(err)
	}

	var object_back []map[string]interface{}
	err = json.Unmarshal([]byte(body_back), &object_back)
	rtime := object_back[0]["rtime"].(string)
	res_2xx := object_back[0]["hrsp_2xx"].(string)
	res_4xx := object_back[0]["hrsp_4xx"].(string)
	res_5xx := object_back[0]["hrsp_5xx"].(string)
	rtime_int, _ := strconv.ParseInt(rtime, 10, 64)
	res2xx_int, _ := strconv.ParseInt(res_2xx, 10, 64)
	res4xx_int, _ := strconv.ParseInt(res_4xx, 10, 64)
	res5xx_int, _ := strconv.ParseInt(res_5xx, 10, 64)
	if err == nil {
	} else {
		fmt.Println(err)
	}

	fmt.Println("rps: ", rps, ", rtime: ", rtime)
	// find the cpu avarage of application cpu usage
	average_cpu := summary_cpu / float32(len(container_ids))
	fmt.Println("avg_cpu : ", average_cpu)
	fmt.Println("summary_cpu : ", summary_cpu)
	// Cal Avg Mem usage
	var avgMem uint64
	for i := 0; i < podNum; i++ {
		avgMem += memory_bundle[i].MemUsageInBytes
	}
	avgMem = avgMem / uint64(podNum)
	avgMem = avgMem / uint64(1024*1024) // MB

	// create appliction object
	app_metric := thoth.AppMetric{
		App:         appName,
		Cpu:         average_cpu,
		Memory:      int64(avgMem),
		Request:     rps_int,
		Response:    rtime_int,
		Response2xx: res2xx_int,
		Response4xx: res4xx_int,
		Response5xx: res5xx_int,
	}

	app_json, err := json.MarshalIndent(app_metric, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Fprint(w, string(app_json))

}

/**
	pull new image form dockerhub
**/
func PullDockerhub(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// get app Name
	docker_repo := vars["dockerhub_repo"]
	postUrl := "http://localhost:4243/images/create?fromImage=" + docker_repo
	fmt.Println(string(postUrl))
	req, err := http.NewRequest("POST", postUrl, nil)
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// defer for ensure
	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(response))
}

func CreateRc(w http.ResponseWriter, r *http.Request) {
	var rc thoth.SendRC
	// limits json post request for prevent overflow attack.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	// catch error from close reader
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// get request json information
	if err := json.Unmarshal(body, &rc); err != nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	fmt.Println(rc)

	ports := map[string]interface{}{
		"containerPort": rc.Port,
	}

	containers := map[string]interface{}{
		"name":  rc.Name,
		"image": rc.Image,
		"ports": []map[string]interface{}{ports},
		"resources": map[string]interface{}{
			"limits": map[string]interface{}{
				"memory": "200Mi",
				"cpu":    "1000m",
			},
		},
	}

	objReq := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ReplicationController",
		"metadata": map[string]interface{}{
			"name":      rc.Name,
			"namespace": rc.Namespace,
		},
		"spec": map[string]interface{}{
			"replicas": rc.Replicas,
			"selector": map[string]interface{}{
				"app": rc.Name,
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app": rc.Name,
					},
					"name": rc.Name,
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{containers},
				},
			},
		},
	}

	jsonReq, err := json.MarshalIndent(objReq, "", "\t")
	if err != nil {
		panic(err)
	}

	//	fmt.Println("you sent ", string(jsonReq))
	// post json to kubernete api server

	rcResCode, rcResBody := postJson(thoth.KubeApi+"/api/v1/namespaces/"+rc.Namespace+"/replicationcontrollers", jsonReq)

	if rcResCode/200 == 2 { // RC Create fail
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, rcResBody["message"].(string))
	} else {
		// RC Create success
		// Random Port (30,000-31,000) for SVC

		nodePort := RandomTCPPort()
		if nodePort == -1 {
			// Dont have port in specifi range to bind
			fmt.Println("Can't bind port")
		} else {
			// Creat SVC
			svc := thoth.Svc{
				APIVersion: "v1",
				Kind:       "Service",
				Metadata: thoth.Metadata{
					Name:      rc.Name,
					Namespace: rc.Namespace,
				},
				Spec: thoth.Spec{
					Ports: []thoth.Port{
						thoth.Port{
							NodePort:   nodePort,
							Port:       80,
							TargetPort: rc.Port,
						}},
					Selector: thoth.Selector{
						App: rc.Name,
					},
					Type: "LoadBalancer",
				},
				Status: thoth.SvcStatus{
					LoadBalancer: thoth.SvcLoadBalancer{
						Ingress: []thoth.SvcIngress{
							thoth.SvcIngress{
								IP: "10.0.1.17",
							}},
					},
				},
			}
			jsonSvc, err := json.MarshalIndent(svc, "", "\t")
			if err != nil {
				panic(err)
			}
			//		fmt.Println(svc)
			//		fmt.Println(string(jsonSvc))
			svcResCode, svcResBody := postJson(thoth.KubeApi+"/api/v1/namespaces/"+rc.Namespace+"/services", jsonSvc)
			fmt.Println(svcResCode)
			fmt.Println(svcResBody)
			fmt.Fprint(w, "{\"port\":", nodePort, "}")
		}
	}

}

func postJson(url string, data []byte) (int, map[string]interface{}) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// defer for ensure
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		panic(err)
	}
	return resp.StatusCode, response
}

func SpeechRecog(w http.ResponseWriter, r *http.Request) {
	//curl -u 0648b905-6758-4ae0-9c15-5cc53fadfa24:9e90DAlqVwoK -X POST --header "Content-Type: audio/wav" --header "Transfer-Encoding: chunked" --data-binary @hello.wav "https://stream.watsonplatform.net/speech-to-text/api/v1/recognize?continuous=true"
	// Basic-Auth
	// POST
	// Content-Type
	// TransferEncoding
	// DataBinary
	// URL Test https://sj8vwegp7eew.runscope.net

	//fmt.Println(w, r)

	url := "https://stream.watsonplatform.net/speech-to-text/api/v1/recognize?continuous=true"
	//url := "https://sj8vwegp7eew.runscope.net"

	req, err := http.NewRequest("POST", url, r.Body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "curl/7.46.0")
	req.Header.Set("Content-Type", "audio/wav")
	req.SetBasicAuth("0648b905-6758-4ae0-9c15-5cc53fadfa24", "9e90DAlqVwoK")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// defer for ensure
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//	fmt.Println(string(body))

	fmt.Fprint(w, string(body))
	/*
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			panic(err)
		}
		fmt.Println(response)
	*/
}

// The error application
func getErrorApp(w http.ResponseWriter, r *http.Request) {
	var username string = "thoth"
	var password string = "thoth"

	// Connect to InfluxDB
	c, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr:     thoth.InfluxdbApi,
		Username: username,
		Password: password,
	})
	// TODO: iterate to query for each app
	queryRes, err := profil.QueryDB(c, fmt.Sprint("SELECT last(code5xx) FROM thoth"))
	if err != nil {
		log.Fatal(err)
	}
	errorNum5xx := queryRes[0].Series[0].Values[0][1]
	fmt.Println(errorNum5xx)
	// TODO: need to implement some algorithm searching error application
}

// IsTCPPortAvailable returns a flag indicating whether or not a TCP port is
// available.
func IsTCPPortAvailable(port int) bool {
	conn, err := tcp.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// RandomTCPPort gets a free, random TCP port between 1025-65535. If no free
// ports are available -1 is returned.
func RandomTCPPort() int {
	tcpPortRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		p := 30000 + tcpPortRand.Intn(100)
		fmt.Println("RndPort", p)
		if IsTCPPortAvailable(p) {
			return p
		}
	}
	return -1
}

// Return Array of Succesfull running pod
func GetRunStatus(w http.ResponseWriter, r *http.Request) {
	reqParam := mux.Vars(r)
	namespace := reqParam["namespace"]
	fmt.Println(namespace)
	runningPodStatus := profil.GetRunningPodStatus(namespace)
	res, err := json.Marshal(runningPodStatus)
	fmt.Println(err)
	fmt.Fprint(w, string(res))
}
