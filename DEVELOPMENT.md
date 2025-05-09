# JsonServer Operator Development Guide

Document steps performed to build a Kubernetes operator for managing [json-server](https://github.com/typicode/json-server) instances with CustomResourceDefinitions (CRDs), Admission Webhooks, and Controllers.

## Prerequisites

The following dependencies needed to be installed:

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [kubebuilder](https://book.kubebuilder.io/) (+ [kubebuilder prerequisites](https://book.kubebuilder.io/quick-start.html#prerequisites))
- [Go development environment](https://go.dev/doc/install)

## Create Kubebuilder Project

<https://book.kubebuilder.io/quick-start#create-a-project>

```bash
mkdir -p jsonserver-operator && cd jsonserver-operator

go mod init jsonserver-operator

kubebuilder init --domain example.com

kubebuilder create api --group example --version v1 --kind JsonServer
```
