//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
)

type ArgoWorkflows mg.Namespace

type ArgoWorkflowsConfig struct {
	Namespace       string
	Version         string
	PortForwardPort string
}

var (
	ArgoWFConfig = ArgoWorkflowsConfig{
		Namespace:       "argocd",
		Version:         "v2.11.3", // use `stable` for the latest version
		PortForwardPort: "8080",
	}
)

// Install argocd workflows
func (ArgoWorkflows) Install() error {
	// Create the ArgoCD namespace
	output, err := createNamespace(ArgoWFConfig.Namespace)
	if err != nil {
		return fmt.Errorf("unable to create argocd namespace. ERROR: %s", err)
	}
	fmt.Println(output)

	// Deploy Argo on the cluster
	output, err = run(fmt.Sprintf("kubectl apply -n %s -f https://github.com/argoproj/argo-workflows/releases/download/%s/install.yaml", ArgoWFConfig.Namespace, ArgoWFConfig.Version))
	if err != nil {
		return fmt.Errorf("unable to deploy argocd. ERROR: %s", err)
	}
	fmt.Println(output)

	return nil
}

// Port-forward the argocd workflows server
func (ArgoWorkflows) PortForward() error {
	fmt.Println(fmt.Sprintf("Argo can be accessed at:\nhttps://localhost:%s", ArgoWFConfig.PortForwardPort))
	// Port forward the argo-server
	_, err := run(fmt.Sprintf("kubectl port-forward svc/argo-server -n %s %s:%s", ArgoWFConfig.Namespace, ArgoWFConfig.PortForwardPort, ArgoWFConfig.PortForwardPort))
	if err != nil {
		return fmt.Errorf("unable to port-forward svc/argo-server. ERROR: %s", err)
	}

	return nil
}
