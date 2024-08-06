//go:build mage
// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
)

type ArgoCD mg.Namespace

type ArgoCdConfig struct {
	Namespace       string
	Version         string
	PortForwardPort string
	SSHKeyPath      string
}

var (
	argocdHost          = "localhost"
	gitOpsRepo          = "https://github.com/dm0275/argo-gitops.git"
	gitOpsRepoSsh       = "git@github.com:dm0275/argo-gitops.git"
	applicationManifest = "applications/application.yaml"
	ArgoCDConfig        = ArgoCdConfig{
		Namespace:       "argocd",
		Version:         "v2.11.3", // use `stable` for the latest version
		SSHKeyPath:      "~/.ssh/id_rsa",
		PortForwardPort: "8080",
	}
)

// Creates the argocd namespace and installs argocd
func (ArgoCD) Install() error {
	// Create the ArgoCD namespace
	output, err := createNamespace(ArgoCDConfig.Namespace)
	if err != nil {
		return fmt.Errorf("unable to create argocd namespace. ERROR: %s", err)
	}
	fmt.Println(output)

	// Deploy Argo on the cluster
	output, err = run(fmt.Sprintf("kubectl apply -n %s -f https://raw.githubusercontent.com/argoproj/argo-cd/%s/manifests/install.yaml", ArgoCDConfig.Namespace, ArgoCDConfig.Version))
	if err != nil {
		return fmt.Errorf("unable to deploy argocd. ERROR: %s", err)
	}
	fmt.Println(output)

	return nil
}

// Port-forward the argocd gitops server
func (ArgoCD) PortForward() error {
	fmt.Println(fmt.Sprintf("Argo can be accessed at:\nhttps://localhost:%s", ArgoCDConfig.PortForwardPort))
	// Port forward the argo-server
	_, err := run(fmt.Sprintf("kubectl port-forward svc/argocd-server -n %s %s:443", ArgoCDConfig.Namespace, ArgoCDConfig.PortForwardPort))
	if err != nil {
		return fmt.Errorf("unable to port-forward svc/argocd-server. ERROR: %s", err)
	}

	return nil
}

// Get the initial ArgoCD admin password
func (ArgoCD) GetAdminPassword() error {
	// Fetching admin password
	output, err := run(fmt.Sprintf("argocd admin initial-password -n %s | head -n 1", ArgoCDConfig.Namespace))
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Login to argo via the cli (requires the argocd service to be accessible)
func (ArgoCD) Login() error {
	// Fetching admin password
	adminPass, err := run("argocd admin initial-password -n argocd | head -n 1")
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	// Running argocd login using admin pass
	output, err := run(fmt.Sprintf("argocd login --username admin --password %s --insecure localhost:%s", adminPass, ArgoCDConfig.PortForwardPort))
	if err != nil {
		return fmt.Errorf("unable to login. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Add github ssh cert
func (ArgoCD) AddKnownHosts() error {
	// Add github ssh cert
	output, err := run("ssh-keyscan github.com | argocd cert add-ssh --batch")
	if err != nil {
		return fmt.Errorf("unable add github ssh cert. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Add Argo repo credentials
func (ArgoCD) AddRepoCreds() error {
	// Add repocreds
	output, err := run(fmt.Sprintf("argocd repocreds add git@github.com --ssh-private-key-path %s", ArgoCDConfig.SSHKeyPath))
	if err != nil {
		return fmt.Errorf("unable add repocreds. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Add HTTP repository to Argo
func (ArgoCD) AddRepo() error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --server %s", gitOpsRepo, argocdHost))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Add SSH repository to Argo
func (ArgoCD) AddRepoSsh() error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --ssh-private-key-path %s --server %s", gitOpsRepoSsh, ArgoCDConfig.SSHKeyPath, argocdHost))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Created new application via argocd cli
func (ArgoCD) CreateAppCli() error {
	// Add new app via argocd cli
	output, err := run(fmt.Sprintf("argocd app create app1 --repo %s --path applications/1-directory --dest-server https://kubernetes.default.svc --dest-namespace app1", gitOpsRepoSsh))
	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Created new application via manifest
func (ArgoCD) CreateAppManifest() error {
	// Add new app via manifest
	output, err := run(fmt.Sprintf("kubectl apply -f %s", applicationManifest))
	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}
