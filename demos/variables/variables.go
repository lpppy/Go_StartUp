// Package variables 演示 Go 语言的变量声明与基本数据类型。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
package variables

import (
	"fmt"
	"unsafe"
)

// Demo 演示所有变量与类型相关内容。
// 注意：函数末尾会故意触发 nil 指针解引用 panic，由 Runner 的 recover 捕获。
func Demo() {
	demoVarDeclarations()
	demoShortDeclaration()
	demoZeroValues()
	demoIntegerTypes()
	demoFloatAndComplex()
	demoBoolStringByteRune()
	demoConst()
	demoIota()
	demoTypeConversion()
	demoStringLiterals()
	demoNilPanic() // 放在最后，会触发 panic
}

// -----------------------------------------------------------------------------
// 1. var 三种声明形式
// C 差异：Go 变量声明语法为 var name type，类型在变量名之后；
//         C 为 type name，类型在前。
// -----------------------------------------------------------------------------
func demoVarDeclarations() {
	fmt.Println("\n--- 1. var 声明形式 ---")

	// 形式一：仅声明类型，不赋值（变量自动获得零值）
	// C 等价：int a;  （但 C 局部变量不保证零值，Go 保证）
	var a int
	fmt.Printf("仅声明类型: var a int => a = %d\n", a)

	// 形式二：声明并赋值
	// C 等价：int b = 42;
	var b int = 42
	fmt.Printf("声明并赋值: var b int = 42 => b = %d\n", b)

	// 形式三：批量声明 var ( ... )
	// C 没有直接等价语法，通常逐行声明
	var (
		name    string  = "Gopher"
		age     int     = 10
		version float64 = 1.21
	)
	fmt.Printf("批量声明: name=%s, age=%d, version=%.2f\n", name, age, version)
}

// -----------------------------------------------------------------------------
// 2. := 短变量声明
// C 差异：C 没有 := 语法；Go 的 := 只能在函数体内使用，
//         不能用于包级别变量（包级别必须用 var）。
// -----------------------------------------------------------------------------
func demoShortDeclaration() {
	fmt.Println("\n--- 2. := 短变量声明（只能在函数体内使用）---")

	// := 自动推断类型，等价于 var x int = 10
	x := 10
	y := "hello"
	z := 3.14

	fmt.Printf(":= 短声明: x=%d (int), y=%q (string), z=%f (float64)\n", x, y, z)

	// 多变量短声明
	i, j := 1, 2
	fmt.Printf("多变量短声明: i=%d, j=%d\n", i, j)

	// 重新赋值：只要左侧至少有一个新变量，:= 就合法
	i, k := 100, 200
	fmt.Printf("重新赋值（i 已存在，k 是新变量）: i=%d, k=%d\n", i, k)
}

// -----------------------------------------------------------------------------
// 3. 零值机制
// C 差异：C 局部变量未初始化时值不确定（未定义行为）；
//         Go 保证所有变量在声明时自动初始化为零值，消除了未初始化变量的安全隐患。
// -----------------------------------------------------------------------------
func demoZeroValues() {
	fmt.Println("\n--- 3. 零值机制 ---")

	var i int
	var f float64
	var s string
	var b bool
	var p *int        // 指针零值
	var sl []int      // 切片零值
	var m map[string]int // map 零值
	var ch chan int   // channel 零值

	fmt.Printf("int    零值: %d\n", i)
	fmt.Printf("float64 零值: %f\n", f)
	fmt.Printf("string 零值: %q\n", s)
	fmt.Printf("bool   零值: %t\n", b)
	fmt.Printf("*int   零值: %v (nil)\n", p)
	fmt.Printf("[]int  零值: %v (nil)\n", sl)
	fmt.Printf("map    零值: %v (nil)\n", m)
	fmt.Printf("chan   零值: %v (nil)\n", ch)
}

// -----------------------------------------------------------------------------
// 4. 整数类型
// C 差异：
//   - C 的 int 宽度依赖编译器和平台（通常 32 位），Go 的 int 明确为平台原生宽度
//     （32 位平台为 32 位，64 位平台为 64 位）。
//   - Go 没有 unsigned 修饰符，无符号类型有独立名称（uint, uint8 等）。
//   - Go 不同整数类型之间不能隐式转换，必须显式转换。
// -----------------------------------------------------------------------------
func demoIntegerTypes() {
	fmt.Println("\n--- 4. 整数类型 ---")

	var i8 int8 = 127          // -128 ~ 127
	var i16 int16 = 32767      // -32768 ~ 32767
	var i32 int32 = 2147483647 // -2^31 ~ 2^31-1
	var i64 int64 = 9223372036854775807

	var u8 uint8 = 255
	var u16 uint16 = 65535
	var u32 uint32 = 4294967295
	var u64 uint64 = 18446744073709551615

	// int 的宽度与平台相关（32 位或 64 位）
	var n int = 42
	fmt.Printf("int 在当前平台的宽度: %d 字节 (%d 位)\n",
		unsafe.Sizeof(n), unsafe.Sizeof(n)*8)

	fmt.Printf("int8:  %d\n", i8)
	fmt.Printf("int16: %d\n", i16)
	fmt.Printf("int32: %d\n", i32)
	fmt.Printf("int64: %d\n", i64)
	fmt.Printf("uint8:  %d\n", u8)
	fmt.Printf("uint16: %d\n", u16)
	fmt.Printf("uint32: %d\n", u32)
	fmt.Printf("uint64: %d\n", u64)

	// uintptr：足以存储指针值的无符号整数，用于底层指针运算
	var ptr uintptr = uintptr(unsafe.Pointer(&n))
	fmt.Printf("uintptr（指针地址）: 0x%x\n", ptr)
}

// -----------------------------------------------------------------------------
// 5. 浮点类型与复数类型
// C 差异：
//   - Go 的 float32/float64 对应 C 的 float/double。
//   - Go 内置复数类型 complex64/complex128，C 需要 <complex.h>。
//   - Go 使用 real()/imag() 内置函数提取实部和虚部。
// -----------------------------------------------------------------------------
func demoFloatAndComplex() {
	fmt.Println("\n--- 5. 浮点与复数类型 ---")

	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793

	fmt.Printf("float32: %f (精度约 7 位有效数字)\n", f32)
	fmt.Printf("float64: %f (精度约 15 位有效数字，推荐使用)\n", f64)

	// 复数：complex64 = float32 实部 + float32 虚部
	//       complex128 = float64 实部 + float64 虚部
	var c64 complex64 = 1 + 2i
	var c128 complex128 = complex(3.0, 4.0) // 使用内置 complex() 构造

	fmt.Printf("complex64:  %v, 实部=%v, 虚部=%v\n", c64, real(c64), imag(c64))
	fmt.Printf("complex128: %v, 实部=%v, 虚部=%v\n", c128, real(c128), imag(c128))
}

// -----------------------------------------------------------------------------
// 6. bool, string, byte, rune
// C 差异：
//   - Go 的 bool 只有 true/false，不能与整数互转（C 中 0 为假，非 0 为真）。
//   - Go 的 string 是不可变的值类型（底层为只读字节序列 + 长度），
//     不是以 '\0' 结尾的字符数组（C 风格字符串）。
//   - byte 是 uint8 的别名，表示单个字节。
//   - rune 是 int32 的别名，表示一个 Unicode 码点（类似 C 的 wchar_t，但固定 32 位）。
// -----------------------------------------------------------------------------
func demoBoolStringByteRune() {
	fmt.Println("\n--- 6. bool, string, byte, rune ---")

	// bool
	var t bool = true
	var f bool = false
	fmt.Printf("bool: t=%t, f=%t\n", t, f)

	// string：不可变值类型，len() 返回字节数而非字符数
	s := "Hello, 世界"
	fmt.Printf("string: %q\n", s)
	fmt.Printf("len(%q) = %d 字节（注意：中文字符占多个字节）\n", s, len(s))

	// byte (uint8 别名)：访问字符串的单个字节
	var b byte = s[0] // 'H' 的 ASCII 值
	fmt.Printf("byte: s[0] = %d (%c)\n", b, b)

	// rune (int32 别名)：表示 Unicode 码点
	var r rune = '世' // Unicode 码点 U+4E16
	fmt.Printf("rune: '世' = %d (U+%04X)\n", r, r)

	// 遍历字符串时用 range 获取 rune（正确处理多字节字符）
	fmt.Print("range 遍历 rune: ")
	for i, ch := range "Go世界" {
		fmt.Printf("[%d]%c ", i, ch)
	}
	fmt.Println()
}

// -----------------------------------------------------------------------------
// 7. const 常量
// C 差异：
//   - Go 的 const 支持"无类型常量"（untyped constant），在使用时才确定类型，
//     比 C 的 #define 宏更类型安全，比 C 的 const 更灵活。
//   - Go 没有 #define 预处理器，常量统一用 const。
// -----------------------------------------------------------------------------
func demoConst() {
	fmt.Println("\n--- 7. const 常量 ---")

	// 有类型常量
	const Pi float64 = 3.14159265358979323846
	const MaxSize int = 1024

	// 无类型常量（untyped constant）：没有固定类型，使用时根据上下文推断
	// 这使得 UntypedInt 可以赋给任意整数类型，而不需要显式转换
	const UntypedInt = 42
	const UntypedFloat = 3.14
	const UntypedStr = "hello"

	fmt.Printf("有类型常量: Pi=%.5f, MaxSize=%d\n", Pi, MaxSize)
	fmt.Printf("无类型常量: UntypedInt=%d, UntypedFloat=%.2f, UntypedStr=%q\n",
		UntypedInt, UntypedFloat, UntypedStr)

	// 无类型常量的灵活性：可以赋给不同精度的类型
	var x int32 = UntypedInt  // 无需 int32(42)
	var y int64 = UntypedInt  // 无需 int64(42)
	fmt.Printf("无类型常量赋给不同类型: x(int32)=%d, y(int64)=%d\n", x, y)

	// 批量常量声明
	const (
		StatusOK    = 200
		StatusNotFound = 404
		StatusError = 500
	)
	fmt.Printf("批量常量: OK=%d, NotFound=%d, Error=%d\n", StatusOK, StatusNotFound, StatusError)
}

// -----------------------------------------------------------------------------
// 8. iota 枚举
// C 差异：
//   - C 使用 enum 关键字，Go 使用 const + iota 实现枚举。
//   - iota 在每个 const 块中从 0 开始，每行自动递增 1。
//   - Go 的 iota 比 C enum 更灵活，支持位移、跳值等复杂模式。
// -----------------------------------------------------------------------------

// Weekday 演示基本 iota 枚举（从 0 开始）
type Weekday int

const (
	Sunday Weekday = iota // 0
	Monday                // 1
	Tuesday               // 2
	Wednesday             // 3
	Thursday              // 4
	Friday                // 5
	Saturday              // 6
)

// ByteSize 演示位移枚举（1 << iota）
type ByteSize float64

const (
	_           = iota             // 用 _ 跳过第一个值（iota=0）
	KB ByteSize = 1 << (10 * iota) // 1 << 10 = 1024
	MB                             // 1 << 20
	GB                             // 1 << 30
	TB                             // 1 << 40
)

// LogLevel 演示跳值枚举
type LogLevel int

const (
	LevelDebug LogLevel = iota + 1 // 从 1 开始（跳过 0）
	LevelInfo                      // 2
	LevelWarn                      // 3
	LevelError                     // 4
)

func demoIota() {
	fmt.Println("\n--- 8. iota 枚举 ---")

	fmt.Printf("基本 iota: Sunday=%d, Monday=%d, Saturday=%d\n",
		Sunday, Monday, Saturday)

	fmt.Printf("位移枚举: KB=%.0f, MB=%.0f, GB=%.0f, TB=%.0f\n",
		float64(KB), float64(MB), float64(GB), float64(TB))

	fmt.Printf("跳值枚举（从1开始）: Debug=%d, Info=%d, Warn=%d, Error=%d\n",
		LevelDebug, LevelInfo, LevelWarn, LevelError)

	// iota 在每个独立的 const 块中重置为 0
	const (
		A = iota // 0
		B        // 1
	)
	const (
		C = iota // 0（新 const 块，iota 重置）
		D        // 1
	)
	fmt.Printf("iota 在新 const 块中重置: A=%d, B=%d, C=%d, D=%d\n", A, B, C, D)
}

// -----------------------------------------------------------------------------
// 9. 显式类型转换 T(v)
// C 差异：
//   - Go 没有隐式类型转换，所有类型转换必须显式进行（C 允许隐式数值转换）。
//   - Go 的类型转换语法为 T(v)，类似 C 的强制转换 (T)v，但语义更严格。
//   - 字符串与 []byte/[]rune 的互转会产生数据拷贝（string 是不可变的）。
// -----------------------------------------------------------------------------
func demoTypeConversion() {
	fmt.Println("\n--- 9. 显式类型转换 T(v) ---")

	// 整数与浮点互转
	var i int = 42
	var f float64 = float64(i) // 必须显式转换，C 中可以隐式转换
	var u uint = uint(f)
	fmt.Printf("int->float64->uint: %d -> %f -> %d\n", i, f, u)

	// 精度损失（截断，不是四舍五入）
	pi := 3.99
	truncated := int(pi) // 截断为 3，不是 4
	fmt.Printf("float64->int（截断）: %.2f -> %d\n", pi, truncated)

	// 字符串与 []byte 互转（产生拷贝）
	s := "Hello, Go"
	b := []byte(s)         // string -> []byte
	b[0] = 'h'             // 修改字节切片（不影响原字符串，因为 string 不可变）
	s2 := string(b)        // []byte -> string
	fmt.Printf("string->[]byte->string: %q -> %v -> %q\n", s, b, s2)

	// 字符串与 []rune 互转（正确处理 Unicode）
	unicode := "Go世界"
	runes := []rune(unicode)  // string -> []rune，每个元素是一个 Unicode 码点
	fmt.Printf("string->[]rune: %q -> %v (长度=%d 个码点)\n", unicode, runes, len(runes))
	s3 := string(runes)
	fmt.Printf("[]rune->string: %v -> %q\n", runes, s3)

	// 整数转字符串（注意：int->string 转换的是 Unicode 码点，不是数字字符串）
	r := rune(65)
	charStr := string(r) // 65 是 'A' 的 Unicode 码点
	fmt.Printf("rune->string: %d -> %q（Unicode 码点转字符）\n", r, charStr)
}

// -----------------------------------------------------------------------------
// 10. 字符串字面量：多行反引号 vs 转义双引号
// C 差异：
//   - C 没有原始字符串字面量，多行字符串需要用 \ 续行或字符串拼接。
//   - Go 的反引号字符串（raw string literal）不处理任何转义序列，
//     可以直接包含换行、反斜杠等字符，适合正则表达式、JSON 模板等场景。
// -----------------------------------------------------------------------------
func demoStringLiterals() {
	fmt.Println("\n--- 10. 字符串字面量 ---")

	// 双引号字符串：支持转义序列
	escaped := "第一行\n第二行\t制表符\n反斜杠: \\"
	fmt.Println("双引号（转义字符串）:")
	fmt.Println(escaped)

	// 反引号字符串（raw string literal）：不处理转义，可以跨多行
	raw := `第一行
第二行（无需 \n）
反斜杠直接写: \n 不会被转义
适合正则表达式: \d+\.\d+
适合 JSON 模板: {"key": "value"}`
	fmt.Println("反引号（多行原始字符串）:")
	fmt.Println(raw)

	// 对比：同样内容的两种写法
	a := "line1\nline2"
	b := `line1
line2`
	fmt.Printf("两种写法内容相同: %t\n", a == b)
}

// -----------------------------------------------------------------------------
// 11. nil 指针解引用触发 panic
// C 差异：
//   - C 中解引用 NULL 指针是未定义行为，可能导致段错误（segfault）且难以调试。
//   - Go 中解引用 nil 指针会触发 panic，并附带清晰的错误信息和调用栈，
//     由 Runner 的 recover 机制捕获，不会导致整个程序崩溃。
//   - 这演示了 Go 的 nil 安全机制的重要性：使用指针前应检查是否为 nil。
// -----------------------------------------------------------------------------
func demoNilPanic() {
	fmt.Println("\n--- 11. nil 指针解引用（将触发 panic，由 Runner recover 捕获）---")
	fmt.Println("即将解引用 nil 指针...")

	var p *int // p 是 nil 指针
	fmt.Println(*p) // 解引用 nil 指针，触发 panic: runtime error: invalid memory address or nil pointer dereference
}
