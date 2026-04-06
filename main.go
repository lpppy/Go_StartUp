package main

import (
	"fmt"

	"go-basics-demo/demos/collections"
	"go-basics-demo/demos/concurrency"
	"go-basics-demo/demos/controlflow"
	errors_demo "go-basics-demo/demos/errors"
	"go-basics-demo/demos/functions"
	"go-basics-demo/demos/generics"
	"go-basics-demo/demos/interfaces"
	"go-basics-demo/demos/packages"
	"go-basics-demo/demos/pointers"
	"go-basics-demo/demos/stdlib"
	"go-basics-demo/demos/structs"
	"go-basics-demo/demos/testing_demo"
	"go-basics-demo/demos/variables"
)

// demoEntry 表示一个演示模块的注册项
type demoEntry struct {
	name string
	fn   func()
}

// demos 是所有演示模块的注册表
var demos = []demoEntry{
	{"变量与基本类型", variables.Demo},
	{"控制流", controlflow.Demo},
	{"函数", functions.Demo},
	{"集合类型", collections.Demo},
	{"结构体", structs.Demo},
	{"接口", interfaces.Demo},
	{"错误处理", errors_demo.Demo},
	{"并发", concurrency.Demo},
	{"指针", pointers.Demo},
	{"包管理", packages.Demo},
	{"泛型", generics.Demo},
	{"测试", testing_demo.Demo},
	{"标准库", stdlib.Demo},
}

// runSafe 安全地执行演示函数，使用 defer/recover 捕获 panic
func runSafe(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PANIC in %s]: %v\n", name, r)
		}
	}()
	fn()
}

func main() {
	for i, d := range demos {
		fmt.Printf("\n========== %d. %s ==========\n", i+1, d.name)
		runSafe(d.name, d.fn)
	}
}
