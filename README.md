![Go](https://github.com/supernova106/kubestorm/workflows/Go/badge.svg?branch=master)

# kubestorm

- A RESTful API to interact with Kubernetes resources
- Easy to customize and integrate with UI/Chatbot/Automation

It consists of 3 main components:

- Backend - Rest API (Go)
- Frontend (Vue.js) (BEING DEVELOPED)
- Persistent Layer to store configurations

## Features

- To support cross platform Kubernetes Management
- Simple UI to view Kubernetes clusters

## Docs

- Update swagger docs: [https://github.com/swaggo/swag](https://github.com/swaggo/swag)

```sh
swag init
```

## Development

Clone the repository

```sh
git clone https://github.com/supernova106/kubestorm.git
```

Dependency management tools are built into go 1.11+ in the form of go modules. These are used by the main Kubernetes repo (>= v1.15.0) and client-go (>= kubernetes-1.15.0) to manage dependencies. If you are using go 1.11 or 1.12 and are working with a project located within \$GOPATH, you must opt into using go modules:

```sh
export GO111MODULE=on
```

Start the API locally

```sh
cd kubestorm/
go run main.go
```

Or using the pre-built binary

```sh
./kubestorm
```

### Generate Auth TOKEN for your Kubernetes cluster

- create the following `kubestorm.yaml`

```yaml
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kubestorm-role
rules:
  - apiGroups: ["", "extensions", "apps", "batch", "events", "resourcequotas"]
    resources: ["*"]
    verbs: ["list", "watch", "get"]
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kubestorm-user
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kubestorm-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kubestorm-role
subjects:
  - kind: ServiceAccount
    name: kubestorm-user
    namespace: kube-system
```

```sh
kubectl apply -f kubestorm.yaml
```

Get `server URL`

```sh
export CURRENT_CONTEXT=$(kubectl config current-context) && export CURRENT_CLUSTER=$(kubectl config view -o go-template="{{\$curr_context := \"$CURRENT_CONTEXT\" }}{{range .contexts}}{{if eq .name \$curr_context}}{{.context.cluster}}{{end}}{{end}}") && echo $(kubectl config view -o go-template="{{\$cluster_context := \"$CURRENT_CLUSTER\"}}{{range .clusters}}{{if eq .name \$cluster_context}}{{.cluster.server}}{{end}}{{end}}")
```

Get `serverCADataString`

```sh
echo $(kubectl get secret -n kube-system -o go-template='{{index .data "ca.crt" }}' $(kubectl get sa kubestorm-user -n kube-system -o go-template="{{range .secrets}}{{.name}}{{end}}"))
```

Get `token`

```sh
echo $(kubectl get secret -n kube-system -o go-template='{{index .data "token" }}' $(kubectl get sa kubestorm-user -n kube-system -o go-template="{{range .secrets}}{{.name}}{{end}}")) | base64 --decode
```

### ADD cluster auth

Execute the following script to add cluster

```sh
./scripts/add-cluster.sh <cluster_name>
```

### GET cluster auth

```sh
curl --location --request GET 'http://localhost:8080/v1/auth/CLUSTER_NAME'
```

### GET nodes from a cluster

```sh
curl --location --request GET 'http://localhost:8080/v1/nodes/CLUSTER_NAME'
```

### DELETE cluster auth

```sh
curl --location --request DELETE 'http://localhost:8080/v1/auth/CLUSTER_NAME'
```

## Release

- Create new release

```sh
make release
```

## Todo

- Support S3 to store configurations
- Showing Information about Kubernets objects: `nodes`, `services`, `namespace`, `deployment`, `events`

## Contributors

- Binh Nguyen
