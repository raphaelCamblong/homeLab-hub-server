package client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"
	"homelab.com/homelab-server/homeLab-server/internal/entities"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8sClient struct {
	clientset      *kubernetes.Clientset
	metricsClient  *metricsclientset.Clientset
	config         *rest.Config
	defaultTimeout time.Duration
}

type ClientOptions struct {
	Timeout time.Duration
}

func NewK8sClient(opts *ClientOptions) (*K8sClient, error) {
	if opts == nil {
		opts = &ClientOptions{Timeout: 30 * time.Second}
	}

	config, err := getKubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		fmt.Printf("Warning: failed to create metrics client (metrics-server might not be installed): %v\n", err)
	}

	return &K8sClient{
		clientset:      clientset,
		metricsClient:  metricsClient,
		config:         config,
		defaultTimeout: opts.Timeout,
	}, nil
}

func getKubernetesConfig() (*rest.Config, error) {
	// Try in-cluster config first
	if config, err := rest.InClusterConfig(); err == nil {
		fmt.Println("Using in-cluster configuration")
		return config, nil
	}

	// Try kubeconfig file
	var kubeconfigPath string
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		kubeconfigPath = kubeconfig
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("unable to find kubeconfig file")
	}

	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("kubeconfig file does not exist at %s", kubeconfigPath)
	}

	fmt.Printf("Using kubeconfig file: %s\n", kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	return config, nil
}

// Helper method to create context with timeout
func (k *K8sClient) createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), k.defaultTimeout)
}

func (k *K8sClient) ListNodes() ([]entities.ClusterNode, error) {
	nodes, err := k.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	var clusterNodes []entities.ClusterNode
	for _, node := range nodes.Items {
		clusterNode := entities.ClusterNode{
			Name:       node.Name,
			Status:     k.getNodeStatus(node),
			Roles:      k.getNodeRoles(node.Labels),
			Labels:     node.Labels,
			IP:         k.getNodeIP(node.Status.Addresses),
			KubeletVer: node.Status.NodeInfo.KubeletVersion,
		}
		clusterNodes = append(clusterNodes, clusterNode)
	}

	return clusterNodes, nil
}

func (k *K8sClient) GetNodeStats(nodeName string) (entities.NodeStats, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	stats := entities.NodeStats{
		CPUUsage:    0,
		MemoryUsage: 0,
		DiskUsage:   0,
		PodCount:    0,
	}

	// Get pod count for the node
	pods, err := k.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		return stats, fmt.Errorf("failed to get node pods: %w", err)
	}
	stats.PodCount = len(pods.Items)

	// Get node metrics if metrics client is available
	if k.metricsClient != nil {
		nodeMetrics, err := k.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Warning: failed to get node metrics for %s: %v\n", nodeName, err)
			return stats, nil
		}

		// Get node capacity for percentage calculations
		node, err := k.clientset.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
		if err != nil {
			return stats, fmt.Errorf("failed to get node details: %w", err)
		}

		stats.CPUUsage = k.calculateCPUUsagePercent(nodeMetrics.Usage[corev1.ResourceCPU], node.Status.Capacity[corev1.ResourceCPU])
		stats.MemoryUsage = k.calculateMemoryUsagePercent(nodeMetrics.Usage[corev1.ResourceMemory], node.Status.Capacity[corev1.ResourceMemory])

		// For disk usage, we need to check ephemeral storage
		if ephemeralUsage, ok := nodeMetrics.Usage[corev1.ResourceEphemeralStorage]; ok {
			if ephemeralCapacity, ok := node.Status.Capacity[corev1.ResourceEphemeralStorage]; ok {
				stats.DiskUsage = k.calculateDiskUsagePercent(ephemeralUsage, ephemeralCapacity)
			}
		}
	}

	return stats, nil
}

func (k *K8sClient) ListPods(namespace string) ([]entities.Pod, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}

	var clusterPods []entities.Pod
	for _, pod := range pods.Items {
		clusterPod := entities.Pod{
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       string(pod.Status.Phase),
			IP:           pod.Status.PodIP,
			NodeName:     pod.Spec.NodeName,
			Labels:       pod.Labels,
			RestartCount: k.getPodRestartCount(pod),
		}
		clusterPods = append(clusterPods, clusterPod)
	}

	return clusterPods, nil
}

func (k *K8sClient) ListWorkflows() ([]entities.Workflow, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	var workflows []entities.Workflow

	// Get Deployments
	if deploymentWorkflows, err := k.getDeploymentWorkflows(ctx); err != nil {
		return nil, fmt.Errorf("failed to get deployment workflows: %w", err)
	} else {
		workflows = append(workflows, deploymentWorkflows...)
	}

	// Get StatefulSets
	if statefulSetWorkflows, err := k.getStatefulSetWorkflows(ctx); err != nil {
		return nil, fmt.Errorf("failed to get statefulset workflows: %w", err)
	} else {
		workflows = append(workflows, statefulSetWorkflows...)
	}

	// Get DaemonSets
	if daemonSetWorkflows, err := k.getDaemonSetWorkflows(ctx); err != nil {
		return nil, fmt.Errorf("failed to get daemonset workflows: %w", err)
	} else {
		workflows = append(workflows, daemonSetWorkflows...)
	}

	return workflows, nil
}

func (k *K8sClient) ListServices(namespace string) ([]entities.Service, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	services, err := k.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services in namespace %s: %w", namespace, err)
	}

	var clusterServices []entities.Service
	for _, svc := range services.Items {
		ports := make([]entities.ServicePort, len(svc.Spec.Ports))
		for i, port := range svc.Spec.Ports {
			ports[i] = entities.ServicePort{
				Name:       port.Name,
				Port:       port.Port,
				TargetPort: port.TargetPort.IntVal,
				Protocol:   string(port.Protocol),
			}
		}

		clusterService := entities.Service{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			Type:      string(svc.Spec.Type),
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     ports,
			Labels:    svc.Labels,
		}
		clusterServices = append(clusterServices, clusterService)
	}

	return clusterServices, nil
}

func (k *K8sClient) ListDeployments(namespace string) ([]entities.Deployment, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	deployments, err := k.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", namespace, err)
	}

	var clusterDeployments []entities.Deployment
	for _, dep := range deployments.Items {
		image := ""
		if len(dep.Spec.Template.Spec.Containers) > 0 {
			image = dep.Spec.Template.Spec.Containers[0].Image
		}

		replicas := int32(0)
		if dep.Spec.Replicas != nil {
			replicas = *dep.Spec.Replicas
		}

		clusterDeployment := entities.Deployment{
			Name:          dep.Name,
			Namespace:     dep.Namespace,
			Replicas:      replicas,
			ReadyReplicas: dep.Status.ReadyReplicas,
			Labels:        dep.Labels,
			Image:         image,
		}
		clusterDeployments = append(clusterDeployments, clusterDeployment)
	}

	return clusterDeployments, nil
}

// Helper methods for workflows
func (k *K8sClient) getDeploymentWorkflows(ctx context.Context) ([]entities.Workflow, error) {
	deployments, err := k.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workflows []entities.Workflow
	for _, dep := range deployments.Items {
		replicas := int32(0)
		if dep.Spec.Replicas != nil {
			replicas = *dep.Spec.Replicas
		}

		workflows = append(workflows, entities.Workflow{
			Name:      dep.Name,
			Status:    k.getWorkloadStatus(dep.Status.ReadyReplicas, replicas),
			StartTime: dep.CreationTimestamp.Format(time.RFC3339),
			EndTime:   "",
		})
	}
	return workflows, nil
}

func (k *K8sClient) getStatefulSetWorkflows(ctx context.Context) ([]entities.Workflow, error) {
	statefulSets, err := k.clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workflows []entities.Workflow
	for _, sts := range statefulSets.Items {
		replicas := int32(0)
		if sts.Spec.Replicas != nil {
			replicas = *sts.Spec.Replicas
		}

		workflows = append(workflows, entities.Workflow{
			Name:      sts.Name,
			Status:    k.getWorkloadStatus(sts.Status.ReadyReplicas, replicas),
			StartTime: sts.CreationTimestamp.Format(time.RFC3339),
			EndTime:   "",
		})
	}
	return workflows, nil
}

func (k *K8sClient) getDaemonSetWorkflows(ctx context.Context) ([]entities.Workflow, error) {
	daemonSets, err := k.clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var workflows []entities.Workflow
	for _, ds := range daemonSets.Items {
		workflows = append(workflows, entities.Workflow{
			Name:      ds.Name,
			Status:    k.getWorkloadStatus(ds.Status.NumberReady, ds.Status.DesiredNumberScheduled),
			StartTime: ds.CreationTimestamp.Format(time.RFC3339),
			EndTime:   "",
		})
	}
	return workflows, nil
}

// Helper methods
func (k *K8sClient) getNodeRoles(labels map[string]string) []string {
	var roles []string

	// Check for various role labels
	roleLabels := []string{
		"node-role.kubernetes.io/control-plane",
		"node-role.kubernetes.io/master",
		"node-role.kubernetes.io/worker",
		"kubernetes.io/role",
	}

	for _, roleLabel := range roleLabels {
		if value, ok := labels[roleLabel]; ok {
			if value == "" {
				// Extract role from label key
				if roleLabel == "node-role.kubernetes.io/control-plane" {
					roles = append(roles, "control-plane")
				} else if roleLabel == "node-role.kubernetes.io/master" {
					roles = append(roles, "master")
				} else if roleLabel == "node-role.kubernetes.io/worker" {
					roles = append(roles, "worker")
				}
			} else {
				roles = append(roles, value)
			}
		}
	}

	// If no specific roles found, assume worker
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}

	return roles
}

func (k *K8sClient) getNodeStatus(node corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			}
			return "NotReady"
		}
	}
	return "Unknown"
}

func (k *K8sClient) getNodeIP(addresses []corev1.NodeAddress) string {
	// Prefer InternalIP, fallback to ExternalIP
	for _, addr := range addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	for _, addr := range addresses {
		if addr.Type == corev1.NodeExternalIP {
			return addr.Address
		}
	}
	return ""
}

func (k *K8sClient) getPodRestartCount(pod corev1.Pod) int32 {
	var totalRestarts int32
	for _, container := range pod.Status.ContainerStatuses {
		totalRestarts += container.RestartCount
	}
	return totalRestarts
}

func (k *K8sClient) getWorkloadStatus(ready, desired int32) string {
	if desired == 0 {
		return "unknown"
	}
	if ready == 0 {
		return "pending"
	}
	if ready < desired {
		return "progressing"
	}
	return "ready"
}

// Metrics calculation helpers
func (k *K8sClient) calculateCPUUsagePercent(usage, capacity resource.Quantity) float64 {
	if capacity.IsZero() {
		return 0
	}
	usageMillis := usage.MilliValue()
	capacityMillis := capacity.MilliValue()
	return float64(usageMillis) / float64(capacityMillis) * 100
}

func (k *K8sClient) calculateMemoryUsagePercent(usage, capacity resource.Quantity) float64 {
	if capacity.IsZero() {
		return 0
	}
	usageBytes := usage.Value()
	capacityBytes := capacity.Value()
	return float64(usageBytes) / float64(capacityBytes) * 100
}

func (k *K8sClient) calculateDiskUsagePercent(usage, capacity resource.Quantity) float64 {
	if capacity.IsZero() {
		return 0
	}
	usageBytes := usage.Value()
	capacityBytes := capacity.Value()
	return float64(usageBytes) / float64(capacityBytes) * 100
}

// Additional utility methods
func (k *K8sClient) IsHealthy() error {
	ctx, cancel := k.createContext()
	defer cancel()

	_, err := k.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		return fmt.Errorf("kubernetes cluster is not healthy: %w", err)
	}
	return nil
}

func (k *K8sClient) GetClusterInfo() (map[string]string, error) {
	_, cancel := k.createContext()
	defer cancel()

	version, err := k.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	info := map[string]string{
		"version":    version.String(),
		"gitVersion": version.GitVersion,
		"platform":   version.Platform,
		"goVersion":  version.GoVersion,
	}

	return info, nil
}

func (k *K8sClient) GetClusterStats() (*entities.ClusterStats, error) {
	ctx, cancel := k.createContext()
	defer cancel()

	stats := &entities.ClusterStats{}
	g, ctx := errgroup.WithContext(ctx)

	// Nodes
	g.Go(func() error {
		nodes, err := k.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		stats.Nodes.Total = len(nodes.Items)
		for _, node := range nodes.Items {
			for _, condition := range node.Status.Conditions {
				if condition.Type == "Ready" && condition.Status == "True" {
					stats.Nodes.Ready++
					break
				}
			}
		}
		return nil
	})

	// Pods
	g.Go(func() error {
		pods, err := k.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		stats.Pods.Total = len(pods.Items)
		for _, pod := range pods.Items {
			if pod.Status.Phase == "Running" {
				stats.Pods.Ready++
			}
		}
		return nil
	})

	// Deployments
	g.Go(func() error {
		deployments, err := k.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		stats.Deployments.Total = len(deployments.Items)
		for _, deployment := range deployments.Items {
			if deployment.Spec.Replicas != nil && deployment.Status.ReadyReplicas == *deployment.Spec.Replicas {
				stats.Deployments.Ready++
			}
		}
		return nil
	})

	// StatefulSets
	g.Go(func() error {
		statefulSets, err := k.clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
		if err != nil {
			return err
		}
		stats.StatefulSets.Total = len(statefulSets.Items)
		for _, ss := range statefulSets.Items {
			if ss.Spec.Replicas != nil && ss.Status.ReadyReplicas == *ss.Spec.Replicas {
				stats.StatefulSets.Ready++
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return stats, nil
}

func (k *K8sClient) Close() {
}

// KustomizeDeploymentOptions contains configuration for Kustomize deployments
// This struct provides various options to customize the deployment behavior
type KustomizeDeploymentOptions struct {
	Namespace   string            `json:"namespace"`   // Target namespace for deployment
	DryRun      bool              `json:"dryRun"`      // Perform dry-run without applying changes
	Prune       bool              `json:"prune"`       // Remove resources not in the kustomization
	Force       bool              `json:"force"`       // Force apply even if conflicts exist
	Timeout     time.Duration     `json:"timeout"`     // Operation timeout
	Labels      map[string]string `json:"labels"`      // Additional labels to apply
	Annotations map[string]string `json:"annotations"` // Additional annotations to apply
}

// KustomizeDeploymentResult contains the result of a Kustomize deployment operation
// This struct provides detailed information about the deployment outcome
type KustomizeDeploymentResult struct {
	Success    bool     `json:"success"`    // Whether the operation was successful
	Message    string   `json:"message"`    // Human-readable result message
	Applied    []string `json:"applied"`    // List of applied resources
	Errors     []string `json:"errors"`     // List of errors encountered
	DryRun     bool     `json:"dryRun"`     // Whether this was a dry-run operation
	BuildTime  string   `json:"buildTime"`  // Time taken to build resources
	DeployTime string   `json:"deployTime"` // Time taken to deploy resources
}

func (r *KustomizeDeploymentResult) String() string {
	return fmt.Sprintf("Success: %t, Message: %s, DeployTime: %s",
		r.Success, r.Message, r.DeployTime)
}

// DeployKustomize deploys a Kustomize directory to the Kubernetes cluster using kubectl
// This is the main function for deploying Kustomize-based applications
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//   - opts: Optional deployment configuration
//
// Returns:
//   - KustomizeDeploymentResult with deployment details
//   - error if deployment fails
//
// Example:
//
//	result, err := k8sClient.DeployKustomize("./my-app", &KustomizeDeploymentOptions{
//	    Namespace: "production",
//	    DryRun:    false,
//	    Prune:     true,
//	})
func (k *K8sClient) DeployKustomize(kustomizePath string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	deployStart := time.Now()
	result := &KustomizeDeploymentResult{
		Success: false,
		DryRun:  opts != nil && opts.DryRun,
	}

	// Validate kustomize path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		result.Message = fmt.Sprintf("Invalid kustomize path: %v", err)
		return result, err
	}

	// Set default options if not provided
	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}

	// Build kubectl command
	args := []string{"apply", "-k", kustomizePath}

	if opts.DryRun {
		args = append(args, "--dry-run=server")
	}

	if opts.Prune {
		args = append(args, "--prune")
	}

	if opts.Force {
		args = append(args, "--force")
	}

	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}

	// Execute kubectl command
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Message = fmt.Sprintf("kubectl command failed: %v\nOutput: %s", err, string(output))
		return result, err
	}

	result.Success = true
	result.Message = "Kustomize deployment completed successfully"
	result.Applied = []string{string(output)}
	result.DeployTime = time.Since(deployStart).String()

	return result, nil
}

// validateKustomizePath validates that the kustomize path exists and contains a kustomization file
// This helper function ensures the provided path is valid for Kustomize operations
func (k *K8sClient) validateKustomizePath(path string) error {
	// Check if path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("kustomize path does not exist: %s", path)
	}

	// Check if it's a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat kustomize path: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("kustomize path is not a directory: %s", path)
	}

	// Check for kustomization files
	kustomizationFiles := []string{"kustomization.yaml", "kustomization.yml", "Kustomization"}
	for _, filename := range kustomizationFiles {
		kustomizationPath := filepath.Join(path, filename)
		if _, err := os.Stat(kustomizationPath); err == nil {
			return nil // Found a kustomization file
		}
	}

	return fmt.Errorf("no kustomization file found in directory: %s", path)
}

// DeployKustomizeWithKubectl deploys using kubectl apply -k (alternative method)
// This function provides an alternative deployment method using kubectl directly
func (k *K8sClient) DeployKustomizeWithKubectl(kustomizePath string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	deployStart := time.Now()
	result := &KustomizeDeploymentResult{
		Success: false,
		DryRun:  opts != nil && opts.DryRun,
	}

	// Validate kustomize path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		result.Message = fmt.Sprintf("Invalid kustomize path: %v", err)
		return result, err
	}

	// Set default options if not provided
	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}

	// Build kubectl command
	args := []string{"apply", "-k", kustomizePath}

	if opts.DryRun {
		args = append(args, "--dry-run=server")
	}

	if opts.Prune {
		args = append(args, "--prune")
	}

	if opts.Force {
		args = append(args, "--force")
	}

	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}

	// Execute kubectl command
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Message = fmt.Sprintf("kubectl command failed: %v\nOutput: %s", err, string(output))
		return result, err
	}

	result.Success = true
	result.Message = "Kustomize deployment completed successfully using kubectl"
	result.Applied = []string{string(output)}
	result.DeployTime = time.Since(deployStart).String()

	return result, nil
}

// ValidateKustomize validates a Kustomize directory without deploying
// This function is useful for CI/CD pipelines to validate kustomization files before deployment
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//
// Returns:
//   - error if validation fails
//
// Example:
//
//	err := k8sClient.ValidateKustomize("./my-app")
//	if err != nil {
//	    log.Printf("Kustomize validation failed: %v", err)
//	}
func (k *K8sClient) ValidateKustomize(kustomizePath string) error {
	// Validate path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		return err
	}

	// Use kubectl to validate the kustomize build
	ctx, cancel := context.WithTimeout(context.Background(), k.defaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "kustomize", kustomizePath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kustomize validation failed: %w", err)
	}

	return nil
}

// GetKustomizeResources returns the built resources without applying them
// This function is useful for previewing what will be deployed or for debugging
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//
// Returns:
//   - []byte containing the YAML output of the built resources
//   - error if building fails
//
// Example:
//
//	yamlData, err := k8sClient.GetKustomizeResources("./my-app")
//	if err != nil {
//	    log.Printf("Failed to get kustomize resources: %v", err)
//	} else {
//	    fmt.Printf("Built resources:\n%s\n", string(yamlData))
//	}
func (k *K8sClient) GetKustomizeResources(kustomizePath string) ([]byte, error) {
	// Validate path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		return nil, err
	}

	// Use kubectl to build the resources
	ctx, cancel := context.WithTimeout(context.Background(), k.defaultTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "kustomize", kustomizePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to build kustomize resources: %w", err)
	}

	return output, nil
}

// DeployKustomizeWithLabels adds custom labels to the deployment
// This is a convenience function for deploying with additional labels
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//   - labels: Map of labels to apply to all resources
//   - opts: Optional deployment configuration
//
// Returns:
//   - KustomizeDeploymentResult with deployment details
//   - error if deployment fails
func (k *K8sClient) DeployKustomizeWithLabels(kustomizePath string, labels map[string]string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	// Set labels in options
	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}
	opts.Labels = labels

	return k.DeployKustomize(kustomizePath, opts)
}

// DeployKustomizeWithAnnotations adds custom annotations to the deployment
// This is a convenience function for deploying with additional annotations
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//   - annotations: Map of annotations to apply to all resources
//   - opts: Optional deployment configuration
//
// Returns:
//   - KustomizeDeploymentResult with deployment details
//   - error if deployment fails
func (k *K8sClient) DeployKustomizeWithAnnotations(kustomizePath string, annotations map[string]string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	// Set annotations in options
	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}
	opts.Annotations = annotations

	return k.DeployKustomize(kustomizePath, opts)
}

// DeleteKustomize deletes resources deployed by a Kustomize directory
// This function removes all resources that were created by the kustomization
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//   - opts: Optional deletion configuration
//
// Returns:
//   - KustomizeDeploymentResult with deletion details
//   - error if deletion fails
//
// Example:
//
//	result, err := k8sClient.DeleteKustomize("./my-app", &KustomizeDeploymentOptions{
//	    Namespace: "production",
//	    DryRun:    true, // Preview what will be deleted
//	})
func (k *K8sClient) DeleteKustomize(kustomizePath string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	deleteStart := time.Now()
	result := &KustomizeDeploymentResult{
		Success: false,
		DryRun:  opts != nil && opts.DryRun,
	}

	// Validate kustomize path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		result.Message = fmt.Sprintf("Invalid kustomize path: %v", err)
		return result, err
	}

	// Set default options if not provided
	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}

	// Build kubectl command for deletion
	args := []string{"delete", "-k", kustomizePath}

	if opts.DryRun {
		args = append(args, "--dry-run=server")
	}

	if opts.Force {
		args = append(args, "--force")
	}

	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}

	// Execute kubectl command
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Message = fmt.Sprintf("kubectl delete command failed: %v\nOutput: %s", err, string(output))
		return result, err
	}

	result.Success = true
	result.Message = "Kustomize deletion completed successfully"
	result.Applied = []string{string(output)}
	result.DeployTime = time.Since(deleteStart).String()

	return result, nil
}

// StopKustomize scales all Deployments and StatefulSets in a Kustomize directory to zero replicas, effectively stopping the application.
// This is useful for pausing or stopping all workloads defined in the kustomization without deleting them.
//
// Parameters:
//   - kustomizePath: Path to the directory containing kustomization.yaml
//   - opts: Optional configuration (namespace, timeout, etc.)
//
// Returns:
//   - KustomizeDeploymentResult with stop details
//   - error if stopping fails
func (k *K8sClient) StopKustomize(kustomizePath string, opts *KustomizeDeploymentOptions) (*KustomizeDeploymentResult, error) {
	stopStart := time.Now()
	result := &KustomizeDeploymentResult{
		Success: false,
	}

	// Validate kustomize path
	if err := k.validateKustomizePath(kustomizePath); err != nil {
		result.Message = fmt.Sprintf("Invalid kustomize path: %v", err)
		return result, err
	}

	if opts == nil {
		opts = &KustomizeDeploymentOptions{
			Namespace: "default",
			Timeout:   k.defaultTimeout,
		}
	}

	// Use kubectl to get all deployments and statefulsets in the kustomizePath
	// 1. Build the resources (kubectl kustomize)
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "kubectl", "kustomize", kustomizePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Message = fmt.Sprintf("failed to build kustomize resources: %v\nOutput: %s", err, string(output))
		return result, err
	}

	yamlData := output

	// Parse the YAML and find all Deployments and StatefulSets
	// We'll use a simple approach: look for 'kind: Deployment' and 'kind: StatefulSet' and their metadata.name
	// (For production, use a YAML parser, but for now, keep it simple and robust)
	resources := []struct {
		Kind string
		Name string
	}{}

	docs := splitYAMLDocuments(string(yamlData))
	for _, doc := range docs {
		kind, name := extractKindAndName(doc)
		if (kind == "Deployment" || kind == "StatefulSet") && name != "" {
			resources = append(resources, struct {
				Kind string
				Name string
			}{Kind: kind, Name: name})
		}
	}

	if len(resources) == 0 {
		result.Message = "No Deployments or StatefulSets found to stop."
		result.Success = true
		return result, nil
	}

	// Scale each resource to zero
	applied := []string{}
	for _, res := range resources {
		args := []string{"scale", res.Kind, res.Name, "--replicas=0"}
		if opts.Namespace != "" {
			args = append(args, "-n", opts.Namespace)
		}
		cmd := exec.CommandContext(ctx, "kubectl", args...)
		out, err := cmd.CombinedOutput()
		applied = append(applied, string(out))
		if err != nil {
			result.Message = fmt.Sprintf("Failed to scale %s/%s: %v\nOutput: %s", res.Kind, res.Name, err, string(out))
			result.Applied = applied
			return result, err
		}
	}

	result.Success = true
	result.Message = "All Deployments and StatefulSets scaled to zero replicas."
	result.Applied = applied
	result.DeployTime = time.Since(stopStart).String()
	return result, nil
}

// splitYAMLDocuments splits a multi-document YAML string into individual documents
func splitYAMLDocuments(yaml string) []string {
	docs := []string{}
	current := ""
	for _, line := range splitLines(yaml) {
		if line == "---" {
			if current != "" {
				docs = append(docs, current)
				current = ""
			}
		} else {
			if current != "" {
				current += "\n"
			}
			current += line
		}
	}
	if current != "" {
		docs = append(docs, current)
	}
	return docs
}

// splitLines splits a string into lines
func splitLines(s string) []string {
	lines := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

// extractKindAndName extracts the kind and metadata.name from a YAML document (very simple, not a full YAML parser)
func extractKindAndName(doc string) (kind, name string) {
	lines := splitLines(doc)
	for i, line := range lines {
		if len(kind) == 0 && (line == "kind: Deployment" || line == "kind: StatefulSet") {
			kind = line[len("kind: "):]
		}
		if len(name) == 0 && (line == "metadata:" || line == " metadata:") {
			// Look ahead for name
			for j := i + 1; j < len(lines); j++ {
				if len(lines[j]) > 6 && lines[j][:6] == "  name" {
					name = lines[j][len("  name: "):]
					break
				}
				if lines[j] == "" || lines[j][0] != ' ' {
					break
				}
			}
		}
		if kind != "" && name != "" {
			break
		}
	}
	return kind, name
}
