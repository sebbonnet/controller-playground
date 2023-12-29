# Controller playground

## Bootstrap

1. Initialize project
```
go mod init github.com/sebbonnet/controller-playground
```

2. Add main.go from https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

3. Import all dependencies

```
go mod tidy
```

## Usage

```
IMG=<repo>/my-controller NAMESPACE=${NAMESPACE} make docker-build docker-push deploy
```

```
kubectl -n ${NAMESPACE} logs -f -lapp=my-controller
```
