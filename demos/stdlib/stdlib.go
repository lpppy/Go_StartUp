// Package stdlib 演示 Go 标准库常用包。
// 面向有 C 语言基础的开发者，每个演示均附有与 libc 的关键差异说明。
//
// 标准库设计哲学：
//   - Go 标准库覆盖面广（网络、加密、JSON、HTTP 等），C 的 libc 相对基础；
//   - Go 标准库接口一致，错误处理统一（返回 error），C 用 errno 和返回值；
//   - Go 标准库充分利用接口（io.Reader/io.Writer），组合性强；
//   - 无需第三方库即可完成大多数常见任务。
package stdlib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// FormatTime 使用 Go 参考时间格式化时间。
// Go 的时间格式化使用参考时间：Mon Jan 2 15:04:05 MST 2006
// 这是 Go 特有的设计：用一个具体的时间点作为格式模板，而非 %Y/%m/%d 这样的占位符。
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// CountWords 统计字符串中的单词数（以空白字符分隔）。
func CountWords(s string) int {
	fields := strings.Fields(s)
	return len(fields)
}

// Demo 演示所有标准库相关内容。
func Demo() {
	demoFmt()
	demoStrings()
	demoStrconv()
	demoOS()
	demoBufio()
	demoJSON()
	demoHTTP()
	demoTime()
	demoRand()
	demoRegexp()
	demoSort()
}

// -----------------------------------------------------------------------------
// 1. fmt 格式化动词
// C 差异：
//   - Go 的 fmt 包类似 C 的 printf/scanf，但更安全（类型检查）；
//   - Go 特有的动词：%v（默认格式）、%+v（含字段名）、%#v（Go 语法格式）、%T（类型）；
//   - Go 的 fmt.Sprintf 返回字符串，fmt.Fprintf 写入 io.Writer；
//   - Go 没有 sprintf 的缓冲区溢出风险（字符串自动扩展）。
// -----------------------------------------------------------------------------
func demoFmt() {
	fmt.Println("\n--- 1. fmt 格式化动词 ---")

	type Point struct{ X, Y int }
	p := Point{X: 10, Y: 20}
	ptr := &p

	// %v：默认格式
	fmt.Printf("%%v（默认格式）: %v\n", p)
	// %+v：含字段名
	fmt.Printf("%%+v（含字段名）: %+v\n", p)
	// %#v：Go 语法格式（可直接复制到代码中）
	fmt.Printf("%%#v（Go 语法）: %#v\n", p)
	// %T：类型名
	fmt.Printf("%%T（类型）: %T\n", p)
	// %p：指针地址
	fmt.Printf("%%p（指针地址）: %p\n", ptr)

	// 数值格式
	n := 255
	fmt.Printf("%%d（十进制）: %d\n", n)
	fmt.Printf("%%b（二进制）: %b\n", n)
	fmt.Printf("%%o（八进制）: %o\n", n)
	fmt.Printf("%%x（十六进制小写）: %x\n", n)
	fmt.Printf("%%X（十六进制大写）: %X\n", n)
	fmt.Printf("%%08d（补零）: %08d\n", n)
	fmt.Printf("%%-10d（左对齐）: %-10d|\n", n)

	// 浮点格式
	f := 3.14159265
	fmt.Printf("%%f（默认精度）: %f\n", f)
	fmt.Printf("%%.2f（2位小数）: %.2f\n", f)
	fmt.Printf("%%e（科学计数）: %e\n", f)
	fmt.Printf("%%g（紧凑格式）: %g\n", f)

	// 字符串格式
	s := "Hello, Go"
	fmt.Printf("%%s（字符串）: %s\n", s)
	fmt.Printf("%%q（带引号）: %q\n", s)
	fmt.Printf("%%10s（右对齐）: %10s|\n", s)
	fmt.Printf("%%-10s（左对齐）: %-10s|\n", s)

	// Sprintf 返回格式化字符串
	formatted := fmt.Sprintf("Point(%d, %d)", p.X, p.Y)
	fmt.Printf("Sprintf 结果: %q\n", formatted)
}

// -----------------------------------------------------------------------------
// 2. strings 包
// C 差异：
//   - C 的字符串操作分散在 <string.h>（strstr、strchr 等），Go 集中在 strings 包；
//   - Go 的 strings 包操作不可变字符串，返回新字符串；
//   - strings.Builder 类似 C 的动态字符串缓冲区，避免频繁内存分配。
// -----------------------------------------------------------------------------
func demoStrings() {
	fmt.Println("\n--- 2. strings 包 ---")

	s := "  Hello, Go World!  "

	fmt.Printf("原始字符串: %q\n", s)
	fmt.Printf("Contains(\"Go\"): %v\n", strings.Contains(s, "Go"))
	fmt.Printf("HasPrefix(\"  Hello\"): %v\n", strings.HasPrefix(s, "  Hello"))
	fmt.Printf("HasSuffix(\"!  \"): %v\n", strings.HasSuffix(s, "!  "))
	fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(s))
	fmt.Printf("ToUpper: %q\n", strings.ToUpper(strings.TrimSpace(s)))
	fmt.Printf("ToLower: %q\n", strings.ToLower(strings.TrimSpace(s)))
	fmt.Printf("Replace(\"Go\",\"Golang\"): %q\n", strings.Replace(s, "Go", "Golang", 1))
	fmt.Printf("Count(\"o\"): %d\n", strings.Count(s, "o"))
	fmt.Printf("Index(\"Go\"): %d\n", strings.Index(s, "Go"))

	// Split 和 Join
	csv := "apple,banana,cherry,date"
	parts := strings.Split(csv, ",")
	fmt.Printf("Split(%q, \",\"): %v\n", csv, parts)
	joined := strings.Join(parts, " | ")
	fmt.Printf("Join(parts, \" | \"): %q\n", joined)

	// strings.Builder：高效字符串拼接
	// C 差异：C 需要手动管理缓冲区，Go 的 Builder 自动扩展
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		fmt.Fprintf(&sb, "item%d", i)
		if i < 4 {
			sb.WriteString(", ")
		}
	}
	fmt.Printf("strings.Builder 结果: %q\n", sb.String())

	// CountWords 演示
	text := "  Hello   Go   World  "
	fmt.Printf("CountWords(%q): %d\n", text, CountWords(text))
}

// -----------------------------------------------------------------------------
// 3. strconv 包：数值与字符串互转
// C 差异：
//   - C 用 atoi/atof/sprintf 进行转换，Go 用 strconv 包；
//   - Go 的转换函数返回 error，比 C 的 atoi（失败返回 0）更安全；
//   - strconv.Itoa 是 strconv.FormatInt(n, 10) 的简写。
// -----------------------------------------------------------------------------
func demoStrconv() {
	fmt.Println("\n--- 3. strconv 包 ---")

	// int -> string
	n := 42
	s := strconv.Itoa(n) // 等价于 strconv.FormatInt(int64(n), 10)
	fmt.Printf("Itoa(%d): %q\n", n, s)

	// string -> int
	n2, err := strconv.Atoi("123")
	if err == nil {
		fmt.Printf("Atoi(\"123\"): %d\n", n2)
	}

	// 错误处理：无效输入
	_, err = strconv.Atoi("abc")
	fmt.Printf("Atoi(\"abc\") 错误: %v\n", err)

	// float64 -> string
	f := 3.14159
	fs := strconv.FormatFloat(f, 'f', 2, 64) // 格式 f，2位小数，64位精度
	fmt.Printf("FormatFloat(%.5f, 'f', 2): %q\n", f, fs)

	// string -> float64
	f2, err := strconv.ParseFloat("2.718", 64)
	if err == nil {
		fmt.Printf("ParseFloat(\"2.718\"): %f\n", f2)
	}

	// bool 转换
	fmt.Printf("FormatBool(true): %q\n", strconv.FormatBool(true))
	b, _ := strconv.ParseBool("true")
	fmt.Printf("ParseBool(\"true\"): %v\n", b)

	// 不同进制
	fmt.Printf("FormatInt(255, 16): %q（十六进制）\n", strconv.FormatInt(255, 16))
	fmt.Printf("FormatInt(255, 2): %q（二进制）\n", strconv.FormatInt(255, 2))
	n3, _ := strconv.ParseInt("ff", 16, 64)
	fmt.Printf("ParseInt(\"ff\", 16): %d\n", n3)
}

// -----------------------------------------------------------------------------
// 4. os 包：文件读写、环境变量、命令行参数
// C 差异：
//   - C 用 fopen/fread/fwrite/fclose，Go 用 os.ReadFile/os.WriteFile（更简洁）；
//   - C 用 getenv，Go 用 os.Getenv（相同语义）；
//   - C 用 argc/argv，Go 用 os.Args（切片，更方便）；
//   - Go 的文件操作统一返回 error，C 用 errno。
// -----------------------------------------------------------------------------
func demoOS() {
	fmt.Println("\n--- 4. os 包 ---")

	// 命令行参数（os.Args）
	// C 等价：int main(int argc, char *argv[])
	fmt.Printf("os.Args[0]（程序名）: %s\n", os.Args[0])
	fmt.Printf("os.Args 长度: %d\n", len(os.Args))

	// 环境变量
	// C 等价：getenv("PATH")
	path := os.Getenv("PATH")
	if len(path) > 50 {
		path = path[:50] + "..."
	}
	fmt.Printf("os.Getenv(\"PATH\"): %s\n", path)

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = "（未设置）"
	}
	fmt.Printf("os.Getenv(\"GOPATH\"): %s\n", gopath)

	// 文件写入和读取
	tmpFile := filepath.Join(os.TempDir(), "go_stdlib_demo.txt")
	content := "Hello, Go stdlib!\n第二行内容\n第三行内容\n"

	// os.WriteFile：写入文件（自动创建，0644 权限）
	// C 等价：fopen + fwrite + fclose
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		fmt.Printf("WriteFile 错误: %v\n", err)
		return
	}
	fmt.Printf("os.WriteFile 写入: %q\n", tmpFile)

	// os.ReadFile：读取整个文件
	// C 等价：fopen + fread + fclose（需要先获取文件大小）
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		fmt.Printf("ReadFile 错误: %v\n", err)
		return
	}
	fmt.Printf("os.ReadFile 读取（%d 字节）:\n%s", len(data), data)

	// 清理临时文件
	os.Remove(tmpFile)
	fmt.Println("临时文件已删除")
}

// -----------------------------------------------------------------------------
// 5. bufio.Scanner：逐行读取
// C 差异：
//   - C 用 fgets 逐行读取，Go 用 bufio.Scanner（更简洁）；
//   - bufio.Scanner 可以从任何 io.Reader 读取（文件、网络、字符串等）；
//   - strings.NewReader 将字符串包装为 io.Reader，方便测试。
// -----------------------------------------------------------------------------
func demoBufio() {
	fmt.Println("\n--- 5. bufio.Scanner 逐行读取 ---")

	text := "第一行：Hello\n第二行：Go\n第三行：World\n第四行：bufio"

	// strings.NewReader 将字符串包装为 io.Reader
	reader := strings.NewReader(text)

	// bufio.Scanner 逐行扫描
	scanner := bufio.NewScanner(reader)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		fmt.Printf("  行 %d: %s\n", lineNum, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner 错误: %v\n", err)
	}
	fmt.Printf("共读取 %d 行\n", lineNum)
}

// -----------------------------------------------------------------------------
// 6. encoding/json：JSON 序列化与反序列化
// C 差异：
//   - C 没有内置 JSON 支持，需要第三方库（cJSON、jansson 等）；
//   - Go 标准库内置 encoding/json，通过反射自动处理；
//   - json 标签控制字段名映射（见 structs 模块的详细演示）。
// -----------------------------------------------------------------------------
func demoJSON() {
	fmt.Println("\n--- 6. encoding/json ---")

	// 嵌套结构体
	type Address struct {
		City    string `json:"city"`
		Country string `json:"country"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Email   string  `json:"email,omitempty"` // 零值时省略
		Address Address `json:"address"`
		Tags    []string `json:"tags,omitempty"`
	}

	// 序列化（struct -> JSON）
	p := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{City: "Beijing", Country: "China"},
		Tags: []string{"Go", "developer"},
	}

	data, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Marshal 错误: %v\n", err)
		return
	}
	fmt.Printf("json.Marshal: %s\n", data)

	// 格式化输出
	pretty, _ := json.MarshalIndent(p, "", "  ")
	fmt.Printf("MarshalIndent:\n%s\n", pretty)

	// 反序列化（JSON -> struct）
	jsonStr := `{"name":"Bob","age":25,"address":{"city":"Shanghai","country":"China"},"tags":["Python","AI"]}`
	var p2 Person
	err = json.Unmarshal([]byte(jsonStr), &p2)
	if err != nil {
		fmt.Printf("Unmarshal 错误: %v\n", err)
		return
	}
	fmt.Printf("json.Unmarshal: %+v\n", p2)

	// map 的 JSON 序列化
	m := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}
	mapData, _ := json.Marshal(m)
	fmt.Printf("map JSON: %s\n", mapData)
}

// -----------------------------------------------------------------------------
// 7. net/http：发起 HTTP GET 请求
// C 差异：
//   - C 没有内置 HTTP 客户端，需要 libcurl 等第三方库；
//   - Go 标准库内置完整的 HTTP 客户端和服务器；
//   - 网络请求可能失败，需要处理错误并设置超时。
// -----------------------------------------------------------------------------
func demoHTTP() {
	fmt.Println("\n--- 7. net/http ---")

	// 创建带超时的 HTTP 客户端
	client := &http.Client{
		Timeout: 3 * time.Second, // 3 秒超时
	}

	// 发起 GET 请求
	url := "https://httpbin.org/get"
	fmt.Printf("发起 GET 请求: %s\n", url)

	resp, err := client.Get(url)
	if err != nil {
		// 网络请求可能因网络问题失败，打印错误后继续演示
		fmt.Printf("HTTP GET 失败（网络问题或超时）: %v\n", err)
		fmt.Println("注意：net/http 内置完整 HTTP 客户端，C 需要 libcurl 等第三方库")
		return
	}
	defer resp.Body.Close()

	fmt.Printf("响应状态: %s\n", resp.Status)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// 读取响应体（限制大小）
	buf := make([]byte, 200)
	n, _ := resp.Body.Read(buf)
	fmt.Printf("响应体（前 %d 字节）: %s...\n", n, buf[:n])

	fmt.Println("注意：net/http 内置完整 HTTP 客户端，C 需要 libcurl 等第三方库")
}

// -----------------------------------------------------------------------------
// 8. time 包：时间操作
// C 差异：
//   - C 用 time.h（time_t、struct tm、strftime），Go 用 time 包；
//   - Go 的时间格式化使用参考时间（2006-01-02 15:04:05），而非 %Y/%m/%d；
//   - Go 的 time.Duration 是纳秒精度的整数，类型安全；
//   - Go 的时间操作更直观（Add、Sub、Since、Until 等）。
// -----------------------------------------------------------------------------
func demoTime() {
	fmt.Println("\n--- 8. time 包 ---")

	// 当前时间
	now := time.Now()
	fmt.Printf("当前时间: %v\n", now)

	// 格式化（Go 参考时间：2006-01-02 15:04:05）
	// C 等价：strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &tm)
	fmt.Printf("格式化（参考时间）: %s\n", FormatTime(now))
	fmt.Printf("自定义格式: %s\n", now.Format("2006年01月02日 15:04:05"))
	fmt.Printf("RFC3339 格式: %s\n", now.Format(time.RFC3339))
	fmt.Printf("仅日期: %s\n", now.Format("2006-01-02"))
	fmt.Printf("仅时间: %s\n", now.Format("15:04:05"))

	// 解析时间字符串
	t, err := time.Parse("2006-01-02", "2024-01-15")
	if err == nil {
		fmt.Printf("Parse(\"2024-01-15\"): %v\n", t)
	}

	// 时间计算
	future := now.Add(24 * time.Hour)   // 加 24 小时
	past := now.Add(-7 * 24 * time.Hour) // 减 7 天
	fmt.Printf("明天: %s\n", future.Format("2006-01-02"))
	fmt.Printf("7天前: %s\n", past.Format("2006-01-02"))

	// Duration（时间段）
	d := 2*time.Hour + 30*time.Minute + 15*time.Second
	fmt.Printf("Duration: %v\n", d)
	fmt.Printf("Duration 秒数: %.0f 秒\n", d.Seconds())
	fmt.Printf("Duration 分钟数: %.1f 分钟\n", d.Minutes())

	// 时间差
	start := time.Now()
	time.Sleep(1 * time.Millisecond) // 短暂休眠
	elapsed := time.Since(start)
	fmt.Printf("time.Since: %v（短暂休眠后）\n", elapsed)

	// 时间戳
	fmt.Printf("Unix 时间戳（秒）: %d\n", now.Unix())
	fmt.Printf("Unix 时间戳（纳秒）: %d\n", now.UnixNano())
}

// -----------------------------------------------------------------------------
// 9. math/rand：随机数生成
// C 差异：
//   - C 用 rand()/srand()，Go 用 math/rand 包；
//   - Go 1.20+ 的 rand 包自动使用随机种子，无需手动 srand；
//   - Go 还有 crypto/rand 包用于密码学安全的随机数。
// -----------------------------------------------------------------------------
func demoRand() {
	fmt.Println("\n--- 9. math/rand ---")

	// Go 1.20+ 自动使用随机种子，每次运行结果不同
	fmt.Printf("rand.Intn(100): %d\n", rand.Intn(100))
	fmt.Printf("rand.Intn(100): %d\n", rand.Intn(100))
	fmt.Printf("rand.Float64(): %.4f\n", rand.Float64())
	fmt.Printf("rand.Float64(): %.4f\n", rand.Float64())

	// 生成随机切片
	nums := make([]int, 5)
	for i := range nums {
		nums[i] = rand.Intn(100)
	}
	fmt.Printf("随机切片: %v\n", nums)

	// 固定种子（可重现的随机序列，用于测试）
	r := rand.New(rand.NewSource(42))
	fmt.Printf("固定种子(42) 序列: %d %d %d\n", r.Intn(100), r.Intn(100), r.Intn(100))

	fmt.Println("注意：crypto/rand 用于密码学安全的随机数（如生成密钥）")
}

// -----------------------------------------------------------------------------
// 10. regexp：正则表达式
// C 差异：
//   - C 标准库没有正则表达式（需要 POSIX regex.h 或第三方库）；
//   - Go 标准库内置 regexp 包，使用 RE2 语法（保证线性时间复杂度）；
//   - regexp.MustCompile 在编译时验证正则，适合静态正则表达式。
// -----------------------------------------------------------------------------
func demoRegexp() {
	fmt.Println("\n--- 10. regexp ---")

	// 编译正则表达式
	// MustCompile：编译失败时 panic（适合静态正则，在程序启动时验证）
	emailRe := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	numRe := regexp.MustCompile(`\d+`)
	wordRe := regexp.MustCompile(`\b[A-Z][a-z]+\b`)

	// 匹配
	fmt.Printf("MatchString(email): %v\n", emailRe.MatchString("user@example.com"))
	fmt.Printf("MatchString(invalid): %v\n", emailRe.MatchString("not-an-email"))

	// 查找
	text := "Go 1.21 was released in 2023, and Go 1.22 in 2024"
	nums := numRe.FindAllString(text, -1)
	fmt.Printf("FindAllString(数字): %v\n", nums)

	// 查找并返回位置
	loc := numRe.FindStringIndex(text)
	fmt.Printf("FindStringIndex(第一个数字位置): %v\n", loc)

	// 捕获组
	dateRe := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)
	dateStr := "今天是 2024-01-15，明天是 2024-01-16"
	matches := dateRe.FindAllStringSubmatch(dateStr, -1)
	for _, m := range matches {
		fmt.Printf("日期: %s（年=%s 月=%s 日=%s）\n", m[0], m[1], m[2], m[3])
	}

	// 替换
	result := numRe.ReplaceAllString(text, "NUM")
	fmt.Printf("ReplaceAllString(数字->NUM): %s\n", result)

	// 大写单词查找
	sentence := "Hello World from Go Programming"
	words := wordRe.FindAllString(sentence, -1)
	fmt.Printf("FindAllString(大写开头单词): %v\n", words)
}

// -----------------------------------------------------------------------------
// 11. sort 包：排序
// C 差异：
//   - C 用 qsort（需要函数指针和类型转换），Go 用 sort 包（类型安全）；
//   - Go 1.21+ 推荐使用 slices.Sort（泛型，更简洁）；
//   - sort.Slice 用闭包定义比较函数，比 C 的 qsort 更直观。
// -----------------------------------------------------------------------------
func demoSort() {
	fmt.Println("\n--- 11. sort 包 ---")

	// 整数排序
	nums := []int{5, 2, 8, 1, 9, 3, 7, 4, 6}
	fmt.Printf("排序前: %v\n", nums)
	sort.Ints(nums)
	fmt.Printf("sort.Ints 后: %v\n", nums)

	// 字符串排序
	words := []string{"banana", "apple", "cherry", "date", "elderberry"}
	sort.Strings(words)
	fmt.Printf("sort.Strings 后: %v\n", words)

	// 自定义排序：sort.Slice（用闭包定义比较函数）
	// C 等价：qsort(arr, n, sizeof(Person), compare_func)
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"Diana", 28},
	}

	// 按年龄升序
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Printf("按年龄升序: ")
	for _, p := range people {
		fmt.Printf("%s(%d) ", p.Name, p.Age)
	}
	fmt.Println()

	// 按名字降序
	sort.Slice(people, func(i, j int) bool {
		return people[i].Name > people[j].Name
	})
	fmt.Printf("按名字降序: ")
	for _, p := range people {
		fmt.Printf("%s(%d) ", p.Name, p.Age)
	}
	fmt.Println()

	// 检查是否已排序
	sorted := []int{1, 2, 3, 4, 5}
	unsorted := []int{3, 1, 4, 1, 5}
	fmt.Printf("sort.IntsAreSorted([1,2,3,4,5]): %v\n", sort.IntsAreSorted(sorted))
	fmt.Printf("sort.IntsAreSorted([3,1,4,1,5]): %v\n", sort.IntsAreSorted(unsorted))

	// 二分查找（在已排序切片中）
	idx := sort.SearchInts(sorted, 3)
	fmt.Printf("sort.SearchInts([1,2,3,4,5], 3): 索引=%d\n", idx)

	fmt.Println("\nC 对比：C 的 qsort 需要函数指针和 void* 转换，Go 的 sort.Slice 用闭包更直观")
	fmt.Println("Go 1.21+ 推荐使用 slices.Sort（泛型，见 generics 模块）")
}
