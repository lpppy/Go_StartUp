# go-basics-demo

面向有 C 语言基础的开发者的 Go 语言入门示例项目。通过 13 个独立演示模块，系统覆盖 Go 核心语法与惯用法，每个模块均附有与 C 语言的详细对比注释。

## 前置要求

- Go 1.21 及以上版本（[下载地址](https://go.dev/dl/)）
- 验证安装：`go version`

## 快速开始

```bash
# 克隆项目后进入目录
cd go-basics-demo

# 下载依赖
go mod download

# 运行所有演示模块（按顺序输出每个模块的演示结果）
go run .
```

## 运行测试

```bash
# 运行所有测试
go test ./...

# 详细输出（显示每个测试用例）
go test -v ./demos/testing_demo/

# 运行基准测试
go test -bench=. ./demos/testing_demo/

# 查看测试覆盖率
go test -cover ./...
```

## 模块说明

按推荐学习顺序排列，每个模块对应 `demos/` 下的一个子包。

| # | 目录 | 主题 | 核心知识点 |
|---|------|------|-----------|
| 1 | `demos/variables` | 变量与基本类型 | var / := / 零值 / iota / 类型转换 |
| 2 | `demos/controlflow` | 控制流 | for / if / switch / range / label |
| 3 | `demos/functions` | 函数 | 多返回值 / 闭包 / defer / 可变参数 |
| 4 | `demos/collections` | 数组、切片与 Map | slice 底层结构 / append / copy / map |
| 5 | `demos/structs` | 结构体与方法 | 值接收者 / 指针接收者 / 嵌入 / JSON |
| 6 | `demos/interfaces` | 接口 | 隐式实现 / 多态 / 类型断言 / nil 陷阱 |
| 7 | `demos/errors` | 错误处理 | error 接口 / errors.Is / errors.As / panic/recover |
| 8 | `demos/concurrency` | 并发 | goroutine / channel / select / sync / context |
| 9 | `demos/pointers` | 指针 | & / * / new / 逃逸分析 / 无指针运算 |
| 10 | `demos/packages` | 包与模块 | import / 导出规则 / init / go.mod |
| 11 | `demos/generics` | 泛型（Go 1.18+） | 类型参数 / 约束 / ~T / 泛型结构体 |
| 12 | `demos/testing_demo` | 测试框架 | 表驱动测试 / 子测试 / 基准测试 / 示例测试 |
| 13 | `demos/stdlib` | 标准库常用包 | fmt / strings / os / json / http / time / regexp |

## 项目结构

```
go-basics-demo/
├── main.go              # 入口：Runner 按顺序执行所有模块，recover 捕获 panic
├── go.mod               # 模块声明（go 1.21）
├── go.sum               # 依赖校验
├── README.md
└── demos/
    ├── variables/
    │   └── variables.go
    ├── controlflow/
    │   └── controlflow.go
    ├── functions/
    │   └── functions.go
    ├── collections/
    │   └── collections.go
    ├── structs/
    │   └── structs.go
    ├── interfaces/
    │   └── interfaces.go
    ├── errors/
    │   └── errors.go
    ├── concurrency/
    │   └── concurrency.go
    ├── pointers/
    │   └── pointers.go
    ├── packages/
    │   └── packages.go
    ├── generics/
    │   └── generics.go
    ├── testing_demo/
    │   ├── testing_demo.go
    │   └── testing_demo_test.go
    └── stdlib/
        └── stdlib.go
```

## 设计说明

**Runner 机制**：`main.go` 中的 Runner 用 `defer/recover` 包裹每个模块的调用，部分模块（如 `collections`、`variables`）会故意触发 panic 来演示语言行为，Runner 会捕获并打印错误后继续执行下一个模块，不会中断整个程序。

**与 C 的对比**：每个模块的源码注释中均包含与 C 语言的详细对比，涵盖语法差异、内存模型、错误处理、并发模型等方面，帮助有 C 基础的开发者快速建立 Go 语言思维。

**Go 与 C 的核心差异速览**：

| 特性 | C | Go |
|------|---|----|
| 内存管理 | 手动 malloc/free | GC 自动管理 |
| 错误处理 | errno / 返回码 | error 接口 / 多返回值 |
| 并发 | pthread（重量级线程） | goroutine（轻量级，初始栈 2KB） |
| 类型转换 | 隐式转换 | 必须显式转换 |
| 数组 | 退化为指针，无边界检查 | 值类型，越界触发 panic |
| 字符串 | `\0` 结尾的字符数组 | 不可变值类型，UTF-8 |
| 面向对象 | 无 | 结构体 + 方法 + 接口（无 class/继承） |
| 泛型 | void* / 宏 | 类型参数（Go 1.18+） |

## 依赖

- [`pgregory.net/rapid`](https://github.com/flyingmutant/rapid) — 属性测试库（用于 `testing_demo` 模块）
