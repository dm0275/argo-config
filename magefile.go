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

var (
	argocdHost    = "localhost:8080"
	gitOpsRepo    = "https://github.com/dm0275/argo-gitops.git"
	gitOpsRepoSsh = "git@github.com:dm0275/argo-gitops.git"
	config        = ArgoConfig{
		Version: "v2.11.3", // use `stable` for the latest version
	}
)

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

// Add github ssh cert
func AddKnownHosts() error {
	// Add github ssh cert
	output, err := run("ssh-keyscan github.com | argocd cert add-ssh --batch")
	if err != nil {
		return fmt.Errorf("unable add github ssh cert. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Add Argo repo credentials
func AddRepoCreds() error {
	// Add repocreds
	output, err := run("argocd repocreds add git@github.com --ssh-private-key-path ~/.ssh/id_rsa")
	if err != nil {
		return fmt.Errorf("unable add repocreds. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func AddRepo() error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --server %s", gitOpsRepo, argocdHost))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func AddRepoSsh() error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --ssh-private-key-path ~/.ssh/id_rsa --server %s", gitOpsRepoSsh, argocdHost))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func CreateAppCli() error {
	// Add new app via argocd cli
	output, err := run(fmt.Sprintf("argocd app create app1 --repo %s --path applications/1-directory --dest-server https://kubernetes.default.svc --dest-namespace app1", gitOpsRepoSsh))
	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func CreateAppManifest() error {
	// Add new app via manifest
	output, err := run("kubectl apply -f applications/application.yaml")
	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

func run(command string) (string, error) {
	return sh.Output("bash", "-c", command)
}
