//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
)

type ArgoConfig struct {
	Version string
}

var config = ArgoConfig{
	Version: "v2.11.3", // use `stable` for the latest version
}

// Creates the argocd namespace and installs argocd
func InstallArgo() error {
	// Create the ArgoCD namespace
	output, err := run("kubectl create namespace argocd")
	if err != nil {
		return fmt.Errorf("unable to create argocd namespace. ERROR: %s", err)
	}
	fmt.Println(output)

	// Deploy Argo on the cluster
	output, err = run(fmt.Sprintf("kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/%s/manifests/install.yaml", config.Version))
	if err != nil {
		return fmt.Errorf("unable to deploy argocd. ERROR: %s", err)
	}
	fmt.Println(output)

	return nil
}

// Port-forward the argocd server
func PortForward() error {
	fmt.Println("Argo can be accessed at:\nhttps://localhost:8080")
	// Port forward the argo-server
	_, err := run("kubectl port-forward svc/argocd-server -n argocd 8080:443")
	if err != nil {
		return fmt.Errorf("unable to port-forward svc/argocd-server. ERROR: %s", err)
	}

	return nil
}

// Get the initial admin password
func GetAdminPassword() error {
	// Fetching admin password
	output, err := run("argocd admin initial-password -n argocd | head -n 1")
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Login to argo via the cli (requires the argocd service to be accessible)
func ArgoLogin() error {
	// Fetching admin password
	adminPass, err := run("argocd admin initial-password -n argocd | head -n 1")
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	// Running argocd login using admin pass
	output, err := run(fmt.Sprintf("argocd login --username admin --password %s --insecure localhost:8080", adminPass))
	if err != nil {
		return fmt.Errorf("unable to login. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func run(command string) (string, error) {
	return sh.Output("bash", "-c", command)
}
