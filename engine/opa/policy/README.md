# OPA 授权策略

当前文件夹下面包含了OPA的授权策略定义。
有关OPA的信息，请参见：

- [OPA - Github](https://github.com/open-policy-agent/opa/)
- [OPA - 网站](https://www.openpolicyagent.org/)

## 使用 REPL （交互式解释器）

[下载 OPA](https://www.openpolicyagent.org/docs/get-started.html#prerequisites) 

运行以下命令：

```shell
opa run -w authz.rego common.rego policies:../example/policies.json
```

进入交互式解释器：

```shell
OPA 0.9.2 (commit 9fbff4c3, built at 2018-09-24T16:12:26Z)

> data.authz.authorized
false
# This matches against an action/resource from a statement in a policy.
> data.authz.authorized with input as { "subjects": [ "team:local:admins" ], "action": "iam:teams:create", "resource": "iam:teams" }
true
> data.authz.authorized with input as { "subjects": [ "team:local:admins" ], "action": "iam:teams:create", "resource": "iam:users" }
false
# This matches against an action/resource from a statement in a policy.
> data.authz.authorized with input as { "subjects": [ "team:local:admins" ], "action": "infra:nodes:delete", "resource": "infra:nodes" }
true
>
```

## 运行OPA单元测试

运行测试，只获取最终结果：

```shell
opa test authz.rego common.rego authz_test.rego
```

运行测试，并打印详情：

```shell
opa test -v authz.rego common.rego authz_test.rego
```
