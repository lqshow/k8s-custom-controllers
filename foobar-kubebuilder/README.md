# Overview

Kubebuilder 节省大量工作，方便用户从零开始开发 CRDs，Controllers 和 Admission Webhooks，让扩展 K8s 变得更简单

## Installation

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -sL https://go.kubebuilder.io/dl/2.0.1/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.0.1_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

## Step-by-step write CR

### 1. 初始化项目

```bash
cd $GOPATH/src/github.com/lqshow/k8s-custom-controllers/foobar-kubebuilder
kubebuilder init --domain basebit.me  --owner "LQ"
```

### 2. 创建 API

```bash
kubebuilder create api --group samplecrd --version v1alpha2 --kind FooBar
```

### 3. 安装 CRD 并启动 controller

```bash
make install
make run
```

```bash
➜  kubectl get crd foobars.samplecrd.basebit.me
NAME                           CREATED AT
foobars.samplecrd.basebit.me   2019-10-09T15:28:44Z
```

```bash
kubectl apply -f config/samples/samplecrd_v1_foobar.yaml

➜   kubectl get foobars.samplecrd.basebit.me
NAME            AGE
foobar-sample   3m
```

## References

- [The Kubebuilder Book](https://book.kubebuilder.io/introduction.html)
- [进阶 K8s 高级玩家必备 | Kubebuilder：让编写 CRD 变得更简单](https://mp.weixin.qq.com/s/UzEcj2eXKM0m8f4XzZCYAA)