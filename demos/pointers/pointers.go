// Package pointers 演示 Go 语言的指针用法。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// 关键差异总结：
//   - Go 指针不支持指针运算（p++、p+1、p-q 均非法），与 C 的最大区别之一；
//   - Go 由 GC 自动管理内存，无需 malloc/free，但需注意避免意外持有大对象引用；
//   - Go 的 & 和 * 语义与 C 相同，但更安全（nil 解引用触发 panic 而非 segfault）。
package pointers

import "fmt"

// Increment 通过指针将外部变量加 1。
// 演示：通过指针修改函数外部变量（类似 C 的 void increment(int *p)）。
// C 差异：语义完全相同，但 Go 不允许对指针做算术运算（p++ 非法）。
func Increment(p *int) {
	*p++ // 解引用后自增：等价于 (*p)++
}

// newInt 演示函数返回局部变量指针（逃逸分析）。
// C 差异：C 中返回局部变量指针是未定义行为（栈帧销毁后指针悬空）；
// Go 编译器通过逃逸分析自动将该变量分配到堆上，安全且无需手动 malloc。
func newInt(v int) *int {
	x := v   // 局部变量，但因为被取地址返回，编译器将其分配到堆
	return &x // 安全：Go 保证返回后 x 仍然有效
}

// Demo 演示所有指针相关内容。
func Demo() {
	demoBasicPointer()
	demoIncrement()
	demoNewVsAddressOf()
	demoNilCheck()
	demoEscapeAnalysis()
	demoStructPointer()
	demoNoPointerArithmetic()
}

// -----------------------------------------------------------------------------
// 1. 指针基本操作：声明、取地址、解引用
// C 差异：
//   - 语法与 C 基本相同：* 声明指针类型，& 取地址，* 解引用；
//   - Go 指针零值为 nil（C 为 NULL，本质相同）；
//   - Go 不允许指针运算（p++、p+1 均编译错误），C 允许。
// -----------------------------------------------------------------------------
func demoBasicPointer() {
	fmt.Println("\n--- 1. 指针基本操作：声明、取地址、解引用 ---")

	// 声明指针变量（零值为 nil）
	// C 等价：int *p = NULL;
	var p *int
	fmt.Printf("声明指针（零值）: p = %v (nil)\n", p)

	// 取地址：& 运算符
	// C 等价：int x = 42; int *p = &x;
	x := 42
	p = &x
	fmt.Printf("取地址: x = %d, p = %p, &x = %p\n", x, p, &x)
	fmt.Printf("p == &x: %v\n", p == &x)

	// 解引用：* 运算符读取指针指向的值
	// C 等价：int val = *p;
	val := *p
	fmt.Printf("解引用: *p = %d\n", val)

	// 通过指针修改值
	// C 等价：*p = 100;
	*p = 100
	fmt.Printf("通过指针修改: *p = 100, x = %d（x 也变了，因为 p 指向 x）\n", x)
}

// -----------------------------------------------------------------------------
// 2. 通过指针修改函数外部变量（Increment 函数）
// C 差异：
//   - 语义与 C 完全相同，Go 同样通过传递指针实现"传引用"效果；
//   - Go 没有引用类型（&T 参数），只有指针，语义更明确；
//   - C++ 有引用（int &x），Go 没有，只用指针。
// -----------------------------------------------------------------------------
func demoIncrement() {
	fmt.Println("\n--- 2. 通过指针修改函数外部变量（Increment）---")

	n := 10
	fmt.Printf("调用 Increment 前: n = %d\n", n)

	Increment(&n) // 传递 n 的地址
	fmt.Printf("调用 Increment 后: n = %d\n", n)

	Increment(&n)
	Increment(&n)
	fmt.Printf("再调用两次后: n = %d\n", n)
	fmt.Println("注意：Go 通过指针实现传引用效果，C++ 有引用类型，Go 只有指针")
}

// -----------------------------------------------------------------------------
// 3. new(T) 与 var x T; &x 的等价性
// C 差异：
//   - new(T) 类似 C 的 malloc(sizeof(T))，但返回已清零的指针，且由 GC 管理；
//   - Go 的 new 不需要 free，GC 自动回收；
//   - var x T; &x 在函数内部时，若指针逃逸到函数外，编译器自动分配到堆（逃逸分析）；
//   - 两种方式功能等价，实践中更常用 &T{...} 字面量形式。
// -----------------------------------------------------------------------------
func demoNewVsAddressOf() {
	fmt.Println("\n--- 3. new(T) 与 var x T; &x 的等价性 ---")

	// 方式一：new(T) 分配零值并返回指针
	// C 等价：int *p1 = calloc(1, sizeof(int));（已清零）
	p1 := new(int)
	fmt.Printf("new(int): p1 = %p, *p1 = %d（零值）\n", p1, *p1)
	*p1 = 42
	fmt.Printf("赋值后: *p1 = %d\n", *p1)

	// 方式二：var x T; &x（取局部变量地址）
	// 若 &x 逃逸到函数外，编译器自动将 x 分配到堆
	var x int // 零值为 0
	p2 := &x
	fmt.Printf("var x int; &x: p2 = %p, *p2 = %d（零值）\n", p2, *p2)
	*p2 = 42
	fmt.Printf("赋值后: *p2 = %d, x = %d（同一内存）\n", *p2, x)

	// 两种方式等价：都返回指向零值的指针
	fmt.Println("结论：new(T) 与 var x T; &x 功能等价，都返回指向零值的指针")
	fmt.Println("注意：Go 的 new 无需 free，GC 自动管理内存")

	// 结构体的 new 与字面量取地址
	type Point struct{ X, Y int }
	p3 := new(Point)           // 零值 Point
	p4 := &Point{X: 1, Y: 2}  // 字面量取地址（更常用）
	fmt.Printf("new(Point): %v\n", *p3)
	fmt.Printf("&Point{1,2}: %v\n", *p4)
}

// -----------------------------------------------------------------------------
// 4. nil 指针检查：使用前检查是否为 nil
// C 差异：
//   - C 中解引用 NULL 指针是未定义行为（通常导致 segfault）；
//   - Go 中解引用 nil 指针触发 panic（有明确错误信息和调用栈）；
//   - 两者都应在使用前检查，但 Go 的错误更易调试；
//   - 良好实践：函数接收指针参数时，应在开头检查 nil。
// -----------------------------------------------------------------------------
func demoNilCheck() {
	fmt.Println("\n--- 4. nil 指针检查 ---")

	// 安全的 nil 检查模式
	safeDeref := func(p *int) int {
		if p == nil {
			fmt.Println("  指针为 nil，返回默认值 0")
			return 0
		}
		return *p
	}

	var nilPtr *int
	val := 99
	validPtr := &val

	fmt.Printf("safeDeref(nil): %d\n", safeDeref(nilPtr))
	fmt.Printf("safeDeref(&99): %d\n", safeDeref(validPtr))

	// nil 指针比较
	fmt.Printf("nilPtr == nil: %v\n", nilPtr == nil)
	fmt.Printf("validPtr == nil: %v\n", validPtr == nil)

	fmt.Println("良好实践：函数接收指针参数时，应在开头检查 nil，避免 panic")
	fmt.Println("注意：Go 的 nil 解引用触发 panic（有调用栈），C 的 NULL 解引用是 UB（segfault）")
}

// -----------------------------------------------------------------------------
// 5. 函数返回局部变量指针（逃逸分析）
// C 差异：
//   - C 中返回局部变量指针是严重错误（悬空指针，未定义行为）；
//   - Go 编译器通过逃逸分析（escape analysis）自动检测：
//     若局部变量的地址被返回或存储到堆上，编译器自动将其分配到堆；
//   - 程序员无需关心栈/堆分配，GC 负责回收；
//   - 可用 go build -gcflags="-m" 查看逃逸分析结果。
// -----------------------------------------------------------------------------
func demoEscapeAnalysis() {
	fmt.Println("\n--- 5. 函数返回局部变量指针（逃逸分析）---")

	// newInt 返回局部变量的指针，Go 自动将其分配到堆
	p := newInt(42)
	fmt.Printf("newInt(42) 返回的指针: p = %p, *p = %d\n", p, *p)

	// 修改堆上的值
	*p = 100
	fmt.Printf("修改后: *p = %d\n", *p)

	fmt.Println("注意：C 中返回局部变量指针是 UB（悬空指针）")
	fmt.Println("      Go 编译器通过逃逸分析自动将变量分配到堆，安全且无需 malloc")
	fmt.Println("      可用 go build -gcflags=\"-m\" 查看哪些变量发生了逃逸")
}

// -----------------------------------------------------------------------------
// 6. 结构体指针字段访问语法糖：p.Field 等价于 (*p).Field
// C 差异：
//   - C 通过指针访问结构体字段需要 -> 运算符（p->Field）；
//   - Go 统一使用 . 运算符，编译器自动处理解引用（p.Field 等价于 (*p).Field）；
//   - 这使得 Go 代码更简洁，无需区分值和指针的访问语法。
// -----------------------------------------------------------------------------
func demoStructPointer() {
	fmt.Println("\n--- 6. 结构体指针字段访问语法糖 ---")

	type Person struct {
		Name string
		Age  int
	}

	p := &Person{Name: "Alice", Age: 30}

	// Go 语法糖：p.Name 等价于 (*p).Name
	// C 等价：p->Name（Go 没有 -> 运算符）
	fmt.Printf("p.Name（语法糖）: %q\n", p.Name)
	fmt.Printf("(*p).Name（显式解引用）: %q\n", (*p).Name)
	fmt.Printf("两者相等: %v\n", p.Name == (*p).Name)

	// 通过指针修改字段
	p.Age = 31 // 等价于 (*p).Age = 31
	fmt.Printf("p.Age = 31 后: (*p).Age = %d\n", (*p).Age)

	// 方法调用也适用
	fmt.Println("注意：Go 用 . 统一访问值和指针的字段，C 需要区分 . 和 ->")
}

// -----------------------------------------------------------------------------
// 7. Go 指针不支持指针运算（与 C 的关键差异）
// C 差异：
//   - C 允许指针运算：p++（移动到下一个元素）、p+n、p1-p2（指针差）；
//   - Go 完全禁止指针运算，p++、p+1、p-q 均为编译错误；
//   - 这是 Go 内存安全的重要保证之一，防止越界访问；
//   - 如需底层指针操作，可使用 unsafe.Pointer（不推荐，破坏类型安全）；
//   - Go 通过切片（slice）提供安全的数组遍历，无需指针运算。
//
// GC 内存管理注意事项：
//   - Go 由 GC 自动管理内存，无需 malloc/free；
//   - 但需注意避免意外持有大对象引用（导致 GC 无法回收，类似内存泄漏）；
//   - 常见场景：全局变量、长生命周期的 goroutine、缓存中持有大切片的子切片；
//   - 解决方案：及时将不再需要的引用置为 nil，或使用 copy 而非子切片。
// -----------------------------------------------------------------------------
func demoNoPointerArithmetic() {
	fmt.Println("\n--- 7. Go 指针不支持指针运算（与 C 的关键差异）---")

	arr := [5]int{10, 20, 30, 40, 50}
	p := &arr[0]
	fmt.Printf("&arr[0] = %p, *p = %d\n", p, *p)

	// Go 中遍历数组/切片应使用 range 或索引，而非指针运算
	// C 中可以写：for (int *q = arr; q < arr+5; q++) { ... }
	// Go 中必须写：
	fmt.Print("Go 安全遍历（range）: ")
	for _, v := range arr {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	fmt.Print("Go 安全遍历（索引）: ")
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%d ", arr[i])
	}
	fmt.Println()

	// 以下代码在 Go 中是编译错误（演示注释）：
	// p++        // 编译错误：invalid operation: p++ (non-numeric type *int)
	// p = p + 1  // 编译错误：invalid operation: p + 1 (mismatched types *int and untyped int)
	fmt.Println("注意：p++、p+1 在 Go 中均为编译错误（C 允许，Go 禁止）")
	fmt.Println("注意：Go 通过切片提供安全的数组遍历，无需指针运算")

	fmt.Println("\n--- GC 内存管理注意事项 ---")
	fmt.Println("✓ Go 由 GC 自动管理内存，无需 malloc/free")
	fmt.Println("✓ 无需担心悬空指针（dangling pointer）")
	fmt.Println("⚠ 注意：避免意外持有大对象引用（导致 GC 无法回收）")
	fmt.Println("  常见场景：全局变量持有大切片、goroutine 泄漏持有大对象")
	fmt.Println("  解决方案：不再使用时将引用置为 nil，或使用 copy 而非子切片")
}
