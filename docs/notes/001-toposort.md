# 001 Toposort（拓扑排序）

## 目标

给 DAG 输出一个执行顺序：对每条依赖 `B depends on A`，必须保证 `A` 在 `B` 之前。

## 我在 FlowMesh 里做了什么

- 新增 `Graph.TopoSort() ([]string, error)`
- 先调用 `Validate()`，保证错误类型与校验一致：
  - 缺失依赖：`*MissingDependencyError`
  - 有环：`*CycleDetectedError`
- 正常情况下用 Kahn 算法输出拓扑序。

## Kahn 算法要点（依赖边 dep -> node）

1. 计算每个节点入度 `indegree[node]`：
  - 对每个节点 `node`，对每个依赖 `dep`，`indegree[node]++`
2. 建出边表 `out[dep]`：表示 `dep` 被哪些节点依赖（`dep -> node`）
3. 把所有 `indegree == 0` 的节点入队
4. 循环： 
  - 出队一个节点加入 `order`
  - 遍历它的 `out` 邻居，把邻居入度减1；减到0就入队
5. 若最终 `order` 长度 != 节点数：说明有环（本项目里理论上先被 `ValidateAcyclic()` 拦住）

## 测试怎么写（避免不稳定）

- 不断言“固定顺序”，只断言“相对顺序”：
  - 用 `pos[id]` 记录节点在结果里的位置
  - 断言 `pos[A] < pos[B]` 这类约束

## 今天踩的坑

- 先写测试会编译不通过：需要先加一个 `TopoSort()` 空壳让测试能跑，再逐步实现。
- HTTPS push 会被 reset: 改用 SSH remote 更稳。

