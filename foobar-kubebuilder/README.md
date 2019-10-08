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

## Create a Project

```bash
cd $GOPATH/src
kubebuilder init --domain basebit.me  --owner "LQ"
```

## Create an API

```bash
kubebuilder create api --group enigma --version v1 --kind DagNodeRunner
```



## References

- [The Kubebuilder Book](https://book.kubebuilder.io/introduction.html)
- [进阶 K8s 高级玩家必备 | Kubebuilder：让编写 CRD 变得更简单](https://mp.weixin.qq.com/s/UzEcj2eXKM0m8f4XzZCYAA)