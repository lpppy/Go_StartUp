// Package interfaces 演示 Go 语言的接口机制。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// Go 的接口是隐式实现的（鸭子类型）：只要一个类型实现了接口的所有方法，
// 它就自动满足该接口，无需任何 implements 关键字或显式声明。
// 这与 C 的函数指针结构体（手动模拟多态）有本质区别。
package interfaces

import (
	"fmt"
	"io"
	"math"
	"strings"
)

// -----------------------------------------------------------------------------
// 类型定义
// C 差异：
//   - C 模拟多态需要定义含函数指针的结构体（如 struct Shape { double (*area)(void*); }）；
//   - Go 的接口是纯方法集合，任何实现了这些方法的类型都自动满足接口；
//   - 无需 implements 关键字，无需注册，编译器自动检查——这就是"鸭子类型"；
//   - 接口与实现完全解耦：接口定义者和实现者可以在不同包中，互不依赖。
// -----------------------------------------------------------------------------

// Shape 接口定义了几何形状的通用行为。
// C 差异：C 需要手动定义函数指针结构体来模拟接口，Go 直接用 interface 关键字。
type Shape interface {
	Area() float64      // 计算面积
	Perimeter() float64 // 计算周长
}

// Circle 表示圆形，隐式实现 Shape 接口。
// C 差异：C 需要手动将函数指针赋值给结构体字段，Go 只需实现方法即可。
type Circle struct {
	Radius float64 // 半径
}

// Area 返回圆的面积（π * r²）。
// 值接收者方法：Circle 自动满足 Shape 接口（无需任何声明）。
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter 返回圆的周长（2 * π * r）。
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// String 实现 fmt.Stringer 接口，自定义 fmt.Println 输出格式。
// C 差异：C 需要手动调用格式化函数，Go 通过接口自动调用。
func (c Circle) String() string {
	return fmt.Sprintf("Circle{Radius: %.2f}", c.Radius)
}

// Rectangle 表示矩形，隐式实现 Shape 接口。
type Rectangle struct {
	Width  float64 // 宽
	Height float64 // 高
}

// Area 返回矩形的面积（宽 * 高）。
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter 返回矩形的周长（2 * (宽 + 高)）。
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// String 实现 fmt.Stringer 接口。
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{Width: %.2f, Height: %.2f}", r.Width, r.Height)
}

// -----------------------------------------------------------------------------
// 接口组合相关类型
// C 差异：C 没有接口组合的概念，需要手动合并函数指针结构体。
// -----------------------------------------------------------------------------

// Reader 演示接口组合用的简单读接口。
type Reader interface {
	Read() string
}

// Writer 演示接口组合用的简单写接口。
type Writer interface {
	Write(s string)
}

// ReadWriter 通过嵌入 Reader 和 Writer 组合成新接口。
// C 差异：C 需要手动将两个函数指针结构体的字段合并到一个新结构体中。
// Go 直接用接口嵌入，ReadWriter 自动包含 Read() 和 Write() 两个方法。
type ReadWriter interface {
	Reader
	Writer
}

// Buffer 实现了 ReadWriter 接口（同时实现了 Reader 和 Writer）。
type Buffer struct {
	data strings.Builder
}

func (b *Buffer) Read() string {
	return b.data.String()
}

func (b *Buffer) Write(s string) {
	b.data.WriteString(s)
}

// -----------------------------------------------------------------------------
// 自定义错误类型（实现 error 接口）
// C 差异：C 用 errno 全局变量或返回负数表示错误，Go 用 error 接口。
// error 接口只有一个方法：Error() string
// -----------------------------------------------------------------------------

// ShapeError 是自定义错误类型，实现了 error 接口。
type ShapeError struct {
	Op  string  // 操作名称
	Val float64 // 导致错误的值
}

// Error 实现 error 接口（只需实现 Error() string 方法）。
func (e *ShapeError) Error() string {
	return fmt.Sprintf("shape error: %s with value %.2f", e.Op, e.Val)
}

// NewCircle 创建 Circle，若半径为负则返回错误。
// 演示 error 接口的实际用法。
func NewCircle(radius float64) (Circle, error) {
	if radius < 0 {
		return Circle{}, &ShapeError{Op: "NewCircle", Val: radius}
	}
	return Circle{Radius: radius}, nil
}

// -----------------------------------------------------------------------------
// 自定义 io.Reader 实现
// C 差异：C 的 FILE* 是具体类型，Go 的 io.Reader 是接口，任何实现了 Read 方法的类型都可以用。
// -----------------------------------------------------------------------------

// StringReader 实现 io.Reader 接口，从字符串中读取数据。
type StringReader struct {
	data string
	pos  int
}

// Read 实现 io.Reader 接口（签名必须与接口完全一致）。
func (sr *StringReader) Read(p []byte) (n int, err error) {
	if sr.pos >= len(sr.data) {
		return 0, io.EOF
	}
	n = copy(p, sr.data[sr.pos:])
	sr.pos += n
	return n, nil
}

// Demo 演示所有接口相关内容。
func Demo() {
	demoInterfaceDefinition()
	demoImplicitImplementation()
	demoPolymorphism()
	demoInterfaceComposition()
	demoEmptyInterface()
	demoTypeAssertion()
	demoTypeSwitch()
	demoNilInterface()
	demoStandardInterfaces()
}

// -----------------------------------------------------------------------------
// 1. 接口定义
// C 差异：
//   - C 模拟接口需要定义含函数指针的结构体，并手动初始化每个函数指针；
//   - Go 的 interface 关键字直接定义方法集合，语法简洁；
//   - Go 接口只包含方法签名，不包含数据字段（与 C 的函数指针结构体不同）；
//   - 接口是引用类型，零值为 nil（表示"没有具体类型和值"）。
// -----------------------------------------------------------------------------
func demoInterfaceDefinition() {
	fmt.Println("\n--- 1. 接口定义 ---")

	// Shape 接口定义：
	//   type Shape interface {
	//       Area() float64
	//       Perimeter() float64
	//   }
	//
	// C 等价（函数指针结构体模拟）：
	//   typedef struct {
	//       double (*area)(void *self);
	//       double (*perimeter)(void *self);
	//   } ShapeVTable;
	//   typedef struct { ShapeVTable *vtable; void *data; } Shape;

	fmt.Println("Shape 接口定义了两个方法：Area() float64 和 Perimeter() float64")
	fmt.Println("任何实现了这两个方法的类型都自动满足 Shape 接口（无需声明）")

	// 接口变量的零值是 nil
	var s Shape
	fmt.Printf("接口变量零值: s == nil → %v\n", s == nil)
	fmt.Println("注意：nil 接口表示没有绑定任何具体类型，调用方法会 panic")
	fmt.Println("C 差异：Go 接口是隐式实现（鸭子类型），C 需要手动绑定函数指针")
}

// -----------------------------------------------------------------------------
// 2. 结构体隐式实现接口
// C 差异：
//   - C++ 需要 class Circle : public Shape（显式继承）；
//   - Go 只需实现接口要求的所有方法，编译器自动验证——无需任何声明；
//   - 这种设计使得接口定义者和实现者完全解耦，可以在不同包中独立演化；
//   - 如果漏实现某个方法，编译时会报错（cannot use Circle as type Shape）。
// -----------------------------------------------------------------------------
func demoImplicitImplementation() {
	fmt.Println("\n--- 2. 结构体隐式实现接口（无需 implements 关键字）---")

	c := Circle{Radius: 5.0}
	r := Rectangle{Width: 4.0, Height: 3.0}

	// Circle 和 Rectangle 都实现了 Shape 接口（隐式，无需声明）
	// C++ 等价：class Circle : public Shape { ... }（显式继承）
	// Go：只需实现 Area() 和 Perimeter() 方法，自动满足 Shape 接口

	fmt.Printf("Circle{Radius:5}   → Area=%.4f, Perimeter=%.4f\n",
		c.Area(), c.Perimeter())
	fmt.Printf("Rectangle{4×3}     → Area=%.4f, Perimeter=%.4f\n",
		r.Area(), r.Perimeter())

	// 将具体类型赋值给接口变量（隐式转换）
	var s Shape = c // Circle 自动满足 Shape，无需任何转换声明
	fmt.Printf("var s Shape = Circle → s.Area()=%.4f\n", s.Area())

	s = r // Rectangle 也自动满足 Shape
	fmt.Printf("var s Shape = Rectangle → s.Area()=%.4f\n", s.Area())

	fmt.Println("C 差异：Go 无需 implements 关键字，编译器自动检查方法集是否满足接口")
}

// -----------------------------------------------------------------------------
// 3. 接口变量的多态用法
// C 差异：
//   - C 通过函数指针结构体 + void* 实现多态，需要手动管理类型信息；
//   - Go 的接口变量内部包含（类型, 值）两个字段，运行时自动分发方法调用；
//   - 同一接口变量可以持有不同具体类型，调用相同方法产生不同行为；
//   - 这是面向对象编程中"多态"的核心思想，Go 通过接口优雅实现。
// -----------------------------------------------------------------------------
func demoPolymorphism() {
	fmt.Println("\n--- 3. 接口变量的多态用法 ---")

	// 创建不同形状的切片（统一用 Shape 接口存储）
	// C 差异：C 需要 Shape* shapes[]（指针数组）+ 手动管理类型标签
	shapes := []Shape{
		Circle{Radius: 3.0},
		Rectangle{Width: 5.0, Height: 2.0},
		Circle{Radius: 1.5},
		Rectangle{Width: 10.0, Height: 4.0},
	}

	// 多态调用：同一接口方法，不同类型产生不同行为
	fmt.Println("遍历所有形状（多态调用）：")
	totalArea := 0.0
	for _, s := range shapes {
		area := s.Area()
		perimeter := s.Perimeter()
		totalArea += area
		fmt.Printf("  %T → Area=%.4f, Perimeter=%.4f\n", s, area, perimeter)
	}
	fmt.Printf("所有形状总面积: %.4f\n", totalArea)

	// 函数接受接口参数（多态函数）
	// C 差异：C 需要 void printShape(Shape *s)（传指针 + 函数指针调用）
	printShapeInfo := func(s Shape) {
		fmt.Printf("  形状信息 → 面积=%.4f, 周长=%.4f\n", s.Area(), s.Perimeter())
	}

	fmt.Println("通过接口参数调用（多态函数）：")
	printShapeInfo(Circle{Radius: 2.0})
	printShapeInfo(Rectangle{Width: 3.0, Height: 4.0})
	fmt.Println("C 差异：Go 接口自动分发方法，C 需要手动通过函数指针调用")
}

// -----------------------------------------------------------------------------
// 4. 接口组合
// C 差异：
//   - C 没有接口组合的概念，需要手动将多个函数指针结构体的字段合并；
//   - Go 通过在接口中嵌入其他接口，自动组合方法集；
//   - 标准库大量使用接口组合：io.ReadWriter = io.Reader + io.Writer；
//   - 接口组合遵循最小接口原则（interface segregation），每个接口只做一件事。
// -----------------------------------------------------------------------------
func demoInterfaceComposition() {
	fmt.Println("\n--- 4. 接口组合（嵌入多个接口定义新接口）---")

	// ReadWriter 接口 = Reader + Writer（接口组合）
	// 标准库等价：io.ReadWriter = io.Reader + io.Writer
	//
	// C 差异：C 需要手动合并：
	//   typedef struct { char* (*read)(void*); void (*write)(void*, char*); } ReadWriter;

	var rw ReadWriter = &Buffer{}
	rw.Write("Hello, ")
	rw.Write("Go Interfaces!")
	fmt.Printf("ReadWriter.Read() = %q\n", rw.Read())

	// 接口组合的好处：可以单独使用子接口
	var w Writer = &Buffer{}
	w.Write("只写不读")
	fmt.Println("接口组合允许单独使用子接口（Reader 或 Writer）")

	// 标准库接口组合示例
	fmt.Println("\n标准库接口组合示例：")
	fmt.Println("  io.Reader    = Read(p []byte) (n int, err error)")
	fmt.Println("  io.Writer    = Write(p []byte) (n int, err error)")
	fmt.Println("  io.ReadWriter = io.Reader + io.Writer（组合）")
	fmt.Println("  io.ReadWriteCloser = io.ReadWriter + io.Closer（再次组合）")
	fmt.Println("C 差异：Go 接口组合通过嵌入实现，C 需要手动合并函数指针结构体")
}

// -----------------------------------------------------------------------------
// 5. 空接口 interface{} / any（Go 1.18+）
// C 差异：
//   - C 用 void* 作为通用指针，但丢失了类型信息，需要手动管理类型标签；
//   - Go 的空接口 interface{} 可以持有任何类型的值，且保留了类型信息；
//   - Go 1.18 引入 any 作为 interface{} 的别名，语义完全相同；
//   - 空接口是 Go 泛型出现前的通用容器，现在推荐用泛型替代（类型安全）。
// -----------------------------------------------------------------------------
func demoEmptyInterface() {
	fmt.Println("\n--- 5. 空接口 interface{} / any（通用容器）---")

	// 空接口可以持有任何类型的值
	// C 差异：C 用 void* 实现，但丢失类型信息；Go 的 interface{} 保留类型信息
	var anything interface{}
	fmt.Printf("interface{} 零值: %v (nil=%v)\n", anything, anything == nil)

	anything = 42
	fmt.Printf("持有 int:    %v (类型=%T)\n", anything, anything)

	anything = "hello"
	fmt.Printf("持有 string: %v (类型=%T)\n", anything, anything)

	anything = Circle{Radius: 3.0}
	fmt.Printf("持有 Circle: %v (类型=%T)\n", anything, anything)

	anything = []int{1, 2, 3}
	fmt.Printf("持有 []int:  %v (类型=%T)\n", anything, anything)

	// any 是 interface{} 的别名（Go 1.18+）
	var val any = 3.14
	fmt.Printf("\nany（interface{} 别名）: %v (类型=%T)\n", val, val)

	// 空接口切片作为通用容器
	// C 差异：C 需要 void* 数组 + 类型标签数组，Go 直接用 []interface{}
	mixed := []interface{}{1, "two", 3.0, true, Circle{Radius: 1}}
	fmt.Println("\n[]interface{} 混合类型切片：")
	for i, v := range mixed {
		fmt.Printf("  [%d] %v (类型=%T)\n", i, v, v)
	}

	// fmt.Println 接受 ...interface{}，这就是它能打印任何类型的原因
	fmt.Println("注意：fmt.Println 的参数类型就是 ...interface{}（接受任意类型）")
	fmt.Println("C 差异：Go interface{} 保留类型信息，C 的 void* 丢失类型信息")
}

// -----------------------------------------------------------------------------
// 6. 类型断言（type assertion）
// C 差异：
//   - C 的 void* 强制转换没有运行时类型检查，转换错误会导致未定义行为；
//   - Go 的类型断言在运行时检查接口的动态类型，安全形式不会 panic；
//   - 安全形式：v, ok := i.(T)，ok 为 false 时 v 为零值，不 panic；
//   - 不安全形式：v := i.(T)，类型不匹配时直接 panic（类似 C 的错误转换）。
// -----------------------------------------------------------------------------
func demoTypeAssertion() {
	fmt.Println("\n--- 6. 类型断言（type assertion）---")

	var s Shape = Circle{Radius: 5.0}

	// 安全形式：v, ok := i.(ConcreteType)
	// C 差异：C 没有安全的类型转换，Go 的 ok 形式避免 panic
	if c, ok := s.(Circle); ok {
		fmt.Printf("安全断言成功: Circle.Radius=%.2f\n", c.Radius)
	}

	// 断言为错误类型：ok=false，v 为零值，不 panic
	if r, ok := s.(Rectangle); ok {
		fmt.Printf("断言为 Rectangle: %v\n", r)
	} else {
		fmt.Printf("安全断言失败: s 不是 Rectangle（ok=false，r=%v，不 panic）\n", r)
	}

	// 空接口的类型断言
	var i interface{} = "hello, world"
	str, ok := i.(string)
	fmt.Printf("interface{} 断言为 string: %q, ok=%v\n", str, ok)

	num, ok := i.(int)
	fmt.Printf("interface{} 断言为 int:    %d, ok=%v（失败时返回零值）\n", num, ok)

	// 不安全形式：类型不匹配时 panic（由 main.go 的 recover 捕获）
	// C 差异：类似 C 的错误类型转换，但 Go 会 panic 而不是静默产生未定义行为
	fmt.Println("\n演示不安全断言 panic（将被 recover 捕获）：")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到 panic: %v\n", r)
				fmt.Println("  注意：不安全断言 i.(T) 在类型不匹配时 panic，应使用 v, ok := i.(T)")
			}
		}()
		var x interface{} = "not an int"
		n := x.(int) // 类型不匹配，panic: interface conversion
		fmt.Println(n) // 不会执行到这里
	}()

	fmt.Println("C 差异：Go 类型断言有运行时检查，C 的强制转换无检查（UB）")
}

// -----------------------------------------------------------------------------
// 7. 类型选择（type switch）
// C 差异：
//   - C 需要手动维护类型标签（enum）并用 switch 判断，然后强制转换；
//   - Go 的 type switch 语法简洁，编译器自动处理类型断言和变量绑定；
//   - switch v := i.(type) 中，每个 case 分支的 v 自动具有对应的具体类型；
//   - 这是处理多种接口实现类型的惯用模式。
// -----------------------------------------------------------------------------
func demoTypeSwitch() {
	fmt.Println("\n--- 7. 类型选择（type switch）---")

	// describe 函数用 type switch 处理多种类型
	// C 差异：C 需要 switch(tag) { case TYPE_INT: ... case TYPE_STRING: ... }
	describe := func(i interface{}) string {
		switch v := i.(type) {
		case int:
			return fmt.Sprintf("整数: %d（值*2=%d）", v, v*2)
		case float64:
			return fmt.Sprintf("浮点数: %.4f（平方=%.4f）", v, v*v)
		case string:
			return fmt.Sprintf("字符串: %q（长度=%d）", v, len(v))
		case bool:
			return fmt.Sprintf("布尔值: %v", v)
		case Circle:
			return fmt.Sprintf("Circle（半径=%.2f，面积=%.4f）", v.Radius, v.Area())
		case Rectangle:
			return fmt.Sprintf("Rectangle（%.2f×%.2f，面积=%.4f）", v.Width, v.Height, v.Area())
		case nil:
			return "nil 值"
		default:
			// v 的类型是 interface{}，可以用 %T 打印实际类型
			return fmt.Sprintf("未知类型: %T = %v", v, v)
		}
	}

	values := []interface{}{
		42,
		3.14,
		"hello",
		true,
		Circle{Radius: 2.0},
		Rectangle{Width: 3.0, Height: 4.0},
		nil,
		[]int{1, 2, 3},
	}

	fmt.Println("type switch 处理多种类型：")
	for _, v := range values {
		fmt.Printf("  %v\n", describe(v))
	}

	// type switch 也可以用于接口类型
	fmt.Println("\ntype switch 处理 Shape 接口：")
	shapes := []Shape{Circle{Radius: 1}, Rectangle{Width: 2, Height: 3}}
	for _, s := range shapes {
		switch v := s.(type) {
		case Circle:
			fmt.Printf("  Circle: 半径=%.2f\n", v.Radius)
		case Rectangle:
			fmt.Printf("  Rectangle: %.2f×%.2f\n", v.Width, v.Height)
		}
	}
	fmt.Println("C 差异：Go type switch 自动绑定具体类型变量，C 需要手动强制转换")
}

// -----------------------------------------------------------------------------
// 8. 接口值的内部结构：nil 接口 vs 持有 nil 指针的接口值
// C 差异：
//   - C 的 void* 为 NULL 就是 NULL，没有"持有 NULL 的非 NULL 指针"的概念；
//   - Go 的接口值内部是（类型, 值）对：nil 接口两者都为 nil；
//   - 持有 nil 指针的接口值：类型不为 nil，值为 nil——接口本身不为 nil！
//   - 这是 Go 最著名的陷阱之一，尤其在错误处理中容易踩坑。
// -----------------------------------------------------------------------------
func demoNilInterface() {
	fmt.Println("\n--- 8. 接口值内部结构：nil 接口 vs 持有 nil 指针的接口值 ---")

	// 情况 1：nil 接口——类型和值都为 nil
	var s1 Shape // 零值：(type=nil, value=nil)
	fmt.Printf("nil 接口:              s1 == nil → %v\n", s1 == nil)
	fmt.Printf("nil 接口内部:          (type=nil, value=nil)\n")

	// 情况 2：持有 nil 指针的接口值——类型不为 nil，值为 nil
	// 这是 Go 的经典陷阱！
	var c *Circle = nil          // c 是 nil 指针
	var s2 Shape = c             // s2 持有 (*Circle, nil)——接口不为 nil！
	fmt.Printf("\n持有 nil 指针的接口:   s2 == nil → %v（陷阱！）\n", s2 == nil)
	fmt.Printf("持有 nil 指针的接口内部: (type=*Circle, value=nil)\n")
	fmt.Println("注意：s2 != nil，因为接口内部的类型字段不为 nil！")

	// 演示陷阱场景：函数返回接口类型时的常见错误
	fmt.Println("\n演示经典陷阱（函数返回接口）：")

	// 错误写法：返回具体类型的 nil 指针，赋值给接口后不为 nil
	badGetShape := func() Shape {
		var c *Circle = nil
		return c // 返回 (*Circle, nil)，不是 nil 接口！
	}

	// 正确写法：直接返回 nil（nil 接口）
	goodGetShape := func() Shape {
		return nil // 返回 (nil, nil)，是真正的 nil 接口
	}

	bad := badGetShape()
	good := goodGetShape()
	fmt.Printf("  badGetShape()  == nil → %v（陷阱：返回了持有 nil 指针的接口）\n", bad == nil)
	fmt.Printf("  goodGetShape() == nil → %v（正确：返回了 nil 接口）\n", good == nil)

	// 调用持有 nil 指针的接口方法会 panic（如果方法解引用了指针）
	fmt.Println("\n演示调用持有 nil 指针的接口方法（将被 recover 捕获）：")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到 panic: %v\n", r)
			}
		}()
		var c *Circle = nil
		var s Shape = c
		// 调用 Area()：(*Circle).Area() 会解引用 nil 指针，panic
		fmt.Println(s.Area())
	}()

	fmt.Println("C 差异：C 的 NULL 指针就是 NULL，Go 的接口有类型和值两个字段")
	fmt.Println("经验：函数返回接口类型时，应直接返回 nil，而非具体类型的 nil 指针")
}

// -----------------------------------------------------------------------------
// 9. 标准库常用接口：fmt.Stringer、error、io.Reader/io.Writer
// C 差异：
//   - C 没有标准接口的概念，每个库定义自己的回调函数类型；
//   - Go 标准库定义了一套通用接口，任何类型都可以实现这些接口；
//   - 实现 fmt.Stringer 后，fmt.Println 自动调用 String() 方法；
//   - 实现 error 接口后，可以作为错误值传递；
//   - io.Reader/io.Writer 是 Go I/O 系统的基础，所有 I/O 操作都基于这两个接口。
// -----------------------------------------------------------------------------
func demoStandardInterfaces() {
	fmt.Println("\n--- 9. 标准库常用接口：fmt.Stringer、error、io.Reader/io.Writer ---")

	// fmt.Stringer 接口：String() string
	// 实现后，fmt.Println/Printf(%v/%s) 自动调用 String() 方法
	// C 差异：C 需要手动调用 animal_to_string()，Go 通过接口自动调用
	fmt.Println("【fmt.Stringer 接口】")
	c := Circle{Radius: 3.0}
	r := Rectangle{Width: 4.0, Height: 5.0}
	fmt.Printf("  fmt.Println(Circle):    %v\n", c)    // 自动调用 c.String()
	fmt.Printf("  fmt.Println(Rectangle): %v\n", r)    // 自动调用 r.String()
	fmt.Printf("  fmt.Sprintf(%%s):        %s\n", c)   // %s 也调用 String()
	fmt.Println("  注意：实现 String() string 方法后，fmt 包自动调用（无需手动）")

	// error 接口：Error() string
	// C 差异：C 用 errno 全局变量或返回负数，Go 用 error 接口（更安全、更清晰）
	fmt.Println("\n【error 接口】")
	_, err := NewCircle(-1.0)
	if err != nil {
		fmt.Printf("  NewCircle(-1) 错误: %v\n", err)
		// 类型断言获取具体错误信息
		if se, ok := err.(*ShapeError); ok {
			fmt.Printf("  ShapeError.Op=%q, Val=%.2f\n", se.Op, se.Val)
		}
	}
	c2, err2 := NewCircle(5.0)
	if err2 == nil {
		fmt.Printf("  NewCircle(5) 成功: %v\n", c2)
	}
	fmt.Println("  注意：error 接口只有 Error() string 一个方法，任何类型都可以实现")
	fmt.Println("  C 差异：Go 的 error 接口比 C 的 errno 更安全（类型安全，可携带上下文）")

	// io.Reader 接口：Read(p []byte) (n int, err error)
	// C 差异：C 的 fread/read 是具体函数，Go 的 io.Reader 是接口（任何类型都可以实现）
	fmt.Println("\n【io.Reader 接口】")
	sr := &StringReader{data: "Hello, Go Interfaces!"}
	buf := make([]byte, 8)
	fmt.Println("  从 StringReader 读取数据（每次 8 字节）：")
	for {
		n, err := sr.Read(buf)
		if n > 0 {
			fmt.Printf("    读取 %d 字节: %q\n", n, buf[:n])
		}
		if err == io.EOF {
			fmt.Println("    读取完毕（io.EOF）")
			break
		}
		if err != nil {
			fmt.Printf("    读取错误: %v\n", err)
			break
		}
	}

	// io.Reader 的强大之处：任何实现了 Read 方法的类型都可以用于 io.Copy 等函数
	sr2 := &StringReader{data: "io.Reader 是 Go I/O 的基础接口"}
	var sb strings.Builder
	// io.Copy 接受 io.Writer 和 io.Reader，与具体类型无关
	n, _ := io.Copy(&sb, sr2)
	fmt.Printf("  io.Copy 复制了 %d 字节: %q\n", n, sb.String())

	fmt.Println("\n【接口设计哲学总结】")
	fmt.Println("  1. 接口越小越好（最小接口原则）：io.Reader 只有 1 个方法")
	fmt.Println("  2. 接口在使用方定义（而非实现方）：解耦依赖")
	fmt.Println("  3. 接受接口，返回具体类型（Accept interfaces, return structs）")
	fmt.Println("  4. 空接口 interface{} 失去类型安全，优先使用泛型（Go 1.18+）")
	fmt.Println("  C 差异：Go 接口是隐式实现的鸭子类型，C 需要手动绑定函数指针")
}
