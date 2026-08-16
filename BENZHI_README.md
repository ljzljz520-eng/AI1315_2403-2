# 泥间陶艺体验馆多页站

这是一个单一 Go module 的陶艺体验馆演示项目。它提供 12 个可直接访问的页面、一个内存预约记录和一组用于验收操作员确认流程的 HTTP API。数据来自固定 fixture，启动时不需要数据库、网络服务、时钟或随机数。

## 运行

要求 Go 1.22.12 或更高的兼容版本，以及 Node.js 20。

```bash
go run ./cmd/pottery-house
```

浏览器打开 `http://localhost:8080/`。预约页会提交到固定记录 `陶艺预约-2026-001`，也可以直接调用：

```bash
curl -X POST 'http://localhost:8080/api/records/%E9%99%B6%E8%89%BA%E9%A2%84%E7%BA%A6-2026-001/confirm' \
  -H 'Content-Type: application/json' \
  -d '{"operator":"周老师","field":"泥料","value":"白陶泥"}'
```

## 页面

首页、课程、作品、釉色、体验项目总览、拉坯入门、手捏杯子、盘中风景、预约、门店环境、活动和联系门店分别位于 `web/public/`。页面使用清晰的 header、section、article 和 form 分区，便于检查 HTML 与 CSS。

前端构建命令：

```bash
cd web
npm install
npm run build
```

构建会把 `web/public/` 整理到 `web/dist/`；`dist/` 是临时产物，不纳入源树。Go 服务优先读取 `web/dist/`，不存在时读取 `web/public/`。

## 业务验收

```bash
go test -count=1 ./...
```

并发验收使用 `CGO_ENABLED=1`，普通全量测试可使用 `CGO_ENABLED=0`。并发业务链路会让两名操作员同时确认同一预约的不同字段，并检查汇总是否同时包含两项内容。当前项目按题目要求保留一个可稳定复现的并发更新缺陷，因此该验收用例会失败，失败内容就是汇总中缺少其中一项确认。

## 目录

```text
cmd/pottery-house/       可执行入口
internal/booking/        fixture、内存仓库、更新服务与业务测试
internal/site/           HTTP API 和静态页面路由
web/public/              12 个页面、样式与交互脚本
web/scripts/             Node.js 20 构建脚本
```
