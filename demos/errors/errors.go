// Package errors 演示 Go 语言的错误处理机制。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// Go 的错误处理哲学：错误是值（error is a value），通过返回值传递，
// 而非异常机制（try/catch）。这与 C 的 errno 全局变量类似，但更安全、更灵活。
package errors

import (
	"errors"
	"fmt"
)

// -----------------------------------------------------------------------------
// 导出的哨兵错误（sentinel errors）
// C 差异：
//   - C 用 errno 全局变量（ENOENT、EACCES 等）表示预定义错误；
//   - Go 用包级别的导出变量作为哨兵错误，调用方用 errors.Is 检查；
//   - 哨兵错误是不可变的值，用 == 或 errors.Is 比较（支持错误链）；
//   - 命名约定：以 Err 开头（如 ErrNotFound、io.EOF）。
// -----------------------------------------------------------------------------

// ErrNotFound 表示资源未找到的哨兵错误。
// C 差异：类似 C 的 ENOENT（errno.h），但 Go 是类型安全的包级变量。
var ErrNotFound = errors.New("not found")

// ErrDivisionByZero 表示除以零的哨兵错误。
// C 差异：C 中除以零是未定义行为（UB），Go 通过返回 error 明确处理。
var ErrDivisionByZero = errors.New("division by zero")

// -----------------------------------------------------------------------------
// 自定义错误类型
// C 差异：
//   - C 的 errno 只是一个整数，无法携带额外上下文信息；
//   - Go 可以定义实现 error 接口的结构体，携带任意上下文（字段、方法）；
//   - error 接口定义：type error interface { Error() string }
//   - 自定义错误类型可以通过 errors.As 从错误链中提取，获取详细信息。
// -----------------------------------------------------------------------------

// ValidationError 是自定义错误类型，携带字段名和错误消息。
// C 差异：C 无法在 errno 中携带字段名等上下文，Go 的自定义错误类型可以。
type ValidationError struct {
	Field   string // 验证失败的字段名
	Message string // 错误描述
}

// Error 实现 error 接口（只需实现 Error() string 方法）。
// C 差异：C 用 strerror(errno) 获取错误字符串，Go 通过接口方法直接调用。
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %q - %s", e.Field, e.Message)
}

// -----------------------------------------------------------------------------
// 导出函数
// -----------------------------------------------------------------------------

// Divide 执行除法，除数为零时返回 ErrDivisionByZero。
// C 差异：
//   - C 中整数除以零是未定义行为（UB），浮点除以零得到 +Inf/-Inf/NaN；
//   - Go 惯用模式：返回 (result, error)，调用方必须检查 error；
//   - 这比 C 的 errno 更安全：error 是返回值，不会被忽略（编译器警告）。
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// -----------------------------------------------------------------------------
// Demo 函数：演示所有错误处理相关内容
// -----------------------------------------------------------------------------

// Demo 演示 Go 语言的错误处理机制。
func Demo() {
	demoErrorInterface()
	demoReturnErrorPattern()
	demoErrorsNew()
	demoErrorWrapping()
	demoCustomErrorType()
	demoErrorsIs()
	demoErrorsAs()
	demoPanic()
	demoRecover()
	demoSentinelErrors()
}

// -----------------------------------------------------------------------------
// 1. error 接口定义
// C 差异：
//   - C 没有接口概念，错误通过 errno 全局变量或返回负整数表示；
//   - Go 的 error 是内置接口：type error interface { Error() string }
//   - 任何实现了 Error() string 方法的类型都自动满足 error 接口；
//   - nil error 表示"无错误"，与 C 的返回 0 表示成功类似。
// -----------------------------------------------------------------------------
func demoErrorInterface() {
	fmt.Println("\n--- 1. error 接口定义 ---")

	// error 接口定义（内置于 Go 语言）：
	//   type error interface {
	//       Error() string
	//   }
	//
	// C 差异：
	//   C 用 errno（全局整数）+ strerror() 获取错误字符串：
	//     errno = ENOENT;
	//     printf("%s\n", strerror(errno)); // "No such file or directory"
	//   Go 用 error 接口，类型安全，可携带任意上下文。

	fmt.Println("error 接口定义：")
	fmt.Println("  type error interface {")
	fmt.Println("      Error() string")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("任何实现了 Error() string 方法的类型都自动满足 error 接口")
	fmt.Println("nil error 表示无错误（类似 C 函数返回 0 表示成功）")
	fmt.Println("C 差异：Go 用接口而非 errno 全局变量，类型安全，可携带上下文")

	// 演示 error 接口变量
	var err error // 零值为 nil
	fmt.Printf("\nerror 零值: err == nil → %v\n", err == nil)

	err = errors.New("something went wrong")
	fmt.Printf("赋值后: err == nil → %v, err.Error() = %q\n", err == nil, err.Error())
}

// -----------------------------------------------------------------------------
// 2. 函数返回 (result, error) 的惯用模式
// C 差异：
//   - C 通常用返回值表示成功/失败（0/-1），用 errno 或输出参数传递错误；
//   - Go 惯用模式：返回 (result, error)，调用方立即检查 error；
//   - Go 没有异常（try/catch），错误处理是显式的、强制的；
//   - 标准写法：if err != nil { return ..., err }（错误向上传播）。
// -----------------------------------------------------------------------------
func demoReturnErrorPattern() {
	fmt.Println("\n--- 2. 函数返回 (result, error) 的惯用模式 ---")

	// 标准调用模式：立即检查 error
	// C 差异：
	//   C 写法：int result = divide(10, 2); if (result < 0) { /* 错误 */ }
	//   Go 写法：result, err := Divide(10, 2); if err != nil { /* 错误 */ }
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Printf("Divide(10, 2) 错误: %v\n", err)
	} else {
		fmt.Printf("Divide(10, 2) = %.2f\n", result)
	}

	// 除以零的情况
	result, err = Divide(5, 0)
	if err != nil {
		fmt.Printf("Divide(5, 0) 错误: %v\n", err)
	} else {
		fmt.Printf("Divide(5, 0) = %.2f\n", result)
	}

	fmt.Println()
	fmt.Println("惯用模式总结：")
	fmt.Println("  result, err := someFunc()")
	fmt.Println("  if err != nil {")
	fmt.Println("      // 处理错误：记录日志、返回、包装后返回")
	fmt.Println("      return ..., fmt.Errorf(\"context: %w\", err)")
	fmt.Println("  }")
	fmt.Println("  // 使用 result")
	fmt.Println("C 差异：Go 无异常机制，错误处理是显式的返回值检查")
}

// -----------------------------------------------------------------------------
// 3. errors.New 创建简单错误
// C 差异：
//   - C 用 errno 预定义常量（ENOENT=2、EACCES=13 等）表示错误；
//   - Go 用 errors.New("message") 创建不可变的错误值；
//   - errors.New 返回的错误每次调用都是不同的指针（即使消息相同）；
//   - 包级别的哨兵错误（var ErrXxx = errors.New(...)）是单例，可用 == 比较。
// -----------------------------------------------------------------------------
func demoErrorsNew() {
	fmt.Println("\n--- 3. errors.New 创建简单错误 ---")

	// 创建简单错误
	// C 差异：C 用 errno = ENOENT，Go 用 errors.New("not found")
	err1 := errors.New("file not found")
	err2 := errors.New("file not found") // 消息相同，但是不同的指针

	fmt.Printf("err1 = %v\n", err1)
	fmt.Printf("err2 = %v\n", err2)
	fmt.Printf("err1 == err2 → %v（不同指针，即使消息相同）\n", err1 == err2)

	// 包级别哨兵错误是单例
	fmt.Printf("\nErrNotFound = %v\n", ErrNotFound)
	fmt.Printf("ErrNotFound == ErrNotFound → %v（同一指针）\n", ErrNotFound == ErrNotFound)

	// 演示哨兵错误的用法
	findUser := func(id int) error {
		if id <= 0 {
			return ErrNotFound // 返回哨兵错误
		}
		return nil
	}

	if err := findUser(0); err != nil {
		fmt.Printf("\nfindUser(0) 错误: %v\n", err)
		fmt.Printf("err == ErrNotFound → %v\n", err == ErrNotFound)
	}

	fmt.Println("C 差异：Go 的 errors.New 比 C 的 errno 更灵活，可携带任意消息")
}

// -----------------------------------------------------------------------------
// 4. fmt.Errorf("%w", err) 包装错误（error wrapping）
// C 差异：
//   - C 没有错误链的概念，errno 被覆盖后原始错误丢失；
//   - Go 1.13+ 引入 %w 动词，将原始错误包装进新错误，形成错误链；
//   - 包装后可用 errors.Is/errors.As 遍历整个错误链；
//   - 这类似于异常链（exception chaining），但通过值而非异常实现。
// -----------------------------------------------------------------------------
func demoErrorWrapping() {
	fmt.Println("\n--- 4. fmt.Errorf(\"%w\", err) 包装错误（error wrapping）---")

	// 模拟多层调用中的错误包装
	// C 差异：C 中每层函数通常直接返回 errno，丢失调用栈上下文
	//         Go 用 fmt.Errorf("context: %w", err) 保留原始错误并添加上下文

	// 底层操作
	dbQuery := func(id int) error {
		if id == 0 {
			return ErrNotFound // 原始哨兵错误
		}
		return nil
	}

	// 中间层：包装错误，添加上下文
	getUser := func(id int) error {
		if err := dbQuery(id); err != nil {
			// %w 包装原始错误，保留错误链
			return fmt.Errorf("getUser(id=%d): %w", id, err)
		}
		return nil
	}

	// 上层：再次包装
	handleRequest := func(id int) error {
		if err := getUser(id); err != nil {
			return fmt.Errorf("handleRequest: %w", err)
		}
		return nil
	}

	err := handleRequest(0)
	fmt.Printf("包装后的错误: %v\n", err)
	fmt.Println("错误链：handleRequest → getUser → ErrNotFound")

	// %v 和 %w 的区别
	original := ErrNotFound
	wrapped := fmt.Errorf("layer1: %w", original)
	notWrapped := fmt.Errorf("layer1: %v", original) // %v 不包装，只格式化字符串

	fmt.Printf("\n%%w 包装: %v\n", wrapped)
	fmt.Printf("%%v 不包装: %v\n", notWrapped)
	fmt.Printf("errors.Is(wrapped, ErrNotFound)    → %v（%%w 保留错误链）\n",
		errors.Is(wrapped, ErrNotFound))
	fmt.Printf("errors.Is(notWrapped, ErrNotFound) → %v（%%v 不保留错误链）\n",
		errors.Is(notWrapped, ErrNotFound))

	fmt.Println("C 差异：C 的 errno 被覆盖后原始错误丢失，Go 的 %w 保留完整错误链")
}

// -----------------------------------------------------------------------------
// 5. 自定义错误类型 ValidationError
// C 差异：
//   - C 的 errno 只是整数，无法携带字段名、消息等上下文信息；
//   - Go 可以定义结构体实现 error 接口，携带任意上下文；
//   - 自定义错误类型可以添加额外方法，提供更丰富的错误信息；
//   - 调用方可以用 errors.As 提取具体错误类型，访问额外字段。
// -----------------------------------------------------------------------------
func demoCustomErrorType() {
	fmt.Println("\n--- 5. 自定义错误类型 ValidationError ---")

	// ValidationError 实现了 error 接口
	// C 差异：C 无法在 errno 中携带字段名，Go 的自定义错误类型可以
	validate := func(name string, age int) error {
		if name == "" {
			return &ValidationError{
				Field:   "name",
				Message: "cannot be empty",
			}
		}
		if age < 0 || age > 150 {
			return &ValidationError{
				Field:   "age",
				Message: fmt.Sprintf("must be between 0 and 150, got %d", age),
			}
		}
		return nil
	}

	// 测试空名称
	if err := validate("", 25); err != nil {
		fmt.Printf("validate(\"\", 25) 错误: %v\n", err)
	}

	// 测试无效年龄
	if err := validate("Alice", -1); err != nil {
		fmt.Printf("validate(\"Alice\", -1) 错误: %v\n", err)
		// 类型断言获取具体错误类型
		if ve, ok := err.(*ValidationError); ok {
			fmt.Printf("  字段: %q, 消息: %q\n", ve.Field, ve.Message)
		}
	}

	// 测试有效输入
	if err := validate("Bob", 30); err == nil {
		fmt.Println("validate(\"Bob\", 30) 成功（无错误）")
	}

	fmt.Println("C 差异：Go 自定义错误类型可携带任意上下文，C 的 errno 只是整数")
}

// -----------------------------------------------------------------------------
// 6. errors.Is 检查错误链中是否包含特定哨兵错误
// C 差异：
//   - C 用 errno == ENOENT 直接比较，但 errno 可能被中间调用覆盖；
//   - Go 的 errors.Is 遍历整个错误链（通过 Unwrap() 方法），即使错误被包装也能找到；
//   - errors.Is 等价于：err == target || errors.Is(err.Unwrap(), target)（递归）；
//   - 对于哨兵错误，应始终用 errors.Is 而非 == 比较（支持包装场景）。
// -----------------------------------------------------------------------------
func demoErrorsIs() {
	fmt.Println("\n--- 6. errors.Is 检查错误链中是否包含特定哨兵错误 ---")

	// 直接比较（未包装）
	err1 := ErrNotFound
	fmt.Printf("errors.Is(ErrNotFound, ErrNotFound) → %v\n",
		errors.Is(err1, ErrNotFound))

	// 包装后仍能检测到（遍历错误链）
	// C 差异：C 的 errno == ENOENT 在错误被包装后无法检测，Go 的 errors.Is 可以
	wrapped1 := fmt.Errorf("layer1: %w", ErrNotFound)
	wrapped2 := fmt.Errorf("layer2: %w", wrapped1)

	fmt.Printf("\n单层包装: errors.Is(wrapped1, ErrNotFound) → %v\n",
		errors.Is(wrapped1, ErrNotFound))
	fmt.Printf("双层包装: errors.Is(wrapped2, ErrNotFound) → %v\n",
		errors.Is(wrapped2, ErrNotFound))

	// 不同哨兵错误
	fmt.Printf("\nerrors.Is(ErrNotFound, ErrDivisionByZero) → %v（不同错误）\n",
		errors.Is(ErrNotFound, ErrDivisionByZero))

	// 实际使用场景
	_, err := Divide(10, 0)
	wrappedDivErr := fmt.Errorf("calculate: %w", err)
	fmt.Printf("\nDivide(10,0) 包装后 errors.Is(err, ErrDivisionByZero) → %v\n",
		errors.Is(wrappedDivErr, ErrDivisionByZero))

	fmt.Println("C 差异：errors.Is 遍历错误链，C 的 errno == X 只检查当前值")
}

// -----------------------------------------------------------------------------
// 7. errors.As 从错误链中提取特定类型的错误
// C 差异：
//   - C 没有从错误链中提取特定类型的机制；
//   - Go 的 errors.As 遍历错误链，找到第一个可赋值给目标类型的错误；
//   - errors.As 等价于：类型断言 + 递归 Unwrap()；
//   - 用于提取自定义错误类型的额外字段（如 ValidationError.Field）。
// -----------------------------------------------------------------------------
func demoErrorsAs() {
	fmt.Println("\n--- 7. errors.As 从错误链中提取特定类型的错误 ---")

	// 创建包含 ValidationError 的错误链
	originalErr := &ValidationError{Field: "email", Message: "invalid format"}
	wrappedErr := fmt.Errorf("user registration: %w", originalErr)
	doubleWrapped := fmt.Errorf("api handler: %w", wrappedErr)

	fmt.Printf("错误链: %v\n", doubleWrapped)

	// errors.As 从错误链中提取 *ValidationError
	// C 差异：C 没有等价机制，需要手动解析错误信息字符串
	var ve *ValidationError
	if errors.As(doubleWrapped, &ve) {
		fmt.Printf("\nerrors.As 成功提取 *ValidationError:\n")
		fmt.Printf("  Field:   %q\n", ve.Field)
		fmt.Printf("  Message: %q\n", ve.Message)
	}

	// errors.As 与类型断言的区别
	// 类型断言只检查当前层，errors.As 遍历整个错误链
	_, ok := doubleWrapped.(*ValidationError) // 类型断言：失败（当前层是 *fmt.wrapError）
	fmt.Printf("\n直接类型断言 doubleWrapped.(*ValidationError) → ok=%v（只检查当前层）\n", ok)
	fmt.Printf("errors.As(doubleWrapped, &ve) → %v（遍历整个错误链）\n",
		errors.As(doubleWrapped, &ve))

	// errors.As 找不到目标类型时返回 false
	var divErr *ValidationError
	_, divideErr := Divide(1, 0)
	fmt.Printf("\nerrors.As(ErrDivisionByZero, &ValidationError) → %v（类型不匹配）\n",
		errors.As(divideErr, &divErr))

	fmt.Println("C 差异：errors.As 遍历错误链提取具体类型，C 没有等价机制")
}

// -----------------------------------------------------------------------------
// 8. panic 的触发场景
// C 差异：
//   - C 中数组越界、空指针解引用等是未定义行为（UB），可能静默损坏内存；
//   - Go 的 panic 在运行时检测到非法操作时触发，立即停止当前 goroutine；
//   - panic 会打印堆栈跟踪，比 C 的 segfault 更容易调试；
//   - panic 应用于程序员错误（不可恢复的逻辑错误），不用于正常错误处理。
// -----------------------------------------------------------------------------
func demoPanic() {
	fmt.Println("\n--- 8. panic 的触发场景 ---")

	fmt.Println("panic 的两种触发方式：")
	fmt.Println("  1. 运行时自动触发（程序员错误）：索引越界、nil 指针解引用、类型断言失败等")
	fmt.Println("  2. 主动调用 panic(value)：表示程序遇到不可恢复的错误")

	// 演示主动 panic（由 defer+recover 捕获）
	fmt.Println("\n演示主动 panic（将被 recover 捕获）：")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到 panic: %v\n", r)
			}
		}()
		panic("something went terribly wrong") // 主动触发 panic
	}()

	// 演示索引越界 panic（运行时自动触发）
	fmt.Println("\n演示索引越界 panic（将被 recover 捕获）：")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到 panic: %v\n", r)
			}
		}()
		s := []int{1, 2, 3}
		_ = s[10] // 索引越界，运行时 panic
	}()

	// 演示 nil map 写入 panic
	fmt.Println("\n演示 nil map 写入 panic（将被 recover 捕获）：")
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("  捕获到 panic: %v\n", r)
			}
		}()
		var m map[string]int // nil map
		m["key"] = 1         // 写入 nil map，panic
	}()

	fmt.Println("\npanic 使用原则：")
	fmt.Println("  - 用于程序员错误（不应该发生的情况）：如参数校验失败、不变量被破坏")
	fmt.Println("  - 不用于正常错误处理（用 error 返回值代替）")
	fmt.Println("  - 库代码应避免 panic 传播到调用方（用 recover 转换为 error）")
	fmt.Println("C 差异：C 的 UB 静默损坏内存，Go 的 panic 立即停止并打印堆栈")
}

// -----------------------------------------------------------------------------
// 9. recover 在 defer 中捕获 panic，转换为 error（SafeRun 模式）
// C 差异：
//   - C 没有 panic/recover 机制，只能用 setjmp/longjmp 模拟（复杂且危险）；
//   - Go 的 recover 只能在 defer 函数中调用，捕获当前 goroutine 的 panic；
//   - SafeRun 模式：将可能 panic 的函数包装为返回 error 的安全函数；
//   - 这是库代码的常见模式：内部用 panic 表示逻辑错误，对外暴露 error 接口。
// -----------------------------------------------------------------------------
func demoRecover() {
	fmt.Println("\n--- 9. recover 在 defer 中捕获 panic，转换为 error（SafeRun 模式）---")

	// SafeRun 将可能 panic 的函数包装为返回 error 的安全函数
	// C 差异：C 需要用 setjmp/longjmp 实现类似功能，复杂且容易出错
	safeRun := func(fn func()) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// 将 panic 值转换为 error
				switch v := r.(type) {
				case error:
					err = fmt.Errorf("recovered panic: %w", v)
				default:
					err = fmt.Errorf("recovered panic: %v", v)
				}
			}
		}()
		fn() // 执行可能 panic 的函数
		return nil
	}

	// 测试 SafeRun：正常函数
	err := safeRun(func() {
		fmt.Println("  正常函数执行成功")
	})
	fmt.Printf("safeRun(正常函数) → err=%v\n", err)

	// 测试 SafeRun：panic 函数
	err = safeRun(func() {
		panic("unexpected state")
	})
	fmt.Printf("safeRun(panic函数) → err=%v\n", err)

	// 测试 SafeRun：索引越界
	err = safeRun(func() {
		s := []int{1, 2, 3}
		_ = s[100]
	})
	fmt.Printf("safeRun(索引越界) → err=%v\n", err)

	// recover 的关键规则
	fmt.Println("\nrecover 的关键规则：")
	fmt.Println("  1. recover 只能在 defer 函数中调用（否则返回 nil）")
	fmt.Println("  2. recover 只捕获当前 goroutine 的 panic（不能跨 goroutine）")
	fmt.Println("  3. recover 后程序继续执行（从 defer 返回后）")
	fmt.Println("  4. 捕获 panic 后应记录日志或转换为 error，不要静默忽略")
	fmt.Println("C 差异：C 的 setjmp/longjmp 复杂危险，Go 的 defer+recover 简洁安全")
}

// -----------------------------------------------------------------------------
// 10. 哨兵错误（sentinel error）模式
// C 差异：
//   - C 的 errno 是全局整数常量（ENOENT=2、EACCES=13 等），所有库共享同一命名空间；
//   - Go 的哨兵错误是包级别的导出变量，每个包有自己的命名空间（如 io.EOF、sql.ErrNoRows）；
//   - 哨兵错误是不可变的值，用 errors.Is 比较（支持错误链）；
//   - 命名约定：以 Err 开头（ErrNotFound、ErrDivisionByZero）。
// -----------------------------------------------------------------------------
func demoSentinelErrors() {
	fmt.Println("\n--- 10. 哨兵错误（sentinel error）模式 ---")

	// 演示 ErrNotFound 哨兵错误
	// C 差异：类似 C 的 ENOENT，但 Go 是包级别变量，不共享全局命名空间
	findItem := func(id int) (string, error) {
		items := map[int]string{1: "apple", 2: "banana", 3: "cherry"}
		if item, ok := items[id]; ok {
			return item, nil
		}
		return "", ErrNotFound // 返回哨兵错误
	}

	// 演示 ErrDivisionByZero 哨兵错误
	fmt.Println("【ErrNotFound 哨兵错误】")
	for _, id := range []int{1, 5, 2} {
		item, err := findItem(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				fmt.Printf("  findItem(%d) → ErrNotFound（资源不存在）\n", id)
			} else {
				fmt.Printf("  findItem(%d) → 未知错误: %v\n", id, err)
			}
		} else {
			fmt.Printf("  findItem(%d) → %q\n", id, item)
		}
	}

	fmt.Println("\n【ErrDivisionByZero 哨兵错误】")
	for _, pair := range [][2]float64{{10, 2}, {5, 0}, {9, 3}} {
		result, err := Divide(pair[0], pair[1])
		if err != nil {
			if errors.Is(err, ErrDivisionByZero) {
				fmt.Printf("  Divide(%.0f, %.0f) → ErrDivisionByZero\n", pair[0], pair[1])
			}
		} else {
			fmt.Printf("  Divide(%.0f, %.0f) = %.2f\n", pair[0], pair[1], result)
		}
	}

	// 标准库哨兵错误示例
	fmt.Println("\n标准库常见哨兵错误：")
	fmt.Println("  io.EOF              — 读取到文件末尾")
	fmt.Println("  io.ErrUnexpectedEOF — 意外的文件末尾")
	fmt.Println("  os.ErrNotExist      — 文件不存在（类似 C 的 ENOENT）")
	fmt.Println("  os.ErrPermission    — 权限不足（类似 C 的 EACCES）")
	fmt.Println("  sql.ErrNoRows       — 查询无结果")
	fmt.Println("  context.Canceled    — 上下文被取消")
	fmt.Println("  context.DeadlineExceeded — 超时")

	fmt.Println("\n哨兵错误最佳实践：")
	fmt.Println("  1. 用 errors.Is 比较（支持错误链），不用 ==")
	fmt.Println("  2. 命名以 Err 开头（ErrNotFound、ErrTimeout）")
	fmt.Println("  3. 只为调用方需要区分处理的错误定义哨兵错误")
	fmt.Println("  4. 避免过多哨兵错误（优先用自定义错误类型携带上下文）")
	fmt.Println("C 差异：Go 哨兵错误是包级变量（有命名空间），C 的 errno 是全局整数常量")
}
