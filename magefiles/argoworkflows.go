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
	workflowsRepo    = "https://github.com/dm0275/argo-gitops.git"
	workflowsRepoSsh = "git@github.com:dm0275/argo-gitops.git"
	ArgoWFConfig     = ArgoWorkflowsConfig{
		Namespace:       "argo",
		Version:         "v3.5.10", // use `stable` for the latest version
		PortForwardPort: "2746",
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

// Start the Argo Server
func (ArgoWorkflows) ArgoServer() error {
	fmt.Println(fmt.Sprintf("Argo can be accessed at:\nhttps://localhost:%s", ArgoWFConfig.PortForwardPort))
	// Port forward the argo-server
	_, err := run("argo server --auth-mode=server")
	if err != nil {
		return fmt.Errorf("unable to start the argo-server. ERROR: %s", err)
	}

	return nil
}
