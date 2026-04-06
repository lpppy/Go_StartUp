// Package testing_demo 演示 Go 测试框架的各种用法。
//
// 测试文件命名约定：
//   - 测试文件必须以 _test.go 结尾（如 testing_demo_test.go）
//   - 测试文件与被测代码在同一包中（也可以用 package xxx_test 做黑盒测试）
//   - 运行命令：go test ./demos/testing_demo/
//   - 详细输出：go test -v ./demos/testing_demo/
//   - 运行基准测试：go test -bench=. ./demos/testing_demo/
//
// 与 C 测试框架的对比：
//   - C 需要第三方框架（CUnit、Check 等），Go 内置 testing 包；
//   - C 测试通常需要手写 main 函数，Go 的 go test 自动发现并运行测试；
//   - Go 的测试函数签名固定：TestXxx(t *testing.T)，无需注册；
//   - Go 内置基准测试和示例测试，C 框架通常不支持。
package testing_demo

import (
	"fmt"
	"testing"
)

// =============================================================================
// 1. 基本测试：TestXxx(t *testing.T) 签名
// t.Error/t.Errorf：报告失败但继续执行
// t.Fatal/t.Fatalf：报告失败并立即停止当前测试函数
// =============================================================================

// TestAdd 演示基本测试函数的写法。
// 测试函数必须以 Test 开头，参数为 *testing.T，无返回值。
func TestAdd(t *testing.T) {
	result := Add(2, 3)
	expected := 5

	// t.Errorf：格式化错误信息，测试失败但继续执行后续代码
	if result != expected {
		t.Errorf("Add(2, 3) = %d, 期望 %d", result, expected)
	}

	// t.Error：不格式化，测试失败但继续执行
	if Add(0, 0) != 0 {
		t.Error("Add(0, 0) 应该返回 0")
	}

	// t.Fatal：测试失败并立即停止当前测试函数（后续代码不执行）
	// 适用于：后续测试依赖此步骤的结果时
	if Add(-1, 1) != 0 {
		t.Fatalf("Add(-1, 1) = %d, 期望 0，此错误是致命的", Add(-1, 1))
	}

	// t.Fatalf：格式化版本的 t.Fatal
	// 如果上面的 Fatal 触发，这行不会执行
	if Add(100, -100) != 0 {
		t.Fatalf("Add(100, -100) = %d, 期望 0", Add(100, -100))
	}
}

// =============================================================================
// 2. 表驱动测试（Table-Driven Test）
// Go 最推荐的测试模式：用结构体切片定义测试用例，循环执行
// 优点：易于添加新用例，测试逻辑集中，输出清晰
// =============================================================================

// TestFibonacci 演示表驱动测试模式。
func TestFibonacci(t *testing.T) {
	// 定义测试用例表：结构体切片
	// C 对比：C 通常用数组或手动列举，Go 的结构体切片更清晰
	tests := []struct {
		name     string // 测试用例名称（用于 t.Run 和错误信息）
		input    int    // 输入
		expected int    // 期望输出
	}{
		{"F(0)", 0, 0},
		{"F(1)", 1, 1},
		{"F(2)", 2, 1},
		{"F(3)", 3, 2},
		{"F(4)", 4, 3},
		{"F(5)", 5, 5},
		{"F(6)", 6, 8},
		{"F(7)", 7, 13},
		{"F(10)", 10, 55},
		{"负数", -1, 0},
	}

	// 循环执行所有测试用例
	for _, tt := range tests {
		// 使用 tt.name 作为子测试名称（见下面的 t.Run 演示）
		result := Fibonacci(tt.input)
		if result != tt.expected {
			t.Errorf("Fibonacci(%d) = %d, 期望 %d（用例：%s）",
				tt.input, result, tt.expected, tt.name)
		}
	}
}

// TestIsPrime 演示表驱动测试（含边界情况）。
func TestIsPrime(t *testing.T) {
	tests := []struct {
		n        int
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{11, true},
		{12, false},
		{13, true},
		{97, true},
		{100, false},
	}

	for _, tt := range tests {
		result := IsPrime(tt.n)
		if result != tt.expected {
			t.Errorf("IsPrime(%d) = %v, 期望 %v", tt.n, result, tt.expected)
		}
	}
}

// =============================================================================
// 3. 子测试 t.Run("name", func(t *testing.T) {...})
// 优点：
//   - 可以单独运行某个子测试：go test -run TestAdd/正数相加
//   - 失败的子测试不影响其他子测试
//   - 输出更有层次感
// =============================================================================

// TestAddWithSubtests 演示子测试的用法。
func TestAddWithSubtests(t *testing.T) {
	// t.Run 创建子测试，第一个参数是子测试名称
	t.Run("正数相加", func(t *testing.T) {
		if Add(1, 2) != 3 {
			t.Errorf("Add(1, 2) 应该等于 3")
		}
	})

	t.Run("负数相加", func(t *testing.T) {
		if Add(-1, -2) != -3 {
			t.Errorf("Add(-1, -2) 应该等于 -3")
		}
	})

	t.Run("正负相加", func(t *testing.T) {
		if Add(5, -3) != 2 {
			t.Errorf("Add(5, -3) 应该等于 2")
		}
	})

	t.Run("零值", func(t *testing.T) {
		if Add(0, 0) != 0 {
			t.Errorf("Add(0, 0) 应该等于 0")
		}
		if Add(42, 0) != 42 {
			t.Errorf("Add(42, 0) 应该等于 42")
		}
	})
}

// =============================================================================
// 4. 并行子测试 t.Parallel()
// 用途：加速独立测试的执行（并发运行）
// 注意：并行测试中不能使用循环变量（需要捕获副本）
// =============================================================================

// TestFibonacciParallel 演示并行子测试。
func TestFibonacciParallel(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"F(5)", 5, 5},
		{"F(8)", 8, 21},
		{"F(10)", 10, 55},
		{"F(12)", 12, 144},
	}

	for _, tt := range tests {
		tt := tt // 重要：捕获循环变量副本，避免并行测试中的数据竞争
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // 标记此子测试可以并行运行

			result := Fibonacci(tt.input)
			if result != tt.expected {
				t.Errorf("Fibonacci(%d) = %d, 期望 %d", tt.input, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 5. 基准测试 BenchmarkXxx(b *testing.B)
// 运行命令：go test -bench=. ./demos/testing_demo/
// 运行特定基准：go test -bench=BenchmarkFibonacci ./demos/testing_demo/
// b.N：测试框架自动调整的迭代次数，确保测试运行足够长时间
// =============================================================================

// BenchmarkAdd 基准测试 Add 函数的性能。
func BenchmarkAdd(b *testing.B) {
	// b.N 由测试框架自动设置，确保基准测试运行足够长时间
	for i := 0; i < b.N; i++ {
		Add(100, 200)
	}
}

// BenchmarkFibonacci 基准测试 Fibonacci 函数（递归实现较慢）。
func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20) // 测试 F(20)，递归调用次数较多
	}
}

// BenchmarkFibonacciSmall 对比小输入的性能。
func BenchmarkFibonacciSmall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

// BenchmarkIsPrime 基准测试 IsPrime 函数。
func BenchmarkIsPrime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(97) // 97 是质数，需要完整检查
	}
}

// =============================================================================
// 6. 示例测试 ExampleXxx()
// 特点：
//   - 函数名以 Example 开头
//   - 通过 // Output: 注释指定期望输出
//   - go test 会运行示例并验证输出是否匹配
//   - 同时作为文档（go doc 会显示示例）
// =============================================================================

// ExampleAdd 演示 Add 函数的用法（示例测试）。
func ExampleAdd() {
	fmt.Println(Add(2, 3))
	fmt.Println(Add(-1, 1))
	fmt.Println(Add(0, 0))
	// Output:
	// 5
	// 0
	// 0
}

// ExampleFibonacci 演示 Fibonacci 函数的用法。
func ExampleFibonacci() {
	fmt.Println(Fibonacci(0))
	fmt.Println(Fibonacci(1))
	fmt.Println(Fibonacci(10))
	// Output:
	// 0
	// 1
	// 55
}

// ExampleIsPrime 演示 IsPrime 函数的用法。
func ExampleIsPrime() {
	fmt.Println(IsPrime(2))
	fmt.Println(IsPrime(4))
	fmt.Println(IsPrime(97))
	// Output:
	// true
	// false
	// true
}
