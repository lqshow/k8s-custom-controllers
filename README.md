# Overview

1. Custom controllers 允许用户自定义控制器的逻辑，基于已有的资源定义更高阶的控制器，实现 Kubernetes 集群原生不支持的功能
2. Custom controllers 的逻辑其实是很简单的：watch CRD 实例（以及关联的资源）的 CRUD 事件，然后开始执行相应的业务逻辑


## Controller model

> 一个无限循环不断地对比实际状态和期望状态，如果有出入则进行调谐逻辑将实际状态调整为期望状态，最终达到与申明一致

- 实际状态: 来自于 Kubernetes 集群本身
- 期望状态: 来自于 用户提交的 YAML 文件


```golang
for {
  实际状态 := 获取集群中对象 X 的实际状态（Actual State）
  期望状态 := 获取集群中对象 X 的期望状态（Desired State）
  if 实际状态 == 期望状态{
    // do nothing
  } else {
    // reconcile loop (执行编排动作，将实际状态调整为期望状态)
  }
}
```

![image](https://user-images.githubusercontent.com/8086910/66443576-3baee600-ea72-11e9-994c-28db8ce74ce7.png)
![image](https://user-images.githubusercontent.com/8086910/66454556-c1df2280-ea9a-11e9-8a0a-de325513378e.png)

### Informer

1. Informer 与 API 对象是一一对应的
2. Informer 其实就是一个带有本地缓存和索引机制的、可以注册 EventHandler 的 client。它是自定义控制器跟 APIServer 进行数据同步的重要组件。
3. Informer 通过一种叫作 ListAndWatch 的方法，把 APIServer 中的 API 对象缓存在了本地，并负责更新和维护这个缓存。

## Workflow

1. 编写 CRD
2. 编写 Custom Controller

### CRD

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: foobars.k8s.io
spec:
  group: k8s.io
  names:
    kind: FooBar
    listKind: FooBarList
    plural: foobars
    singular: foobar
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: FooBar is the Schema for the foobars API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: FooBarSpec defines the desired state of FooBar
          type: object
        status:
          description: FooBarStatus defines the observed state of FooBar
          type: object
      type: object
  version: v1
  versions:
    - name: v1
      served: true
      storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
```

### CR

```yaml
apiVersion: k8s.io/v1
kind: FooBar
metadata:
  name: foobar-sample
spec:
  foo: bar
  barz: "false"
  command: "echo foobar"
```

## References

- [Writing Controllers](https://github.com/kubernetes/community/blob/8decfe4/contributors/devel/controllers.md)
- [Example Kubernetes controller: the cloud native at command](https://github.com/programming-kubernetes/cnat)
- [浅析 Kubernetes 控制器的工作原理](https://www.yangcs.net/posts/a-deep-dive-into-kubernetes-controllers/)