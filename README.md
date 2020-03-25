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

Create service account
```bash
kubectl -n monitoring create -f deploy/serviceaccount.yaml
```

Install helm chart
```bash
helm3 upgrade node-health-agent funix/generic \
  --namespace monitoring \
  --install \
  --values deploy/helm-values.yaml \
  --dry-run --debug
```

The ingress expects hostname `node-health-agent` which isn't configured on the load balancers.
To test if node-health-agent is running correctly you can execute a curl request from the mgtm server:
```bash
ssh of-mgmt-gen-001
curl -i -H Host:node-health-agent of-kube-gen-301.funix.nl/?host=of-kube-gen-301.funix.nl
```
