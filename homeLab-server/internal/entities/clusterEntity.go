package entities

type Node struct {
	ID      uint    `gorm:"primary_key"`
	Role    string  `gorm:"not null"`
	Label   string  `gorm:"not null"`
	Version string  `gorm:"not null"`
	Status  string  `gorm:"not null"`
	CPU     int     `gorm:"not null"`
	Memory  int     `gorm:"not null"`
	Disk    int     `gorm:"not null"`
	GPU     int     `gorm:"not null"`
	Alert   *string `gorm:"not null"`
}

type Cluster struct {
	Nodes       []Node `gorm:"many2many:cluster_nodes;"`
	Pods        []Pod  `gorm:"many2many:cluster_pods;"`
	MemoryUsage int    `gorm:"not null"`
	CPUUsage    int    `gorm:"not null"`
}

type Pod struct {
	ID     uint   `gorm:"primary_key"`
	Label  string `gorm:"not null"`
	Status string `gorm:"not null"`
}

type Workload struct {
	Name      string `gorm:"not null"`
	Namespace string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Ready     int    `gorm:"not null"`
	State     string `gorm:"not null"`
	Age       string `gorm:"not null"`
}
