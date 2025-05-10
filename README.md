# 🔄 jsonserver-operator

[![Build, Scan and Push to ttl.sh](https://github.com/chickenbeef/jsonserver/actions/workflows/build-push-ttl.yml/badge.svg)](https://github.com/chickenbeef/jsonserver/actions/workflows/build-push-ttl.yml)

Kubernetes operator for managing [json-server](https://github.com/typicode/json-server) instances with CustomResourceDefinitions (CRDs), Admission Webhooks, and Controllers.

For development steps see [DEVELOPMENT.md](DEVELOPMENT.md)

## 🚀 Getting Started

### 🔧 Prerequisites

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)

### 📦 Deploy to local cluster

1. Set a cluster name:

    ```sh
    CLUSTER_NAME=jsonserver-cluster
    ```

1. Create a Kind cluster:

    ```sh
    kind create cluster --name $CLUSTER_NAME
    ```

1. Install cert-manager (required for webhooks):

    - <https://book.kubebuilder.io/cronjob-tutorial/cert-manager>
    - <https://cert-manager.io/docs/installation/>

    ```sh
    kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.17.2/cert-manager.yaml
    kubectl wait --for=condition=Available --timeout=300s deployment/cert-manager -n cert-manager
    ```

1. Install CRDs for JsonServer objects:

    ```sh
    make install
    ```

1. Build and push the `jsonserver-operator` image to your own registry:

    Set the image repository and tag details then build and push it:

    ```sh
    IMG=<some-registry>/jsonserver-operator:tag
    make docker-build docker-push IMG=$IMG
    ```

    **Example:**

    ```sh
    IMG=chickenbeef/jsonserver-operator:latest
    make docker-build docker-push IMG=$IMG
    ```

1. Load image into Kind cluster:

    ```sh
    kind load docker-image $IMG --name $CLUSTER_NAME
    ```

1. Deploy

    ```sh
    make deploy IMG=$IMG

    kubectl wait --for=condition=Available --timeout=300s deployment/jsonserver-operator-controller-manager -n jsonserver-operator-system
    ```

#### 🔍 Testing the Operator

1. Test the name validation (should be enforced by the webhook):

    This should fail due to the name not matching the `app-*` convention.

    ```sh
    kubectl apply -f - <<EOF
    apiVersion: example.example.com/v1
    kind: JsonServer
    metadata:
      name: invalid-name
      namespace: default
    spec:
      replicas: 1
      jsonConfig: |
        { "test": [ { "id": 1, "name": "Test" } ] }
    EOF
    ```

    Expected error:

    ```sh
    Error from server (Forbidden): error when creating "STDIN": admission webhook "vjsonserver-v1.kb.io" denied the request: JsonServer name must follow the convention 'app-${name}'
    ```

1. Test the JSON validation (should be handled by the controller, not blocked by webhook):

    ```sh
    kubectl apply -f - <<EOF
    apiVersion: example.example.com/v1
    kind: JsonServer
    metadata:
      name: app-invalid-json
      namespace: default
    spec:
      replicas: 1
      jsonConfig: |
        { "invalid json here }
    EOF
    ```

    Expected output:

    > jsonserver.example.example.com/app-invalid-json created

    Check the status - should show `Error` state:

    ```sh
    kubectl get jsonserver app-invalid-json -o jsonpath='{.status}' | jq
    ```

    Expected output:

    ```json
    {
      "message": "Error: spec.jsonConfig is not a valid json object",
      "state": "Error"
    }
    ```

    Check that no resources have been created:

    ```sh
    kubectl get configmap,deployment,service,pods -l app=app-invalid-json
    ```

    Cleanup:

    ```sh
    kubectl delete jsonserver app-invalid-json
    ```

1. Create a valid JsonServer instance:

    ```sh
    kubectl apply -f - <<EOF
    apiVersion: example.example.com/v1
    kind: JsonServer
    metadata:
      name: app-my-server
      namespace: default
    spec:
      replicas: 2
      jsonConfig: |
        { "people": [
            { "id": 1, "name": "Person A" },
            { "id": 2, "name": "Person B" }
          ]
        }
    EOF
    ```

1. Verify the resources were created:

    ```sh
    watch -d kubectl get configmap,deployment,service,pods -l app=app-my-server
    ```

1. Test accessing the JSON server:

    ```sh
    kubectl port-forward svc/app-my-server 8080:3000
    ```

    In another tab:

    ```sh
    curl http://localhost:8080/people
    ```

    Expected output:

    ```json
    [
      {
        "id": "1",
        "name": "Person A"
      },
      {
        "id": "2",
        "name": "Person B"
      }
    ]
    ```

1. (Bonus) Test scaling

    Scale up:

    ```bash
    kubectl scale jsonserver app-my-server --replicas 5

    # Check new pods created:
    kubectl get pods -l app=app-my-server
    ```

    Scale back down:

    ```bash
    kubectl scale jsonserver app-my-server --replicas 1

    # Check pods terminated:
    kubectl get pods -l app=app-my-server
    ```

1. Cleanup

    Delete the test `jsonserver` object:

    ```sh
    kubectl delete jsonserver app-my-server

    # Check all resources deleted
    kubectl get configmap,deployment,service,pods -l app=app-my-server
    ```

### 🧹 Cleanup

1. Delete the APIs(CRDs) from the cluster (OR):

    ```sh
    make uninstall
    ```

1. Undeploy the controller *and* CRDs from the cluster:

    ```sh
    make undeploy
    ```

### 🔄 CI/CD with GitHub Actions (Bonus)

This project includes a [GitHub workflow](https://github.com/chickenbeef/jsonserver/actions) that automatically builds, scans, and pushes the Docker image to [ttl.sh](https://ttl.sh/):

- **Triggers**: Runs on pushes to the `main` branch and all pull requests
- **Registry**: Uses ttl.sh as a free, ephemeral Docker registry (images expire after 24 hours)
- **Tags**: Images are tagged as `ttl.sh/jsonserver-operator-{git-sha}:24h`
- **Security Scanning**: [Uses Trivy to scan for vulnerabilities](https://github.com/chickenbeef/jsonserver/security) in the container image

To use the CI-built images:

```sh
# Example of pulling the latest CI-built image (replace with actual SHA)
docker pull ttl.sh/jsonserver-operator-115ad73:24h

# Using in your deployment
make deploy IMG=ttl.sh/jsonserver-operator-115ad73:24h
```

---

### 📚 Alternative Installation Methods

**Note:** *Following sections are auto-generated by Kubebuilder*

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/jsonserver-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

1. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/jsonserver-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

1. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.
