# gokrm : An Object-Kubernetes Resource-Mapper for Go

gokrm (pronounced _go-kay-rem_) is a Go library that provides tooling to build declarative, simple reconcilers for Kubernetes resources, using a single source of truth that can be a custom resource or a configuration value, and reflection, to dynamically introspect the soure type and generate the mapped resources.

## Example

The following is an example of a `SourceType` mapped into a `Deployment` and `Service`:

```go
// Type declaration
type sourceType struct {
  metav1.ObjectMeta

  Image     string `gokrmTarget:"Deployment" gokrmTargetField:"Spec.Template.Spec.Containers[0].Image"`
  Replicas  *int32 `gokrmTarget:"Deployment" gokrmTargetField:"Spec.Replicas"`
  ClusterIP string `gokrmTarget:"Service" gokrmTargetField:"Spec.ClusterIP"`
  Port      int32  `gokrmTarget:"Service" gokrmTargetField:"Spec.Ports[0].Port"`
}
```

> Note the `gokrmTarget` and `gokrmTargetField` tags. They indicate where to map these fields to.

An example value

```go
source := SourceType{
  ObjectMeta: metav1.ObjectMeta{
    Name:      "test",
    Namespace: "default",
  },
  Replicas:  intptr(10),
  ClusterIP: "1.1.1.1",
  Image:     "image:latest",
  Port:      8080,
}
```

Will result in the following resources being generated:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  name: test
  namespace: default
spec:
  replicas: 10
  selector: null
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
    spec:
      containers:
      - image: image:latest
        name: ""
        resources: {}
status: {}
---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  name: test
  namespace: default
spec:
  clusterIP: 1.1.1.1
  ports:
  - port: 8080
    targetPort: 0
status:
  loadBalancer: {}

```

## Usage