// Package structs 演示 Go 语言的结构体、方法与组合。
// 面向有 C 语言基础的开发者，每个演示均附有与 C 的关键差异说明。
//
// Go 没有 class、没有继承，通过结构体+方法+组合实现面向对象编程风格。
package structs

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// -----------------------------------------------------------------------------
// 类型定义
// C 差异：
//   - C 的 struct 只能包含数据字段，方法需要通过函数指针模拟；
//   - Go 的结构体可以绑定方法（值接收者或指针接收者），但方法定义在 struct 外部；
//   - Go 没有 class，没有继承，通过嵌入（embedding）实现代码复用；
//   - 结构体标签（struct tag）是 Go 特有的元数据机制，C 没有对应概念。
// -----------------------------------------------------------------------------

// Animal 表示一种动物，含 json 标签用于序列化。
// C 差异：C 的 struct 没有标签机制，JSON 序列化需要手动编写或使用宏。
type Animal struct {
	Name    string `json:"name"`              // json 标签：序列化时使用 "name" 作为键
	Species string `json:"species"`           // 物种
	Age     int    `json:"age"`               // 年龄
	Weight  float64 `json:"weight,omitempty"` // omitempty：零值时不输出该字段
}

// String 是值接收者方法，返回 Animal 的字符串表示。
// 值接收者：方法内操作的是 Animal 的副本，不会修改原始值。
// C 差异：C 需要通过函数指针或普通函数实现，Go 直接绑定到类型上。
func (a Animal) String() string {
	return fmt.Sprintf("Animal{Name:%q, Species:%q, Age:%d, Weight:%.1f}",
		a.Name, a.Species, a.Age, a.Weight)
}

// Birthday 是指针接收者方法，将 Animal 的年龄加 1。
// 指针接收者：方法内可以修改原始值（通过指针访问）。
// C 差异：C 需要显式传递 struct 指针，如 void birthday(Animal *a)。
func (a *Animal) Birthday() {
	a.Age++ // 直接修改原始值，而非副本
}

// Dog 通过嵌入 Animal 实现代码复用（组合，而非继承）。
// C 差异：
//   - C 可以将 struct 作为另一个 struct 的字段，但无法自动提升方法；
//   - Go 的匿名字段（嵌入）会自动提升被嵌入类型的字段和方法；
//   - 这是 Go 实现"继承"效果的方式，但本质是组合，不是继承。
type Dog struct {
	Animal        // 匿名字段（嵌入）：Animal 的字段和方法被提升到 Dog
	Breed  string `json:"breed"`           // Dog 特有字段：品种
	Leash  bool   `json:"leash,omitempty"` // 是否有牵引绳，omitempty 示例
}

// Demo 演示所有结构体相关内容。
func Demo() {
	demoStructInit()
	demoValueType()
	demoPointerPass()
	demoValueReceiver()
	demoPointerReceiver()
	demoAutoAddressability()
	demoEmbedding()
	demoPromotedMethods()
	demoStructTag()
	demoJSON()
}

// -----------------------------------------------------------------------------
// 1. 结构体定义与两种初始化方式
// C 差异：
//   - C 的结构体初始化：struct Animal a = {"Buddy", "Dog", 3, 10.5}（按位置）；
//   - Go 推荐按字段名初始化，可读性更好，字段顺序无关；
//   - Go 的按位置初始化必须提供所有字段，不推荐使用（字段增减时容易出错）；
//   - Go 结构体字段未指定时自动初始化为零值（C 局部变量不自动清零）。
// -----------------------------------------------------------------------------
func demoStructInit() {
	fmt.Println("\n--- 1. 结构体定义与两种初始化方式 ---")

	// 推荐方式：按字段名初始化（字段顺序无关，未指定字段为零值）
	// C 等价（C99 指定初始化器）：struct Animal a = {.Name="Buddy", .Species="Dog", .Age=3};
	a1 := Animal{
		Name:    "Buddy",
		Species: "Dog",
		Age:     3,
		Weight:  10.5,
	}
	fmt.Printf("按字段名初始化: %v\n", a1)

	// 不推荐方式：按位置初始化（必须提供所有字段，顺序必须与定义一致）
	// C 等价：struct Animal a = {"Whiskers", "Cat", 5, 4.2};
	a2 := Animal{"Whiskers", "Cat", 5, 4.2}
	fmt.Printf("按位置初始化:   %v\n", a2)

	// 部分字段初始化：未指定字段自动为零值
	// C 差异：C 局部 struct 未初始化字段是未定义值（UB），Go 保证为零值
	a3 := Animal{Name: "Ghost", Species: "Cat"}
	fmt.Printf("部分初始化（Age=%d, Weight=%.1f 为零值）: %v\n", a3.Age, a3.Weight, a3)

	// 访问字段：使用点号
	fmt.Printf("a1.Name=%q, a1.Age=%d\n", a1.Name, a1.Age)
	fmt.Println("注意：Go 推荐按字段名初始化，字段增减时不会破坏现有代码")
}

// -----------------------------------------------------------------------------
// 2. 结构体是值类型：赋值时复制，修改副本不影响原结构体
// C 差异：
//   - C 的 struct 赋值也是值复制（浅拷贝），这一点与 Go 相同；
//   - 但 C 的数组字段在赋值时也会被复制（Go 同理）；
//   - Go 的切片/map 字段赋值时只复制引用（浅拷贝），深层数据共享；
//   - 与 C++ 不同，Go 没有拷贝构造函数，始终是浅拷贝。
// -----------------------------------------------------------------------------
func demoValueType() {
	fmt.Println("\n--- 2. 结构体是值类型（赋值时复制）---")

	original := Animal{Name: "Buddy", Species: "Dog", Age: 3, Weight: 10.5}

	// 赋值：复制整个结构体（值语义）
	// C 等价：struct Animal copy = original;（C struct 赋值也是值复制）
	copyAnimal := original
	copyAnimal.Name = "Copy"
	copyAnimal.Age = 99

	fmt.Printf("original: %v\n", original)
	fmt.Printf("copy:     %v\n", copyAnimal)
	fmt.Println("结论：修改 copy 不影响 original（值类型，赋值时深拷贝）")

	// 函数传参也是值复制
	modifyAnimal := func(a Animal) {
		a.Name = "Modified" // 修改的是副本
		a.Age = 0
	}
	modifyAnimal(original)
	fmt.Printf("传参后 original: %v（函数内修改不影响外部）\n", original)
	fmt.Println("注意：大结构体传参开销大，如需避免复制请传递指针 *Animal")
}

// -----------------------------------------------------------------------------
// 3. 通过指针传递结构体避免复制开销
// C 差异：
//   - C 通过传递 struct 指针避免复制，Go 同理；
//   - C 通过指针访问字段需要 ->（如 p->Name），Go 统一使用 .（自动解引用）；
//   - Go 的 p.Name 等价于 (*p).Name，编译器自动处理解引用。
// -----------------------------------------------------------------------------
func demoPointerPass() {
	fmt.Println("\n--- 3. 通过指针传递结构体避免复制开销 ---")

	a := Animal{Name: "Buddy", Species: "Dog", Age: 3, Weight: 10.5}

	// 通过指针修改结构体
	// C 等价：void rename(Animal *p, const char *name) { p->Name = name; }
	rename := func(p *Animal, name string) {
		p.Name = name // Go 自动解引用：p.Name 等价于 (*p).Name
		// C 中需要写 p->Name = name
	}

	fmt.Printf("修改前: %v\n", a)
	rename(&a, "Max") // 传递指针，避免复制
	fmt.Printf("修改后: %v\n", a)

	// new(T) 分配堆上的结构体，返回指针
	// C 等价：Animal *p = malloc(sizeof(Animal));（但 Go 由 GC 自动释放）
	p := new(Animal)
	p.Name = "NewAnimal"
	p.Species = "Unknown"
	fmt.Printf("new(Animal): %v\n", *p)

	// 结构体字面量取地址（常用方式）
	// C 等价：Animal *p2 = &(Animal){"Kitty", "Cat", 2, 3.5};（C99 复合字面量）
	p2 := &Animal{Name: "Kitty", Species: "Cat", Age: 2, Weight: 3.5}
	fmt.Printf("&Animal{...}: %v\n", *p2)
	fmt.Println("注意：Go 用 . 访问指针字段（自动解引用），C 需要用 ->")
}

// -----------------------------------------------------------------------------
// 4. 值接收者方法：方法内操作的是副本
// C 差异：
//   - C 没有方法，只有函数；Go 的方法是绑定到类型的函数；
//   - 值接收者相当于 C 的 void func(Animal a)（传值，操作副本）；
//   - 值接收者方法不能修改原始值，适合只读操作或需要副本的场景；
//   - 实现 fmt.Stringer 接口（String() string）可以自定义 fmt.Println 输出。
// -----------------------------------------------------------------------------
func demoValueReceiver() {
	fmt.Println("\n--- 4. 值接收者方法（操作副本）---")

	a := Animal{Name: "Buddy", Species: "Dog", Age: 3, Weight: 10.5}

	// 调用值接收者方法 String()
	// C 等价：char* animal_string(Animal a)（传值）
	s := a.String()
	fmt.Printf("a.String() = %s\n", s)

	// 验证：值接收者方法内修改不影响原始值
	// 通过一个演示函数来说明
	demonstrateValueReceiver := func(a Animal) string {
		a.Name = "MODIFIED_IN_METHOD" // 修改的是副本
		return a.Name
	}
	result := demonstrateValueReceiver(a)
	fmt.Printf("方法内修改结果: %q\n", result)
	fmt.Printf("原始值未变: a.Name=%q\n", a.Name)

	// fmt.Println 会自动调用 String() 方法（实现了 fmt.Stringer 接口）
	fmt.Printf("fmt.Printf(%%v): %v\n", a)
	fmt.Println("注意：值接收者方法操作副本，不能修改原始值（适合只读操作）")
}

// -----------------------------------------------------------------------------
// 5. 指针接收者方法：方法内可修改原始值
// C 差异：
//   - 指针接收者相当于 C 的 void func(Animal *a)（传指针，可修改原始值）；
//   - Go 的指针接收者方法通过 *T 绑定，可以修改接收者的字段；
//   - 如果方法需要修改接收者，或接收者是大结构体（避免复制），应使用指针接收者；
//   - 同一类型的方法集应统一使用值接收者或指针接收者，避免混用。
// -----------------------------------------------------------------------------
func demoPointerReceiver() {
	fmt.Println("\n--- 5. 指针接收者方法（可修改原始值）---")

	a := Animal{Name: "Buddy", Species: "Dog", Age: 3, Weight: 10.5}
	fmt.Printf("Birthday() 前: Age=%d\n", a.Age)

	// 调用指针接收者方法 Birthday()
	// C 等价：void birthday(Animal *a) { a->Age++; }
	a.Birthday() // Go 自动取地址：等价于 (&a).Birthday()
	fmt.Printf("Birthday() 后: Age=%d（原始值被修改）\n", a.Age)

	a.Birthday()
	a.Birthday()
	fmt.Printf("再调用两次后: Age=%d\n", a.Age)

	// 通过指针变量调用
	p := &a
	p.Birthday()
	fmt.Printf("通过指针调用 p.Birthday() 后: Age=%d\n", a.Age)
	fmt.Println("注意：指针接收者方法修改原始值，适合需要修改接收者或大结构体的场景")
}

// -----------------------------------------------------------------------------
// 6. Go 自动处理值与指针接收者的调用
// C 差异：
//   - C 需要显式区分传值和传指针，调用时需要手动取地址或解引用；
//   - Go 编译器自动处理：对值变量调用指针接收者方法时自动取地址（&v）；
//   - 对指针变量调用值接收者方法时自动解引用（*p）；
//   - 但接口方法集有严格规定：只有指针类型才包含指针接收者方法。
// -----------------------------------------------------------------------------
func demoAutoAddressability() {
	fmt.Println("\n--- 6. Go 自动处理值与指针接收者的调用 ---")

	// 值变量调用指针接收者方法：Go 自动取地址
	a := Animal{Name: "Buddy", Species: "Dog", Age: 3}
	a.Birthday() // 等价于 (&a).Birthday()，Go 自动取地址
	fmt.Printf("值变量调用指针接收者方法（自动取地址）: Age=%d\n", a.Age)

	// 指针变量调用值接收者方法：Go 自动解引用
	p := &Animal{Name: "Max", Species: "Dog", Age: 5}
	s := p.String() // 等价于 (*p).String()，Go 自动解引用
	fmt.Printf("指针变量调用值接收者方法（自动解引用）: %s\n", s)

	// 注意：字面量不可寻址，不能直接调用指针接收者方法
	// Animal{Name: "X"}.Birthday() // 编译错误：cannot take the address of Animal literal
	fmt.Println("注意：字面量不可寻址，不能直接对字面量调用指针接收者方法")
	fmt.Println("注意：Go 自动取地址/解引用，C 需要手动写 &a 或 *p")
}

// -----------------------------------------------------------------------------
// 7. 结构体嵌套（组合）：将 Animal 嵌入 Dog，实现代码复用
// C 差异：
//   - C 可以将 struct 作为另一个 struct 的字段，但字段名不能省略；
//   - Go 的匿名字段（嵌入）允许省略字段名，直接访问嵌入类型的字段；
//   - 这是 Go 实现"继承"效果的方式，但本质是组合（has-a），不是继承（is-a）；
//   - 嵌入类型的字段和方法被"提升"到外层类型，可以直接访问。
// -----------------------------------------------------------------------------
func demoEmbedding() {
	fmt.Println("\n--- 7. 结构体嵌套（组合）---")

	// 创建 Dog（嵌入了 Animal）
	d := Dog{
		Animal: Animal{Name: "Rex", Species: "Dog", Age: 2, Weight: 25.0},
		Breed:  "German Shepherd",
		Leash:  true,
	}

	// 直接访问嵌入字段（字段提升）
	// C 差异：C 需要 d.Animal.Name，Go 可以直接写 d.Name
	fmt.Printf("d.Name（提升字段）: %q\n", d.Name)       // 等价于 d.Animal.Name
	fmt.Printf("d.Age（提升字段）:  %d\n", d.Age)        // 等价于 d.Animal.Age
	fmt.Printf("d.Breed（Dog 自有字段）: %q\n", d.Breed)

	// 也可以通过完整路径访问
	fmt.Printf("d.Animal.Name（完整路径）: %q\n", d.Animal.Name)

	// 修改嵌入字段
	d.Name = "Max" // 等价于 d.Animal.Name = "Max"
	fmt.Printf("修改后 d.Name=%q, d.Animal.Name=%q（同一字段）\n", d.Name, d.Animal.Name)
	fmt.Println("注意：Go 的嵌入是组合（has-a），不是继承（is-a），但效果类似")
}

// -----------------------------------------------------------------------------
// 8. 匿名字段（嵌入）与方法提升（promoted methods）
// C 差异：
//   - C 没有方法提升的概念，需要手动调用 animal_string(&d.animal)；
//   - Go 的嵌入会自动提升被嵌入类型的所有方法，可以直接调用；
//   - 如果外层类型定义了同名方法，则外层方法优先（遮蔽）；
//   - 方法提升使得 Dog 可以直接实现 Animal 的接口（如 fmt.Stringer）。
// -----------------------------------------------------------------------------
func demoPromotedMethods() {
	fmt.Println("\n--- 8. 匿名字段与方法提升（promoted methods）---")

	d := Dog{
		Animal: Animal{Name: "Rex", Species: "Dog", Age: 2, Weight: 25.0},
		Breed:  "German Shepherd",
	}

	// 直接调用 Animal 的值接收者方法（方法提升）
	// C 差异：C 需要 animal_string(d.animal)，Go 可以直接 d.String()
	s := d.String() // 等价于 d.Animal.String()
	fmt.Printf("d.String()（提升的值接收者方法）: %s\n", s)

	// 直接调用 Animal 的指针接收者方法（方法提升）
	fmt.Printf("Birthday() 前: d.Age=%d\n", d.Age)
	d.Birthday() // 等价于 d.Animal.Birthday()
	fmt.Printf("Birthday() 后: d.Age=%d（提升的指针接收者方法）\n", d.Age)

	// fmt.Println 也会调用提升的 String() 方法
	fmt.Printf("fmt.Printf(%%v): %v\n", d.Animal)
	fmt.Println("注意：方法提升让 Dog 自动拥有 Animal 的所有方法，无需重新实现")
}

// -----------------------------------------------------------------------------
// 9. 结构体标签（struct tag）的定义与通过 reflect 包读取
// C 差异：
//   - C 没有结构体标签机制，元数据通常通过注释或宏实现；
//   - Go 的结构体标签是字符串字面量，紧跟在字段类型后面；
//   - 标签通过 reflect 包在运行时读取，常用于 JSON、数据库 ORM、验证等；
//   - 标签格式：`key:"value" key2:"value2"`，多个标签用空格分隔。
// -----------------------------------------------------------------------------
func demoStructTag() {
	fmt.Println("\n--- 9. 结构体标签（struct tag）与 reflect 读取 ---")

	// 通过 reflect 读取结构体标签
	t := reflect.TypeOf(Animal{})
	fmt.Printf("Animal 结构体共 %d 个字段:\n", t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json") // 读取 json 标签
		fmt.Printf("  字段 %-8s: json 标签=%q\n", field.Name, jsonTag)
	}

	// 读取 Dog 的标签（包括嵌入字段）
	fmt.Println()
	dt := reflect.TypeOf(Dog{})
	fmt.Printf("Dog 结构体共 %d 个字段:\n", dt.NumField())
	for i := 0; i < dt.NumField(); i++ {
		field := dt.Field(i)
		jsonTag := field.Tag.Get("json")
		fmt.Printf("  字段 %-8s: json 标签=%q, 是否匿名=%v\n",
			field.Name, jsonTag, field.Anonymous)
	}
	fmt.Println("注意：reflect 包可以在运行时读取标签，C 没有对应的元数据机制")
}

// -----------------------------------------------------------------------------
// 10. 结构体与 JSON 互转：json.Marshal 和 json.Unmarshal
// C 差异：
//   - C 没有内置 JSON 支持，需要使用第三方库（如 cJSON、jansson）；
//   - Go 标准库 encoding/json 通过反射自动处理序列化/反序列化；
//   - json 标签控制字段名映射和特殊行为（omitempty、-）；
//   - omitempty：字段为零值时不输出该字段（减少 JSON 体积）；
//   - json:"-"：完全忽略该字段（不序列化也不反序列化）。
// -----------------------------------------------------------------------------
func demoJSON() {
	fmt.Println("\n--- 10. 结构体与 JSON 互转 ---")

	// 序列化（Marshal）：struct → JSON
	a := Animal{Name: "Buddy", Species: "Dog", Age: 3, Weight: 10.5}
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Printf("Marshal 错误: %v\n", err)
		return
	}
	fmt.Printf("json.Marshal(Animal): %s\n", data)

	// omitempty 效果：Weight 为零值时不输出
	aNoWeight := Animal{Name: "Ghost", Species: "Cat", Age: 2}
	dataNoWeight, _ := json.Marshal(aNoWeight)
	fmt.Printf("Weight=0 时（omitempty）: %s\n", dataNoWeight)
	fmt.Println("注意：Weight 字段有 omitempty 标签，零值时不出现在 JSON 中")

	// 反序列化（Unmarshal）：JSON → struct
	jsonStr := `{"name":"Max","species":"Dog","age":5,"weight":20.0}`
	var a2 Animal
	err = json.Unmarshal([]byte(jsonStr), &a2)
	if err != nil {
		fmt.Printf("Unmarshal 错误: %v\n", err)
		return
	}
	fmt.Printf("json.Unmarshal 结果: %v\n", a2)

	// Dog 的 JSON 序列化（嵌入字段会被展开）
	d := Dog{
		Animal: Animal{Name: "Rex", Species: "Dog", Age: 2, Weight: 25.0},
		Breed:  "German Shepherd",
		Leash:  true,
	}
	dogData, _ := json.Marshal(d)
	fmt.Printf("json.Marshal(Dog): %s\n", dogData)
	fmt.Println("注意：嵌入的 Animal 字段被展开到 Dog 的 JSON 中（不是嵌套对象）")

	// Leash=false 时（omitempty 效果）
	dNoLeash := Dog{
		Animal: Animal{Name: "Buddy", Species: "Dog", Age: 3},
		Breed:  "Labrador",
	}
	noLeashData, _ := json.Marshal(dNoLeash)
	fmt.Printf("Leash=false 时（omitempty）: %s\n", noLeashData)

	// 格式化输出（MarshalIndent）
	prettyData, _ := json.MarshalIndent(d, "", "  ")
	fmt.Printf("json.MarshalIndent:\n%s\n", prettyData)
	fmt.Println("注意：Go 标准库内置 JSON 支持，C 需要第三方库（如 cJSON）")
}
