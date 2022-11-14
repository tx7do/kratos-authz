# Google Zanzibar

Google Zanzibar是谷歌2016年起上线的一致性全球授权系统。这套系统的主要功能是

1. 储存来自各个服务的访问控制列表(Access Control Lists, ACLs)，也就是所谓的权限(Permission)。
2. 根据储存的ACL，进行权限校验。

这套系统上线后对接的服务有谷歌地图，谷歌图片，谷歌云盘，GCP，以及Youtube等等重要的服务

为了服务如此重要的业务，Zanzibar有着以下特点：

- 一致性：面对并发度如此大的业务场景，Zanzibar在检查权限的同时必须保证按照各个ACL的添加顺序判断。比如A添加了一条规则后立即删除，这两个动作如果没有按照正确的顺序执行，那么会造成权限泄露。
- 灵活性：各个业务场景的鉴权需求都不尽相同，所以Zanzibar灵活地支持不同的权限模式
- 横向扩展：以横向扩展支持数万亿条规则，每秒百万级鉴权
- 性能：95%的请求10毫秒内完成，99%的请求100毫秒内完成
- 可用性：上线三年间保证了99.999%的可用时间

以上的各个特性中除了灵活性之外都是性能或算法上的特点，性能和可靠性上也有很大一部分得益于底层的[Spanner](https://research.google/pubs/pub39966/)数据库。如果有兴趣可以阅读以下这篇论文：[Zanzibar: Google’s Consistent, Global Authorization System对Zanzibar](https://research.google/pubs/pub48190/)进行更深入的了解。下面我们就灵活性这一特点看一下Zanzibar是如何定义鉴权模型的。

## 概念与定义

### 关系元组(Relation Tuples)

Relation Tuples是Zanzibar的核心概念，一条Relation Tuples就对应了一条ACL。关系元组由：**命名空间(Namespace)**，**对象(Object)**，**关系(Relation)** 和 **主体(Subject)** 组成。一条Relation Tuples可以用[BNF语法](https://en.wikipedia.org/wiki/Backus%E2%80%93Naur_form)这样描述：

```ini
<relation-tuple> ::= <object>'#'relation'@'<subject>
<object> ::= namespace':'object_id
<subject> ::= subject_id | <subject_set>
<subject_set> ::= <object>'#'relation
```

一条Relation Tuples写作

```ini
namespace:object#relation@subject
```

意味着`subject`对`namespace`中的`object`有一种`relation`。

换成更有语义的例子：

```ini
videos:cat.mp4#view@felix
```

就意味着felix对videos中的cat.mp4有view的关系。

上述BNF定义的第四条`subject_set`是由`<object>'#'relation`组成的，也就是代表了一群对某种object有relation的subject。举例来说就是

```ini
groups:admin#member@felix
groups:admin#member@john
videos:cat.mp4#view@(groups:admin#member)
```

在这个例子中，felix和john都对groups:admin有member的关系，而对groups:admin有member的关系的subject_set对videos:cat.mp4有view的关系。也就是说felix和john都对videos:cat.mp4有view的关系。这种嵌套的语法可以有很多层，从而达到了整个ACL规则灵活可配的目的。

### 命名空间(Namespaces)， 对象(Object)与主体(Subject)

Zanzibar中的Namespace并不是起隔离作用的，就像上面的那个例子，在编写videosNamespace时也可以引用groupsNamespace。这里的命名空间概念更多是用来将数据分为同质的分块（并应用不同的配置），并且在储存层面上也是分离的。所以在多租户的使用场景中，用租户的UUID作为Namespace并不是一个好的选择，而应该使用tenants作为Namespace，从而实现

```ini
tenants:tenant-id-1#member@felix
tenants:tenant-id-1#member@john
```

这样的Relation Tuples，并且用tenants:tenant-id-1#member作为鉴权的subject_set。在命名方面，一般建议Namespace使用单词的复数形式，而Object和Subject使用UUID。 将Relation Tuples转换为图有助于更好地理解object与subject之间的关系，考虑[Keto官方文档](https://www.ory.sh/keto/docs/concepts/graph-of-relations)上的以下例子

```ini
// user1 has access on dir1
dir1#access@user1
// Have a look on the subjects concept page if you don't know the empty relation.
dir1#parent@(file1#)
// Everyone with access to dir1 has access to file1. This would probably be defined
// through a subject set rewrite that defines this inherited relation globally.
// In this example, we define this tuple explicitly.
file1#access@(dir1#access)
// Direct access on file2 was granted.
file2#access@user1
// user2 is owner of file2
file2#owner@user2
// Owners of file2 have access to it; possibly defined through subject set rewrites.
file2#access@(file2#owner)
```

## 什么是 ReBAC？


## Ory Keto

- [官方网站](https://www.ory.sh/keto/)
- [Github代码库](https://github.com/ory/keto-client-go)
- [官方文档](https://www.ory.sh/docs/keto/sdk/go)

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

## OpenFGA

Docker安装

```shell
docker pull openfga/openfga:latest

docker run -itd --name openfga-server `
  -p 8080:8080 `
  -p 8081:8081 `
  -p 3000:3000 `
  openfga/openfga:latest run
```

其中，8080是GRPC的接口，8081是HTTP的接口。

3000 提供了playground：<http://localhost:3000/playground>


## 参考资料

- [Zanzibar: Google’s Consistent, Global Authorization System](https://research.google/pubs/pub48190/)
- [My Reading on Google Zanzibar: Consistent, Global Authorization System](https://pushpalanka.medium.com/my-reading-on-google-zanzibar-consistent-global-authorization-system-f4a12df85cbb)
- [AuthZ: Carta’s highly scalable permissions system](https://medium.com/building-carta/authz-cartas-highly-scalable-permissions-system-782a7f2c840f)
- [Zanzibar-style ACLs with OPA Rego](https://gruchalski.com/posts/2022-05-07-zanzibar-style-acls-with-opa-rego/)
- [Zanzibar: A Global Authorization System - Presented by Auth0](https://zanzibar.academy/)
- [The Evolution of Ory Keto: A Global Scale Authorization System](https://www.ory.sh/keto-zanzibar-evolution/)
- [Building Zanzibar from Scratch](https://www.osohq.com/post/zanzibar)
- [OpenFGA : Auth0’s an open-source authorization solution](https://openfga.dev/)
- [What is Zanzibar?](https://authzed.com/blog/what-is-zanzibar/)
- [ZANZIBAR与ORY/KETO: 权限管理服务简介](https://chennima.github.io/keto-permission-manager-introduction)
- [What is Relationship Based Access Control (ReBAC)?](https://www.ubisecure.com/access-management/what-is-relationship-based-access-control-rebac/)
- [Relationship-Based Access Control (ReBAC)](https://www.osohq.com/academy/relationship-based-access-control-rebac)
- [详解微服务中的三种授权模式](https://www.infoq.cn/article/rl6g3buvaal8aiwvugdf)
- [Announcing OpenFGA - Auth0’s Open Source Fine Grained Authorization System](https://auth0.com/blog/auth0s-openfga-open-source-fine-grained-authorization-system/)
- [如何使用 Ory Kratos 和 Ory Keto 保护您的烧瓶应用程序](https://devpress.csdn.net/python/62f99ab8c6770329307fef6d.html)
