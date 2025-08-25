package main

import (
	"fmt"
	"log"
	"time"

	"homelab.com/homelab-server/homeLab-server/infrastructure/client"
)

func main() {
	// Create a new Kubernetes client
	k8sClient, err := client.NewK8sClient(&client.ClientOptions{
		Timeout: 60 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	defer k8sClient.Close()

	// Example 1: Basic Kustomize deployment
	fmt.Println("=== Example 1: Basic Kustomize Deployment ===")
	basicDeployment(k8sClient)

	// Example 2: Deployment with custom options
	fmt.Println("\n=== Example 2: Deployment with Custom Options ===")
	deploymentWithOptions(k8sClient)

	// Example 3: Dry-run deployment (preview)
	fmt.Println("\n=== Example 3: Dry-run Deployment (Preview) ===")
	dryRunDeployment(k8sClient)

	// Example 4: Deployment with labels and annotations
	fmt.Println("\n=== Example 4: Deployment with Labels and Annotations ===")
	deploymentWithLabelsAndAnnotations(k8sClient)

	// Example 5: Validation
	fmt.Println("\n=== Example 5: Kustomize Validation ===")
	validateKustomize(k8sClient)

	// Example 6: Get resources without deploying
	fmt.Println("\n=== Example 6: Get Resources Without Deploying ===")
	getResources(k8sClient)

	// Example 7: Delete deployment
	fmt.Println("\n=== Example 7: Delete Deployment ===")
	deleteDeployment(k8sClient)
}

func basicDeployment(k8sClient *client.K8sClient) {
	// Deploy a simple application
	result, err := k8sClient.DeployKustomize("./examples/simple-app", nil)
	if err != nil {
		log.Printf("Basic deployment failed: %v", err)
		return
	}

	fmt.Printf("Deployment result: %+v\n", result)
	if result.Success {
		fmt.Println("✅ Basic deployment completed successfully")
	} else {
		fmt.Printf("❌ Basic deployment failed: %s\n", result.Message)
	}
}

func deploymentWithOptions(k8sClient *client.K8sClient) {
	// Deploy with custom options
	opts := &client.KustomizeDeploymentOptions{
		Namespace: "production",
		DryRun:    false,
		Prune:     true,
		Force:     false,
		Timeout:   120 * time.Second,
	}

	result, err := k8sClient.DeployKustomize("./examples/production-app", opts)
	if err != nil {
		log.Printf("Deployment with options failed: %v", err)
		return
	}

	fmt.Printf("Deployment result: %+v\n", result)
	if result.Success {
		fmt.Println("✅ Deployment with options completed successfully")
		fmt.Printf("⏱️  Deployment time: %s\n", result.DeployTime)
	} else {
		fmt.Printf("❌ Deployment with options failed: %s\n", result.Message)
	}
}

func dryRunDeployment(k8sClient *client.K8sClient) {
	// Perform a dry-run to preview changes
	opts := &client.KustomizeDeploymentOptions{
		Namespace: "staging",
		DryRun:    true,
		Prune:     true,
	}

	result, err := k8sClient.DeployKustomize("./examples/staging-app", opts)
	if err != nil {
		log.Printf("Dry-run deployment failed: %v", err)
		return
	}

	fmt.Printf("Dry-run result: %+v\n", result)
	if result.Success {
		fmt.Println("✅ Dry-run completed successfully")
		fmt.Println("📋 Preview of changes:")
		for _, applied := range result.Applied {
			fmt.Printf("   %s\n", applied)
		}
	} else {
		fmt.Printf("❌ Dry-run failed: %s\n", result.Message)
	}
}

func deploymentWithLabelsAndAnnotations(k8sClient *client.K8sClient) {
	// Deploy with custom labels
	labels := map[string]string{
		"app.kubernetes.io/name":    "my-app",
		"app.kubernetes.io/version": "v1.0.0",
		"app.kubernetes.io/part-of": "home-lab",
		"environment":               "production",
	}

	annotations := map[string]string{
		"deployment.kubernetes.io/revision":         "1",
		"kustomize.config.k8s.io/needs-hash-suffix": "true",
	}

	opts := &client.KustomizeDeploymentOptions{
		Namespace:   "production",
		Labels:      labels,
		Annotations: annotations,
	}

	result, err := k8sClient.DeployKustomize("./examples/labeled-app", opts)
	if err != nil {
		log.Printf("Deployment with labels and annotations failed: %v", err)
		return
	}

	fmt.Printf("Deployment result: %+v\n", result)
	if result.Success {
		fmt.Println("✅ Deployment with labels and annotations completed successfully")
	} else {
		fmt.Printf("❌ Deployment with labels and annotations failed: %s\n", result.Message)
	}
}

func validateKustomize(k8sClient *client.K8sClient) {
	// Validate a Kustomize directory
	err := k8sClient.ValidateKustomize("./examples/valid-app")
	if err != nil {
		fmt.Printf("❌ Kustomize validation failed: %v\n", err)
		return
	}

	fmt.Println("✅ Kustomize validation passed")

	// Try to validate an invalid directory
	err = k8sClient.ValidateKustomize("./examples/invalid-app")
	if err != nil {
		fmt.Printf("✅ Correctly detected invalid kustomize: %v\n", err)
	} else {
		fmt.Println("❌ Should have detected invalid kustomize")
	}
}

func getResources(k8sClient *client.K8sClient) {
	// Get the built resources without deploying
	yamlData, err := k8sClient.GetKustomizeResources("./examples/simple-app")
	if err != nil {
		log.Printf("Failed to get kustomize resources: %v", err)
		return
	}

	fmt.Printf("✅ Retrieved kustomize resources (%d bytes)\n", len(yamlData))
	fmt.Println("📄 Resource preview:")
	fmt.Println("---")
	fmt.Println(string(yamlData[:min(len(yamlData), 500)])) // Show first 500 characters
	if len(yamlData) > 500 {
		fmt.Println("... (truncated)")
	}
}

func deleteDeployment(k8sClient *client.K8sClient) {
	// First, do a dry-run to see what will be deleted
	opts := &client.KustomizeDeploymentOptions{
		Namespace: "production",
		DryRun:    true,
	}

	result, err := k8sClient.DeleteKustomize("./examples/production-app", opts)
	if err != nil {
		log.Printf("Dry-run deletion failed: %v", err)
		return
	}

	fmt.Printf("Dry-run deletion result: %+v\n", result)
	if result.Success {
		fmt.Println("✅ Dry-run deletion completed successfully")
		fmt.Println("📋 Preview of resources to be deleted:")
		for _, applied := range result.Applied {
			fmt.Printf("   %s\n", applied)
		}

		// Now perform the actual deletion
		opts.DryRun = false
		result, err = k8sClient.DeleteKustomize("./examples/production-app", opts)
		if err != nil {
			log.Printf("Actual deletion failed: %v", err)
			return
		}

		if result.Success {
			fmt.Println("✅ Actual deletion completed successfully")
		} else {
			fmt.Printf("❌ Actual deletion failed: %s\n", result.Message)
		}
	} else {
		fmt.Printf("❌ Dry-run deletion failed: %s\n", result.Message)
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
