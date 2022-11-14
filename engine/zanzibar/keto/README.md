# Ory Keto

- [官方网站](https://www.ory.sh/keto/)
- [Github代码库](https://github.com/ory/keto-client-go)
- [官方文档](https://www.ory.sh/docs/keto/sdk/go)

支持的数据库：PostgreSQL、MySQL、CockroachDB、SQLite。

其中，开发时可以使用SQLite，而运营时最好不要使用。

### 安装部署Keto服务

具体文档可见：<https://www.ory.sh/docs/keto/install>

最基本的keto.yml

```yaml
version: v0.10.0-alpha.0

log:
  level: debug

namespaces:
  - id: 0
    name: app

serve:
  read:
    host: 0.0.0.0
    port: 4466
  write:
    host: 0.0.0.0
    port: 4467

dsn: memory
```

需要注意的是，新的版本当中，必须要有namespaces的定义，不然启动不了。

### Docker

#### 直接docker run启动

```powershell
docker pull oryd/keto:latest

docker run -itd --name keto-server `
    -p 4466:4466 -p 4467:4467 `
    -v /d/keto.yml:/home/ory/keto.yml `
    oryd/keto:latest serve -c /home/ory/keto.yml
```

需要注意的是，我把宿主的keto.yml直接挂载上去了，不然启动不了。

#### docker-compose启动

```yaml
version: "3"

services:
  keto:
    image: oryd/keto:v0.10.0-alpha.0
    ports:
      - "4466:4466"
      - "4467:4467"
    command: serve -c /home/ory/keto.yml
    restart: on-failure
    volumes:
      - type: bind
        source: .
        target: /home/ory
```

### Linux

```shell
bash <(curl https://raw.githubusercontent.com/ory/meta/master/install.sh) -d -b . keto v0.10.0-alpha.0
./keto help
```

### macOS

```shell
brew install ory/tap/keto
keto help
```

### Windows

```shell
irm get.scoop.sh | iex

scoop bucket add ory https://github.com/ory/scoop.git
scoop install keto

keto help
```

我尝试了使用sqlite启动，结果说没有支持：`could not create new connection: sqlite3 support was not compiled into the binary stack_trace`

### Kubernetes

```shell
helm repo add ory https://k8s.ory.sh/helm/charts
helm repo update
```

### 安装SDK

#### 安装gRPC API

```shell
go get github.com/ory/keto/proto@v0.10.0-alpha.0
```

#### 安装REST API

```shell
go get github.com/ory/keto-client-go@v0.10.0-alpha.0
```

## 参考资料

- [Zanzibar: Google’s Consistent, Global Authorization System](https://research.google/pubs/pub48190/)
- [My Reading on Google Zanzibar: Consistent, Global Authorization System](https://pushpalanka.medium.com/my-reading-on-google-zanzibar-consistent-global-authorization-system-f4a12df85cbb)
- [AuthZ: Carta’s highly scalable permissions system](https://medium.com/building-carta/authz-cartas-highly-scalable-permissions-system-782a7f2c840f)
- [Zanzibar-style ACLs with OPA Rego](https://gruchalski.com/posts/2022-05-07-zanzibar-style-acls-with-opa-rego/)
- [Zanzibar: A Global Authorization System - Presented by Auth0](https://zanzibar.academy/)
- [The Evolution of Ory Keto: A Global Scale Authorization System](https://www.ory.sh/keto-zanzibar-evolution/)
- [Building Zanzibar from Scratch](https://www.osohq.com/post/zanzibar)
- [What is Zanzibar?](https://authzed.com/blog/what-is-zanzibar/)
- [ZANZIBAR与ORY/KETO: 权限管理服务简介](https://chennima.github.io/keto-permission-manager-introduction)
- [What is Relationship Based Access Control (ReBAC)?](https://www.ubisecure.com/access-management/what-is-relationship-based-access-control-rebac/)
- [Relationship-Based Access Control (ReBAC)](https://www.osohq.com/academy/relationship-based-access-control-rebac)
- [详解微服务中的三种授权模式](https://www.infoq.cn/article/rl6g3buvaal8aiwvugdf)
- [如何使用 Ory Kratos 和 Ory Keto 保护您的烧瓶应用程序](https://devpress.csdn.net/python/62f99ab8c6770329307fef6d.html)
