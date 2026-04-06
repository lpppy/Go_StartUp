// Package concurrency 演示 Go 语言的并发模型。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// Go 并发哲学："不要通过共享内存来通信，而要通过通信来共享内存。"
// （"Do not communicate by sharing memory; instead, share memory by communicating."）
//
// C 并发对比：
//   - C 用 pthread_create/pthread_join 创建线程，线程栈默认约 8MB；
//   - Go 的 goroutine 初始栈约 2KB，可动态增长，由 Go 运行时调度（M:N 模型）；
//   - C 用 mutex/semaphore/condition variable 同步，容易死锁；
//   - Go 提供 channel 作为首选同步原语，同时也支持 sync.Mutex 等传统方式。
package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MergeChannels 将多个只读 channel 合并为一个 channel（扇入模式）。
// 泛型参数 T 支持任意类型。
//
// 实现原理：为每个输入 channel 启动一个 goroutine，将数据转发到输出 channel；
// 使用 sync.WaitGroup 等待所有转发 goroutine 完成后关闭输出 channel。
//
// C 差异：C 没有原生的 channel 概念，需要用 pipe/socket/共享内存+mutex 实现类似功能。
func MergeChannels[T any](cs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup

	// 为每个输入 channel 启动一个转发 goroutine
	for _, c := range cs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch { // range 会在 channel 关闭后退出
				out <- v
			}
		}(c)
	}

	// 等待所有转发 goroutine 完成后关闭输出 channel
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// Demo 演示 Go 语言的并发模型，涵盖 goroutine、channel、select、sync 包和 context。
func Demo() {
	demoGoroutine()
	demoWaitGroup()
	demoUnbufferedChannel()
	demoBufferedChannel()
	demoSelect()
	demoChannelCloseAndRange()
	demoDirectionalChannels()
	demoMutex()
	demoSyncOnce()
	demoContext()
	demoGoroutineLeak()
}

// -----------------------------------------------------------------------------
// 1. go 关键字启动 goroutine
// C 差异：
//   - C 用 pthread_create(&tid, NULL, func, arg) 创建线程，线程栈默认约 8MB；
//   - Go 用 go func() 启动 goroutine，初始栈约 2KB，可动态增长到 1GB；
//   - Go 运行时使用 M:N 调度模型（M 个 goroutine 映射到 N 个 OS 线程）；
//   - 可以轻松启动数十万个 goroutine，而 C 线程受系统资源限制（通常几千个）；
//   - goroutine 没有 ID，不能直接等待某个 goroutine（需要 channel 或 WaitGroup）。
// -----------------------------------------------------------------------------
func demoGoroutine() {
	fmt.Println("\n--- 1. go 关键字启动 goroutine ---")
	fmt.Println("goroutine 轻量级特性：初始栈约 2KB（C 线程默认约 8MB）")
	fmt.Println("Go 运行时使用 M:N 调度：多个 goroutine 复用少量 OS 线程")

	// 用 channel 同步，等待 goroutine 完成
	// C 差异：C 用 pthread_join(tid, NULL) 等待线程，Go 用 channel 或 WaitGroup
	done := make(chan struct{})

	go func() {
		// 这是一个 goroutine，与主 goroutine 并发执行
		// C 差异：C 线程函数签名固定为 void* func(void*)，Go 可以是任意无参函数
		fmt.Println("  [goroutine] 我在独立的 goroutine 中运行")
		fmt.Println("  [goroutine] goroutine 初始栈约 2KB，按需动态增长")
		done <- struct{}{} // 通知主 goroutine 已完成
	}()

	<-done // 等待 goroutine 完成（阻塞直到收到信号）
	fmt.Println("  [main] goroutine 已完成")
	fmt.Println("C 差异：go 关键字比 pthread_create 简洁，goroutine 比线程轻量 4000 倍")
}

// -----------------------------------------------------------------------------
// 2. sync.WaitGroup 等待多个 goroutine 完成
// C 差异：
//   - C 用 pthread_join 逐个等待线程，需要保存所有 tid；
//   - Go 的 WaitGroup 可以等待任意数量的 goroutine，无需保存 goroutine 标识；
//   - WaitGroup 三个方法：Add(n) 增加计数、Done() 减少计数（等价 Add(-1)）、Wait() 阻塞直到计数为 0；
//   - 必须在启动 goroutine 之前调用 Add，避免 Wait 在 goroutine 启动前返回。
// -----------------------------------------------------------------------------
func demoWaitGroup() {
	fmt.Println("\n--- 2. sync.WaitGroup 等待多个 goroutine 完成 ---")

	var wg sync.WaitGroup
	results := make([]int, 5) // 预分配结果切片（每个 goroutine 写不同索引，无 data race）

	for i := 0; i < 5; i++ {
		wg.Add(1) // 必须在 go 语句之前调用 Add
		go func(idx int) {
			defer wg.Done() // goroutine 完成时减少计数（defer 确保一定执行）
			results[idx] = idx * idx
			fmt.Printf("  [goroutine %d] 计算 %d² = %d\n", idx, idx, results[idx])
		}(i) // 传入 i 的副本，避免闭包捕获循环变量（C 差异：Go 闭包捕获变量引用）
	}

	wg.Wait() // 阻塞直到所有 goroutine 调用 Done（计数归零）
	fmt.Printf("  [main] 所有 goroutine 完成，结果: %v\n", results)
	fmt.Println("C 差异：WaitGroup 比 pthread_join 更灵活，无需保存线程 ID")
}

// -----------------------------------------------------------------------------
// 3. 无缓冲 channel：发送方和接收方必须同时就绪（同步）
// C 差异：
//   - C 没有原生 channel，需要用 pipe(2) 或 POSIX 消息队列实现；
//   - 无缓冲 channel 类似 pipe 容量为 0：发送方阻塞直到接收方就绪（反之亦然）；
//   - 这种"握手"机制保证了同步：发送完成意味着接收方已经收到数据；
//   - C 的 pipe 是字节流，Go 的 channel 是类型安全的值传递。
// -----------------------------------------------------------------------------
func demoUnbufferedChannel() {
	fmt.Println("\n--- 3. 无缓冲 channel：同步通信 ---")
	fmt.Println("make(chan T) 创建无缓冲 channel，发送和接收必须同时就绪")

	// 无缓冲 channel：发送方阻塞直到接收方就绪
	ch := make(chan string) // 无缓冲 channel，容量为 0

	go func() {
		// 发送方：发送后阻塞，直到接收方接收
		fmt.Println("  [sender] 准备发送消息...")
		ch <- "Hello from goroutine" // 阻塞直到接收方就绪
		fmt.Println("  [sender] 消息已被接收，继续执行")
	}()

	// 接收方：接收前阻塞，直到发送方发送
	time.Sleep(10 * time.Millisecond) // 模拟接收方稍后就绪
	msg := <-ch                       // 阻塞直到发送方发送
	fmt.Printf("  [main] 收到消息: %q\n", msg)
	fmt.Println("  无缓冲 channel 保证：发送完成 = 接收方已收到（同步握手）")
	fmt.Println("C 差异：类似 pipe 但类型安全，无缓冲 channel 提供同步保证")
}

// -----------------------------------------------------------------------------
// 4. 有缓冲 channel：缓冲满时发送阻塞
// C 差异：
//   - 有缓冲 channel 类似有容量的消息队列（POSIX mq_open 或 ring buffer）；
//   - 发送方在缓冲未满时不阻塞（异步），缓冲满时阻塞（背压机制）；
//   - 接收方在缓冲非空时不阻塞，缓冲空时阻塞；
//   - 有缓冲 channel 可以解耦生产者和消费者的速度差异。
// -----------------------------------------------------------------------------
func demoBufferedChannel() {
	fmt.Println("\n--- 4. 有缓冲 channel：异步通信 ---")
	fmt.Println("make(chan T, n) 创建容量为 n 的有缓冲 channel")

	// 有缓冲 channel：缓冲未满时发送不阻塞
	ch := make(chan int, 3) // 容量为 3 的有缓冲 channel

	// 发送 3 个值（不阻塞，因为缓冲未满）
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("  发送 3 个值后，len=%d, cap=%d\n", len(ch), cap(ch))

	// 第 4 次发送会阻塞（缓冲已满），用 goroutine 演示
	go func() {
		fmt.Println("  [goroutine] 尝试发送第 4 个值（缓冲已满，将阻塞）...")
		ch <- 4 // 阻塞直到有空间
		fmt.Println("  [goroutine] 第 4 个值已发送")
	}()

	// 接收数据，腾出缓冲空间
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 4; i++ {
		v := <-ch
		fmt.Printf("  接收: %d\n", v)
	}

	fmt.Println("C 差异：有缓冲 channel 类似消息队列，提供生产者-消费者解耦")
}

// -----------------------------------------------------------------------------
// 5. select 语句：同时监听多个 channel，带 default 的非阻塞 select
// C 差异：
//   - C 用 select(2)/poll(2)/epoll 监听多个文件描述符（fd），语法复杂；
//   - Go 的 select 语法简洁，直接操作 channel，类型安全；
//   - select 随机选择一个就绪的 case（公平调度，避免饥饿）；
//   - 带 default 的 select 是非阻塞的：没有 case 就绪时立即执行 default；
//   - select 可以同时监听发送和接收操作。
// -----------------------------------------------------------------------------
func demoSelect() {
	fmt.Println("\n--- 5. select 语句：多路复用 channel ---")

	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	// 向两个 channel 发送数据
	ch1 <- "来自 ch1 的消息"
	ch2 <- "来自 ch2 的消息"

	// select 随机选择一个就绪的 case
	fmt.Println("  select 随机选择就绪的 case（两个 channel 都有数据）：")
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  case ch1: %q\n", msg)
		case msg := <-ch2:
			fmt.Printf("  case ch2: %q\n", msg)
		}
	}

	// 带 default 的非阻塞 select
	// C 差异：类似 select(2) 的 timeout=0（立即返回），但语法更简洁
	fmt.Println("\n  带 default 的非阻塞 select（channel 为空时执行 default）：")
	ch3 := make(chan int, 1)

	select {
	case v := <-ch3:
		fmt.Printf("  收到: %d\n", v)
	default:
		fmt.Println("  default: ch3 为空，非阻塞返回")
	}

	ch3 <- 42
	select {
	case v := <-ch3:
		fmt.Printf("  收到: %d\n", v)
	default:
		fmt.Println("  default: ch3 为空")
	}

	// select 实现超时
	fmt.Println("\n  select 实现超时（time.After）：")
	slowCh := make(chan string)
	select {
	case msg := <-slowCh:
		fmt.Printf("  收到: %q\n", msg)
	case <-time.After(50 * time.Millisecond):
		fmt.Println("  超时：50ms 内未收到数据")
	}

	fmt.Println("C 差异：Go select 比 C select(2)/epoll 语法简洁，类型安全")
}

// -----------------------------------------------------------------------------
// 6. channel 的关闭与 range 遍历
// C 差异：
//   - C 的 pipe 关闭后，读端读到 EOF（read 返回 0）；
//   - Go 关闭 channel 后，接收方可以继续接收已缓冲的值，然后收到零值；
//   - 用 v, ok := <-ch 检查 channel 是否已关闭（ok=false 表示已关闭且为空）；
//   - range ch 自动在 channel 关闭后退出循环（最常用的遍历方式）；
//   - 只有发送方应该关闭 channel（接收方关闭会导致发送方 panic）；
//   - 关闭已关闭的 channel 会 panic。
// -----------------------------------------------------------------------------
func demoChannelCloseAndRange() {
	fmt.Println("\n--- 6. channel 关闭与 range 遍历 ---")

	// 生产者：发送数据后关闭 channel
	ch := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i
			fmt.Printf("  [producer] 发送: %d\n", i)
		}
		close(ch) // 关闭 channel，通知消费者没有更多数据
		fmt.Println("  [producer] channel 已关闭")
	}()

	// 消费者：用 range 遍历直到 channel 关闭
	// C 差异：类似 while((n = read(fd, buf, size)) > 0) 读到 EOF
	fmt.Println("  [consumer] 用 range 遍历 channel（自动在关闭后退出）：")
	for v := range ch { // range 在 channel 关闭且为空时自动退出
		fmt.Printf("  [consumer] 收到: %d\n", v)
	}
	fmt.Println("  [consumer] range 循环结束（channel 已关闭）")

	// 用 v, ok 检查 channel 是否关闭
	fmt.Println("\n  用 v, ok := <-ch 检查 channel 状态：")
	ch2 := make(chan string, 2)
	ch2 <- "hello"
	close(ch2)

	v1, ok1 := <-ch2
	fmt.Printf("  第一次接收: v=%q, ok=%v（channel 有数据）\n", v1, ok1)
	v2, ok2 := <-ch2
	fmt.Printf("  第二次接收: v=%q, ok=%v（channel 已关闭且为空）\n", v2, ok2)

	fmt.Println("C 差异：close(ch) 类似 close(fd)，range 自动处理 EOF，比 C 更安全")
}

// -----------------------------------------------------------------------------
// 7. 单向 channel 类型：chan<- T（只写）和 <-chan T（只读）
// C 差异：
//   - C 的 pipe 通过关闭读端或写端来限制方向，但没有类型系统强制；
//   - Go 的单向 channel 类型在编译期强制方向限制，提高代码安全性；
//   - 双向 channel 可以隐式转换为单向 channel（反之不行）；
//   - 函数参数使用单向 channel 类型，明确表达意图（只发送或只接收）。
// -----------------------------------------------------------------------------

// producer 只向 channel 发送数据（chan<- int 只写）
// C 差异：C 无法在类型系统中表达"只写"约束，Go 编译器强制检查
func producer(out chan<- int, n int) {
	for i := 0; i < n; i++ {
		out <- i
	}
	close(out) // 发送方负责关闭
}

// consumer 只从 channel 接收数据（<-chan int 只读）
// C 差异：C 无法在类型系统中表达"只读"约束
func consumer(in <-chan int) []int {
	var results []int
	for v := range in {
		results = append(results, v)
	}
	return results
}

func demoDirectionalChannels() {
	fmt.Println("\n--- 7. 单向 channel 类型 ---")
	fmt.Println("  chan<- T：只写 channel（发送方使用）")
	fmt.Println("  <-chan T：只读 channel（接收方使用）")

	// 双向 channel 可以隐式转换为单向 channel
	ch := make(chan int, 5) // 双向 channel

	// 传给 producer 时隐式转换为 chan<- int（只写）
	// 传给 consumer 时隐式转换为 <-chan int（只读）
	go producer(ch, 5) // ch 隐式转换为 chan<- int
	results := consumer(ch) // ch 隐式转换为 <-chan int

	fmt.Printf("  producer 发送 5 个值，consumer 收到: %v\n", results)

	// 编译期安全：单向 channel 不能反向操作
	// var readOnly <-chan int = ch
	// readOnly <- 1  // 编译错误：cannot send to receive-only channel

	fmt.Println("  单向 channel 在编译期强制方向限制，防止误操作")
	fmt.Println("C 差异：Go 类型系统强制 channel 方向，C 无法在类型层面表达此约束")
}

// -----------------------------------------------------------------------------
// 8. sync.Mutex 保护共享数据，对比 channel 方案
// C 差异：
//   - C 用 pthread_mutex_lock/pthread_mutex_unlock 保护共享数据；
//   - Go 的 sync.Mutex 用法类似，但配合 defer 更安全（不会忘记解锁）；
//   - Go 并发哲学：优先用 channel 通信，但保护共享状态时 Mutex 更合适；
//   - sync.RWMutex 支持多读单写（读多写少场景性能更好）；
//   - 注意：Mutex 不能复制（应通过指针传递）。
// -----------------------------------------------------------------------------
func demoMutex() {
	fmt.Println("\n--- 8. sync.Mutex 保护共享数据 ---")

	// 方案一：sync.Mutex 保护共享计数器
	// C 差异：类似 pthread_mutex_t，但 defer mu.Unlock() 确保不会忘记解锁
	fmt.Println("  【方案一：sync.Mutex】")
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()         // 加锁：C 对应 pthread_mutex_lock(&mu)
			defer mu.Unlock() // 解锁：defer 确保即使 panic 也会解锁
			counter++         // 临界区：安全访问共享数据
		}()
	}
	wg.Wait()
	fmt.Printf("  Mutex 保护的计数器（100 个 goroutine）: %d\n", counter)

	// 方案二：channel 方案（Go 并发哲学）
	// 用 channel 传递数据所有权，避免共享状态
	fmt.Println("\n  【方案二：channel 方案（Go 并发哲学）】")
	counterCh := make(chan int, 1)
	counterCh <- 0 // 初始值放入 channel

	var wg2 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			v := <-counterCh // 取出（获得所有权）
			v++
			counterCh <- v // 放回（释放所有权）
		}()
	}
	wg2.Wait()
	fmt.Printf("  channel 方案的计数器（100 个 goroutine）: %d\n", <-counterCh)

	fmt.Println("\n  Mutex vs Channel 选择原则：")
	fmt.Println("  - 保护共享状态（缓存、计数器）→ Mutex 更直观")
	fmt.Println("  - 传递数据所有权、协调工作流 → Channel 更符合 Go 哲学")
	fmt.Println("C 差异：Go 的 defer mu.Unlock() 比 C 的 pthread_mutex_unlock 更安全")
}

// -----------------------------------------------------------------------------
// 9. sync.Once 实现只执行一次的初始化
// C 差异：
//   - C 用 pthread_once(&once_control, init_func) 实现一次性初始化；
//   - Go 的 sync.Once 用法更简洁，Do(f) 保证 f 只执行一次（即使并发调用）；
//   - sync.Once 常用于延迟初始化单例（懒加载）；
//   - 如果 f 发生 panic，Once 认为 f 已执行（不会再次调用）。
// -----------------------------------------------------------------------------
func demoSyncOnce() {
	fmt.Println("\n--- 9. sync.Once 实现只执行一次的初始化 ---")

	var once sync.Once
	initCount := 0

	// 模拟多个 goroutine 竞争初始化
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Do 保证 f 只执行一次，即使 10 个 goroutine 同时调用
			// C 差异：类似 pthread_once，但语法更简洁
			once.Do(func() {
				initCount++
				fmt.Printf("  [goroutine %d] 执行初始化（只会执行一次）\n", id)
			})
		}(i)
	}
	wg.Wait()

	fmt.Printf("  初始化执行次数: %d（期望: 1）\n", initCount)

	// 典型用法：单例模式（懒加载）
	fmt.Println("\n  典型用法：单例懒加载")
	fmt.Println("  var instance *MyService")
	fmt.Println("  var once sync.Once")
	fmt.Println("  func GetInstance() *MyService {")
	fmt.Println("      once.Do(func() { instance = &MyService{} })")
	fmt.Println("      return instance")
	fmt.Println("  }")
	fmt.Println("C 差异：sync.Once 比 pthread_once 更简洁，无需声明 pthread_once_t 变量")
}

// -----------------------------------------------------------------------------
// 10. context.Context 控制 goroutine 的取消与超时
// C 差异：
//   - C 没有原生的取消机制，通常用全局标志变量或信号（SIGTERM）；
//   - Go 的 context.Context 提供统一的取消、超时、截止时间和值传递机制；
//   - context 通过函数参数传递（约定第一个参数），不存储在结构体中；
//   - context.WithCancel 手动取消，context.WithTimeout 超时自动取消；
//   - goroutine 应定期检查 ctx.Done() channel 来响应取消信号。
// -----------------------------------------------------------------------------
func demoContext() {
	fmt.Println("\n--- 10. context.Context 控制取消与超时 ---")

	// 演示 context.WithCancel：手动取消
	fmt.Println("  【context.WithCancel：手动取消】")
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done(): // 收到取消信号
				fmt.Printf("  [worker] 收到取消信号: %v\n", ctx.Err())
				return
			default:
				// 模拟工作
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	time.Sleep(30 * time.Millisecond)
	cancel() // 手动取消，触发 ctx.Done()
	wg.Wait()
	fmt.Println("  [main] worker 已停止")

	// 演示 context.WithTimeout：超时自动取消
	// C 差异：C 需要用 alarm(2) 或 setitimer 实现超时，Go 用 context 更优雅
	fmt.Println("\n  【context.WithTimeout：超时自动取消】")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2() // 即使超时触发，也应调用 cancel 释放资源

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-time.After(200 * time.Millisecond): // 模拟耗时操作（200ms）
			fmt.Println("  [worker2] 操作完成")
		case <-ctx2.Done(): // 50ms 超时先触发
			fmt.Printf("  [worker2] 超时取消: %v\n", ctx2.Err())
		}
	}()
	wg.Wait()

	// context 传递值（谨慎使用，仅用于请求范围的元数据）
	fmt.Println("\n  【context.WithValue：传递请求元数据】")
	type ctxKey string
	ctx3 := context.WithValue(context.Background(), ctxKey("requestID"), "req-12345")
	requestID := ctx3.Value(ctxKey("requestID"))
	fmt.Printf("  从 context 获取 requestID: %v\n", requestID)
	fmt.Println("  注意：context.Value 仅用于请求范围的元数据（如 requestID、traceID）")
	fmt.Println("C 差异：context 提供统一的取消/超时机制，C 需要手动管理全局标志或信号")
}

// -----------------------------------------------------------------------------
// 11. goroutine 泄漏的典型场景及避免方式
// C 差异：
//   - C 线程泄漏（未 join 的 detached 线程）会消耗系统资源直到进程退出；
//   - Go 的 goroutine 泄漏同样消耗内存和 CPU，且 GC 无法回收被阻塞的 goroutine；
//   - 常见泄漏场景：goroutine 永久阻塞在 channel 发送/接收、等待永不触发的条件；
//   - 避免方式：使用 context 取消、确保 channel 会被关闭、设置超时。
// -----------------------------------------------------------------------------
func demoGoroutineLeak() {
	fmt.Println("\n--- 11. goroutine 泄漏场景及避免方式 ---")

	// 场景一：泄漏示例（仅演示，实际已通过 context 修复）
	// 错误写法（会泄漏）：
	//   go func() {
	//       result := <-neverClosedCh  // 永久阻塞，goroutine 泄漏！
	//       process(result)
	//   }()
	fmt.Println("  【泄漏场景】goroutine 永久阻塞在 channel 接收（无人发送且不关闭）")
	fmt.Println("  错误写法：go func() { result := <-neverClosedCh; ... }()")
	fmt.Println("  问题：neverClosedCh 永不关闭，goroutine 永久阻塞，内存无法释放")

	// 修复方式一：通过关闭 channel 通知 goroutine 退出
	fmt.Println("\n  【修复方式一：关闭 channel 通知退出】")
	done := make(chan struct{})
	resultCh := make(chan int, 1)

	go func() {
		select {
		case v := <-resultCh:
			fmt.Printf("  [worker] 收到结果: %d\n", v)
		case <-done: // 监听退出信号
			fmt.Println("  [worker] 收到退出信号，正常退出（无泄漏）")
		}
	}()

	// 模拟：不发送结果，而是关闭 done channel 通知退出
	time.Sleep(10 * time.Millisecond)
	close(done) // 关闭 done，所有监听 done 的 goroutine 都会收到信号
	time.Sleep(10 * time.Millisecond)

	// 修复方式二：使用 context 取消（推荐方式）
	fmt.Println("\n  【修复方式二：context 取消（推荐）】")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		neverSendCh := make(chan int) // 模拟永不发送的 channel
		select {
		case v := <-neverSendCh:
			fmt.Printf("  [worker2] 收到: %d\n", v)
		case <-ctx.Done(): // context 超时或取消时退出
			fmt.Printf("  [worker2] context 取消，正常退出: %v\n", ctx.Err())
		}
	}()
	wg.Wait()

	fmt.Println("\n  goroutine 泄漏预防原则：")
	fmt.Println("  1. 每个 goroutine 都应有明确的退出条件")
	fmt.Println("  2. 使用 context 传递取消信号（推荐）")
	fmt.Println("  3. 确保 channel 发送方最终会关闭 channel 或发送数据")
	fmt.Println("  4. 使用 select + done channel 或 ctx.Done() 实现可取消的等待")
	fmt.Println("  5. 工具：go tool pprof 的 goroutine profile 可检测泄漏")
	fmt.Println("C 差异：C 线程泄漏同样危险，Go 提供 context 等更优雅的取消机制")

	// 演示 MergeChannels 泛型函数
	fmt.Println("\n  【MergeChannels 泛型扇入演示】")
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)
	ch3 := make(chan int, 3)

	ch1 <- 1
	ch1 <- 4
	close(ch1)
	ch2 <- 2
	ch2 <- 5
	close(ch2)
	ch3 <- 3
	ch3 <- 6
	close(ch3)

	merged := MergeChannels(
		(<-chan int)(ch1),
		(<-chan int)(ch2),
		(<-chan int)(ch3),
	)

	var mergedResults []int
	for v := range merged {
		mergedResults = append(mergedResults, v)
	}
	fmt.Printf("  MergeChannels 合并 3 个 channel，收到 %d 个值（顺序不定）\n", len(mergedResults))
	fmt.Println("  MergeChannels 使用泛型 [T any]，支持任意类型的 channel 合并")
}
