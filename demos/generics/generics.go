// Package generics 演示 Go 1.18+ 泛型特性。
// 面向有 C 语言基础的开发者，每个演示均附有与 C/C++ 的关键差异说明。
//
// 与 C/C++ 的对比：
//   - C 通常用 void* 实现泛型（类型不安全，需要手动转换）；
//   - C 用宏（#define）实现泛型算法（类型不安全，调试困难）；
//   - C++ 用模板（template）实现泛型（强大但语法复杂，编译错误难读）；
//   - Go 的泛型通过类型参数（type parameters）实现，语法简洁，错误信息清晰。
package generics

import (
	"fmt"
	"slices"
)

// Map 将切片中的每个元素通过函数 f 转换为新类型，返回新切片。
// 类型参数：T（输入元素类型）、U（输出元素类型），约束为 any（任意类型）。
func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Filter 返回切片中满足谓词函数 f 的元素组成的新切片。
func Filter[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Contains 检查切片中是否包含元素 v。
// 类型参数：T，约束为 comparable（支持 == 和 != 运算符）。
func Contains[T comparable](s []T, v T) bool {
	for _, elem := range s {
		if elem == v {
			return true
		}
	}
	return false
}

// Number 是自定义约束接口，表示所有数值类型。
// ~ 前缀表示"底层类型为 T"，匹配所有以该类型为底层类型的自定义类型。
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// Sum 计算数值切片的总和。
func Sum[T Number](s []T) T {
	var total T
	for _, v := range s {
		total += v
	}
	return total
}

// Stack 是泛型栈，支持任意类型 T。
// C 对比：C 需要为每种类型写单独的栈，或用 void*（不安全）；
//         C++ 用 template<typename T> class Stack 实现。
type Stack[T any] struct {
	items []T
}

// Push 将元素压入栈顶。
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop 弹出栈顶元素，若栈为空则返回零值和 false。
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top, true
}

// Peek 查看栈顶元素但不弹出，若栈为空则返回零值和 false。
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len 返回栈中元素数量。
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty 返回栈是否为空。
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Demo 演示所有泛型相关内容。
func Demo() {
	demoGenericFunctions()
	demoConstraints()
	demoUnderlyingTypeConstraint()
	demoGenericStack()
	demoTypeInference()
	demoSlicesPackage()
	demoComparison()
}

func demoGenericFunctions() {
	fmt.Println("\n--- 1. 泛型函数定义语法 ---")

	nums := []int{1, 2, 3, 4, 5}
	strs := Map(nums, func(n int) string {
		return fmt.Sprintf("item_%d", n)
	})
	fmt.Printf("Map([]int -> []string): %v\n", strs)

	words := []string{"Go", "泛型", "很好用"}
	lengths := Map(words, func(s string) int { return len([]rune(s)) })
	fmt.Printf("Map([]string -> []int 长度): %v\n", lengths)

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter(偶数): %v\n", evens)

	longWords := Filter(words, func(s string) bool { return len([]rune(s)) > 2 })
	fmt.Printf("Filter(长度>2的字符串): %v\n", longWords)

	fmt.Printf("Contains([1,2,3,4,5], 3): %v\n", Contains(nums, 3))
	fmt.Printf("Contains([1,2,3,4,5], 6): %v\n", Contains(nums, 6))
	fmt.Printf("Contains([\"Go\",\"泛型\",\"很好用\"], \"Go\"): %v\n", Contains(words, "Go"))

	fmt.Println("\n泛型函数语法：func Name[TypeParam Constraint](params) ReturnType")
}

func demoConstraints() {
	fmt.Println("\n--- 2. 类型参数约束 ---")

	fmt.Println("any 约束：接受任意类型")
	mixedMap := Map([]any{1, "hello", true, 3.14}, func(v any) string {
		return fmt.Sprintf("%v", v)
	})
	fmt.Printf("  Map(any 切片): %v\n", mixedMap)

	fmt.Println("\ncomparable 约束：支持 == 和 != 的类型")
	fmt.Printf("  Contains([]int, 42): %v\n", Contains([]int{10, 20, 42, 50}, 42))
	fmt.Printf("  Contains([]string, \"Go\"): %v\n", Contains([]string{"Python", "Go", "Rust"}, "Go"))

	fmt.Println("\nNumber 自定义约束：~int | ~int8 | ... | ~float64")
	fmt.Printf("  Sum([]int{1,2,3,4,5}): %d\n", Sum([]int{1, 2, 3, 4, 5}))
	fmt.Printf("  Sum([]float64{1.1,2.2,3.3}): %.1f\n", Sum([]float64{1.1, 2.2, 3.3}))

	fmt.Println("\n约束类型总结:")
	fmt.Println("  any         → 任意类型（最宽松）")
	fmt.Println("  comparable  → 支持 == 的类型（int, string, struct 等，不含 slice/map）")
	fmt.Println("  自定义接口  → 通过接口定义允许的类型集合")
}

// MyInt 是基于 int 的自定义类型，底层类型为 int
type MyInt int

// MyFloat 是基于 float64 的自定义类型
type MyFloat float64

func demoUnderlyingTypeConstraint() {
	fmt.Println("\n--- 3. ~T 底层类型约束 ---")

	myInts := []MyInt{1, 2, 3, 4, 5}
	total := Sum(myInts)
	fmt.Printf("Sum([]MyInt{1,2,3,4,5}): %v（MyInt 底层类型为 int，满足 ~int）\n", total)

	myFloats := []MyFloat{1.5, 2.5, 3.0}
	floatTotal := Sum(myFloats)
	fmt.Printf("Sum([]MyFloat{1.5,2.5,3.0}): %v（MyFloat 底层类型为 float64）\n", floatTotal)

	fmt.Println("\n~ 前缀的含义：")
	fmt.Println("  ~int  → 匹配所有底层类型为 int 的类型（int 本身 + type MyInt int 等）")
	fmt.Println("  int   → 只匹配 int 本身（不含自定义类型）")
}

func demoGenericStack() {
	fmt.Println("\n--- 4. 泛型结构体 Stack[T any] ---")

	var intStack Stack[int]
	fmt.Printf("初始状态: IsEmpty=%v, Len=%d\n", intStack.IsEmpty(), intStack.Len())

	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	fmt.Printf("Push 10, 20, 30: Len=%d\n", intStack.Len())

	if top, ok := intStack.Peek(); ok {
		fmt.Printf("Peek: %d（不弹出）\n", top)
	}

	for !intStack.IsEmpty() {
		if v, ok := intStack.Pop(); ok {
			fmt.Printf("Pop: %d\n", v)
		}
	}
	fmt.Printf("全部弹出后: IsEmpty=%v\n", intStack.IsEmpty())

	_, ok := intStack.Pop()
	fmt.Printf("空栈 Pop: ok=%v\n", ok)

	var strStack Stack[string]
	strStack.Push("Go")
	strStack.Push("泛型")
	strStack.Push("真好用")
	fmt.Printf("\nstring 栈 Len=%d\n", strStack.Len())
	if top, ok := strStack.Peek(); ok {
		fmt.Printf("string 栈 Peek: %q\n", top)
	}

	fmt.Println("\n泛型结构体语法：type Stack[T any] struct { items []T }")
	fmt.Println("方法语法：func (s *Stack[T]) Push(v T) { ... }")
}

func demoTypeInference() {
	fmt.Println("\n--- 5. 类型推断 ---")

	nums := []int{1, 2, 3}

	result1 := Map[int, string](nums, func(n int) string { return fmt.Sprintf("%d", n) })
	fmt.Printf("显式类型参数 Map[int, string]: %v\n", result1)

	result2 := Map(nums, func(n int) string { return fmt.Sprintf("%d", n) })
	fmt.Printf("类型推断 Map（无需指定类型）: %v\n", result2)

	fmt.Printf("Contains 类型推断: %v\n", Contains(nums, 2))
	fmt.Printf("Sum 类型推断: %d\n", Sum(nums))

	fmt.Println("\n类型推断：编译器从函数参数的类型自动推断类型参数")
	fmt.Println("大多数情况下无需显式指定类型参数")
}

func demoSlicesPackage() {
	fmt.Println("\n--- 6. 标准库 slices 包（Go 1.21+）---")

	nums := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	strs := []string{"banana", "apple", "cherry", "date"}

	fmt.Printf("slices.Contains([3,1,4,...], 5): %v\n", slices.Contains(nums, 5))
	fmt.Printf("slices.Contains([3,1,4,...], 7): %v\n", slices.Contains(nums, 7))

	sortedNums := make([]int, len(nums))
	copy(sortedNums, nums)
	slices.Sort(sortedNums)
	fmt.Printf("slices.Sort([]int): %v\n", sortedNums)

	sortedStrs := make([]string, len(strs))
	copy(sortedStrs, strs)
	slices.Sort(sortedStrs)
	fmt.Printf("slices.Sort([]string): %v\n", sortedStrs)

	idx := slices.Index(nums, 9)
	fmt.Printf("slices.Index([3,1,4,...], 9): %d\n", idx)

	rev := []int{1, 2, 3, 4, 5}
	slices.Reverse(rev)
	fmt.Printf("slices.Reverse([1,2,3,4,5]): %v\n", rev)

	fmt.Printf("slices.Max: %d\n", slices.Max(nums))
	fmt.Printf("slices.Min: %d\n", slices.Min(nums))

	fmt.Println("\nslices 包（Go 1.21+）提供了大量泛型切片操作函数")
}

func demoComparison() {
	fmt.Println("\n--- 7. 与 C/C++ 的对比 ---")

	fmt.Println("C void* 方式（类型不安全）：需要手动转换类型，运行时才能发现类型错误")
	fmt.Println("C 宏方式（预处理器）：没有类型检查，调试困难，可能有副作用")
	fmt.Println("C++ 模板方式：语法复杂，编译错误信息难读，编译时间长")
	fmt.Println("Go 泛型方式：语法简洁，编译时类型检查，错误信息清晰")
	fmt.Println("\nGo 泛型限制：不支持泛型方法（只能在函数级别使用类型参数）")
	fmt.Println("            不支持特化（C++ 模板特化），约束系统相对简单")
}
