# Overview

1. Custom controllers 允许用户自定义控制器的逻辑，基于已有的资源定义更高阶的控制器，实现 Kubernetes 集群原生不支持的功能
2. Custom controllers 的逻辑其实是很简单的：watch CRD 实例（以及关联的资源）的 CRUD 事件，然后开始执行相应的业务逻辑


## Controller model

> 一个无限循环不断地对比实际状态和期望状态，如果有出入则进行调谐逻辑将实际状态调整为期望状态，最终达到与申明一致

实际状态: 来自于 Kubernetes 集群本身
期望状态: 来自于 用户提交的 YAML 文件


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

## Workflow

1. 编写 CRD
2. 编写 Custom Controller


## References

- [Writing Controllers](https://github.com/kubernetes/community/blob/8decfe4/contributors/devel/controllers.md)
- [浅析 Kubernetes 控制器的工作原理](https://www.yangcs.net/posts/a-deep-dive-into-kubernetes-controllers/)

