// Package packages 演示 Go 语言的包与模块系统。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// 与 C 头文件的对比：
//   - C 使用 #include <header.h> 引入声明，链接时再找实现；
//   - Go 使用 import "pkg/path" 直接引入包，编译器同时处理声明和实现；
//   - Go 没有头文件，没有预处理器，包系统更简洁统一；
//   - Go 的包名与目录名通常一致，但不强制要求（包名是 package 声明的名字）。
package packages

import (
	// 标准库包：直接用包名导入
	"fmt"
	"math"

	// 别名导入：import alias "pkg/path"
	// 用途：避免包名冲突，或缩短较长的包名
	str "strings" // 给 strings 包起别名 str

	// 同项目内部包：使用模块路径（go.mod 中的 module 名）+ 子目录路径
	// 例如：import "go-basics-demo/demos/variables"
	// 注意：这里为了演示，我们不实际导入（避免循环依赖），
	// 而是通过注释和 fmt.Println 说明用法。

	// 空白导入：import _ "pkg"
	// 用途：仅执行包的 init 函数，不使用包中的任何标识符
	// 常见场景：注册数据库驱动、注册图片格式解码器等
	// 例如：import _ "github.com/lib/pq"  // 注册 PostgreSQL 驱动
	// 例如：import _ "image/png"          // 注册 PNG 图片解码器
	// 注意：这里不实际使用空白导入（避免引入不必要的依赖）
)

// initExecuted 记录 init 函数是否已执行（包级变量）。
// init 函数在包被导入时自动执行，且只执行一次。
var initExecuted bool

// initMessage 记录 init 函数执行时的消息。
var initMessage string

// init 函数：包初始化函数。
// 特点：
//   - 函数名固定为 init，无参数，无返回值；
//   - 在包被导入时自动调用，无需（也不能）手动调用；
//   - 一个包可以有多个 init 函数（甚至在同一文件中），按声明顺序执行；
//   - 执行时机：所有包级变量初始化完成后，main 函数执行前；
//   - 执行顺序：被导入包的 init 先执行，然后是当前包的 init；
//
// C 差异：
//   - C 没有自动执行的初始化函数（C++ 有全局对象构造函数）；
//   - C 通常用 __attribute__((constructor)) 或 DllMain 实现类似功能；
//   - Go 的 init 机制更简洁、更可预测。
func init() {
	initExecuted = true
	initMessage = "packages 包的 init 函数已执行（在 Demo() 调用之前）"
}

// Demo 演示所有包与模块相关内容。
func Demo() {
	demoImports()
	demoExportRules()
	demoInitFunction()
	demoGoMod()
	demoBlankImport()
}

// -----------------------------------------------------------------------------
// 1. 包导入：标准库、内部包、别名导入
// C 差异：
//   - C 的 #include 是文本替换（预处理器），Go 的 import 是编译器级别的包引用；
//   - C 需要分别 #include 头文件和链接库，Go 的 import 一步完成；
//   - Go 不允许循环导入（A 导入 B，B 导入 A），C 通过头文件保护可以处理；
//   - Go 的未使用导入是编译错误（C 只是警告），强制保持代码整洁。
// -----------------------------------------------------------------------------
func demoImports() {
	fmt.Println("\n--- 1. 包导入 ---")

	// 标准库包使用
	fmt.Printf("math.Pi = %.5f\n", math.Pi)
	fmt.Printf("math.Sqrt(2) = %.5f\n", math.Sqrt(2))

	// 使用别名导入的包（str 是 strings 的别名）
	result := str.ToUpper("hello, go")
	fmt.Printf("str.ToUpper（别名导入）: %q\n", result)

	// 正常使用 strings 包（与别名导入同一个包，演示两种方式）
	fmt.Printf("strings.Contains: %v\n", str.Contains("Go语言", "Go"))

	// 内部包导入示例（注释说明，不实际执行以避免循环依赖）
	fmt.Println("\n内部包导入示例（代码注释）:")
	fmt.Println(`  import "go-basics-demo/demos/variables"`)
	fmt.Println(`  // 然后可以调用 variables.Demo()`)
	fmt.Println("  // 模块路径 = go.mod 中的 module 名 + 子目录路径")

	// 别名导入的常见用途
	fmt.Println("\n别名导入常见用途:")
	fmt.Println(`  import myfmt "fmt"           // 避免与本地变量名冲突`)
	fmt.Println(`  import yaml "gopkg.in/yaml.v3" // 缩短包名`)
	fmt.Println(`  import . "fmt"               // 点导入：直接用 Println 而非 fmt.Println（不推荐）`)

	fmt.Println("\nC 对比：#include <stdio.h> 是文本替换，Go import 是编译器级别的包引用")
	fmt.Println("注意：Go 不允许未使用的导入（编译错误），C 只是警告")
}

// -----------------------------------------------------------------------------
// 2. 导出规则：首字母大写 vs 小写
// C 差异：
//   - C 没有语言级别的访问控制，通常用命名约定（如 _private）区分；
//   - Go 通过标识符首字母大小写实现访问控制，是语言规范的一部分；
//   - 大写开头（如 Demo、Increment）：导出（exported），包外可访问，相当于 public；
//   - 小写开头（如 demoImports、initExecuted）：未导出（unexported），包内私有，相当于private；
//   - 这一规则适用于：函数、类型、变量、常量、结构体字段、接口方法。
// -----------------------------------------------------------------------------

// ExportedConst 是导出常量（首字母大写，包外可访问）
const ExportedConst = "我是导出常量，包外可以访问"

// unexportedConst 是未导出常量（首字母小写，仅包内可访问）
const unexportedConst = "我是未导出常量，仅包内可访问"

// ExportedType 是导出类型
type ExportedType struct {
	PublicField  string // 导出字段（大写）
	privateField string // 未导出字段（小写），包外无法直接访问
}

// unexportedType 是未导出类型（包外无法直接使用）
type unexportedType struct {
	value int
}

// ExportedFunc 是导出函数
func ExportedFunc() string {
	return "我是导出函数"
}

// unexportedFunc 是未导出函数
func unexportedFunc() string {
	return "我是未导出函数"
}

func demoExportRules() {
	fmt.Println("\n--- 2. 导出规则：首字母大写 vs 小写 ---")

	fmt.Printf("导出常量: %q\n", ExportedConst)
	fmt.Printf("未导出常量（包内访问）: %q\n", unexportedConst)

	fmt.Printf("导出函数: %q\n", ExportedFunc())
	fmt.Printf("未导出函数（包内访问）: %q\n", unexportedFunc())

	// 导出类型的使用
	et := ExportedType{
		PublicField:  "公开字段",
		privateField: "私有字段（包内可访问）",
	}
	fmt.Printf("导出类型: PublicField=%q\n", et.PublicField)
	fmt.Printf("未导出字段（包内访问）: privateField=%q\n", et.privateField)

	// 未导出类型（包内可用）
	ut := unexportedType{value: 42}
	fmt.Printf("未导出类型（包内访问）: value=%d\n", ut.value)

	fmt.Println("\n规则总结:")
	fmt.Println("  大写开头 → 导出（public）：包外可访问")
	fmt.Println("  小写开头 → 未导出（package-private）：仅包内可访问")
	fmt.Println("  适用于：函数、类型、变量、常量、结构体字段、接口方法")
	fmt.Println("C 对比：C 没有语言级访问控制，Go 通过命名规范强制实现封装")
}

// -----------------------------------------------------------------------------
// 3. init 函数：包初始化
// C 差异：
//   - C 没有自动执行的初始化函数（C++ 有全局构造函数）；
//   - Go 的 init 在 import 时自动调用，无需手动调用；
//   - 多个 init 函数按声明顺序执行（同一文件内），不同文件按文件名字母序；
//   - init 函数不能被显式调用（编译错误）。
// -----------------------------------------------------------------------------

// secondInit 用于演示同一包中多个 init 函数
var secondInitExecuted bool

func init() {
	// 同一包可以有多个 init 函数，按声明顺序执行
	secondInitExecuted = true
}

func demoInitFunction() {
	fmt.Println("\n--- 3. init 函数：包初始化 ---")

	fmt.Printf("第一个 init 已执行: %v\n", initExecuted)
	fmt.Printf("init 消息: %q\n", initMessage)
	fmt.Printf("第二个 init 已执行: %v\n", secondInitExecuted)

	fmt.Println("\ninit 函数特点:")
	fmt.Println("  1. 函数名固定为 init，无参数，无返回值")
	fmt.Println("  2. 在包被导入时自动调用，无需手动调用")
	fmt.Println("  3. 一个包可以有多个 init 函数，按声明顺序执行")
	fmt.Println("  4. 执行时机：包级变量初始化完成后，main() 执行前")
	fmt.Println("  5. 不能显式调用 init()（编译错误）")
	fmt.Println("  6. 常用于：注册驱动、初始化全局状态、验证配置")

	fmt.Println("\n执行顺序示例:")
	fmt.Println("  导入包 A → A 的包级变量初始化 → A 的 init() 执行")
	fmt.Println("  导入包 B → B 的包级变量初始化 → B 的 init() 执行")
	fmt.Println("  当前包的包级变量初始化 → 当前包的 init() 执行")
	fmt.Println("  main() 执行")

	fmt.Println("\nC 对比：C 没有自动初始化函数，C++ 有全局对象构造函数（顺序不确定）")
}

// -----------------------------------------------------------------------------
// 4. go.mod 文件结构说明
// C 差异：
//   - C 没有官方的包管理系统，通常依赖系统包管理器或手动管理；
//   - Go 模块系统（Go Modules）是官方的依赖管理方案，从 Go 1.11 引入；
//   - go.mod 类似于 Node.js 的 package.json 或 Rust 的 Cargo.toml。
// -----------------------------------------------------------------------------
func demoGoMod() {
	fmt.Println("\n--- 4. go.mod 文件结构说明 ---")

	fmt.Println("go.mod 文件示例：")
	fmt.Println(`  module go-basics-demo    // 模块路径（唯一标识符，通常是仓库路径）`)
	fmt.Println(``)
	fmt.Println(`  go 1.21                  // 最低 Go 版本要求`)
	fmt.Println(``)
	fmt.Println(`  require (`)
	fmt.Println(`      pgregory.net/rapid v1.2.0  // 直接依赖`)
	fmt.Println(`  )`)

	fmt.Println("\ngo.mod 各部分说明:")
	fmt.Println("  module：模块路径，是所有内部包的导入路径前缀")
	fmt.Println("          例如：import \"go-basics-demo/demos/variables\"")
	fmt.Println("  go：指定最低 Go 版本，影响语言特性和标准库可用性")
	fmt.Println("  require：列出所有直接和间接依赖及其版本")
	fmt.Println("  replace：将某个依赖替换为本地路径或其他版本（调试用）")
	fmt.Println("  exclude：排除特定版本（有安全漏洞时使用）")

	fmt.Println("\ngo.sum 文件说明:")
	fmt.Println("  go.sum 记录每个依赖的加密哈希值（SHA-256）")
	fmt.Println("  用于验证下载的依赖未被篡改（安全保证）")
	fmt.Println("  应提交到版本控制（与 go.mod 一起）")

	fmt.Println("\n常用 go 命令:")
	fmt.Println("  go mod init <module>  // 初始化新模块")
	fmt.Println("  go get <pkg>@<version> // 添加/更新依赖")
	fmt.Println("  go mod tidy           // 清理未使用的依赖")
	fmt.Println("  go mod download       // 下载所有依赖到本地缓存")
	fmt.Println("  go mod vendor         // 将依赖复制到 vendor 目录")

	fmt.Println("\nC 对比：C 没有官方包管理，Go Modules 类似 npm/cargo，统一管理依赖版本")
}

// -----------------------------------------------------------------------------
// 5. 空白导入 import _ "pkg" 的用途
// C 差异：
//   - C 没有对应概念，通常通过链接特定库来触发初始化代码；
//   - Go 的空白导入是一种显式的"仅执行 init"机制，意图清晰；
//   - 最常见用途：注册数据库驱动（database/sql 驱动模型）。
// -----------------------------------------------------------------------------
func demoBlankImport() {
	fmt.Println("\n--- 5. 空白导入 import _ \"pkg\" ---")

	fmt.Println("空白导入语法：import _ \"pkg/path\"")
	fmt.Println("用途：仅执行包的 init 函数，不使用包中的任何标识符")
	fmt.Println("编译器不会报「未使用的导入」错误（_ 表示明确忽略）")

	fmt.Println("\n常见使用场景：")

	fmt.Println("\n1. 注册数据库驱动（最典型用途）:")
	fmt.Println(`   import _ "github.com/lib/pq"           // PostgreSQL 驱动`)
	fmt.Println(`   import _ "github.com/go-sql-driver/mysql" // MySQL 驱动`)
	fmt.Println(`   // 驱动的 init() 会调用 sql.Register() 注册自己`)
	fmt.Println(`   // 之后可以用 sql.Open("postgres", ...) 使用`)

	fmt.Println("\n2. 注册图片格式解码器:")
	fmt.Println(`   import _ "image/png"   // 注册 PNG 解码器`)
	fmt.Println(`   import _ "image/jpeg"  // 注册 JPEG 解码器`)
	fmt.Println(`   // 之后 image.Decode() 就能处理这些格式`)

	fmt.Println("\n3. 注册 HTTP 处理器或中间件:")
	fmt.Println(`   import _ "net/http/pprof"  // 注册 pprof 性能分析 HTTP 端点`)

	fmt.Println("\n4. 触发包的副作用初始化（如全局注册表）:")
	fmt.Println(`   import _ "myapp/plugins/auth"  // 注册认证插件`)

	fmt.Println("\n工作原理：")
	fmt.Println("  包的 init() 函数在 import 时自动执行")
	fmt.Println("  空白导入确保 init() 执行，但不引入任何名称到当前包的命名空间")

	fmt.Println("\nC 对比：C 通过链接特定 .o 文件触发初始化，Go 的空白导入更显式、更安全")
}
