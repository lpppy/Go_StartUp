// Package testing_demo 演示 Go 语言的测试框架。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 测试框架的关键差异说明。
//
// 与 C 测试框架的对比：
//   - C 通常使用第三方框架（如 CUnit、Check、cmocka）或手写断言宏；
//   - Go 内置 testing 包，无需第三方依赖，go test 命令直接运行测试；
//   - Go 的测试文件以 _test.go 结尾，测试函数以 Test 开头，签名固定；
//   - Go 支持基准测试（Benchmark）和示例测试（Example），C 框架通常不内置。
package testing_demo

import "fmt"

// Add 返回两个整数之和。
// 这是一个简单的可测试函数，用于演示单元测试。
func Add(a, b int) int {
	return a + b
}

// Fibonacci 返回第 n 个斐波那契数（递归实现）。
// F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2)
// 注意：递归实现简洁但效率低（指数级时间复杂度），仅用于演示。
func Fibonacci(n int) int {
	if n < 0 {
		return 0 // 负数返回 0
	}
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// IsPrime 判断 n 是否为质数。
// 质数：大于 1 且只能被 1 和自身整除的正整数。
func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// Demo 演示测试相关概念（通过打印说明，实际测试在 _test.go 文件中）。
func Demo() {
	fmt.Println("测试模块的实际演示在 testing_demo_test.go 中")
	fmt.Println("运行方式：")
	fmt.Println("  go test ./demos/testing_demo/          # 运行所有测试")
	fmt.Println("  go test -v ./demos/testing_demo/       # 详细输出")
	fmt.Println("  go test -bench=. ./demos/testing_demo/ # 基准测试")
	fmt.Printf("Add(2, 3) = %d\n", Add(2, 3))
	fmt.Printf("Fibonacci(10) = %d\n", Fibonacci(10))
	fmt.Printf("IsPrime(97) = %v\n", IsPrime(97))
}
