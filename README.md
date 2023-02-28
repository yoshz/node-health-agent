# Node Health Agent

An agent that returns the health for a specific node by making sure node taints are healthy.

When deployed as an DaemonSet it can be used as health check endpoint in your load balancer.

For example when you (cordon) taint a node unhealthy, the health endpoint will return unhealthy and the node is not used by the load balancer anymore.

# Installation

```bash
helm repo add node-health-agent https://yoshz.github.io/node-health-agent/
helm install --namespace kube-system node-health-agent node-health-agent/node-health-agent
```

# Development

## Usage

```bash
# install dependencies
go get

# start server
go run main.go --kubeconfig $KUBECONFIG

# get current health for node
curl -i localhost:8991/?host=k8s-master-01

# cordon node
kubectl cordon k8s-master-01

# health check returns 500
curl -i localhost:8991/?host=k8s-master-01

# uncordon node
kubectl uncordon k8s-master-01

# health check returns 200 again
curl -i localhost:8991/?host=k8s-master-01
```

## Build

Build the binary locally:
```bash
go build
```

Build Docker image:
```bash
docker build -t yoshz/node-health-agent .
```

Create Helm package:
```
helm package charts/node-health-agent
git checkout gh-pages
helm repo index . --url https://yoshz.github.io/node-health-agent/
```

