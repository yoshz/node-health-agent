# Node Health Agent

Agent that returns health for a specific node

## Prerequisites

```bash
# install golang
sudo yum install golang

# configure environment
export GOPATH=$HOME/go
export GO111MODULE=on
```

## Usage

```bash
# ensure kubeconfig is configured
kubectx of

# start server
go run main.go

# get current health for node
curl -i localhost:8090/?host=of-kube-gen-001.funix.nl

# cordon node
kubectl cordon of-kube-gen-001.funix.nl

# health check returns 500
curl -i localhost:8090/?host=of-kube-gen-001.funix.nl

# uncordon node
kubectl uncordon of-kube-gen-001.funix.nl

# health check returns 200 again
curl -i localhost:8090/?host=of-kube-gen-001.funix.nl
```

## Build

To build locally:

```bash
go build
```

To build Docker image:

```bash
docker build -t docker-registry.funix.nl/operations/node-health-agent .
```

## Deployment

node-health-agent is deployed from the https://gitlab.funix.nl/operations/k8s-addons repository.
