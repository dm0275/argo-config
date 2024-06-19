.PHONY: help
.DEFAULT_GOAL := help

ARGOCD_HOST := localhost:8080
GITOPS_REPO := https://github.com/dm0275/argo-gitops.git
GITOPS_REPO_SSH := git@github.com:dm0275/argo-gitops.git

install-argocd: ## Install Argo
	# Create the ArgoCD namespace
	kubectl create namespace argocd
	# Deploy Argo on the cluster
	# Use this URL for the latest https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.11.3/manifests/install.yaml

port-forward: ## Port-Forward Argo UI
	@echo "Argo can be accessed at:"
	@echo "https://localhost:8080"
	kubectl port-forward svc/argocd-server -n argocd 8080:443

get-admin-password: ## Get Admin password
	@argocd admin initial-password -n argocd | head -n 1

login: ## Login to Argo
	@export argo_pass=$$(argocd admin initial-password -n argocd | head -n 1) ;\
    argocd login --username admin --password $$argo_pass --insecure localhost:8080

add-repo: ## Add repository
	@argocd repo add $(GITOPS_REPO) --server $(ARGOCD_HOST)

add-repo-ssh: ## Add repository
	@argocd repo add $(GITOPS_REPO_SSH) --ssh-private-key-path ~/.ssh/id_rsa --server $(ARGOCD_HOST) # --insecure --insecure-ignore-host-key

create-app-cli: ## Create application via CLI
	@argocd app create app1 --repo $(GITOPS_REPO) --path applications/1-directory --dest-server https://kubernetes.default.svc --dest-namespace app1

add-known-hosts:
	ssh-keyscan github.com | argocd cert add-ssh --batch

add-repo-creds:
	argocd repocreds add git@github.com --ssh-private-key-path ~/.ssh/id_rsa

create-app-manifest: ## Create application via manifest
	@kubectl apply -f applications/application.yaml

sync-app:
	@argocd app sync app1

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'