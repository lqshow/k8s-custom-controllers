# Overview 

Kubernetes 使用 CRD+Controller 来扩展集群功能，官方提供了 CRD 代码的自动生成器 code-generator

## Kubernetes code-generator

> 其中 informer-gen 和 lister-gen 是构建 controller 的基础

| Code-gen     | Desc                                                         |
| ------------ | ------------------------------------------------------------ |
| deepcopy-gen | 生成`func` `(t *T)` `DeepCopy()` `*T` 和`func` `(t *T)` `DeepCopyInto(*T)`方法 |
| client-gen   | 创建类型化客户端集合(typed client sets)                      |
| informer-gen | 为CR创建一个informer , 当CR有变化的时候, 这个informer可以基于事件接口获取到信息变更 |
| lister-gen   | 为CR创建一个listers , 就是为`GET` and `LIST`请求提供read-only caching layer |


code-generator 脚本

- [generate-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-groups.sh)
- [generate-internal-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-internal-groups.sh)


## Step-by-step write CR

```bash
# 初始化项目
mkdir code-generator-sample && cd code-generator-sample
go mod init github.com/lqshow/k8s-custom-controllers

# 初始化crd资源类型
mkdir -p pkg/api/foobar/v1 && cd pkg/api/foobar/v1
```

## References

- [Generators for kube-like API types](https://github.com/kubernetes/code-generator)
- [Kubernetes Deep Dive: Code Generation for CustomResources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)
