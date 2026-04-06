// Package collections 演示 Go 语言的集合类型：数组、切片、map。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
package collections

import (
	"fmt"
	"sort"
	"unsafe"
)

// Demo 演示所有集合类型相关内容。
// 最后一个演示（nil map 写入）会触发 panic，由 main.go 的 Runner recover 捕获。
func Demo() {
	demoArrayDeclaration()
	demoArrayValueType()
	demoSliceCreation()
	demoSliceAppend()
	demoSliceGrowth()
	demoSliceSharedBacking()
	demoCopy()
	demo2DSlice()
	demoMapBasics()
	demoMapReadWriteDelete()
	demoMapOrdering()
	demoMapZeroValue()
	demoNilMapPanic() // 触发 panic，由 Runner recover 捕获
}

// -----------------------------------------------------------------------------
// 1. 数组三种声明方式
// C 差异：
//   - C 的数组长度可以是运行时变量（VLA，C99），Go 的数组长度必须是编译期常量；
//   - Go 支持 [...] 让编译器根据初始化元素个数推断数组长度；
//   - Go 数组是值类型（赋值/传参时复制整个数组），C 数组退化为指针；
//   - Go 数组的长度是类型的一部分：[3]int 和 [5]int 是不同类型，不可互赋值。
// -----------------------------------------------------------------------------
func demoArrayDeclaration() {
	fmt.Println("\n--- 1. 数组三种声明方式 ---")

	// 方式一：var 声明，元素自动初始化为零值
	// C 等价：int a[5] = {0};（C 局部数组不自动清零，需显式初始化）
	var a [5]int
	fmt.Printf("var a [5]int（零值初始化）: %v\n", a)

	// 方式二：字面量初始化，指定长度
	// C 等价：int b[3] = {1, 2, 3};
	b := [3]int{1, 2, 3}
	fmt.Printf("[3]int{1,2,3}: %v\n", b)

	// 方式三：[...] 让编译器推断长度（编译期确定，不是运行时）
	// C 等价：int c[] = {1, 2, 3};（C 也支持省略长度，但语义略有不同）
	c := [...]int{10, 20, 30, 40}
	fmt.Printf("[...]int{10,20,30,40}（编译器推断长度=%d）: %v\n", len(c), c)

	// 数组长度是类型的一部分：[3]int 和 [4]int 是不同类型
	fmt.Printf("len(a)=%d, len(b)=%d, len(c)=%d\n", len(a), len(b), len(c))
	fmt.Println("注意：[3]int 和 [5]int 是不同类型，不可互赋值（C 数组退化为指针，无此限制）")
}

// -----------------------------------------------------------------------------
// 2. 数组是值类型：赋值时复制整个数组
// C 差异：
//   - C 的数组名在大多数上下文中退化为指向首元素的指针，赋值实际上是指针赋值；
//   - Go 的数组是值类型，赋值（或传参）时复制整个数组，修改副本不影响原数组；
//   - 这意味着大数组传参开销较大，如需避免复制应传递指针 *[N]T 或使用切片。
// -----------------------------------------------------------------------------
func demoArrayValueType() {
	fmt.Println("\n--- 2. 数组是值类型（赋值时复制整个数组）---")

	original := [3]int{1, 2, 3}
	// 赋值：复制整个数组（值语义）
	// C 差异：C 中 int copy[3] = original; 是非法的，数组不能直接赋值
	copyArr := original
	copyArr[0] = 999 // 修改副本

	fmt.Printf("original: %v（未受影响）\n", original)
	fmt.Printf("copyArr:  %v（修改了 [0]）\n", copyArr)
	fmt.Println("结论：Go 数组赋值是深拷贝，修改副本不影响原数组")

	// 验证：通过函数传参也是复制
	modifyArray := func(arr [3]int) {
		arr[0] = 777 // 修改的是副本
	}
	modifyArray(original)
	fmt.Printf("传参后 original: %v（函数内修改不影响外部）\n", original)
	fmt.Println("注意：大数组传参开销大，如需避免复制请传 *[N]T 指针或使用切片")
}

// -----------------------------------------------------------------------------
// 3. 切片三种创建方式
// C 差异：
//   - C 没有内置切片类型，动态数组需要 malloc/realloc 手动管理内存；
//   - Go 的切片是对底层数组的引用（包含指针、长度、容量三个字段），
//     由 GC 自动管理内存，无需手动释放；
//   - 切片是引用类型：赋值传递的是引用（共享底层数组），而非值的副本。
// -----------------------------------------------------------------------------
func demoSliceCreation() {
	fmt.Println("\n--- 3. 切片三种创建方式 ---")

	// 方式一：字面量（最常用）
	// C 等价：int *s1 = (int[]){1, 2, 3};（C99 复合字面量，但无长度/容量信息）
	s1 := []int{1, 2, 3}
	fmt.Printf("字面量 []int{1,2,3}: %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))

	// 方式二：make([]T, len, cap)，指定长度和容量
	// C 等价：int *s2 = calloc(3, sizeof(int));（但 Go 自动清零且有容量概念）
	s2 := make([]int, 3, 5) // 长度 3，容量 5，元素初始化为零值
	fmt.Printf("make([]int, 3, 5): %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 方式三：从数组截取 a[low:high]（左闭右开区间）
	// C 等价：int *s3 = arr + 1;（指针偏移，但无边界检查）
	arr := [5]int{10, 20, 30, 40, 50}
	s3 := arr[1:4] // 包含 arr[1], arr[2], arr[3]，不含 arr[4]
	fmt.Printf("arr[1:4]（从数组截取）: %v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 省略 low 或 high 的简写
	s4 := arr[:3]  // 等价于 arr[0:3]
	s5 := arr[2:]  // 等价于 arr[2:5]
	s6 := arr[:]   // 等价于 arr[0:5]，整个数组的切片
	fmt.Printf("arr[:3]=%v, arr[2:]=%v, arr[:]=%v\n", s4, s5, s6)
	fmt.Println("注意：切片是三元组（指针+长度+容量），Go 自动管理内存，无需 free")
}

// -----------------------------------------------------------------------------
// 4. 切片 append 操作
// C 差异：
//   - C 需要手动 realloc 扩容，容易出现内存泄漏和 use-after-free；
//   - Go 的 append 自动处理扩容，返回新切片（可能指向新底层数组）；
//   - append(a, b...) 语法将切片 b 展开追加到 a，类似 C 的 memcpy 追加。
// -----------------------------------------------------------------------------
func demoSliceAppend() {
	fmt.Println("\n--- 4. 切片 append 操作 ---")

	s := []int{1, 2, 3}
	fmt.Printf("初始切片: %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 追加单个元素
	s = append(s, 4)
	fmt.Printf("append(s, 4): %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 追加多个元素
	s = append(s, 5, 6, 7)
	fmt.Printf("append(s, 5,6,7): %v, len=%d, cap=%d\n", s, len(s), cap(s))

	// 追加另一个切片（使用 ... 展开）
	// C 等价：memcpy(dst + len, src, n * sizeof(int));（手动计算偏移）
	extra := []int{8, 9, 10}
	s = append(s, extra...) // extra... 将切片展开为独立参数
	fmt.Printf("append(s, extra...): %v\n", s)
	fmt.Println("注意：append 可能返回新底层数组的切片，必须用返回值（s = append(s, ...)）")
}

// -----------------------------------------------------------------------------
// 5. 切片扩容行为：对比 append 前后的 cap 和底层指针
// C 差异：
//   - C 的 realloc 可能原地扩容也可能移动内存，需要更新所有指针；
//   - Go 的 append 在容量不足时自动分配新底层数组（通常 2 倍扩容，大切片按比例缩小），
//     旧切片仍指向旧数组，新切片指向新数组，两者互不影响。
// -----------------------------------------------------------------------------
func demoSliceGrowth() {
	fmt.Println("\n--- 5. 切片扩容行为（cap 和底层指针变化）---")

	s := make([]int, 0, 2) // 初始容量 2
	fmt.Printf("初始: len=%d, cap=%d, ptr=%p\n", len(s), cap(s), s)

	for i := 1; i <= 5; i++ {
		oldCap := cap(s)
		oldPtr := fmt.Sprintf("%p", s)
		s = append(s, i)
		newCap := cap(s)
		newPtr := fmt.Sprintf("%p", s)

		if newCap != oldCap {
			fmt.Printf("append(%d): len=%d, cap %d→%d（扩容！），ptr %s→%s（新底层数组）\n",
				i, len(s), oldCap, newCap, oldPtr, newPtr)
		} else {
			fmt.Printf("append(%d): len=%d, cap=%d（未扩容），ptr=%s\n",
				i, len(s), newCap, newPtr)
		}
	}

	// 用 unsafe.Pointer 直接展示底层数组地址（更底层的视角）
	s1 := []int{1, 2, 3}
	s2 := append(s1, 4) // 若 cap 不足，s2 指向新数组
	fmt.Printf("\ns1 底层指针: %v\n", unsafe.Pointer(&s1[0]))
	fmt.Printf("s2 底层指针: %v\n", unsafe.Pointer(&s2[0]))
	if unsafe.Pointer(&s1[0]) != unsafe.Pointer(&s2[0]) {
		fmt.Println("s1 和 s2 指向不同底层数组（append 触发了扩容）")
	} else {
		fmt.Println("s1 和 s2 共享底层数组（append 未触发扩容）")
	}
	fmt.Println("注意：扩容后旧切片仍指向旧数组，新切片指向新数组，两者独立")
}

// -----------------------------------------------------------------------------
// 6. 切片共享底层数组：修改子切片会影响原切片
// C 差异：
//   - C 的指针偏移也共享内存，但没有长度/容量保护，越界访问是 UB；
//   - Go 的切片有边界检查，越界访问触发 panic 而非 UB；
//   - 共享底层数组是 Go 切片的重要特性，既是性能优势也是潜在陷阱；
//   - 使用 copy 或 append([]T{}, src...) 可以创建独立副本，避免意外修改。
// -----------------------------------------------------------------------------
func demoSliceSharedBacking() {
	fmt.Println("\n--- 6. 切片共享底层数组 ---")

	original := []int{1, 2, 3, 4, 5}
	// sub 和 original 共享同一底层数组
	sub := original[1:4] // sub = [2, 3, 4]
	fmt.Printf("original: %v\n", original)
	fmt.Printf("sub = original[1:4]: %v\n", sub)

	// 修改 sub 会影响 original
	sub[0] = 999
	fmt.Printf("sub[0] = 999 后:\n")
	fmt.Printf("  sub:      %v\n", sub)
	fmt.Printf("  original: %v（original[1] 也被修改了！）\n", original)

	// 使用 copy 创建独立副本，避免共享
	original2 := []int{1, 2, 3, 4, 5}
	independent := make([]int, 3)
	copy(independent, original2[1:4]) // 复制到独立切片
	independent[0] = 888
	fmt.Printf("\n使用 copy 创建独立副本:\n")
	fmt.Printf("  independent: %v（修改了 [0]）\n", independent)
	fmt.Printf("  original2:   %v（未受影响）\n", original2)
	fmt.Println("注意：切片赋值（sub := original[1:4]）共享底层数组，用 copy 避免意外修改")
}

// -----------------------------------------------------------------------------
// 7. copy(dst, src) 函数用法及返回值
// C 差异：
//   - C 的 memcpy 不检查边界，需要调用方保证 dst 有足够空间；
//   - Go 的 copy 自动取 min(len(dst), len(src)) 个元素复制，不会越界；
//   - copy 返回实际复制的元素个数，可用于判断是否完整复制；
//   - copy 支持 src 和 dst 重叠（类似 C 的 memmove，而非 memcpy）。
// -----------------------------------------------------------------------------
func demoCopy() {
	fmt.Println("\n--- 7. copy 函数用法及返回值 ---")

	src := []int{1, 2, 3, 4, 5}

	// dst 比 src 短：只复制 dst 能容纳的部分
	dst1 := make([]int, 3)
	n1 := copy(dst1, src)
	fmt.Printf("copy(dst[3], src[5]): 复制了 %d 个元素, dst=%v\n", n1, dst1)

	// dst 比 src 长：复制全部 src 元素
	dst2 := make([]int, 7)
	n2 := copy(dst2, src)
	fmt.Printf("copy(dst[7], src[5]): 复制了 %d 个元素, dst=%v\n", n2, dst2)

	// dst 和 src 长度相同
	dst3 := make([]int, len(src))
	n3 := copy(dst3, src)
	fmt.Printf("copy(dst[5], src[5]): 复制了 %d 个元素, dst=%v\n", n3, dst3)

	// copy 也可以复制字符串到 []byte
	bs := make([]byte, 5)
	n4 := copy(bs, "Hello, Go!")
	fmt.Printf("copy([]byte, string): 复制了 %d 个字节, dst=%v (%q)\n", n4, bs, bs)

	fmt.Println("注意：copy 返回 min(len(dst), len(src))，不会越界（C 的 memcpy 无此保护）")
}

// -----------------------------------------------------------------------------
// 8. 二维切片的创建与访问
// C 差异：
//   - C 的二维数组 int a[3][4] 在内存中是连续的；
//   - Go 的二维切片是切片的切片（[][]T），每行可以有不同长度（锯齿数组）；
//   - 内存不一定连续，每行是独立分配的切片；
//   - 如需连续内存的二维数组，可以用一维切片模拟：a[i*cols+j]。
// -----------------------------------------------------------------------------
func demo2DSlice() {
	fmt.Println("\n--- 8. 二维切片的创建与访问 ---")

	rows, cols := 3, 4

	// 创建二维切片（每行独立分配）
	// C 等价：int **matrix = malloc(rows * sizeof(int*)); for each row: malloc(cols * sizeof(int))
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = i*cols + j // 填充值
		}
	}

	fmt.Printf("二维切片 %dx%d:\n", rows, cols)
	for i, row := range matrix {
		fmt.Printf("  matrix[%d]: %v\n", i, row)
	}

	// 锯齿数组（每行长度不同）
	// C 差异：C 的二维数组每行长度必须相同
	jagged := [][]int{
		{1},
		{2, 3},
		{4, 5, 6},
	}
	fmt.Println("锯齿数组（每行长度不同）:")
	for i, row := range jagged {
		fmt.Printf("  jagged[%d]: %v（len=%d）\n", i, row, len(row))
	}
	fmt.Println("注意：Go 二维切片每行可以有不同长度，内存不连续（C 二维数组内存连续）")
}

// -----------------------------------------------------------------------------
// 9. map 创建：make 和字面量
// C 差异：
//   - C 没有内置 map/哈希表，需要手动实现或使用第三方库；
//   - Go 的 map 是内置类型，键可以是任意可比较类型（comparable），
//     值可以是任意类型；
//   - map 是引用类型：赋值传递的是引用，修改会影响原 map；
//   - 未初始化的 map（nil map）可以读取（返回零值），但写入会 panic。
// -----------------------------------------------------------------------------
func demoMapBasics() {
	fmt.Println("\n--- 9. map 创建方式 ---")

	// 方式一：make(map[K]V)
	// C 等价：需要手动实现哈希表或使用 uthash 等库
	m1 := make(map[string]int)
	m1["alice"] = 90
	m1["bob"] = 85
	fmt.Printf("make(map[string]int): %v\n", m1)

	// 方式二：字面量初始化
	m2 := map[string]int{
		"alice": 90,
		"bob":   85,
		"carol": 92,
	}
	fmt.Printf("字面量 map: %v\n", m2)

	// map 是引用类型：m3 和 m2 共享同一底层数据
	m3 := m2
	m3["alice"] = 100
	fmt.Printf("m3[\"alice\"]=100 后, m2[\"alice\"]=%d（引用类型，共享数据）\n", m2["alice"])
	fmt.Println("注意：map 是引用类型，赋值不会复制数据（与数组的值语义相反）")
}

// -----------------------------------------------------------------------------
// 10. map 读写、删除（delete）、存在性检查
// C 差异：
//   - C 的哈希表查找通常返回指针（NULL 表示不存在），需要区分"不存在"和"值为零"；
//   - Go 的 map 查找支持双返回值 v, ok := m[k]，ok 明确表示键是否存在；
//   - 访问不存在的键返回值类型的零值（不 panic，不返回 NULL）；
//   - delete(m, k) 删除键值对，删除不存在的键是安全的（无操作）。
// -----------------------------------------------------------------------------
func demoMapReadWriteDelete() {
	fmt.Println("\n--- 10. map 读写、删除、存在性检查 ---")

	m := map[string]int{"alice": 90, "bob": 85, "carol": 92}

	// 读取
	fmt.Printf("m[\"alice\"] = %d\n", m["alice"])

	// 写入（新键）
	m["dave"] = 78
	fmt.Printf("写入 m[\"dave\"]=78 后: len=%d\n", len(m))

	// 更新（已有键）
	m["alice"] = 95
	fmt.Printf("更新 m[\"alice\"]=95 后: m[\"alice\"]=%d\n", m["alice"])

	// 删除
	// C 等价：哈希表删除操作，通常需要标记删除或重新哈希
	delete(m, "bob")
	fmt.Printf("delete(m, \"bob\") 后: len=%d\n", len(m))

	// 删除不存在的键是安全的（无操作，不 panic）
	delete(m, "nonexistent")
	fmt.Println("delete(m, \"nonexistent\"): 安全，无操作")

	// 存在性检查：双返回值形式
	// C 差异：C 需要通过返回值是否为 NULL 来判断，无法区分"不存在"和"值为零"
	if v, ok := m["alice"]; ok {
		fmt.Printf("m[\"alice\"] 存在，值=%d\n", v)
	}
	if v, ok := m["bob"]; !ok {
		fmt.Printf("m[\"bob\"] 不存在（已删除），零值=%d\n", v)
	}
	fmt.Println("注意：v, ok := m[k] 中 ok 明确区分「键不存在」和「值为零值」")
}

// -----------------------------------------------------------------------------
// 11. 遍历 map 时键的顺序是随机的，演示如何有序遍历
// C 差异：
//   - C 的哈希表遍历顺序也是不确定的；
//   - Go 故意在每次运行时随机化 map 遍历顺序，防止代码依赖特定顺序；
//   - 如需有序遍历，需要先提取所有键，排序后再按序访问。
// -----------------------------------------------------------------------------
func demoMapOrdering() {
	fmt.Println("\n--- 11. map 遍历顺序随机性与有序遍历 ---")

	m := map[string]int{
		"banana": 3,
		"apple":  5,
		"cherry": 1,
		"date":   2,
	}

	// 直接遍历：顺序随机（每次运行可能不同）
	fmt.Println("直接遍历（顺序随机）:")
	for k, v := range m {
		fmt.Printf("  %s: %d\n", k, v)
	}

	// 有序遍历：先提取键，排序，再按序访问
	// C 差异：C 也需要类似操作，但需要手动实现排序
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 对键排序

	fmt.Println("有序遍历（排序 key 后）:")
	for _, k := range keys {
		fmt.Printf("  %s: %d\n", k, m[k])
	}
	fmt.Println("注意：Go 故意随机化 map 遍历顺序，防止代码依赖特定顺序（这是设计决策）")
}

// -----------------------------------------------------------------------------
// 12. 访问 map 中不存在的键返回零值（不 panic）
// C 差异：
//   - C 的哈希表访问不存在的键通常返回 NULL 或特殊错误码；
//   - Go 访问不存在的键返回值类型的零值（int→0, string→"", bool→false 等），
//     不会 panic，也不会报错；
//   - 这可能导致隐蔽的 bug：无法区分"键不存在"和"值恰好是零值"，
//     需要用双返回值形式 v, ok := m[k] 来明确区分。
// -----------------------------------------------------------------------------
func demoMapZeroValue() {
	fmt.Println("\n--- 12. 访问不存在的键返回零值（不 panic）---")

	scores := map[string]int{"alice": 90}

	// 访问存在的键
	fmt.Printf("scores[\"alice\"] = %d（存在）\n", scores["alice"])

	// 访问不存在的键：返回零值，不 panic
	// C 差异：C 的哈希表通常返回 NULL，需要判断
	fmt.Printf("scores[\"bob\"] = %d（不存在，返回 int 零值 0）\n", scores["bob"])

	// 零值的实际用途：计数器 map（无需初始化）
	wordCount := make(map[string]int)
	words := []string{"go", "is", "go", "fun", "go"}
	for _, w := range words {
		wordCount[w]++ // 不存在的键自动从零值 0 开始累加
	}
	fmt.Printf("词频统计（利用零值）: %v\n", wordCount)

	// 零值的陷阱：无法区分"不存在"和"值为 0"
	m := map[string]int{"zero": 0, "one": 1}
	fmt.Printf("m[\"zero\"]=%d, m[\"missing\"]=%d（无法区分！）\n", m["zero"], m["missing"])
	fmt.Println("解决方案：用 v, ok := m[k] 明确区分键是否存在")

	// 其他类型的零值
	boolMap := map[string]bool{}
	strMap := map[string]string{}
	fmt.Printf("bool map 零值: %v, string map 零值: %q\n", boolMap["x"], strMap["x"])
}

// -----------------------------------------------------------------------------
// 13. 向未初始化（nil）的 map 写入数据触发 panic
// C 差异：
//   - C 向 NULL 指针写入是未定义行为（UB），可能导致段错误（segfault）；
//   - Go 向 nil map 写入会触发 panic（运行时错误），有明确的错误信息；
//   - nil map 可以安全读取（返回零值），但不能写入；
//   - 此 panic 由 main.go 的 Runner（defer+recover）捕获，程序不会崩溃。
// -----------------------------------------------------------------------------
func demoNilMapPanic() {
	fmt.Println("\n--- 13. nil map 写入触发 panic（由 Runner recover 捕获）---")

	// nil map 读取是安全的（返回零值）
	var nilMap map[string]int // 未初始化，值为 nil
	fmt.Printf("nil map 读取: nilMap[\"key\"] = %d（安全，返回零值）\n", nilMap["key"])
	fmt.Printf("nil map 判断: nilMap == nil → %v\n", nilMap == nil)

	fmt.Println("即将向 nil map 写入数据，触发 panic...")
	fmt.Println("（此 panic 将由 main.go 的 runSafe 函数通过 defer+recover 捕获）")

	// 向 nil map 写入：触发 panic: assignment to entry in nil map
	// C 差异：C 的 NULL 指针写入是 UB，Go 的 nil map 写入是明确的 panic
	nilMap["key"] = 1 // panic: assignment to entry in nil map
}
