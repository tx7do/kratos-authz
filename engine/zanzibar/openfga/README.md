# OpenFGA

- [官方网站](https://openfga.dev/)
- [Github代码库](https://github.com/openfga)
- [官方文档](https://openfga.dev/docs/authorization-and-openfga)

支持的数据存储引擎：PostgreSQL、MySQL。

开发时支持内存存储引擎，但是一旦服务关机重启，数据将会丢失。

Docker安装

```shell
docker pull openfga/openfga:latest

docker run -itd --name openfga-server `
  -p 8080:8080 `
  -p 8081:8081 `
  -p 3000:3000 `
  openfga/openfga:latest run
```

- 8080 是GRPC的接口
- 8081 是HTTP的接口
- 3000 提供了playground：<http://localhost:3000/playground>
- 3001 提供了性能探查器

## 参考资料

- [OpenFGA : Auth0’s an open-source authorization solution](https://openfga.dev/)
- [Announcing OpenFGA - Auth0’s Open Source Fine Grained Authorization System](https://auth0.com/blog/auth0s-openfga-open-source-fine-grained-authorization-system/)
