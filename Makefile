# kq Development Makefile

CLUSTER_NAME := kq-dev
BINARY := kq
PLUGIN_BINARY := kubectl-kq
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/edgarvarela24/kq/cmd.Version=$(VERSION)"

# === Build ===

.PHONY: build
build: ## Build the kq binary
	go build $(LDFLAGS) -o $(BINARY) .

.PHONY: build-plugin
build-plugin: ## Build as kubectl plugin (kubectl-kq)
	go build $(LDFLAGS) -o $(PLUGIN_BINARY) .

.PHONY: install-plugin
install-plugin: build-plugin ## Install kubectl-kq to ~/bin
	mkdir -p ~/bin
	cp $(PLUGIN_BINARY) ~/bin/
	@echo "Installed $(PLUGIN_BINARY) to ~/bin/"
	@echo "Make sure ~/bin is in your PATH"

.PHONY: test
test: ## Run all tests
	go test ./... -v

# === Development Cluster ===

.PHONY: cluster-create
cluster-create: ## Create a kind cluster with sample workloads
	kind create cluster --name $(CLUSTER_NAME)
	@echo "Waiting for cluster to be ready..."
	kubectl wait --for=condition=Ready nodes --all --timeout=60s --context kind-$(CLUSTER_NAME)
	@$(MAKE) cluster-populate

.PHONY: cluster-populate
cluster-populate: ## Add sample workloads to the cluster
	kubectl create deployment nginx --image=nginx --replicas=3 --context kind-$(CLUSTER_NAME)
	kubectl create deployment redis --image=redis --replicas=2 --context kind-$(CLUSTER_NAME)
	kubectl create deployment busybox --image=busybox --context kind-$(CLUSTER_NAME) -- sleep 3600
	kubectl create namespace staging --context kind-$(CLUSTER_NAME)
	kubectl create deployment web --image=nginx -n staging --context kind-$(CLUSTER_NAME)
	@echo "Sample workloads created!"

.PHONY: cluster-delete
cluster-delete: ## Delete the kind cluster
	kind delete cluster --name $(CLUSTER_NAME)

.PHONY: cluster-restart
cluster-restart: ## Restart the cluster (delete + create)
	@$(MAKE) cluster-delete || true
	@$(MAKE) cluster-create

.PHONY: cluster-status
cluster-status: ## Show cluster status and pods
	@kubectl cluster-info --context kind-$(CLUSTER_NAME) 2>/dev/null || echo "Cluster not running"
	@echo ""
	@kubectl get pods -A --context kind-$(CLUSTER_NAME) 2>/dev/null || true

# === Helpers ===

.PHONY: run
run: build ## Build and run kq
	./$(BINARY)

.PHONY: clean
clean: ## Remove built binary
	rm -f $(BINARY)

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
