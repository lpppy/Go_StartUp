# 实现计划：go-basics-demo

## 概述

将 go-basics-demo 项目拆解为可逐步执行的编码任务。每个任务对应一个或多个演示模块的实现，最终通过 main.go 的 Runner 将所有模块串联起来。属性测试使用 `pgregory.net/rapid` 库。

## 任务列表

- [x] 1. 初始化项目结构
  - 执行 `go mod init go-basics-demo`，生成 `go.mod`（Go 版本 1.21）
  - 创建 `demos/` 目录及各子包目录（variables、controlflow、functions、collections、structs、interfaces、errors、concurrency、pointers、packages、generics、testing_demo、stdlib）
  - 执行 `go get pgregory.net/rapid` 添加属性测试依赖
  - 创建 `README.md`，说明运行方式与各模块学习顺序
  - _需求：1.4, 1.5, 1.6_

- [x] 2. 实现 main.go 入口与 Runner
  - 创建 `main.go`，定义 `demoEntry` 结构体和 `demos` 注册表切片
  - 实现 `runSafe(name string, fn func())` 函数，使用 `defer/recover` 捕获 panic
  - 实现 `main()` 函数，遍历 demos 列表，打印分隔线和模块名后调用 `runSafe`
  - _需求：1.1, 1.2, 1.3, 1.7_

- [x] 3. 实现 demos/variables 模块
  - [x] 3.1 实现 `variables.go`
    - 演示 `var` 三种声明形式、`:=` 短声明、零值机制
    - 演示整数/浮点/复数/bool/string/byte/rune 类型
    - 演示 `const`、`iota`（跳值、位移枚举）
    - 演示显式类型转换 `T(v)`、字符串与 `[]byte`/`[]rune` 互转
    - 演示多行字面量（反引号）与转义字符串
    - 演示 nil 指针解引用触发 panic（由 Runner recover 捕获）
    - 通过注释说明与 C 的差异
    - _需求：2.1–2.12_
  - [ ]* 3.2 为零值机制编写属性测试
    - **属性 2：Go 零值机制**
    - **验证需求：2.3**
  - [ ]* 3.3 为字符串与 []byte 互转编写属性测试
    - **属性 3：字符串与 []byte 互转 round-trip**
    - **验证需求：2.9**

- [x] 4. 实现 demos/controlflow 模块
  - [x] 4.1 实现 `controlflow.go`
    - 演示 `for` 三种形式（三段式、条件式、无限循环）
    - 演示 `if/else`、带初始化语句的 `if`
    - 演示 `switch` 自动 break、`fallthrough`、无表达式形式、type switch
    - 演示 `range` 遍历数组/切片/map/字符串/channel
    - 演示 `break`/`continue` 配合标签、`goto` 语句
    - 导出 `TypeSwitchDemo(i interface{}) string` 供测试调用
    - 通过注释说明与 C 的差异
    - _需求：3.1–3.10_
  - [ ]* 4.2 为 range 遍历完整性编写属性测试
    - **属性 4：range 遍历完整性**
    - **验证需求：3.7**

- [x] 5. 实现 demos/functions 模块
  - [x] 5.1 实现 `functions.go`
    - 演示基本函数定义、多返回值、命名返回值与裸返回
    - 实现并导出 `Sum(nums ...int) int`（可变参数）
    - 实现并导出 `MakeAdder(x int) func(int) int`（闭包）
    - 实现并导出 `SafeDivide(a, b int) (result int, err error)`（defer+recover）
    - 演示函数作为一等公民、匿名函数 IIFE、闭包循环变量陷阱及修复
    - 演示 `defer` LIFO 顺序、修改命名返回值、资源释放典型用途
    - 通过注释说明与 C 的差异
    - _需求：4.1–4.11_
  - [ ]* 5.2 为可变参数求和编写属性测试
    - **属性 5：可变参数求和**
    - **验证需求：4.4**
  - [ ]* 5.3 为闭包加法器编写属性测试
    - **属性 6：闭包加法器正确性**
    - **验证需求：4.7**
  - [ ]* 5.4 为 defer LIFO 顺序编写属性测试
    - **属性 7：defer LIFO 执行顺序**
    - **验证需求：4.8**

- [x] 6. 实现 demos/collections 模块
  - [x] 6.1 实现 `collections.go`
    - 演示数组三种声明方式、数组值类型语义（复制验证）
    - 演示切片三种创建方式、`append`（单元素/切片展开）、扩容行为（cap 和指针对比）
    - 演示切片共享底层数组、`copy` 函数用法、二维切片
    - 演示 map 创建/读写/删除/存在性检查、遍历顺序随机性与有序遍历
    - 演示访问不存在 key 返回零值、nil map 写入触发 panic
    - 通过注释说明与 C 的差异
    - _需求：5.1–5.14_
  - [ ]* 6.2 为数组值类型语义编写属性测试
    - **属性 8：数组值类型语义**
    - **验证需求：5.2**
  - [ ]* 6.3 为切片 append 正确性编写属性测试
    - **属性 9：切片 append 正确性**
    - **验证需求：5.4**
  - [ ]* 6.4 为 copy 函数正确性编写属性测试
    - **属性 10：copy 函数正确性**
    - **验证需求：5.7**
  - [ ]* 6.5 为 map 读写 round-trip 编写属性测试
    - **属性 11：map 读写 round-trip**
    - **验证需求：5.10**

- [x] 7. 实现 demos/structs 模块
  - [x] 7.1 实现 `structs.go`
    - 定义并导出 `Animal` 结构体（含 json 标签），实现值接收者 `String()` 和指针接收者 `Birthday()`
    - 定义并导出 `Dog` 结构体（嵌入 Animal），演示方法提升
    - 演示结构体值类型语义、指针传递、自动取地址调用
    - 演示结构体标签与 `reflect` 读取、`json.Marshal`/`json.Unmarshal`（含 `omitempty`）
    - 通过注释说明与 C 的差异
    - _需求：6.1–6.10_
  - [ ]* 7.2 为指针接收者方法修改原始值编写属性测试
    - **属性 12：指针接收者方法修改原始值**
    - **验证需求：6.3, 6.4**

- [x] 8. 实现 demos/interfaces 模块
  - [x] 8.1 实现 `interfaces.go`
    - 定义并导出 `Shape` 接口（`Area()`、`Perimeter()`）
    - 实现 `Circle` 和 `Rectangle` 结构体，隐式实现 `Shape`
    - 演示接口多态、接口组合、空接口 `any`
    - 演示类型断言安全形式与 panic 形式、type switch
    - 演示接口值内部结构（nil 接口 vs 持有 nil 指针的接口）
    - 演示 `fmt.Stringer`、`error`、`io.Reader`/`io.Writer` 标准接口
    - 通过注释说明与 C 的差异
    - _需求：7.1–7.10_
  - [ ]* 8.2 为接口多态正确性编写属性测试
    - **属性 13：接口多态正确性**
    - **验证需求：7.3**
  - [ ]* 8.3 为类型断言 round-trip 编写属性测试
    - **属性 14：类型断言 round-trip**
    - **验证需求：7.6**

- [x] 9. 实现 demos/errors 模块
  - [x] 9.1 实现 `errors.go`
    - 定义并导出 `ErrNotFound`、`ErrDivisionByZero` 哨兵错误
    - 定义并导出 `ValidationError` 自定义错误类型
    - 实现并导出 `Divide(a, b float64) (float64, error)`
    - 演示 `errors.New`、`fmt.Errorf("%w",...)`、`errors.Is`、`errors.As`
    - 演示 panic 触发场景、`defer+recover` 转换为 error 的模式
    - 通过注释说明与 C errno 的对比
    - _需求：8.1–8.11_
  - [ ]* 9.2 为错误处理正确性编写属性测试
    - **属性 15：错误处理正确性**
    - **验证需求：8.2**
  - [ ]* 9.3 为错误链操作正确性编写属性测试
    - **属性 16：错误链操作正确性**
    - **验证需求：8.6, 8.7**

- [x] 10. 实现 demos/concurrency 模块
  - [x] 10.1 实现 `concurrency.go`
    - 演示 `go` 关键字启动 goroutine、`sync.WaitGroup` 等待完成
    - 演示无缓冲/有缓冲 channel 的创建与收发
    - 演示 `select` 语句（含 default 非阻塞分支）
    - 演示 channel 关闭与 `range` 遍历、单向 channel 类型
    - 演示 `sync.Mutex` 保护共享数据、`sync.Once` 单次初始化
    - 演示 `context.Context` 控制取消与超时
    - 演示 goroutine 泄漏场景及避免方式
    - 实现并导出 `MergeChannels[T any](cs ...<-chan T) <-chan T`
    - 通过注释说明并发哲学与 C pthread 对比
    - _需求：9.1–9.12_
  - [ ]* 10.2 为 channel 数据传输 round-trip 编写属性测试
    - **属性 17：channel 数据传输 round-trip**
    - **验证需求：9.2**
  - [ ]* 10.3 为并发计数器安全性编写属性测试
    - **属性 18：并发计数器安全性**
    - **验证需求：9.8**

- [x] 11. 检查点 —— 确保所有已实现模块的测试通过
  - 确保所有测试通过，如有问题请向用户反馈。

- [x] 12. 实现 demos/pointers 模块
  - [x] 12.1 实现 `pointers.go`
    - 演示指针声明、取地址、解引用基本操作
    - 实现并导出 `Increment(p *int)`，演示通过指针修改外部变量
    - 演示 `new(T)` 与 `var x T; &x` 的等价性
    - 演示 nil 指针检查、函数返回局部变量指针（逃逸分析）
    - 演示结构体指针字段访问语法糖 `p.Field`
    - 通过注释说明不支持指针运算、GC 安全性与内存泄漏注意事项
    - _需求：10.1–10.8_
  - [ ]* 12.2 为指针取地址与解引用 round-trip 编写属性测试
    - **属性 19：指针取地址与解引用 round-trip**
    - **验证需求：10.1, 10.2**

- [x] 13. 实现 demos/packages 模块
  - [x] 13.1 实现 `packages.go`
    - 演示包导入（标准库、内部包、别名导入）
    - 演示导出规则（首字母大写 vs 小写）
    - 演示 `init` 函数用途、执行时机与多个 init 的顺序
    - 演示 `go.mod` 结构说明、`go get` 添加依赖、`go.sum` 作用
    - 演示空白导入 `import _ "pkg"` 的用途
    - 通过注释说明与 C 头文件的对比
    - _需求：11.1–11.8_

- [x] 14. 实现 demos/generics 模块
  - [x] 14.1 实现 `generics.go`
    - 实现并导出 `Map[T, U any](s []T, f func(T) U) []U`
    - 实现并导出 `Filter[T any](s []T, f func(T) bool) []T`
    - 实现并导出 `Contains[T comparable](s []T, v T) bool`
    - 定义 `Number` 约束接口（`~int | ~float64` 等），实现 `Sum[T Number](s []T) T`
    - 定义并导出泛型 `Stack[T any]`，实现 `Push`/`Pop`/`Peek`/`Len`/`IsEmpty`
    - 演示 `~T` 底层类型约束、类型推断（无需显式指定类型参数）
    - 演示标准库 `slices` 包（Go 1.21+）中泛型函数的使用
    - 通过注释说明与 C void*/宏/C++ 模板的对比
    - _需求：12.1–12.7_
  - [ ]* 14.2 为泛型 Map 函数正确性编写属性测试
    - **属性 20：泛型 Map 函数正确性**
    - **验证需求：12.1**
  - [ ]* 14.3 为泛型 Stack LIFO 语义编写属性测试
    - **属性 21：泛型 Stack LIFO 语义**
    - **验证需求：12.4**

- [x] 15. 实现 demos/testing_demo 模块
  - [x] 15.1 实现 `testing_demo.go` 与 `testing_demo_test.go`
    - 在 `testing_demo.go` 中实现若干可测试的辅助函数（如 `Add`、`Fibonacci`）
    - 在 `testing_demo_test.go` 中演示测试文件命名约定与 `TestXxx` 签名
    - 演示 `t.Error`/`t.Errorf`/`t.Fatal`/`t.Fatalf` 的用法
    - 演示表驱动测试模式（结构体切片 + 循环）
    - 演示子测试 `t.Run` 与并行子测试 `t.Parallel()`
    - 演示基准测试 `BenchmarkXxx` 与 `b.N`
    - 演示示例测试 `ExampleXxx` 与 `// Output:` 注释
    - 通过注释说明与 C 测试框架的对比
    - _需求：13.1–13.10_

- [x] 16. 实现 demos/stdlib 模块
  - [x] 16.1 实现 `stdlib.go`
    - 演示 `fmt` 格式化动词（`%v`、`%+v`、`%#v`、`%T`、`%p` 等）
    - 演示 `strings` 包常用函数与 `strings.Builder`
    - 演示 `strconv` 包数值字符串互转
    - 演示 `os` 包文件读写、环境变量、命令行参数
    - 演示 `bufio.Scanner` 逐行读取、`io.Copy`、`io.TeeReader`
    - 演示 `encoding/json` 嵌套结构与自定义 Marshal/Unmarshal
    - 演示 `net/http` 发起 GET/POST 请求与启动简单 HTTP 服务器
    - 实现并导出 `FormatTime(t time.Time) string`、`CountWords(s string) int`
    - 演示 `time` 包格式化（Go 参考时间）、时间计算与 `Duration`
    - 演示 `math/rand`、`regexp`、`sort` 包
    - 通过注释说明标准库设计哲学与 libc 对比
    - _需求：14.1–14.12_
  - [ ]* 16.2 为 strings.Split/Join round-trip 编写属性测试
    - **属性 22：strings.Split/Join round-trip**
    - **验证需求：14.2**
  - [ ]* 16.3 为 strconv 数值字符串互转 round-trip 编写属性测试
    - **属性 23：strconv 数值字符串互转 round-trip**
    - **验证需求：14.3**
  - [ ]* 16.4 为时间格式化 round-trip 编写属性测试
    - **属性 24：时间格式化 round-trip**
    - **验证需求：14.8**

- [x] 17. 将所有模块注册到 main.go 的 Runner
  - 在 `main.go` 中导入所有 demos 包
  - 完善 `demos` 注册表，按顺序注册全部 13 个模块
  - 确保 `go run .` 能完整输出所有模块的演示结果
  - _需求：1.1, 1.2, 1.3_

- [x] 18. 最终检查点 —— 确保所有测试通过
  - 确保所有测试通过（`go test ./...`），如有问题请向用户反馈。

## 备注

- 标有 `*` 的子任务为可选任务，可跳过以加快 MVP 进度
- 每个任务均引用了具体的需求条款，便于追溯
- 属性测试使用 `pgregory.net/rapid` 库，每个属性对应设计文档中的属性编号
- 单元测试与属性测试互补，共同提供全面的测试覆盖
