package entities

type ClusterNode struct {
	Name       string            `json:"name"`
	Status     string            `json:"status"`
	Roles      []string          `json:"roles"`
	Labels     map[string]string `json:"labels"`
	IP         string            `json:"ip"`
	KubeletVer string            `json:"kubeletVersion"`
}

type NodeStats struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	PodCount    int     `json:"podCount"`
}

type Pod struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Status       string            `json:"status"`
	IP           string            `json:"ip"`
	NodeName     string            `json:"nodeName"`
	Labels       map[string]string `json:"labels"`
	RestartCount int32             `json:"restartCount"`
}

type Workflow struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Service struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Type      string            `json:"type"`
	ClusterIP string            `json:"clusterIP"`
	Ports     []ServicePort     `json:"ports"`
	Labels    map[string]string `json:"labels"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
	Protocol   string `json:"protocol"`
}

type Deployment struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	ReadyReplicas int32             `json:"readyReplicas"`
	Labels        map[string]string `json:"labels"`
	Image         string            `json:"image"`
}

type ClusterStats struct {
	Nodes struct {
		Total int `json:"total"`
		Ready int `json:"ready"`
	} `json:"nodes"`
	Pods struct {
		Total int `json:"total"`
		Ready int `json:"ready"`
	} `json:"pods"`
	Deployments struct {
		Total int `json:"total"`
		Ready int `json:"ready"`
	} `json:"deployments"`
	StatefulSets struct {
		Total int `json:"total"`
		Ready int `json:"ready"`
	} `json:"statefulSets"`
}
