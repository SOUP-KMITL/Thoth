// TODO : need to change to api package
package thoth

type AppMetric struct {
	App         string  `json:"app"`
	Cpu         float32 `json:"cpu"`
	Memory      int64   `json:"memory"`
	Request     int64
	Response    int64
	Response2xx int64
	Response4xx int64
	Response5xx int64
}

// This is schema of Application Profile
type App struct {
	Name       string `json:"name"`
	ExternalIp string `json:"external_ip"`
	InternalIp string `json:"internal_ip"`
	Image      string `json:"image"`
	Pods       []Pod  `json:"pods"`
}

type Apps []App

type RC struct {
	Namespace string // Namespace = User
	Name      string
}

var KubeApi string = "http://localhost:8080"
var InfluxdbApi string = "http://localhost:8086"
var VampApi string = "http://localhost:10001"
