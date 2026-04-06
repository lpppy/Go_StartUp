// Package functions 演示 Go 语言的函数特性。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
package functions

import (
	"errors"
	"fmt"
)

// Demo 演示所有函数相关内容。
func Demo() {
	demoBasicFunction()
	demoMultipleReturn()
	demoNamedReturn()
	demoVariadic()
	demoFirstClass()
	demoAnonymous()
	demoClosure()
	demoDeferLIFO()
	demoDeferNamedReturn()
	demoDeferUseCases()
}

// -----------------------------------------------------------------------------
// 导出函数（供测试调用）
// -----------------------------------------------------------------------------

// Sum 可变参数求和。
// C 差异：C 的可变参数需要 <stdarg.h> 中的 va_list/va_arg 宏，类型不安全；
// Go 的可变参数 ...T 是类型安全的，在函数内部表现为 []T 切片。
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// MakeAdder 返回一个闭包，该闭包将其参数与 x 相加。
// C 差异：C 没有闭包，函数指针无法捕获外部变量；
// Go 的闭包可以捕获并持有外部作用域的变量，形成独立的执行环境。
func MakeAdder(x int) func(int) int {
	return func(y int) int {
		return x + y // 捕获外部变量 x
	}
}

// SafeDivide 使用 defer+recover 捕获除零 panic，将其转换为 error 返回。
// C 差异：C 的整数除零是未定义行为（UB），通常导致程序崩溃（SIGFPE）；
// Go 的整数除零会触发 panic，可以用 defer+recover 安全捕获并转换为 error。
func SafeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("除法错误: %v", r)
		}
	}()
	result = a / b // 若 b==0 触发 panic，由上方 defer 捕获
	return
}

// -----------------------------------------------------------------------------
// 1. 基本函数定义与调用
// C 差异：
//   - Go 的参数类型后置：func add(a int, b int) int，而 C 是 int add(int a, int b)。
//   - 相邻参数类型相同时可以合并：func add(a, b int) int。
//   - Go 不支持默认参数和函数重载（C 也不支持重载，但 C++ 支持）。
//   - Go 函数是包级别的，不需要前向声明（C 需要在使用前声明或提供头文件）。
// -----------------------------------------------------------------------------
func demoBasicFunction() {
	fmt.Println("\n--- 1. 基本函数定义与调用 ---")

	// 参数类型后置语法：func 函数名(参数名 类型) 返回类型
	// C 等价：int add(int a, int b) { return a + b; }
	add := func(a, b int) int { // 相邻同类型参数可合并
		return a + b
	}

	result := add(3, 4)
	fmt.Printf("add(3, 4) = %d\n", result)
	fmt.Println("注意：Go 参数类型后置（与 C 相反），相邻同类型参数可合并写法：a, b int")
	fmt.Println("注意：Go 不支持默认参数和函数重载，不需要前向声明")
}

// multiply 演示包级别函数，无需前向声明即可在 Demo 中调用。
// C 差异：C 需要在调用前声明函数原型（或在调用前定义函数）。
func multiply(a, b int) int {
	return a * b
}

// -----------------------------------------------------------------------------
// 2. 多返回值函数
// C 差异：
//   - C 函数只能返回一个值，多值返回需要通过指针参数输出或结构体；
//   - Go 原生支持多返回值，惯用法是最后一个返回值为 error。
//   - 调用方必须处理所有返回值（或用 _ 显式忽略），不能像 C 那样直接丢弃。
// -----------------------------------------------------------------------------
func demoMultipleReturn() {
	fmt.Println("\n--- 2. 多返回值函数 ---")

	// divide 同时返回商和 error
	// C 等价需要：int divide(int a, int b, int *result)，通过指针输出结果
	divide := func(a, b int) (int, error) {
		if b == 0 {
			return 0, errors.New("除数不能为零")
		}
		return a / b, nil
	}

	// 调用方同时接收两个返回值
	if result, err := divide(10, 3); err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("divide(10, 3) = %d, err = nil\n", result)
	}

	// 用 _ 忽略不需要的返回值
	_, err := divide(5, 0)
	fmt.Printf("divide(5, 0): err = %v\n", err)

	// 演示包级别函数调用（无需前向声明）
	fmt.Printf("multiply(6, 7) = %d（包级别函数，无需前向声明）\n", multiply(6, 7))
}

// -----------------------------------------------------------------------------
// 3. 命名返回值（named return values）与裸返回（naked return）
// C 差异：
//   - C 没有命名返回值的概念；
//   - Go 可以在函数签名中为返回值命名，命名返回值会被初始化为零值；
//   - 裸返回（naked return）直接写 return 不带值，返回当前命名返回值；
//   - 注意事项：裸返回在长函数中会降低可读性，建议仅在短函数中使用。
// -----------------------------------------------------------------------------
func demoNamedReturn() {
	fmt.Println("\n--- 3. 命名返回值与裸返回 ---")

	// minMax 使用命名返回值，函数签名中声明 min, max
	// 命名返回值会被自动初始化为零值（int 为 0）
	minMax := func(nums []int) (min, max int) {
		if len(nums) == 0 {
			return // 裸返回：返回零值 min=0, max=0
		}
		min, max = nums[0], nums[0]
		for _, n := range nums[1:] {
			if n < min {
				min = n
			}
			if n > max {
				max = n
			}
		}
		return // 裸返回：返回当前 min 和 max 的值
	}

	nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
	min, max := minMax(nums)
	fmt.Printf("minMax(%v) = min:%d, max:%d\n", nums, min, max)
	fmt.Println("注意：命名返回值在函数签名中声明，自动初始化为零值")
	fmt.Println("注意：裸返回（naked return）在长函数中会降低可读性，建议仅在短函数中使用")
}

// -----------------------------------------------------------------------------
// 4. 可变参数函数（variadic function）
// C 差异：
//   - C 的可变参数（<stdarg.h>）是类型不安全的，需要手动解析参数；
//   - Go 的可变参数 ...T 是类型安全的，在函数内部表现为 []T 切片；
//   - 可以用 slice... 语法将切片展开传入可变参数函数（类似 C 的 apply）。
// -----------------------------------------------------------------------------
func demoVariadic() {
	fmt.Println("\n--- 4. 可变参数函数 ---")

	// 直接传入多个参数
	fmt.Printf("Sum(1, 2, 3) = %d\n", Sum(1, 2, 3))
	fmt.Printf("Sum() = %d（空参数，返回零值）\n", Sum())

	// 用 slice... 展开切片传入
	// C 差异：C 没有这种语法，需要手动传递数组指针和长度
	nums := []int{10, 20, 30, 40, 50}
	fmt.Printf("Sum(nums...) = %d（slice... 展开切片）\n", Sum(nums...))
	fmt.Println("注意：nums... 将切片展开为独立参数，等价于 Sum(10, 20, 30, 40, 50)")
}

// -----------------------------------------------------------------------------
// 5. 函数作为一等公民（first-class citizen）
// C 差异：
//   - C 支持函数指针，但语法繁琐：int (*fp)(int, int) = &add；
//   - Go 的函数是一等公民，可以赋值给变量、定义函数类型、作为参数和返回值，语法更简洁；
//   - Go 的函数类型是值类型，可以直接比较是否为 nil（但不能比较两个函数是否相等）。
// -----------------------------------------------------------------------------
func demoFirstClass() {
	fmt.Println("\n--- 5. 函数作为一等公民 ---")

	// 5a. 将函数赋值给变量
	// C 等价：int (*fp)(int, int) = multiply;（函数指针，语法繁琐）
	op := multiply // 函数赋值给变量，类型为 func(int, int) int
	fmt.Printf("函数赋值给变量: op(3, 4) = %d\n", op(3, 4))

	// 5b. 定义函数类型
	// C 等价：typedef int (*BinaryOp)(int, int);
	type BinaryOp func(int, int) int
	var myOp BinaryOp = multiply
	fmt.Printf("定义函数类型 BinaryOp: myOp(5, 6) = %d\n", myOp(5, 6))

	// 5c. 函数作为参数传递（高阶函数）
	// C 等价：int apply(int (*f)(int, int), int a, int b) { return f(a, b); }
	apply := func(f func(int, int) int, a, b int) int {
		return f(a, b)
	}
	fmt.Printf("函数作为参数: apply(multiply, 7, 8) = %d\n", apply(multiply, 7, 8))

	// 5d. 函数作为返回值（工厂函数）
	// C 差异：C 的函数指针可以作为返回值，但无法捕获上下文（无闭包）
	makeMultiplier := func(factor int) func(int) int {
		return func(x int) int {
			return x * factor
		}
	}
	double := makeMultiplier(2)
	triple := makeMultiplier(3)
	fmt.Printf("函数作为返回值: double(5)=%d, triple(5)=%d\n", double(5), triple(5))
}

// -----------------------------------------------------------------------------
// 6. 匿名函数（anonymous function）
// C 差异：
//   - C 没有匿名函数（lambda），只有具名函数指针；
//   - Go 支持匿名函数，可以立即调用（IIFE）或赋值给变量后调用；
//   - 匿名函数可以捕获外部变量（形成闭包），C 的函数指针做不到这一点。
// -----------------------------------------------------------------------------
func demoAnonymous() {
	fmt.Println("\n--- 6. 匿名函数 ---")

	// 6a. 立即调用的匿名函数（IIFE: Immediately Invoked Function Expression）
	// C 差异：C 没有 IIFE 语法
	result := func(a, b int) int {
		return a + b
	}(10, 20) // 定义后立即调用，传入参数 10, 20
	fmt.Printf("IIFE（立即调用匿名函数）: (func(10, 20) int { return a+b })() = %d\n", result)

	// 6b. 匿名函数赋值给变量后调用
	greet := func(name string) string {
		return fmt.Sprintf("你好，%s！", name)
	}
	fmt.Printf("匿名函数赋值给变量: greet(\"Go\") = %s\n", greet("Go"))

	// 6c. 匿名函数作为 goroutine（并发场景常见用法）
	// 注意：此处仅演示语法，实际并发需要同步机制
	fmt.Println("匿名函数常用于 goroutine：go func() { ... }()")
}

// -----------------------------------------------------------------------------
// 7. 闭包（closure）
// C 差异：
//   - C 没有闭包，函数指针无法捕获外部变量；
//   - Go 的闭包捕获的是变量的引用（而非值的副本），多个闭包共享同一变量；
//   - 循环变量捕获陷阱：在循环中创建闭包时，所有闭包共享同一个循环变量，
//     需要通过参数传递或局部变量副本来修复。
// -----------------------------------------------------------------------------
func demoClosure() {
	fmt.Println("\n--- 7. 闭包捕获外部变量 ---")

	// 7a. 闭包捕获外部变量（引用捕获）
	counter := 0
	increment := func() int {
		counter++ // 捕获并修改外部变量 counter
		return counter
	}
	fmt.Printf("闭包捕获外部变量: increment()=%d, increment()=%d, increment()=%d\n",
		increment(), increment(), increment())
	fmt.Printf("外部变量 counter 被修改: counter=%d\n", counter)

	// 7b. MakeAdder 演示闭包作为工厂函数
	add5 := MakeAdder(5)
	add10 := MakeAdder(10)
	fmt.Printf("MakeAdder(5)(3)=%d, MakeAdder(10)(3)=%d（每个闭包独立捕获 x）\n",
		add5(3), add10(3))

	// 7c. 循环变量捕获陷阱（常见 Bug）
	// C 差异：C 没有闭包，不存在此问题
	fmt.Println("\n循环变量捕获陷阱演示:")

	// ❌ 错误写法：所有闭包共享同一个循环变量 i
	// 在 Go 1.22 之前，循环变量在每次迭代中是同一个变量
	funcs := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		i := i // ✅ 修复：用同名局部变量遮蔽循环变量，每次迭代创建新变量
		funcs[i] = func() int {
			return i // 捕获的是局部变量 i（每次迭代独立）
		}
	}
	fmt.Print("修复后（局部变量遮蔽）: ")
	for _, f := range funcs {
		fmt.Printf("%d ", f())
	}
	fmt.Println()

	// 另一种修复方式：通过参数传递
	funcs2 := make([]func() int, 3)
	for i := 0; i < 3; i++ {
		func(val int) { // 通过参数传递，val 是 i 的副本
			funcs2[val] = func() int { return val }
		}(i)
	}
	fmt.Print("修复后（参数传递）: ")
	for _, f := range funcs2 {
		fmt.Printf("%d ", f())
	}
	fmt.Println()
	fmt.Println("注意：Go 1.22+ 修改了循环变量语义，每次迭代创建新变量，陷阱不再存在")
}

// -----------------------------------------------------------------------------
// 8. defer LIFO 执行顺序
// C 差异：
//   - C 没有 defer 机制，资源释放需要手动在每个返回路径前调用；
//   - Go 的 defer 语句将函数调用推入栈，在当前函数返回时按 LIFO（后进先出）顺序执行；
//   - defer 的参数在 defer 语句执行时立即求值（而非在延迟调用时求值）。
// -----------------------------------------------------------------------------
func demoDeferLIFO() {
	fmt.Println("\n--- 8. defer LIFO 执行顺序 ---")

	fmt.Println("注册 defer 顺序：1, 2, 3")
	defer fmt.Println("defer 执行顺序：3（最后注册，最先执行）")
	defer fmt.Println("defer 执行顺序：2")
	defer fmt.Println("defer 执行顺序：1（最先注册，最后执行）")
	fmt.Println("函数体执行完毕，即将执行 defer 栈（LIFO）...")

	// defer 参数立即求值演示
	x := 10
	defer fmt.Printf("defer 参数立即求值: x=%d（defer 注册时 x=10，非返回时的值）\n", x)
	x = 20
	fmt.Printf("x 已修改为: %d\n", x)
}

// -----------------------------------------------------------------------------
// 9. defer 修改命名返回值
// C 差异：
//   - C 没有此机制；
//   - Go 的 defer 可以通过命名返回值在函数返回后修改返回值，
//     这是实现 panic 恢复并返回 error 的关键机制。
// -----------------------------------------------------------------------------
func demoDeferNamedReturn() {
	fmt.Println("\n--- 9. defer 修改命名返回值 ---")

	// double 演示 defer 修改命名返回值
	double := func(n int) (result int) { // result 是命名返回值
		defer func() {
			result *= 2 // defer 在 return 之后执行，可以修改命名返回值
		}()
		result = n
		return // 裸返回：result=n，但 defer 会将其乘以 2
	}

	fmt.Printf("double(5) = %d（defer 将 result 从 5 修改为 10）\n", double(5))
	fmt.Println("注意：defer 在 return 语句执行后、函数真正返回前执行，可修改命名返回值")

	// SafeDivide 演示 defer+recover 修改命名返回值（err）
	result, err := SafeDivide(10, 2)
	fmt.Printf("SafeDivide(10, 2) = %d, err = %v\n", result, err)

	result, err = SafeDivide(10, 0)
	fmt.Printf("SafeDivide(10, 0) = %d, err = %v\n", result, err)
}

// -----------------------------------------------------------------------------
// 10. defer 典型用途：资源释放与 panic 恢复
// C 差异：
//   - C 需要在每个返回路径前手动释放资源，容易遗漏（尤其是错误路径）；
//   - Go 的 defer 确保资源释放代码紧跟资源获取代码，无论函数如何返回都会执行；
//   - defer+recover 是 Go 将 panic 转换为 error 的标准模式，
//     类似 C++ 的 try-catch，但 Go 不鼓励用 panic/recover 做常规错误处理。
// -----------------------------------------------------------------------------
func demoDeferUseCases() {
	fmt.Println("\n--- 10. defer 典型用途 ---")

	// 10a. 资源释放（模拟文件关闭）
	// C 等价：FILE *f = fopen(...); ... fclose(f);（需要在每个 return 前调用）
	openFile := func(name string) string {
		fmt.Printf("  打开文件: %s\n", name)
		return name
	}
	closeFile := func(name string) {
		fmt.Printf("  关闭文件: %s（由 defer 保证，无论函数如何返回）\n", name)
	}

	processFile := func(name string) error {
		f := openFile(name)
		defer closeFile(f) // 紧跟资源获取，确保释放
		// ... 处理文件 ...
		fmt.Printf("  处理文件: %s\n", name)
		return nil
	}

	fmt.Println("资源释放演示（模拟文件操作）:")
	_ = processFile("data.txt")

	// 10b. panic 恢复（SafeDivide 已演示，此处展示通用模式）
	fmt.Println("\npanic 恢复演示（SafeDivide）:")
	result, err := SafeDivide(10, 0)
	fmt.Printf("  SafeDivide(10, 0): result=%d, err=%v\n", result, err)
	fmt.Println("  注意：defer+recover 将 panic 转换为 error，是 Go 的标准模式")
	fmt.Println("  注意：不要用 panic/recover 做常规错误处理，仅用于真正的异常情况")

	// 10c. 互斥锁释放（演示模式，不实际使用 sync.Mutex）
	fmt.Println("\n互斥锁释放模式（伪代码演示）:")
	fmt.Println("  mu.Lock()")
	fmt.Println("  defer mu.Unlock()  // 紧跟 Lock，确保解锁")
	fmt.Println("  // ... 访问共享资源 ...")
}
