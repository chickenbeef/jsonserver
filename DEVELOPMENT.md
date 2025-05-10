# 🛠️ JsonServer Operator Development Guide

This file documents the steps performed to build a Kubernetes operator for managing [json-server](https://github.com/typicode/json-server) instances with CustomResourceDefinitions (CRDs), Admission Webhooks, and Controllers.

## 🔧 Prerequisites

The following dependencies needed to be installed:

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [kubebuilder](https://book.kubebuilder.io/) (+ [kubebuilder prerequisites](https://book.kubebuilder.io/quick-start.html#prerequisites))
- [Go development environment](https://go.dev/doc/install)

## 🚀 Development Steps

### 📦 Create Kubebuilder Project

Reference doc: <https://book.kubebuilder.io/quick-start#create-a-project>

```bash
mkdir -p jsonserver-operator && cd jsonserver-operator

go mod init jsonserver-operator

kubebuilder init --domain example.com

kubebuilder create api --group example --version v1 --kind JsonServer

make manifests
```

### ⚙️ Implement Controller Reconciling Logic

Update generated `internal/controller/jsonserver_controller.go`:

1. [RBAC annotations](https://book.kubebuilder.io/reference/markers/rbac) for permissions to manipulate `Deployments`, `Services` and `ConfigMaps` for the JsonServer CRD:

    ```go
    // +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
    // +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
    // +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
    ```

1. Implement the Reconcile functions to check if the JSON config is valid, naming convention is met and then create the underlying resources

    - Check that the resource name matches `app-*` naming convention
    - Validate supplied `jsonConfig` and create/update ConfigMap with the JSON data else fail fast and skip any remaining steps
      - Create/update Deployment using [backplane/json-server](https://hub.docker.com/r/backplane/json-server) image
      - Create/update Service to expose port 3000 (default port)
    - Update JsonServer CRD status

### 🔒 Implement Admission Webhook

References:

- [Webhook Implementation](https://book.kubebuilder.io/cronjob-tutorial/webhook-implementation)
- [Getting Started to Write Your First Kubernetes Admission Webhook](https://medium.com/trendyol-tech/getting-started-to-write-your-first-kubernetes-admission-webhook-part-2-48d0b0b1780e)
- [Simple Kubernetes Mutating Admission Webhook](https://breuer.dev/blog/kubernetes-webhooks)

```bash
kubebuilder create webhook --group example --version v1 --kind JsonServer --defaulting --programmatic-validation
# Next: implement your new Webhook and generate the manifests with:
# $ make manifests

make manifests
```

### 🔧 Configure Webhooks

References:

- [Running Webhooks](https://book.kubebuilder.io/cronjob-tutorial/running-webhook)

Created webhook [manifest](/config/webhook/manifests.yaml) and enabled cert-manager configuration.

### ✅ Check manifests generated correctly

```bash
bin/kustomize build config/default
```

### 🌟 Bonus: Implement scale

References:

- [Scale subresource example](https://book.kubebuilder.io/reference/generating-crd.html#scale)

### 🌟 Bonus: Push to ttl.sh in CI

References:

- [setup-buildx-action](https://github.com/docker/setup-buildx-action)
- [trivy-action](https://github.com/marketplace/actions/aqua-security-trivy)
