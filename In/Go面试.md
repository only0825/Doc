## Go语言方法和接收器



### 函数式方法

```go
func 函数名(参数列表) (返回参数) {
   // 函数体
}
```

函数一经定义都可被调用，无指定作用对象，无归属感



### “类”方法

```go
// 普通接收器
func (接收器变量 接收器类型) 方法名(参数列表) (返回参数){
	// 方法体
}
// 指针接收器
func (接收器变量 *接收器类型) 方法名(参数列表) (返回参数){
	// 方法体
}
```

| 比较项           | 面向对象类 | go中类型+所属方法=“伪类”                                     |
| :--------------- | :--------- | :----------------------------------------------------------- |
| 方法放置同一文件 | 必须       | 可以不在同一源文件，但必须在同一个包中                       |
| 重载             | 支持       | 不支持，但针对接收器方法，同一包内允许存在同名方法（不同类的同名方法） |



### 函数与方法比较

区别，非指针接收器在运行时，将接收器复制一份，仅在方法内修改成员有效，方法外不会影响原接收器；
若需要改变后的接收器，新返回新生成的接收器

```go
package main

type Bag struct {
	items []int
}

// 普通函数
func Insert(b *Bag,id int) {
	b.items = append(b.items,id)
}

// 指针接收器
func (b *Bag) Insert(id int) {
	b.items = append(b.items,id)
}

func main()  {

	newbag := &Bag{}

	// 函数式调用
	Insert(newbag,1024)

	// “类”方法调用
	newbag.Insert(1024)

}

```

1. 小对象由于值复制时的速度较快，所以适合使用非指针接收器，大对象因为复制性能较低，适合使用指针接收器，在接收器和参数间传递时不进行复制，只是传递指针。

2. 或关注类型的本质，成员是内置类型（int，float…）,引用类型（map，slice…）使用值接受者（非指针接收器），成员是结构类型使用指针接受器。但也是根据具体需而定。（总之遇事不决用指针）



### 如何选择指针接收器或非指针接收器

```go
package main

import "fmt"

type apple struct {
	color string
}

type orange struct {
	color string
}

func (a apple) changeColor(color string) apple {
	a.color = color
	fmt.Printf("非指针接收器，apple方法内修改：%v\n", a)
	return a
}

func (o *orange) changeColor(color string) {
	o.color = color
	fmt.Printf("指针接收器，orange方法内修改：%v\n", o)
}

func main() {

	apple := apple{}
	orange := orange{}
	newApple := apple.changeColor("green")
	orange.changeColor("green")
	fmt.Printf("修改后获取apple:%v\n", apple)
	fmt.Printf("修改后获取orange:%v\n", orange)
	fmt.Printf("修改后获取新生成的apple:%v\n", newApple)
}

```



### 指针类型无法调用非指针接收器

```go
package main

import "fmt"


type orange struct {
	color string
	from string
}


func (o *orange) changeColor(color string) {
	o.color = color
	fmt.Printf("指针接收器，orange方法内修改：%v\n", o)
}

func (o orange) changeFrom(from string) orange {
	o.from = from
	return o
}

func main() {

	// 非指针类型能调用指针方法
	o1 := orange{}
	o1.changeColor("green")
	fmt.Printf("修改后获取orange:%v\n", o1)
	o2 := &orange{}
	o2.changeColor("red")
	fmt.Printf("修改后获取orange:%v\n", o2)

	o1 = o1.changeFrom("north")
	//o2 = o2.changeFrom("south") 无法运行，指针类型无法调用非指针接收器
	fmt.Printf("修改后获取orange:%v\n", o1)
	fmt.Printf("修改后获取orange:%v\n", o2)
}

```



### 接口

鸭子类型：当看到一只鸟走起来像鸭子、游泳起来像鸭子、叫起来也像鸭子，那么这只鸟就可以被称为鸭子。

go采用这种类型实现接口。

使用结构体实现接口：当一个结构体具备接口的所有的方法的时候，它就实现了这个接口

```go
package main

import (
	"fmt"
	"net/http"
)

type server interface {
	route() func(pattern string,HandlerFunc http.HandlerFunc)
	start() func(address string) error
}

// 当一个结构体具备接口的所有的方法的时候，它就实现了这个接口
type webserver struct {
	name string
}


func (w *webserver) route(pattern string,HandlerFunc http.HandlerFunc) {
	http.HandleFunc(pattern,HandlerFunc)
}

func (w *webserver) start(address string) error {
	return http.ListenAndServe(address,nil)
}

func home(w http.ResponseWriter,r *http.Request) {
	fmt.Fprintf(w,"Ciao %s",r.URL.Path[1:])
}

func main() {

	obj := &webserver{
		name: "openresty",
	}
	obj.route("/",home)
	obj.start("localhost:8099")

}

```



## 面试题

### `=` 和 `:=` 的区别？

=是赋值变量，:=是定义变量



### 指针的作用

一个指针可以指向任意变量的地址，它所指向的地址在32位或64位机器上分别**固定**占4或8个字节。指针的作用有：

- 获取变量的值

```go
 import fmt
 
 func main(){
  a := 1
  p := &a//取址&
  fmt.Printf("%d\n", *p);//取值*
 }
```

- 改变变量的值

```go
 // 交换函数
 func swap(a, b *int) {
     *a, *b = *b, *a
 }
```

- 用指针替代值传入函数，比如类的接收器就是这样的。

```go
 type A struct{}
 
 func (a *A) fun(){}
```



### Go 有异常类型吗？

有。Go用error类型代替try...catch语句，这样可以节省资源。同时增加代码可读性：

```text
 _, err := funcDemo()
if err != nil {
    fmt.Println(err)
    return
}
```

也可以用errors.New()来定义自己的异常。errors.Error()会返回异常的字符串表示。只要实现error接口就可以定义自己的异常，

```text
 type errorString struct {
  s string
 }
 
 func (e *errorString) Error() string {
  return e.s
 }
 
 // 多一个函数当作构造函数
 func New(text string) error {
  return &errorString{text}
 }
```



###  什么是协程（Goroutine）

协程是**用户态轻量级线程**，它是**线程调度的基本单位**。通常在函数前加上go关键字就能实现并发。一个Goroutine会以一个很小的栈启动2KB或4KB，当遇到栈空间不足时，栈会**自动伸缩**， 因此可以轻易实现成千上万个goroutine同时启动。



### 如何高效地拼接字符串

拼接字符串的方式有：`+` , `fmt.Sprintf` , `strings.Builder`, `bytes.Buffer`, `strings.Join`

**1 "+"**

使用`+`操作符进行拼接时，会对字符串进行遍历，计算并开辟一个新的空间来存储原来的两个字符串。

**2 fmt.Sprintf**

由于采用了接口参数，必须要用反射获取值，因此有性能损耗。

**3 strings.Builder：**

用WriteString()进行拼接，内部实现是指针+切片，同时String()返回拼接后的字符串，它是直接把[]byte转换为string，从而避免变量拷贝。

**4 bytes.Buffer**

`bytes.Buffer`是一个一个缓冲`byte`类型的缓冲器，这个缓冲器里存放着都是`byte`，

`bytes.buffer`底层也是一个`[]byte`切片。

5 **strings.join**

`strings.join`也是基于`strings.builder`来实现的,并且可以自定义分隔符，在join方法内调用了b.Grow(n)方法，这个是进行初步的容量分配，而前面计算的n的长度就是我们要拼接的slice的长度，因为我们传入切片长度固定，所以提前进行容量分配可以减少内存分配，很高效。

**性能比较**：

strings.Join ≈ strings.Builder > bytes.Buffer > "+" > fmt.Sprintf

5种拼接方法的实例代码

```go
func main(){
	a := []string{"a", "b", "c"}
	//方式1：+
	ret := a[0] + a[1] + a[2]
	//方式2：fmt.Sprintf
	ret := fmt.Sprintf("%s%s%s", a[0],a[1],a[2])
	//方式3：strings.Builder
	var sb strings.Builder
	sb.WriteString(a[0])
	sb.WriteString(a[1])
	sb.WriteString(a[2])
	ret := sb.String()
	//方式4：bytes.Buffer
	buf := new(bytes.Buffer)
	buf.Write(a[0])
	buf.Write(a[1])
	buf.Write(a[2])
	ret := buf.String()
	//方式5：strings.Join
	ret := strings.Join(a,"")
}
```



### 什么是 rune 类型

ASCII 码只需要 7 bit 就可以完整地表示，但只能表示英文字母在内的128个字符，为了表示世界上大部分的文字系统，发明了 Unicode， 它是ASCII的超集，包含世界上书写系统中存在的所有字符，并为每个代码分配一个标准编号（称为Unicode CodePoint），在 Go 语言中称之为 rune，是 int32 类型的别名。

Go 语言中，字符串的底层表示是 byte (8 bit) 序列，而非 rune (32 bit) 序列。

```go
sample := "我爱GO"
runeSamp := []rune(sample)
runeSamp[0] = '你'
fmt.Println(string(runeSamp))  // "你爱GO"
fmt.Println(len(runeSamp))  // 4
```



### 如何判断 map 中是否包含某个 key ？

```go
var sample map[int]int
if _, ok := sample[10]; ok {

} else {

}
```



### Go 支持默认参数或可选参数吗？

不支持。但是可以利用结构体参数，或者...传入参数切片数组。

```go
// 这个函数可以传入任意数量的整型参数
func sum(nums ...int) {
    total := 0
    for _, num := range nums {
        total += num
    }
    fmt.Println(total)
}
```



### defer 的执行顺序

defer执行顺序和调用顺序相反，类似于栈**后进先出**(LIFO)。

**defer在return之后执行，但在函数退出之前，defer可以修改返回值（但要是有名称的返回值，无名称的返回值会会创建临时变量所以不能修改原来的值）。**

下面是一个例子：

```go
func test() int {
	i := 0
	defer func() {
		fmt.Println("defer1")
	}()
	defer func() {
		i += 1
		fmt.Println("defer2")
	}()
	return i
}

func main() {
	fmt.Println("return", test())
}
// defer2
// defer1
// return 0
```

上面这个例子中，test返回值并没有修改，这是由于Go的返回机制决定的，执行Return语句后，Go会创建一个临时变量保存返回值。如果是有名返回（也就是指明返回值`func test() (i int)`）

```go
func test() (i int) {
	i = 0
	defer func() {
		i += 1
		fmt.Println("defer2")
	}()
	return i
}

func main() {
	fmt.Println("return", test())
}
// defer2
// return 1
```

这个例子中，返回值被修改了。对于有名返回值的函数，执行 return 语句时，并不会再创建临时变量保存，因此，defer 语句修改了 i，即对返回值产生了影响。



### 如何交换 2 个变量的值？

对于变量而言`a,b = b,a`； 对于指针而言`*a,*b = *b, *a`





### Go语言的CSP模型

### Go语言中 select 和 switch 的比较

