package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// HistoricalMetrics represents metrics data over time
type HistoricalMetrics struct {
	NodeName    string      `json:"nodeName"`
	Timestamps  []time.Time `json:"timestamps"`
	CPUUsage    []float64   `json:"cpuUsage"`    // Percentage
	MemoryUsage []float64   `json:"memoryUsage"` // Percentage
	DiskUsage   []float64   `json:"diskUsage"`   // Percentage
	NetworkRx   []float64   `json:"networkRx"`   // Bytes per second
	NetworkTx   []float64   `json:"networkTx"`   // Bytes per second
}

// ClusterHistoricalMetrics represents cluster-wide metrics data over time
type ClusterHistoricalMetrics struct {
	ClusterName string      `json:"clusterName"`
	Timestamps  []time.Time `json:"timestamps"`
	AvgCPU      []float64   `json:"avgCpu"`     // Average CPU across all nodes (%)
	AvgMemory   []float64   `json:"avgMemory"`  // Average Memory across all nodes (%)
	AvgDisk     []float64   `json:"avgDisk"`    // Average Disk across all nodes (%)
	TotalNetRx  []float64   `json:"totalNetRx"` // Total Network RX across cluster (bytes/sec)
	TotalNetTx  []float64   `json:"totalNetTx"` // Total Network TX across cluster (bytes/sec)
	NodeCount   []int       `json:"nodeCount"`  // Number of active nodes
}

// PrometheusMetricsClient handles historical metrics via Prometheus
type PrometheusMetricsClient struct {
	baseURL    string
	httpClient *http.Client
}

// Method 1: Prometheus Integration (Recommended)
func NewPrometheusMetricsClient(prometheusURL string) *PrometheusMetricsClient {
	return &PrometheusMetricsClient{
		baseURL: prometheusURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetClusterHistoricalMetrics gets cluster-wide average metrics over time
func (p *PrometheusMetricsClient) GetClusterHistoricalMetrics(duration time.Duration) (*ClusterHistoricalMetrics, error) {
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	step := 10 * time.Minute

	metrics := &ClusterHistoricalMetrics{
		ClusterName: "kubernetes-cluster",
		Timestamps:  []time.Time{},
	}

	// Get average CPU usage across all nodes (percentage)
	avgCpuQuery := `avg(100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))`
	cpuData, err := p.queryRange(avgCpuQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster CPU metrics: %w", err)
	}
	metrics.AvgCPU = cpuData.Values
	if len(cpuData.Timestamps) > 0 {
		metrics.Timestamps = cpuData.Timestamps
	}

	// Get average Memory usage across all nodes (percentage)
	avgMemQuery := `avg(100 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100))`
	memData, err := p.queryRange(avgMemQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster memory metrics: %w", err)
	}
	metrics.AvgMemory = memData.Values

	// Get average Disk usage across all nodes (percentage)
	avgDiskQuery := `avg(100 - (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"} * 100))`
	diskData, err := p.queryRange(avgDiskQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster disk metrics: %w", err)
	}
	metrics.AvgDisk = diskData.Values

	// Get total Network RX across all nodes (bytes per second)
	totalNetRxQuery := `sum(irate(node_network_receive_bytes_total{device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m]))`
	netRxData, err := p.queryRange(totalNetRxQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster network RX metrics: %w", err)
	}
	metrics.TotalNetRx = netRxData.Values

	// Get total Network TX across all nodes (bytes per second)
	totalNetTxQuery := `sum(irate(node_network_transmit_bytes_total{device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m]))`
	netTxData, err := p.queryRange(totalNetTxQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster network TX metrics: %w", err)
	}
	metrics.TotalNetTx = netTxData.Values

	// Get node count over time
	nodeCountQuery := `count(up{job="node-exporter"})`
	nodeData, err := p.queryRange(nodeCountQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get node count metrics: %w", err)
	}
	// Convert float64 to int for node count
	for _, val := range nodeData.Values {
		metrics.NodeCount = append(metrics.NodeCount, int(val))
	}

	return metrics, nil
}

func (p *PrometheusMetricsClient) GetNodeHistoricalMetrics(nodeName string, duration time.Duration) (*HistoricalMetrics, error) {
	endTime := time.Now()
	startTime := endTime.Add(-duration)
	step := 10 * time.Minute

	metrics := &HistoricalMetrics{
		NodeName:   nodeName,
		Timestamps: []time.Time{},
	}

	// Get CPU usage (percentage)
	cpuQuery := fmt.Sprintf(`100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle",instance=~".*%s.*"}[5m])) * 100)`, nodeName)
	cpuData, err := p.queryRange(cpuQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU metrics: %w", err)
	}
	metrics.CPUUsage = cpuData.Values
	if len(cpuData.Timestamps) > 0 {
		metrics.Timestamps = cpuData.Timestamps
	}

	// Get Memory usage (percentage)
	memQuery := fmt.Sprintf(`(1 - (node_memory_MemAvailable_bytes{instance=~".*%s.*"} / node_memory_MemTotal_bytes{instance=~".*%s.*"})) * 100`, nodeName, nodeName)
	memData, err := p.queryRange(memQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory metrics: %w", err)
	}
	metrics.MemoryUsage = memData.Values

	// Get Disk usage (percentage)
	diskQuery := fmt.Sprintf(`(1 - (node_filesystem_avail_bytes{instance=~".*%s.*",mountpoint="/"} / node_filesystem_size_bytes{instance=~".*%s.*",mountpoint="/"})) * 100`, nodeName, nodeName)
	diskData, err := p.queryRange(diskQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk metrics: %w", err)
	}
	metrics.DiskUsage = diskData.Values

	// Get Network RX (bytes per second)
	netRxQuery := fmt.Sprintf(`irate(node_network_receive_bytes_total{instance=~".*%s.*",device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m])`, nodeName)
	netRxData, err := p.queryRange(netRxQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get network RX metrics: %w", err)
	}
	metrics.NetworkRx = netRxData.Values

	// Get Network TX (bytes per second)
	netTxQuery := fmt.Sprintf(`irate(node_network_transmit_bytes_total{instance=~".*%s.*",device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m])`, nodeName)
	netTxData, err := p.queryRange(netTxQuery, startTime, endTime, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get network TX metrics: %w", err)
	}
	metrics.NetworkTx = netTxData.Values

	return metrics, nil
}

type PrometheusQueryResult struct {
	Timestamps []time.Time
	Values     []float64
}

// GetCurrentClusterMetrics gets real-time cluster metrics snapshot
func (p *PrometheusMetricsClient) GetCurrentClusterMetrics() (*ClusterSnapshot, error) {
	now := time.Now()

	snapshot := &ClusterSnapshot{
		Timestamp: now,
	}

	// Current average CPU across cluster
	avgCpuQuery := `avg(100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100))`
	cpuResult, err := p.queryInstant(avgCpuQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current cluster CPU: %w", err)
	}
	snapshot.AvgCPU = cpuResult

	// Current average Memory across cluster
	avgMemQuery := `avg(100 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100))`
	memResult, err := p.queryInstant(avgMemQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current cluster memory: %w", err)
	}
	snapshot.AvgMemory = memResult

	// Current average Disk across cluster
	avgDiskQuery := `avg(100 - (node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"} * 100))`
	diskResult, err := p.queryInstant(avgDiskQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current cluster disk: %w", err)
	}
	snapshot.AvgDisk = diskResult

	// Current total network traffic
	totalNetRxQuery := `sum(irate(node_network_receive_bytes_total{device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m]))`
	netRxResult, err := p.queryInstant(totalNetRxQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current cluster network RX: %w", err)
	}
	snapshot.TotalNetRx = netRxResult

	totalNetTxQuery := `sum(irate(node_network_transmit_bytes_total{device!~"lo|veth.*|docker.*|flannel.*|cali.*|br-.*"}[5m]))`
	netTxResult, err := p.queryInstant(totalNetTxQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current cluster network TX: %w", err)
	}
	snapshot.TotalNetTx = netTxResult

	// Current node count
	nodeCountQuery := `count(up{job="node-exporter"})`
	nodeCountResult, err := p.queryInstant(nodeCountQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get current node count: %w", err)
	}
	snapshot.NodeCount = int(nodeCountResult)

	return snapshot, nil
}

type ClusterSnapshot struct {
	Timestamp  time.Time `json:"timestamp"`
	AvgCPU     float64   `json:"avgCpu"`     // Average CPU across cluster (%)
	AvgMemory  float64   `json:"avgMemory"`  // Average Memory across cluster (%)
	AvgDisk    float64   `json:"avgDisk"`    // Average Disk across cluster (%)
	TotalNetRx float64   `json:"totalNetRx"` // Total Network RX (bytes/sec)
	TotalNetTx float64   `json:"totalNetTx"` // Total Network TX (bytes/sec)
	NodeCount  int       `json:"nodeCount"`  // Number of active nodes
}

func (p *PrometheusMetricsClient) queryInstant(query string) (float64, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("time", strconv.FormatInt(time.Now().Unix(), 10))

	url := fmt.Sprintf("%s/api/v1/query?%s", p.baseURL, params.Encode())

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to query Prometheus: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response: %w", err)
	}

	var promResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &promResp); err != nil {
		return 0, fmt.Errorf("failed to parse Prometheus response: %w", err)
	}

	if promResp.Status != "success" {
		return 0, fmt.Errorf("Prometheus query failed")
	}

	if len(promResp.Data.Result) > 0 && len(promResp.Data.Result[0].Value) >= 2 {
		if val, ok := promResp.Data.Result[0].Value[1].(string); ok {
			if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
				return floatVal, nil
			}
		}
	}

	return 0, nil
}

func (p *PrometheusMetricsClient) queryRange(query string, start, end time.Time, step time.Duration) (*PrometheusQueryResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.Unix(), 10))
	params.Set("end", strconv.FormatInt(end.Unix(), 10))
	params.Set("step", fmt.Sprintf("%.0fs", step.Seconds()))

	url := fmt.Sprintf("%s/api/v1/query_range?%s", p.baseURL, params.Encode())

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var promResp struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"`
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &promResp); err != nil {
		return nil, fmt.Errorf("failed to parse Prometheus response: %w", err)
	}

	if promResp.Status != "success" {
		return nil, fmt.Errorf("Prometheus query failed")
	}

	result := &PrometheusQueryResult{}

	if len(promResp.Data.Result) > 0 {
		values := promResp.Data.Result[0].Values
		for _, value := range values {
			if len(value) >= 2 {
				// Parse timestamp
				if timestamp, ok := value[0].(float64); ok {
					result.Timestamps = append(result.Timestamps, time.Unix(int64(timestamp), 0))
				}
				// Parse value
				if val, ok := value[1].(string); ok {
					if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
						result.Values = append(result.Values, floatVal)
					}
				}
			}
		}
	}

	return result, nil
}

// Method 2: Enhanced K8sClient with Historical Tracking
func (k *K8sClient) StartMetricsCollection(interval time.Duration) *MetricsCollector {
	collector := &MetricsCollector{
		client:      k,
		interval:    interval,
		stopChan:    make(chan struct{}),
		metricsData: make(map[string]*HistoricalMetrics),
	}

	go collector.collectLoop()
	return collector
}

type MetricsCollector struct {
	client      *K8sClient
	interval    time.Duration
	stopChan    chan struct{}
	metricsData map[string]*HistoricalMetrics
}

func (mc *MetricsCollector) collectLoop() {
	ticker := time.NewTicker(mc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.collectCurrentMetrics()
		case <-mc.stopChan:
			return
		}
	}
}

func (mc *MetricsCollector) collectCurrentMetrics() {
	nodes, err := mc.client.ListNodes()
	if err != nil {
		fmt.Printf("Failed to list nodes for metrics collection: %v\n", err)
		return
	}

	now := time.Now()
	for _, node := range nodes {
		stats, err := mc.client.GetNodeStats(node.Name)
		if err != nil {
			fmt.Printf("Failed to get stats for node %s: %v\n", node.Name, err)
			continue
		}

		if _, exists := mc.metricsData[node.Name]; !exists {
			mc.metricsData[node.Name] = &HistoricalMetrics{
				NodeName: node.Name,
			}
		}

		metrics := mc.metricsData[node.Name]
		metrics.Timestamps = append(metrics.Timestamps, now)
		metrics.CPUUsage = append(metrics.CPUUsage, stats.CPUUsage)
		metrics.MemoryUsage = append(metrics.MemoryUsage, stats.MemoryUsage)
		metrics.DiskUsage = append(metrics.DiskUsage, stats.DiskUsage)

		// Keep only last 24 hours of data
		cutoff := now.Add(-24 * time.Hour)
		mc.trimOldData(metrics, cutoff)
	}
}

func (mc *MetricsCollector) trimOldData(metrics *HistoricalMetrics, cutoff time.Time) {
	// Find first index after cutoff
	startIdx := 0
	for i, timestamp := range metrics.Timestamps {
		if timestamp.After(cutoff) {
			startIdx = i
			break
		}
	}

	if startIdx > 0 {
		metrics.Timestamps = metrics.Timestamps[startIdx:]
		metrics.CPUUsage = metrics.CPUUsage[startIdx:]
		metrics.MemoryUsage = metrics.MemoryUsage[startIdx:]
		metrics.DiskUsage = metrics.DiskUsage[startIdx:]
		if len(metrics.NetworkRx) > startIdx {
			metrics.NetworkRx = metrics.NetworkRx[startIdx:]
		}
		if len(metrics.NetworkTx) > startIdx {
			metrics.NetworkTx = metrics.NetworkTx[startIdx:]
		}
	}
}

func (mc *MetricsCollector) GetHistoricalMetrics(nodeName string) *HistoricalMetrics {
	if metrics, exists := mc.metricsData[nodeName]; exists {
		return metrics
	}
	return nil
}

func (mc *MetricsCollector) Stop() {
	close(mc.stopChan)
}

// Method 3: Integration with existing monitoring stack
func (k *K8sClient) GetHistoricalMetricsFromGrafana(grafanaURL, apiKey, nodeName string) (*HistoricalMetrics, error) {
	// This would integrate with Grafana's API to fetch historical data
	// Implementation depends on your Grafana setup and dashboard configuration

	client := &http.Client{Timeout: 30 * time.Second}

	// Example Grafana API call (customize based on your setup)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/datasources/proxy/1/api/v1/query_range", grafanaURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// Add query parameters for your specific metrics
	q := req.URL.Query()
	q.Add("query", fmt.Sprintf(`node_cpu_usage{instance="%s"}`, nodeName))
	q.Add("start", strconv.FormatInt(time.Now().Add(-24*time.Hour).Unix(), 10))
	q.Add("end", strconv.FormatInt(time.Now().Unix(), 10))
	q.Add("step", "600") // 10 minutes
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response and return HistoricalMetrics
	// Implementation depends on your Grafana response format

	return &HistoricalMetrics{NodeName: nodeName}, nil
}

// Utility functions for cluster metrics analysis
func (chm *ClusterHistoricalMetrics) GetClusterAverages() map[string]float64 {
	return map[string]float64{
		"avgCPU":       calculateAverage(chm.AvgCPU),
		"avgMemory":    calculateAverage(chm.AvgMemory),
		"avgDisk":      calculateAverage(chm.AvgDisk),
		"avgNetworkRx": calculateAverage(chm.TotalNetRx),
		"avgNetworkTx": calculateAverage(chm.TotalNetTx),
		"avgNodeCount": calculateAverageInt(chm.NodeCount),
	}
}

func (chm *ClusterHistoricalMetrics) GetClusterPeaks() map[string]float64 {
	return map[string]float64{
		"peakCPU":       findMax(chm.AvgCPU),
		"peakMemory":    findMax(chm.AvgMemory),
		"peakDisk":      findMax(chm.AvgDisk),
		"peakNetworkRx": findMax(chm.TotalNetRx),
		"peakNetworkTx": findMax(chm.TotalNetTx),
		"maxNodeCount":  float64(findMaxInt(chm.NodeCount)),
	}
}

func findMax(f []float64) float64 {
	max := f[0]
	for _, v := range f {
		if v > max {
			max = v
		}
	}
	return max
}

func (chm *ClusterHistoricalMetrics) GetNetworkTotalGB() map[string]float64 {
	totalRxBytes := 0.0
	totalTxBytes := 0.0

	// Calculate total bytes transferred (sum of all intervals * step duration)
	stepDurationSeconds := 600.0 // 10 minutes = 600 seconds

	for _, rx := range chm.TotalNetRx {
		totalRxBytes += rx * stepDurationSeconds
	}

	for _, tx := range chm.TotalNetTx {
		totalTxBytes += tx * stepDurationSeconds
	}

	return map[string]float64{
		"totalRxGB": totalRxBytes / (1024 * 1024 * 1024), // Convert to GB
		"totalTxGB": totalTxBytes / (1024 * 1024 * 1024), // Convert to GB
	}
}

// Utility functions for metrics analysis
func (hm *HistoricalMetrics) GetAverageMetrics() map[string]float64 {
	return map[string]float64{
		"avgCPU":       calculateAverage(hm.CPUUsage),
		"avgMemory":    calculateAverage(hm.MemoryUsage),
		"avgDisk":      calculateAverage(hm.DiskUsage),
		"avgNetworkRx": calculateAverage(hm.NetworkRx),
		"avgNetworkTx": calculateAverage(hm.NetworkTx),
	}
}

func (hm *HistoricalMetrics) GetPeakMetrics() map[string]float64 {
	return map[string]float64{
		"peakCPU":       findMax(hm.CPUUsage),
		"peakMemory":    findMax(hm.MemoryUsage),
		"peakDisk":      findMax(hm.DiskUsage),
		"peakNetworkRx": findMax(hm.NetworkRx),
		"peakNetworkTx": findMax(hm.NetworkTx),
	}
}

func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateAverageInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}

func findMaxInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
