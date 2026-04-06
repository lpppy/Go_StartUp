// Package controlflow 演示 Go 语言的控制流语句。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
package controlflow

import "fmt"

// Demo 演示所有控制流相关内容。
func Demo() {
	demoForLoop()
	demoIfElse()
	demoIfWithInit()
	demoSwitch()
	demoSwitchFallthrough()
	demoSwitchNoExpr()
	demoRange()
	demoLabeledBreakContinue()
	demoGoto()
}

// TypeSwitchDemo 演示 type switch，根据传入值的类型返回类型描述字符串。
// 供测试调用。
//
// C 差异：C 没有运行时类型信息（RTTI）机制，无法在运行时判断变量的具体类型；
// Go 的 interface{} 携带动态类型信息，type switch 可以安全地分支处理多种类型。
func TypeSwitchDemo(i interface{}) string {
	switch v := i.(type) {
	case int:
		return fmt.Sprintf("int: %d", v)
	case float64:
		return fmt.Sprintf("float64: %f", v)
	case string:
		return fmt.Sprintf("string: %q", v)
	case bool:
		return fmt.Sprintf("bool: %t", v)
	case []int:
		return fmt.Sprintf("[]int: len=%d", len(v))
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

// -----------------------------------------------------------------------------
// 1. for 三种形式
// C 差异：
//   - Go 只有 for 关键字，没有 while 和 do-while；
//     仅条件式 for 等价于 C 的 while，无限循环等价于 C 的 while(1) 或 for(;;)。
//   - Go 的 for 条件表达式不需要括号（C 要求括号），但花括号是强制要求的。
// -----------------------------------------------------------------------------
func demoForLoop() {
	fmt.Println("\n--- 1. for 三种形式 ---")

	// 形式一：标准三段式（init; condition; post）
	// C 等价：for (int i = 0; i < 3; i++) { ... }
	fmt.Print("标准三段式 for: ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 形式二：仅条件式（等价于 while）
	// C 等价：while (n < 8) { n *= 2; }
	n := 1
	fmt.Print("仅条件式 for（等价 while）: ")
	for n < 8 {
		fmt.Printf("%d ", n)
		n *= 2
	}
	fmt.Println()

	// 形式三：无限循环（等价于 while(1) 或 for(;;)），用 break 退出
	// C 等价：for (;;) { if (count >= 3) break; count++; }
	count := 0
	fmt.Print("无限循环 for + break: ")
	for {
		if count >= 3 {
			break
		}
		fmt.Printf("%d ", count)
		count++
	}
	fmt.Println()
}

// -----------------------------------------------------------------------------
// 2. if/else 语句
// C 差异：
//   - Go 的 if 条件表达式不需要括号（C 要求括号）。
//   - Go 强制要求花括号，即使只有一条语句（C 允许省略花括号）。
//   - else 必须与上一个 } 在同一行（Go 的分号自动插入规则决定）。
// -----------------------------------------------------------------------------
func demoIfElse() {
	fmt.Println("\n--- 2. if/else 语句 ---")

	x := 42

	// 条件无需括号，花括号强制要求
	// C 等价：if (x > 0) { ... } else if (x < 0) { ... } else { ... }
	if x > 0 {
		fmt.Printf("x=%d 是正数\n", x)
	} else if x < 0 {
		fmt.Printf("x=%d 是负数\n", x)
	} else {
		fmt.Printf("x=%d 是零\n", x)
	}

	// else 必须与 } 在同一行，否则编译错误（Go 自动插入分号）
	// 以下写法是正确的：
	score := 85
	var grade string
	if score >= 90 {
		grade = "A"
	} else if score >= 80 {
		grade = "B"
	} else if score >= 70 {
		grade = "C"
	} else {
		grade = "D"
	}
	fmt.Printf("score=%d, grade=%s\n", score, grade)
}

// doSomething 是一个辅助函数，用于演示带初始化语句的 if。
// 返回 nil 表示成功，模拟实际业务中的错误返回。
func doSomething() error {
	return nil // 模拟成功
}

// doSomethingFail 模拟返回错误的函数。
func doSomethingFail() error {
	return fmt.Errorf("操作失败：权限不足")
}

// -----------------------------------------------------------------------------
// 3. 带初始化语句的 if
// C 差异：
//   - C 没有在 if 条件中声明变量的语法（C99 的 for 循环支持，但 if 不支持）。
//   - Go 的 if 初始化语句声明的变量作用域仅限于该 if/else if/else 块，
//     不会污染外部作用域，这是一种常见的 Go 惯用法（idiom）。
// -----------------------------------------------------------------------------
func demoIfWithInit() {
	fmt.Println("\n--- 3. 带初始化语句的 if ---")

	// 语法：if init; condition { ... }
	// err 的作用域仅限于这个 if/else 块
	if err := doSomething(); err != nil {
		fmt.Printf("操作失败: %v\n", err)
	} else {
		fmt.Println("操作成功（err 变量作用域仅在此 if/else 块内）")
	}
	// 此处 err 不可访问，编译器会报错：undefined: err

	// 实际错误处理示例
	if err := doSomethingFail(); err != nil {
		fmt.Printf("捕获到错误: %v\n", err)
	}

	// 带初始化语句的 if 常用于：
	// 1. 函数调用并检查错误
	// 2. 类型断言并检查 ok
	// 3. map 查找并检查存在性
	m := map[string]int{"a": 1, "b": 2}
	if val, ok := m["a"]; ok {
		fmt.Printf("map 查找（带初始化）: m[\"a\"] = %d\n", val)
	}
	if _, ok := m["z"]; !ok {
		fmt.Println("map 查找（带初始化）: m[\"z\"] 不存在")
	}
}

// -----------------------------------------------------------------------------
// 4. switch 自动 break 行为
// C 差异：
//   - C 的 switch 默认穿透（fall-through），需要显式 break 阻止；
//     Go 的 switch 默认不穿透，每个 case 执行完自动 break，无需显式写 break。
//   - Go 的 switch case 可以包含多个值（逗号分隔），C 需要多个 case 标签。
//   - Go 的 switch 表达式可以是任意可比较类型，不限于整数（C 只支持整数）。
// -----------------------------------------------------------------------------
func demoSwitch() {
	fmt.Println("\n--- 4. switch 自动 break（无需显式 break）---")

	day := "Monday"

	// Go switch 每个 case 执行完自动退出，无需 break
	// C 等价需要在每个 case 末尾加 break
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		// 一个 case 可以匹配多个值（逗号分隔）
		// C 需要写多个 case 标签
		fmt.Printf("%s 是工作日\n", day)
	case "Saturday", "Sunday":
		fmt.Printf("%s 是周末\n", day)
	default:
		fmt.Printf("%s 不是有效的星期名\n", day)
	}

	// switch 也可以用于整数
	code := 404
	switch code {
	case 200:
		fmt.Println("HTTP 200: OK")
	case 404:
		fmt.Println("HTTP 404: Not Found（自动 break，不会穿透到下一个 case）")
	case 500:
		fmt.Println("HTTP 500: Internal Server Error")
	default:
		fmt.Printf("HTTP %d: 未知状态码\n", code)
	}
}

// -----------------------------------------------------------------------------
// 5. switch 使用 fallthrough 显式穿透
// C 差异：
//   - C 默认穿透，需要 break 阻止；Go 默认不穿透，需要 fallthrough 显式穿透。
//   - Go 的 fallthrough 必须是 case 块的最后一条语句。
//   - fallthrough 会无条件执行下一个 case 的代码，不检查下一个 case 的条件。
// -----------------------------------------------------------------------------
func demoSwitchFallthrough() {
	fmt.Println("\n--- 5. switch fallthrough 显式穿透 ---")

	n := 2
	fmt.Printf("n=%d，使用 fallthrough 穿透演示:\n", n)

	switch n {
	case 1:
		fmt.Println("  case 1")
		fallthrough // 显式穿透到 case 2
	case 2:
		fmt.Println("  case 2（匹配）")
		fallthrough // 显式穿透到 case 3（注意：不检查 case 3 的条件）
	case 3:
		fmt.Println("  case 3（因 fallthrough 执行，即使 n != 3）")
	case 4:
		fmt.Println("  case 4（不会执行，case 3 没有 fallthrough）")
	}

	// fallthrough 的典型用途：版本兼容性处理
	version := 2
	fmt.Printf("\nversion=%d 支持的特性:\n", version)
	switch version {
	case 3:
		fmt.Println("  - 特性 C（v3 新增）")
		fallthrough
	case 2:
		fmt.Println("  - 特性 B（v2 新增）")
		fallthrough
	case 1:
		fmt.Println("  - 特性 A（v1 基础特性）")
	}
}

// -----------------------------------------------------------------------------
// 6. switch 无表达式形式（等价于 if-else if 链）
// C 差异：
//   - C 的 switch 必须有表达式；Go 的 switch 可以省略表达式，
//     此时每个 case 包含一个布尔条件，等价于 if-else if 链，但更清晰。
// -----------------------------------------------------------------------------
func demoSwitchNoExpr() {
	fmt.Println("\n--- 6. switch 无表达式（等价于 if-else if 链）---")

	temp := 25.5

	// 无表达式 switch：每个 case 是一个布尔条件
	// 等价于 if temp < 0 { ... } else if temp < 10 { ... } ...
	switch {
	case temp < 0:
		fmt.Printf("%.1f°C: 冰点以下，注意防冻\n", temp)
	case temp < 10:
		fmt.Printf("%.1f°C: 寒冷，需要穿厚衣服\n", temp)
	case temp < 20:
		fmt.Printf("%.1f°C: 凉爽，适合外出\n", temp)
	case temp < 30:
		fmt.Printf("%.1f°C: 舒适，天气宜人\n", temp)
	default:
		fmt.Printf("%.1f°C: 炎热，注意防暑\n", temp)
	}

	// 无表达式 switch 也常用于替代复杂的 if-else if 链
	hour := 14
	switch {
	case hour < 6:
		fmt.Printf("当前 %d 时：深夜\n", hour)
	case hour < 12:
		fmt.Printf("当前 %d 时：上午\n", hour)
	case hour < 18:
		fmt.Printf("当前 %d 时：下午\n", hour)
	default:
		fmt.Printf("当前 %d 时：晚上\n", hour)
	}
}

// -----------------------------------------------------------------------------
// 7. range 遍历
// C 差异：
//   - C 没有 range 关键字，遍历数组需要手动管理索引；
//     Go 的 range 提供了统一的遍历语法，适用于数组、切片、map、字符串、channel。
//   - range 遍历字符串时返回的是 rune（Unicode 码点），而非 byte，
//     能正确处理多字节 UTF-8 字符（C 需要手动处理 UTF-8 编码）。
//   - range 遍历 map 时顺序是随机的（Go 故意随机化，防止依赖顺序）。
// -----------------------------------------------------------------------------
func demoRange() {
	fmt.Println("\n--- 7. range 遍历 ---")

	// 遍历数组（index + value）
	arr := [3]string{"Go", "Python", "Rust"}
	fmt.Print("range 遍历数组: ")
	for i, v := range arr {
		fmt.Printf("[%d]=%s ", i, v)
	}
	fmt.Println()

	// 遍历切片（index + value）
	// 用 _ 忽略不需要的 index
	nums := []int{10, 20, 30, 40, 50}
	fmt.Print("range 遍历切片（仅 value）: ")
	for _, v := range nums {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// 遍历切片（仅 index）
	fmt.Print("range 遍历切片（仅 index）: ")
	for i := range nums {
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// 遍历 map（key + value，顺序随机）
	m := map[string]int{"alice": 90, "bob": 85, "carol": 92}
	fmt.Println("range 遍历 map（key+value，顺序随机）:")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// 遍历字符串（返回 rune，正确处理 Unicode）
	// 注意：index 是字节偏移量，不是字符索引
	s := "Go世界"
	fmt.Print("range 遍历字符串（rune）: ")
	for i, r := range s {
		fmt.Printf("[字节偏移%d]'%c'(U+%04X) ", i, r, r)
	}
	fmt.Println()
}

// -----------------------------------------------------------------------------
// 8. break/continue 配合标签（label）跳出多层循环
// C 差异：
//   - C 的 break/continue 只能作用于最近的一层循环，跳出多层循环需要 goto 或标志变量；
//     Go 支持带标签的 break/continue，可以直接跳出或继续指定的外层循环。
//   - 标签必须紧接在 for/switch/select 语句之前，且必须被使用（否则编译错误）。
// -----------------------------------------------------------------------------
func demoLabeledBreakContinue() {
	fmt.Println("\n--- 8. break/continue 配合标签 ---")

	// 带标签的 break：跳出外层循环
	// C 等价需要 goto 或 found 标志变量
	fmt.Println("带标签的 break（跳出外层循环）:")
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	target := 5
	found := false

OuterBreak:
	for i, row := range matrix {
		for j, val := range row {
			if val == target {
				fmt.Printf("  找到目标值 %d，位置 [%d][%d]，跳出所有循环\n", target, i, j)
				found = true
				break OuterBreak // 直接跳出外层 OuterBreak 标记的循环
			}
		}
	}
	if !found {
		fmt.Printf("  未找到目标值 %d\n", target)
	}

	// 带标签的 continue：跳过外层循环的当前迭代
	fmt.Println("带标签的 continue（跳过外层循环当前迭代）:")
OuterContinue:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				// 跳过外层循环的当前迭代（i 的当前值），继续下一个 i
				fmt.Printf("  i=%d, j=%d: 触发 continue OuterContinue，跳过 i=%d 的剩余迭代\n", i, j, i)
				continue OuterContinue
			}
			fmt.Printf("  i=%d, j=%d\n", i, j)
		}
	}
}

// -----------------------------------------------------------------------------
// 9. goto 语句
// C 差异：
//   - Go 和 C 都支持 goto，但两者都不推荐使用。
//   - Go 的 goto 有限制：不能跳过变量声明，不能跳入内层代码块。
//   - 实际 Go 代码中极少使用 goto，通常用 for/break/continue/return 替代。
//   - 注释说明：goto 会破坏代码的结构化流程，降低可读性和可维护性，
//     仅在极少数性能敏感的底层代码中偶尔使用（如生成的解析器代码）。
// -----------------------------------------------------------------------------
func demoGoto() {
	fmt.Println("\n--- 9. goto 语句（不推荐使用，仅作演示）---")

	i := 0

	// goto 跳转到标签处
	// 注意：这里用 goto 模拟循环，实际应使用 for 循环
	// C 等价：goto loop; ... loop: if (i < 3) { ... goto loop; }
loop:
	if i < 3 {
		fmt.Printf("  goto 模拟循环: i=%d\n", i)
		i++
		goto loop // 跳回 loop 标签
	}

	fmt.Println("  goto 演示结束")
	fmt.Println("  ⚠️  注意：实际开发中应使用 for 循环替代 goto，goto 会降低代码可读性")
}
