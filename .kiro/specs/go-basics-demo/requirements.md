# 需求文档

## 简介

本项目是一个面向有 C 语言基础的开发者的 Go 语言入门示例项目。项目通过一系列独立但相互关联的示例模块，系统地演示 Go 语言的核心概念与惯用法，帮助开发者快速掌握 Go 与 C 的异同，建立 Go 语言思维方式。每个模块均包含可运行的代码示例、详细注释以及与 C 语言的对比说明。

## 词汇表

- **Demo_App**：整个 Go 入门示例应用程序
- **Module**：每个独立的知识点示例模块（如变量模块、函数模块等）
- **Runner**：负责按顺序执行各示例模块并输出结果的主程序
- **Goroutine**：Go 语言的轻量级并发执行单元，由 Go 运行时调度
- **Channel**：Go 语言用于 goroutine 间通信的类型安全管道
- **Interface**：Go 语言的接口类型，定义方法集合，隐式实现
- **Struct**：Go 语言的结构体类型，可附加方法
- **Slice**：Go 语言的动态数组，底层由数组、长度、容量三元组构成
- **Goroutine_Scheduler**：Go 运行时的 M:N 调度器，将 goroutine 映射到 OS 线程
- **Zero_Value**：Go 中变量声明后未赋值时的默认值
- **Type_Assertion**：从接口值中提取具体类型的操作
- **Defer_Stack**：存储 defer 调用的后进先出栈结构

---

## 需求

### 需求 1：项目结构与入口

**用户故事：** 作为一名有 C 语言基础的开发者，我希望有一个清晰的项目入口，能够依次运行所有示例模块，以便我能系统地学习 Go 语言基础。

#### 验收标准

1. THE Demo_App SHALL 提供一个 `main.go` 作为唯一入口文件，位于项目根目录
2. THE Runner SHALL 按顺序调用每个示例模块的演示函数，每个模块调用前打印分隔线和模块名称
3. WHEN 运行 `go run .` 或 `go run main.go` 时，THE Demo_App SHALL 在终端输出所有模块的演示结果
4. THE Demo_App SHALL 使用 `go mod init` 初始化为独立的 Go module，module 名称为 `go-basics-demo`
5. THE Demo_App SHALL 将每个示例模块组织为 `demos/` 子目录下的独立包，例如 `demos/variables`、`demos/functions`
6. THE Demo_App SHALL 提供 `README.md` 说明如何运行项目及各模块的学习顺序
7. WHEN 任意模块的演示函数发生 panic 时，THE Runner SHALL 使用 recover 捕获并打印错误信息，继续执行后续模块

---

### 需求 2：变量与基本类型

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的变量声明与基本数据类型，以便我能理解 Go 与 C 在类型系统上的差异。

#### 验收标准

1. THE Module SHALL 演示使用 `var` 关键字声明变量的三种形式：仅声明类型、声明并赋值、批量声明 `var ( ... )`
2. THE Module SHALL 演示使用 `:=` 短变量声明的方式，并说明其只能在函数体内使用
3. THE Module SHALL 演示 Go 的零值机制：`int` 为 `0`，`float64` 为 `0.0`，`string` 为 `""`，`bool` 为 `false`，指针/切片/map/channel 为 `nil`
4. THE Module SHALL 演示基本整数类型：`int`、`int8`、`int16`、`int32`、`int64` 及对应无符号类型，并说明 `int` 的平台相关宽度（32位或64位）
5. THE Module SHALL 演示浮点类型 `float32` 和 `float64`，以及复数类型 `complex64` 和 `complex128`
6. THE Module SHALL 演示 `bool`、`string`、`byte`（`uint8` 别名）、`rune`（`int32` 别名，表示 Unicode 码点）
7. THE Module SHALL 演示常量 `const` 的声明，包括无类型常量（untyped constant）的灵活性
8. THE Module SHALL 演示 `iota` 枚举用法，包括跳值、位移枚举（`1 << iota`）等场景
9. THE Module SHALL 演示显式类型转换语法 `T(v)`，并演示整数与浮点、字符串与 `[]byte`/`[]rune` 之间的转换
10. THE Module SHALL 通过注释说明与 C 语言的关键差异：Go 无隐式类型转换、无 `unsigned` 修饰符、`string` 是不可变值类型而非字符数组
11. THE Module SHALL 演示字符串的多行字面量（反引号 `` ` ``）与转义字符串（双引号）的区别
12. IF 对 `nil` 指针进行解引用，THEN THE Demo_App SHALL 触发 panic 并由 Runner 捕获，演示 nil 安全的重要性

---

### 需求 3：控制流

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的控制流语句，以便我能快速适应 Go 的语法风格。

#### 验收标准

1. THE Module SHALL 演示 `for` 循环的三种形式：标准三段式 `for i := 0; i < n; i++`、仅条件式 `for condition`（等价于 C 的 `while`）、无限循环 `for {}`
2. THE Module SHALL 演示 `if/else` 语句，条件表达式无需括号，并说明花括号是强制要求的
3. THE Module SHALL 演示带初始化语句的 `if` 形式：`if err := doSomething(); err != nil { ... }`，说明初始化变量的作用域仅限于 if/else 块
4. THE Module SHALL 演示 `switch` 语句的自动 `break` 行为，以及使用 `fallthrough` 显式穿透的用法
5. THE Module SHALL 演示 `switch` 的无表达式形式（等价于 `if-else if` 链）
6. THE Module SHALL 演示 `switch` 对类型的匹配（type switch），结合接口使用
7. THE Module SHALL 演示 `range` 遍历数组、切片（获取 index 和 value）、map（获取 key 和 value）、字符串（获取 rune）、channel
8. THE Module SHALL 演示 `break` 和 `continue` 配合标签（label）跳出多层循环的用法
9. THE Module SHALL 演示 `goto` 语句，并通过注释说明其在 Go 中的限制与不推荐使用的原因
10. THE Module SHALL 通过注释说明 Go 没有 `while` 和 `do-while` 关键字，以及 `for range` 与 C 的 `for` 循环的对应关系

---

### 需求 4：函数

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 函数的特性，以便我能掌握 Go 函数与 C 函数的核心差异。

#### 验收标准

1. THE Module SHALL 演示基本函数定义与调用，说明参数类型后置的语法（与 C 相反）
2. THE Module SHALL 演示多返回值函数，例如同时返回计算结果和 error
3. THE Module SHALL 演示命名返回值（named return values）与裸返回（naked return）的用法及注意事项
4. THE Module SHALL 演示可变参数函数（variadic function）`func sum(nums ...int)`，以及用 `slice...` 展开切片传入的方式
5. THE Module SHALL 演示函数作为一等公民：将函数赋值给变量、定义函数类型、作为参数传递、作为返回值返回
6. THE Module SHALL 演示匿名函数的立即调用（IIFE）与赋值给变量后调用两种形式
7. THE Module SHALL 演示闭包捕获外部变量的行为，包括循环变量捕获的常见陷阱及修复方式
8. THE Module SHALL 演示 `defer` 语句的后进先出执行顺序，通过多个 defer 的示例验证
9. THE Module SHALL 演示 `defer` 在函数返回时修改命名返回值的行为
10. THE Module SHALL 演示 `defer` 的典型用途：资源释放（文件关闭、锁释放）、panic 恢复
11. THE Module SHALL 通过注释说明 Go 函数不支持默认参数和函数重载，以及与 C 函数指针的对比

---

### 需求 5：数组、切片与 Map

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的集合类型，以便我能理解切片与 C 数组的本质区别。

#### 验收标准

1. THE Module SHALL 演示数组的声明方式：`var a [5]int`、`a := [3]int{1,2,3}`、`a := [...]int{1,2,3}`（编译器推断长度）
2. THE Module SHALL 演示数组是值类型：赋值和传参时会复制整个数组，并通过示例验证修改副本不影响原数组
3. THE Module SHALL 演示切片的三种创建方式：字面量 `[]int{1,2,3}`、`make([]int, len, cap)`、从数组或切片截取 `a[low:high]`
4. THE Module SHALL 演示切片的 `append` 操作，包括追加单个元素、追加另一个切片（`append(a, b...)`）
5. THE Module SHALL 演示切片扩容时底层数组的重新分配，通过对比 append 前后的 `cap` 和底层指针说明
6. THE Module SHALL 演示切片共享底层数组的行为：修改子切片会影响原切片，并说明使用 `copy` 避免此问题
7. THE Module SHALL 演示 `copy(dst, src)` 函数的用法及返回值（实际复制的元素数量）
8. THE Module SHALL 演示二维切片的创建与访问
9. THE Module SHALL 演示 `map` 的创建：`make(map[K]V)` 和字面量 `map[K]V{k:v}`
10. THE Module SHALL 演示 map 的读写、删除（`delete`）与存在性检查（`v, ok := m[k]`）
11. THE Module SHALL 演示遍历 map 时键的顺序是随机的，并说明如何获得有序遍历（排序 key 后遍历）
12. IF 访问 map 中不存在的键，THEN THE Module SHALL 演示返回值类型的零值而非 panic
13. IF 向未初始化（nil）的 map 写入数据，THEN THE Module SHALL 演示触发 panic，并说明必须先用 make 初始化
14. THE Module SHALL 通过注释说明切片与 C 数组指针的本质区别：切片携带长度和容量信息，越界访问会 panic 而非未定义行为

---

### 需求 6：结构体与方法

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的结构体与方法，以便我能理解 Go 的面向对象风格与 C struct 的差异。

#### 验收标准

1. THE Module SHALL 演示结构体的定义与两种初始化方式：按字段名初始化和按位置初始化，并说明推荐使用字段名方式
2. THE Module SHALL 演示结构体是值类型，赋值时会复制，并演示通过指针传递结构体以避免复制开销
3. THE Module SHALL 演示为结构体定义值接收者方法 `func (s MyStruct) Method()`，说明方法内操作的是副本
4. THE Module SHALL 演示为结构体定义指针接收者方法 `func (s *MyStruct) Method()`，说明方法内可修改原始值
5. THE Module SHALL 演示 Go 自动处理值与指针接收者的调用：对值变量调用指针接收者方法时自动取地址
6. THE Module SHALL 演示结构体嵌套（组合）：将一个结构体作为另一个结构体的字段，实现代码复用
7. THE Module SHALL 演示匿名字段（嵌入）：`type Dog struct { Animal }` 形式，以及方法提升（promoted methods）
8. THE Module SHALL 演示结构体标签（struct tag）的定义与通过 `reflect` 包读取，说明其在 JSON 序列化中的典型用途
9. THE Module SHALL 演示结构体与 JSON 的互转：`json.Marshal` 和 `json.Unmarshal`，包括字段名映射和 `omitempty` 标签
10. THE Module SHALL 通过注释说明 Go 没有类（class）、继承（inheritance），用结构体+方法+组合实现面向对象，以及与 C struct 的对比

---

### 需求 7：接口

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的接口机制，以便我能掌握 Go 隐式接口实现的设计哲学。

#### 验收标准

1. THE Module SHALL 演示接口的定义：`type Writer interface { Write(p []byte) (n int, err error) }`
2. THE Module SHALL 演示结构体隐式实现接口：只要实现了接口的所有方法即视为实现，无需 `implements` 关键字
3. THE Module SHALL 演示接口变量的多态用法：同一接口变量持有不同具体类型，调用相同方法产生不同行为
4. THE Module SHALL 演示接口组合：通过嵌入多个接口定义新接口，例如 `io.ReadWriter`
5. THE Module SHALL 演示空接口 `interface{}` 或 `any`（Go 1.18+）的用法，以及作为通用容器的场景
6. THE Module SHALL 演示类型断言 `v, ok := i.(ConcreteType)` 的安全形式，以及不带 ok 时 panic 的情况
7. THE Module SHALL 演示类型选择（type switch）`switch v := i.(type)` 处理多种类型的模式
8. THE Module SHALL 演示接口值的内部结构（类型+值对），以及接口值为 nil 与持有 nil 指针的接口值的区别
9. THE Module SHALL 演示标准库中常用接口：`fmt.Stringer`（实现 `String() string`）、`error`、`io.Reader`、`io.Writer`
10. THE Module SHALL 通过注释说明 Go 接口与 C 函数指针结构体（虚函数表）的概念对比，以及鸭子类型的设计哲学

---

### 需求 8：错误处理

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的错误处理模式，以便我能掌握 Go 用返回值处理错误的惯用法。

#### 验收标准

1. THE Module SHALL 演示 `error` 接口的定义：`type error interface { Error() string }`
2. THE Module SHALL 演示函数返回 `(result, error)` 的惯用模式，以及调用方检查 error 的标准写法
3. THE Module SHALL 演示使用 `errors.New("message")` 创建简单错误
4. THE Module SHALL 演示使用 `fmt.Errorf("context: %w", err)` 包装错误（wrapping），保留错误链
5. THE Module SHALL 演示自定义错误类型：定义结构体并实现 `Error() string` 方法，携带额外上下文信息
6. THE Module SHALL 演示 `errors.Is(err, target)` 检查错误链中是否包含特定错误值（哨兵错误）
7. THE Module SHALL 演示 `errors.As(err, &target)` 从错误链中提取特定类型的错误
8. THE Module SHALL 演示 `panic` 的触发场景：程序员错误（如索引越界）与主动调用
9. THE Module SHALL 演示 `recover` 在 `defer` 中捕获 panic，将其转换为 error 返回的模式
10. THE Module SHALL 演示哨兵错误（sentinel error）模式：定义包级别的 `var ErrNotFound = errors.New(...)`
11. THE Module SHALL 通过注释说明 Go 错误处理与 C 返回码（errno）模式的对比，以及 Go 不使用异常的设计原因

---

### 需求 9：并发——Goroutine 与 Channel

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的并发模型，以便我能掌握 goroutine 和 channel 的基本用法。

#### 验收标准

1. THE Module SHALL 演示使用 `go` 关键字启动 goroutine，并说明 goroutine 的轻量级特性（初始栈约 2KB，可动态增长）
2. THE Module SHALL 演示无缓冲 channel 的创建 `make(chan T)` 与同步收发操作，说明发送方和接收方必须同时就绪
3. THE Module SHALL 演示有缓冲 channel `make(chan T, n)` 的创建与使用，说明缓冲满时发送阻塞的行为
4. THE Module SHALL 演示使用 `sync.WaitGroup` 等待多个 goroutine 完成的标准模式
5. THE Module SHALL 演示使用 `select` 语句同时监听多个 channel，包括带 `default` 分支的非阻塞 select
6. THE Module SHALL 演示 channel 的关闭 `close(ch)` 与用 `range` 遍历直到 channel 关闭
7. THE Module SHALL 演示单向 channel 类型 `chan<- T`（只写）和 `<-chan T`（只读）在函数参数中的用法
8. THE Module SHALL 演示使用 `sync.Mutex` 保护共享数据，对比 channel 方案说明各自适用场景
9. THE Module SHALL 演示使用 `sync.Once` 实现只执行一次的初始化
10. THE Module SHALL 演示使用 `context.Context` 控制 goroutine 的取消与超时
11. THE Module SHALL 演示 goroutine 泄漏的典型场景，并说明如何通过 channel 关闭或 context 取消来避免
12. THE Module SHALL 通过注释说明 Go 的并发哲学："不要通过共享内存来通信，而要通过通信来共享内存"，以及与 C pthread 的对比

---

### 需求 10：指针

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的指针用法，以便我能理解 Go 指针与 C 指针的异同。

#### 验收标准

1. THE Module SHALL 演示指针的声明 `var p *int`、取地址 `p = &x`、解引用 `*p` 的基本操作
2. THE Module SHALL 演示通过指针修改函数外部变量的值，对比值传递与指针传递的效果
3. THE Module SHALL 演示 `new(T)` 函数分配零值内存并返回指针，对比与 `var x T; &x` 的等价性
4. THE Module SHALL 演示指针的零值为 `nil`，以及在解引用前检查 nil 的必要性
5. THE Module SHALL 演示函数返回局部变量的指针（逃逸分析），说明 Go 编译器会自动将其分配到堆上
6. THE Module SHALL 演示指向结构体的指针访问字段时可省略解引用：`p.Field` 等价于 `(*p).Field`
7. THE Module SHALL 通过注释说明 Go 指针不支持指针运算（`p++`、`p+1` 均非法），与 C 的关键差异
8. THE Module SHALL 通过注释说明 Go 的垃圾回收机制使指针使用更安全，无需手动 `free`，但需注意避免意外持有大对象引用导致内存泄漏

---

### 需求 11：包与模块

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 的包管理机制，以便我能理解 Go 模块化与 C 头文件的差异。

#### 验收标准

1. THE Demo_App SHALL 将各示例模块组织为 `demos/` 子目录下的独立 Go 包，每个包名与目录名一致
2. THE Module SHALL 演示包的导入：标准库包、同项目内部包，以及使用别名导入 `import alias "pkg/path"`
3. THE Module SHALL 演示 Go 的导出规则：首字母大写的标识符为导出（public），小写为包内私有（package-private）
4. THE Module SHALL 演示 `init` 函数的用途（包初始化）、执行时机（import 时自动调用）及多个 init 函数的执行顺序
5. THE Module SHALL 演示 `go.mod` 文件的结构：module 路径、Go 版本、依赖声明
6. THE Module SHALL 演示使用 `go get` 添加第三方依赖，以及 `go.sum` 文件的作用（依赖校验）
7. THE Module SHALL 演示空白导入 `import _ "pkg"` 的用途（仅执行 init 函数，如注册数据库驱动）
8. THE Module SHALL 通过注释说明 Go 包与 C 头文件的概念对比：Go 包是编译单元，无需头文件，循环依赖在编译时报错

---

### 需求 12：泛型（Generics）

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 1.18 引入的泛型特性，以便我能编写类型安全的通用代码，理解其与 C 宏和 void* 的本质区别。

#### 验收标准

1. THE Module SHALL 演示泛型函数的定义语法：`func Map[T, U any](s []T, f func(T) U) []U`
2. THE Module SHALL 演示类型参数约束（constraint）：使用 `any`、`comparable`，以及自定义约束接口
3. THE Module SHALL 演示使用 `~T` 语法定义底层类型约束，例如约束所有底层类型为 `int` 的类型
4. THE Module SHALL 演示泛型结构体的定义：`type Stack[T any] struct { items []T }`，并实现 Push/Pop 方法
5. THE Module SHALL 演示 `golang.org/x/exp/slices` 或标准库 `slices`（Go 1.21+）中泛型函数的使用
6. THE Module SHALL 演示类型推断：调用泛型函数时编译器自动推断类型参数，无需显式指定
7. THE Module SHALL 通过注释说明 Go 泛型与 C 的 `void*`（类型不安全）和宏（无类型检查）的对比，以及与 C++ 模板的异同

---

### 需求 13：测试

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 内置的测试框架，以便我能为 Go 代码编写单元测试，理解 Go 测试与 C 单元测试框架的差异。

#### 验收标准

1. THE Module SHALL 演示测试文件的命名约定：`xxx_test.go`，以及测试函数的签名 `func TestXxx(t *testing.T)`
2. THE Module SHALL 演示使用 `t.Error`、`t.Errorf`、`t.Fatal`、`t.Fatalf` 报告测试失败
3. THE Module SHALL 演示表驱动测试（table-driven test）模式：使用结构体切片定义测试用例，循环执行
4. THE Module SHALL 演示子测试 `t.Run("name", func(t *testing.T) {...})` 的用法，支持并行子测试 `t.Parallel()`
5. THE Module SHALL 演示基准测试（benchmark）：`func BenchmarkXxx(b *testing.B)`，使用 `b.N` 控制循环次数
6. THE Module SHALL 演示示例测试（example test）：`func ExampleXxx()`，通过 `// Output:` 注释验证输出
7. THE Module SHALL 演示使用 `testing/iotest`、`net/http/httptest` 等标准库测试辅助包
8. WHEN 运行 `go test ./...` 时，THE Demo_App SHALL 所有测试通过
9. WHEN 运行 `go test -bench=. ./...` 时，THE Demo_App SHALL 输出基准测试结果
10. THE Module SHALL 通过注释说明 Go 测试无需第三方框架，内置于工具链，以及与 C 的 CUnit/Check 等框架的对比

---

### 需求 14：标准库常用包

**用户故事：** 作为一名有 C 语言基础的开发者，我希望了解 Go 标准库中最常用的包，以便我能在日常开发中快速找到所需功能，减少对第三方库的依赖。

#### 验收标准

1. THE Module SHALL 演示 `fmt` 包：`Printf`/`Sprintf`/`Fprintf` 的格式化动词（`%v`、`%+v`、`%#v`、`%T`、`%p` 等）
2. THE Module SHALL 演示 `strings` 包：`Contains`、`HasPrefix`、`Split`、`Join`、`TrimSpace`、`Builder`（高效字符串拼接）
3. THE Module SHALL 演示 `strconv` 包：`Itoa`、`Atoi`、`ParseFloat`、`FormatFloat` 等数值与字符串互转
4. THE Module SHALL 演示 `os` 包：读写文件（`os.ReadFile`、`os.WriteFile`）、获取环境变量（`os.Getenv`）、命令行参数（`os.Args`）
5. THE Module SHALL 演示 `io` 和 `bufio` 包：`bufio.Scanner` 逐行读取、`io.Copy`、`io.TeeReader`
6. THE Module SHALL 演示 `encoding/json` 包：结构体与 JSON 互转，包括处理嵌套结构和自定义 `MarshalJSON`/`UnmarshalJSON`
7. THE Module SHALL 演示 `net/http` 包：发起 HTTP GET/POST 请求，以及启动简单 HTTP 服务器（`http.HandleFunc`、`http.ListenAndServe`）
8. THE Module SHALL 演示 `time` 包：获取当前时间、时间格式化（Go 特有的参考时间 `2006-01-02 15:04:05`）、时间计算与 `time.Duration`
9. THE Module SHALL 演示 `math/rand` 包（或 Go 1.20+ 的 `math/rand/v2`）：生成随机数，以及设置随机种子
10. THE Module SHALL 演示 `regexp` 包：编译正则表达式、匹配、查找、替换
11. THE Module SHALL 演示 `sort` 包：对基本类型切片排序，以及使用 `sort.Slice` 自定义排序
12. THE Module SHALL 通过注释说明 Go 标准库的设计哲学：小而精，接口驱动，与 C 标准库（libc）的对比
