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


### code-generator 脚本

- [generate-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-groups.sh)
- [generate-internal-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-internal-groups.sh)

### Tag

Tag 语法

```
// +tag-name 
或
// +tag-name=value
```

Tag 类型

1. Global tags: 全局的 tag, 起到的是全局的代码生成控制的作用, 放在具体版本的 doc.go 文件中
2. Local tags: 本地的 tag, 放在 types.go 文件中的具体的 struct 上.

### 关于注释

+genclient 的意思是：请为下面这个 API 资源类型生成对应的 Client 代码。

+genclient:noStatus 的意思是：这个 API 资源类型定义里，没有 Status 字段。否则，生成的 Client 就会自动带上 UpdateStatus 方法。

+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object 的意思是：在生成 DeepCopy 的时候，实现 Kubernetes 提供的 runtime.Object 接口。否则，在某些版本的 Kubernetes 里，你的这个类型定义会出现编译错误。这是一个固定的操作

## Step-by-step write CR

1. 准备一个目录结构如下的项目
2. 通过 “复制、粘贴、替换、code generator” 可快速生成新的 crd

```bash
➜   tree $GOPATH/src/github.com/lqshow/k8s-custom-controllers/foobar-code-generator
.
├── artifacts
│   └── examples
│       ├── foobar-crd.yaml
│       └── foobar-example.yaml
├── go.mod
├── go.sum
└── pkg
    └── apis
       └── foobar
          ├── register.go
          └── v1
              ├── doc.go
              ├── register.go
              └── types.go 
```

### 1. 初始化项目
```bash

mkdir `PWD`/foobar-code-generator && cd `PWD`/foobar-code-generator
go mod init github.com/lqshow/k8s-custom-controllers/foobar-code-generator
mkdir -p pkg/apis/foobar/v1
```

### 2. 创建 register.go 文件，用来放置全局变量

```go
package foobar

const (
	GroupName = "k8s.io"
	Version   = "v1"
)
```
### 3. 初始化 crd 资源类型

```bash
cd pkg/apis/foobar/v1
```

doc.go
```go
// +k8s:deepcopy-gen=package
// +groupName=k8s.io

// Package v1 is the v1 version of the API.
package v1
```

types.go

1. 定义 FooBar 类型的具体内容
2. +genclient 只需要写在 FooBar 类型上，而不用写在 FooBarList 上。因为 FooBarList 只是一个返回值类型，FooBar 才是“主类型”。


register.go


### 4. 使用代码生成工具，为 FooBar 这个资源类型自动生成 clientset, informer 和 lister。

```bash
# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="github.com/lqshow/k8s-custom-controllers/foobar-code-generator"

# API Group
CUSTOM_RESOURCE_NAME="foobar"

# API Version
CUSTOM_RESOURCE_VERSION="v1"

# 执行代码自动生成，其中 pkg/generated 是生成目标目录，pkg/apis 是类型定义目录
$GOPATH/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
```

```bash
➜  $GOPATH/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
Generating deepcopy funcs
Generating clientset for foobar:v1 at github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/clientset
Generating listers for foobar:v1 at github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/listers
Generating informers for foobar:v1 at github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/generated/informers
```

自动生成的代码
1. pkg/generated/clientset
2. pkg/generated/listers
3. pkg/generated/informers
4. pkg/apis/foobar/v1/zz_generated.deepcopy.go
       
### 5. 编写自定义控制器代码

1. 编写 main 函数
2. 编写自定义控制器的定义
3. 编写控制器的业务逻辑

## Create FooBar Object in kubernetes

```bash
kubectl apply -f `PWD`/artifacts/examples/foobar-crd.yaml
kubectl apply -f `PWD`/artifacts/examples/foobar-example.yaml
```

## TODO

1. 在 GO 工作区目录（例如$GOPATH/src） 执行命令会报以下错误，在当前目录下执行就没问题，很妖

```bash
➜  $GOPATH/src/k8s.io/code-generator/generate-groups.sh all "$ROOT_PACKAGE/pkg/generated" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"
Generating deepcopy funcs
F1008 19:54:25.625104   96258 deepcopy.go:885] Hit an unsupported type invalid type for invalid type, from github.com/lqshow/k8s-custom-controllers/foobar-code-generator/pkg/apis/foobar/v1.FooBar
```

## References

- [Generators for kube-like API types](https://github.com/kubernetes/code-generator)
- [Kubernetes Deep Dive: Code Generation for CustomResources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)
- [Extending Kubernetes: Create Controllers for Core and Custom Resources](https://medium.com/@trstringer/create-kubernetes-controllers-for-core-and-custom-resources-62fc35ad64a3)