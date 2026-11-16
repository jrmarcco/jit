# jit

自用的简易工具集。

## 包索引

| 包 | 说明 |
|---|---|
| `jit` | 指针工具（`Ptr`/`DePtr`）、数值约束（`RealNumber`/`Number`）、`Comparator` |
| `xslice` | 切片工具：查找、集合运算（并/交/差/对称差）、映射、聚合 |
| `xmap` | map 工具（`Keys`/`Vals`/`Merge`/`ToMap`）与 map 实现：`HashMap`（自定义哈希）、`TreeMap`（红黑树）、`MultiMap` |
| `xset` | 集合实现：`MapSet`、`TreeSet` |
| `xlist` | 列表实现：`ArrayList`、`LinkedList`、`CowArrayList`、`ConcurrentList`、`SkipList` |
| `xsync` | 并发原语：泛型 `Pool`/`Map`、支持 context 的 `Cond` |
| `xretry` | 重试循环与退避策略（固定间隔、指数退避） |
| `xpool` | 任务池：动态伸缩 goroutine 的 `BlockTaskPool` |
| `xjwt` | 基于 golang-jwt/jwt/v5 的 jwt 管理器（HMAC/RSA/Ed25519） |
| `xbean/option` | 泛型 option 模式 |
| `xbean/copy` | 基于反射的 struct 拷贝器（字段忽略、字段转换） |

## 开发

```bash
make setup   # 安装 gofumpt / golangci-lint / goimports 等工具与 git hooks
make check   # 格式化 + tidy
make vet     # go vet
make lint    # golangci-lint（CI 基线）
go test -race ./...
```
