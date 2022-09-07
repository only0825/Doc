# NIO基础



## 1. 三大组件

### 1.1 Channel与Buffer

Java NIO（New IO）系统的**核心**在于：**通道(Channel)和缓冲区(Buffer)**。通道表示打开到 IO 设备(例如：文件、套接字)的连接。

若需要使用 NIO 系统，需要获取用于**连接 IO 设备的通道**以及用于**容纳数据的缓冲区**。然后操作缓冲区，对数据进行处理

简而言之，**通道负责传输，缓冲区负责存储**

**常见的Channel有以下四种**，其中FileChannel主要用于文件传输，其余三种用于网络通信

- FileChannel
- DatagramChannel
- SocketChannel
- ServerSocketChannel

**Buffer有以下几种**，其中使用较多的是ByteBuffer

- ByteBuffer
  - MappedByteBuffer
  - DirectByteBuffer
  - HeapByteBuffer
- ShortBuffer
- IntBuffer
- LongBuffer
- FloatBuffer
- DoubleBuffer
- CharBuffer

![img](https://nyimapicture.oss-cn-beijing.aliyuncs.com/img/20210412135510.png)

### 1.2 Selector 

在使用Selector之前，处理socket连接还有以下两种方法 

##### **使用多线程技术** 

为每个连接分别开辟一个线程，分别去处理对应的socke连接

![img](https://nyimapicture.oss-cn-beijing.aliyuncs.com/img/20210418181918.png)

这种方法存在以下几个问题

* 内存占用高
  * 每个线程都需要占用一定的内存，当连接较多时，会开辟大量线程，导致占用大量内存
* 线程上下文切换成本高
* 只适合连接数少的场景
  * 连接数过多，会导致创建很多线程，从而出现问题

**使用线程池技术**

使用线程池，让线程池中的线程去处理连接

![img](https://nyimapicture.oss-cn-beijing.aliyuncs.com/img/20210418181933.png)

这种方法存在以下几个问题

* 阻塞模式下，线程仅能处理一个连接

  * 线程池中的线程获取任务（task）后，只有当其执行完任务之后（断开连接后），才会去获取并执行下一个任务

  * 若socke连接一直未断开，则其对应的线程无法处理其他socke连接

* 仅适合短连接场景

  * 短连接即建立连接发送请求并响应后就立即断开，使得线程池中的线程可以快速处理其他连接

**使用选择器**

**selector 的作用就是配合一个线程来管理多个 channel（fileChannel因为是阻塞式的，所以无法使用selector）**，

获取这些 channel 上发生的**事件**，这些 channel 工作在**非阻塞模式**下，当一个channel中没有执行任务时，可以去执行其他channel中的任务。

**适合连接数多，但流量较少的场景**

![img](https://nyimapicture.oss-cn-beijing.aliyuncs.com/img/20210418181947.png)

若事件未就绪，调用 selector 的 select() 方法会阻塞线程，直到 channel 发生了就绪事件。这些事件就绪后，select 方法就会返回这些事件交给 thread 来处理



## 2. ByteBuffer

### 2.1 ByteBuffer 正确使用姿势

1. 向 buffer 写入数据，例如调用 channel.read(buffer)
2. 调用 flip() 切换至**读模式**
3. 从 buffer 读取数据，例如调用 buffer.get()
4. 调用clear() 或 compact() 切换至**写模式**
5. 重复 1~4 步骤



### 2.2 ByteBuffer 结构

ByteBuffer 有以下重要属性

* capacity
* position
* limit

https://www.bilibili.com/video/BV1py4y1E7oA?p=7

具体还可看netty-demo中test里面的TestByteBufferReadWrite



### 2.3 ByteBuffer 常见方法

#### 分配空间

可以使用 allocate 方法为 ByteBuffer 分配空间， 其他 buffer 类也有该方法

```java
Bytebuffer buf = ByteBuffer.allocate(16);
```



#### 向 buffer 写入数据

有两种办法

* 调用 channel 的 read 方法
* 调用 buffer 自己的 put 方法

``` java
int readBytes = channel.read(buf);
```

和

```java
buf.put((byte)127);
```



#### 从 buffer 读取数据

同样有两种办法

* 调用 channel 的 write 方法
* 调用 buffer 自己的 get 方法

```java
int writeBytes = channel.write(buf)
```

和

```java
byte b = buf.get();
```

get 方法会让 position 读指针向后走，如果想重复读取数据

* 可以调用 rewind 方法将 position 重新置为 0
* 或者调用 **get(int i)** 方法获取索引 i 的内容， 它就不会移动读指针





#### rewind 

从头开始读

```java
// 'a', 'b', 'c', 'd'
buffer.get(new byte[4]);
debugAll(buffer);
buffer.rewind();
System.out.println((char) buffer.get());
```

```
+--------+-------------------- all ------------------------+----------------+
position: [4], limit: [4]
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 61 62 63 64 00 00 00 00 00 00                   |abcd......      |
+--------+-------------------------------------------------+----------------+
a
```



#### mark & reset

mark 做一个标记，记录 position 位置， reset 是将 position 重置到 mark 的位置

```java
// 'a', 'b', 'c', 'd'
System.out.println((char) buffer.get());
System.out.println((char) buffer.get());
buffer.mark(); // 加标记，索引2 的位置
System.out.println((char) buffer.get());
System.out.println((char) buffer.get());
buffer.reset(); // 将 position 重置到索引 2
System.out.println((char) buffer.get());
System.out.println((char) buffer.get());
/*
a
b
c
d
c
d
/*
```

#### get(i) 

不会改变读索引的位置

```java
System.out.println((char) buffer.get(3));
debugAll(buffer);
```

```
d
+--------+-------------------- all ------------------------+----------------+
position: [0], limit: [4]
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 61 62 63 64 00 00 00 00 00 00                   |abcd......      |
+--------+-------------------------------------------------+----------------+
```



#### 字符串与 ByteBuffer 互转

```java
// 字符串转为 ByteBuffer 3种方法
// 1. 原始方法
ByteBuffer buffer1 = ByteBuffer.allocate(16);
buffer1.put("hello".getBytes());
debugAll(buffer1);
// 2. Charset
ByteBuffer buffer2 = StandardCharsets.UTF_8.encode("hello"); // 会自动转为
debugAll(buffer2);
// 3. wrap
ByteBuffer buffer3 = ByteBuffer.wrap("hello".getBytes());
debugAll(buffer3);


// ByteBuffer 转字符串
// 2 和 3 都是直接切换到读模式 再转可以不用filp()切换到读模式
String str1 = StandardCharsets.UTF_8.decode(buffer2).toString();
System.out.println(str1);
// 1方法再转就要用filp()转成读模式
buffer1.flip();
String str2 = StandardCharsets.UTF_8.decode(buffer1).toString();
System.out.println(str2);
```



### 2.4 分散读与集中写

主要是一种思想，可以减少数据再ByteBuffer之间的拷贝，变相的提高了效率

#### 2.4.1 Scattering Reads

分散读取，有一个文本words.txt

```
onetwothree
```

**需求**：读取words.txt中的onetwothree，并输出one、two、three

**思路** (重点)

1. 将数据存入一个ByteBuffer，之后再用别的方法拆分成三组，涉及到数据重新的分割复制
2. 读取时一次读到3个ByteBuffer (分散读取)

```java
// 分散读取
try (FileChannel channel = new RandomAccessFile("words.txt", "r").getChannel()) {
  ByteBuffer b1 = ByteBuffer.allocate(3);
  ByteBuffer b2 = ByteBuffer.allocate(3);
  ByteBuffer b3 = ByteBuffer.allocate(5);
  channel.read(new ByteBuffer[]{b1, b2, b3});
  b1.flip();
  b2.flip();
  b3.flip();
  debugAll(b1);
  debugAll(b2);
  debugAll(b3);
} catch (IOException e) {
  e.printStackTrace();
}
```

```
+--------+-------------------- all ------------------------+----------------+
position: [0], limit: [3]
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 6f 6e 65                                        |one             |
+--------+-------------------------------------------------+----------------+
+--------+-------------------- all ------------------------+----------------+
position: [0], limit: [3]
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 74 77 6f                                        |two             |
+--------+-------------------------------------------------+----------------+
+--------+-------------------- all ------------------------+----------------+
position: [0], limit: [5]
         +-------------------------------------------------+
         |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |
+--------+-------------------------------------------------+----------------+
|00000000| 74 68 72 65 65                                  |three           |
+--------+-------------------------------------------------+----------------+
```



#### 2.4.2 Gathering Writes

需求：三个ByteBuffer写入到一个文件中

思路：

1. 三个ByteBuffer组合到一个大的ByteBuffer中，涉及到数据到多次拷贝
2. 三个ByteBuffer组合到一起，以一个整体写入（集中写）



```java
ByteBuffer b1 = StandardCharsets.UTF_8.encode("hello");
ByteBuffer b2 = StandardCharsets.UTF_8.encode("world");
ByteBuffer b3 = StandardCharsets.UTF_8.encode("你好");

try (FileChannel channel = new RandomAccessFile("words2.txt", "rw").getChannel()) {
    channel.write(new ByteBuffer[]{b1, b2, b3});
} catch (IOException e) {
    e.printStackTrace();
}
```





### 2.5 粘包与半包

#### 现象

网络上有多条数据发送给服务端，数据之间使用 \n 进行分隔
但由于某种原因这些数据在接收时，被进行了重新组合，例如原始数据有3条为

- Hello,world\n
- I’m Nyima\n
- How are you?\n

变成了下面的两个 byteBuffer (粘包，半包)

- Hello,world\nI’m Nyima\nHo
- w are you?\n

#### 出现原因

**粘包**

发送方在发送数据时，并不是一条一条地发送数据，而是**将数据整合在一起**，当数据达到一定的数量后再一起发送。这就会导致多条信息被放在一个缓冲区中被一起发送出去

**半包**

接收方的缓冲区的大小是有限的，当接收方的缓冲区满了以后，就需要**将信息截断**，等缓冲区空了以后再继续放入数据。这就会发生一段完整的数据最后被截断的现象

#### 解决办法

下面只是原始解法，要会，这样才能知道netty帮我们做了哪些事情

```java
   // 解决粘包和半包问题，下面只是原始解法，要会，这样才能知道netty帮我们做了哪些事情
    public static void main(String[] args) {
        /*
            网络上有多条数据发送给服务端，数据之间使用 \n 进行分隔
            但由于某种原因这些数据在接收时，被进行了重新组合，例如原始数据有3条为
            - Hello,world\n
            - I’m Nyima\n
            - How are you?\n
            变成了下面的两个 byteBuffer (粘包，半包)
            - Hello,world\nI’m Nyima\nHo
            - w are you?\n
            现在要求你编写程序，将错误的数据恢复成原始按 \n 分割的数据
         */
        ByteBuffer source = ByteBuffer.allocate(32);
      	// 模拟粘包+半包
        source.put("Hello,world\nI'm zhangsan\nHo".getBytes());
        split(source);
        source.put("w are you ?\n".getBytes());
        split(source);
    }

    public static void split(ByteBuffer source) {
      	// 切换为读模式
        source.flip();
        for (int i = 0; i < source.limit(); i++) {
            // 找到一条完整消息
            if (source.get(i) == '\n') { // get(i)不会移动position
                int length = i + 1 - source.position(); // 当前位置 - 读指针位置
                // 把这条完整消息存入新的 ByteBuffers
                ByteBuffer target = ByteBuffer.allocate(length);
                // 从source 读，向 target 写
                for (int j = 0; j < length; j++) {
                    target.put(source.get()); // get 方法会让 position 读指针向后走
                }
                debugAll(target);
            }
        }
        // 切换为写模式，但是缓冲区可能未读完，这里需要使用compact
        source.compact();
    }
```



## 3. 文件编程

重点是网络编程，文件编程只了解即可

### 3.1 FileChannel

#### :warning: 注意

> FileChannel 只能工作再阻塞模式下  不能配合Selector



#### 获取

不能直接打开 FileChannel，**必须**通过 FileInputStream、FileOutputStream 或者 RandomAccessFile 来获取 FileChannel，它们都有 getChannel 方法

- 通过 FileInputStream 获取的 channel **只能读**
- 通过 FileOutputStream 获取的 channel **只能写**
- 通过 RandomAccessFile 是否能读写**根据构造 RandomAccessFile 时的读写模式决定**



#### 读取

会从 channel 读取数据填充 ByteBuffer，返回值表示读到了多少字节，-1 表示到达了文件的末尾

```java
int readBytes = channel.read(buffer)
```



#### 写入

因为channel也是有大小的，所以 write 方法并不能保证一次将 buffer 中的内容全部写入 channel。必须**需要按照以下规则进行写入**

```java
// 通过hasRemaining()方法查看缓冲区中是否还有数据未写入到通道中
while(buffer.hasRemaining()) {
	channel.write(buffer);
}
```



#### 关闭

channel必须关闭，一般情况通过try-with-resource进行关闭。

调用了 FileInputStream、FileOutputStream 或者 RandomAccessFile 的 close 方法会间接的调用 channel 的 close 方法



#### 位置

**position**

channel也拥有一个保存读取数据位置的属性，即position

```java
long pos = channel.position();
```

可以通过position(int pos)设置channel中position的值

```java
long newPos = ...;
channel.position(newPos);
```

设置当前位置时，如果设置为文件的末尾

- 这时读取会返回 -1
- 这时写入，会追加内容，但要注意如果 position 超过了文件末尾，再写入时在新内容和原末尾之间会有空洞（00）



#### 大小

使用 size() 方法获取文件的大小



#### 强制写入

操作系统出于性能的考虑，会将数据缓存，不是立刻写入磁盘，而是等到缓存满了以后将所有数据一次性的写入磁盘。可以调用 **force(true)** 方法将文件内容和元数据（文件的权限等信息）立刻写入磁盘



### 3.2 两个 Channel 传输数据

transferTo 方法，大小限制2G，但可以多次传输

```java
        try (
                FileChannel from = new FileInputStream("data.txt").getChannel();
                FileChannel to = new FileOutputStream("to.txt").getChannel();
        ) {
            // 效率高，底层会利用操作系统的零拷贝进行优化, 有2G数据限制
            long size = from.size();
            // left 变量代表还剩余多少字节
            for (long left = size; left > 0; ) {
                System.out.println("position:" + (size - left) + " left:" + left);
                left -= from.transferTo((size - left), left, to);
            }
            from.transferTo(0, from.size(), to);
        } catch (IOException e) {
            e.printStackTrace();
        }
```

### 3.3 Path

- Path 用来表示文件路径
- Paths 是工具类，用来获取 Path 实例

```
Path source = Paths.get("1.txt"); // 相对路径 不带盘符 使用 user.dir 环境变量来定位 1.txt

Path source = Paths.get("d:\\1.txt"); // 绝对路径 代表了  d:\1.txt 反斜杠需要转义

Path source = Paths.get("d:/1.txt"); // 绝对路径 同样代表了  d:\1.txt

Path projects = Paths.get("d:\\data", "projects"); // 代表了  d:\data\projectsCopy
```

- `.` 代表了当前路径
- `..` 代表了上一级路径

例如目录结构如下

```
d:
	|- data
		|- projects
			|- a
			|- bCopy
```

代码

```
Path path = Paths.get("d:\\data\\projects\\a\\..\\b");
System.out.println(path);
System.out.println(path.normalize()); // 正常化路径 会去除 . 以及 ..Copy
```

输出结果为

```
d:\data\projects\a\..\b
d:\data\projects\bCopy
```

### 3.4 Files

#### 查找

检查文件是否存在

```
Path path = Paths.get("helloword/data.txt");
System.out.println(Files.exists(path));Copy
```

#### 创建

创建**一级目录**

```
Path path = Paths.get("helloword/d1");
Files.createDirectory(path);Copy
```

- 如果目录已存在，会抛异常 FileAlreadyExistsException
- 不能一次创建多级目录，否则会抛异常 NoSuchFileException

创建**多级目录用**

```
Path path = Paths.get("helloword/d1/d2");
Files.createDirectories(path);Copy
```

#### 拷贝及移动

**拷贝文件**

```
Path source = Paths.get("helloword/data.txt");
Path target = Paths.get("helloword/target.txt");

Files.copy(source, target);Copy
```

- 如果文件已存在，会抛异常 FileAlreadyExistsException

如果希望用 source **覆盖**掉 target，需要用 StandardCopyOption 来控制

```
Files.copy(source, target, StandardCopyOption.REPLACE_EXISTING);Copy
```

移动文件

```
Path source = Paths.get("helloword/data.txt");
Path target = Paths.get("helloword/data.txt");

Files.move(source, target, StandardCopyOption.ATOMIC_MOVE);Copy
```

- **StandardCopyOption.ATOMIC_MOVE 保证文件移动的原子性**

#### 删除

删除文件

```
Path target = Paths.get("helloword/target.txt");

Files.delete(target);Copy
```

- 如果文件不存在，会抛异常 NoSuchFileException

删除目录

```
Path target = Paths.get("helloword/d1");

Files.delete(target);Copy
```

- 如果**目录还有内容**，会抛异常 DirectoryNotEmptyException

#### 删除多级目录

```java
Files.walkFileTree(Paths.get("路径"), new SimpleFileVisitor<Path>() {
    @Override
    public FileVisitResult visitFile(Path file, BasicFileAttributes attrs) throws IOException {
        // 访问目录时删除目录里面的文件
        Files.delete(file);
        return super.visitFile(file, attrs);
    }
    @Override
    public FileVisitResult postVisitDirectory(Path dir, IOException exc) throws IOException {
        // 退出目录时删除目录
        Files.delete(dir);
        return super.postVisitDirectory(dir, exc);
    }
});
```



#### 遍历多级目录

```java
Path path = Paths.get("/Users/wh37/Documents/Doc/Netty/netty-demo");
// 文件目录数目
AtomicInteger dirCount = new AtomicInteger();
// 文件数目
AtomicInteger fileCount = new AtomicInteger();
Files.walkFileTree(path, new SimpleFileVisitor<Path>(){
    @Override
    public FileVisitResult preVisitDirectory(Path dir, BasicFileAttributes attrs) throws IOException {
        System.out.println("===>"+dir);
        // 增加文件目录数
        dirCount.incrementAndGet();
        return super.preVisitDirectory(dir, attrs);
    }

    @Override
    public FileVisitResult visitFile(Path file, BasicFileAttributes attrs) throws IOException {
        System.out.println(file);
        // 增加文件数
        fileCount.incrementAndGet();
        return super.visitFile(file, attrs);
    }

});
// 打印数目
System.out.println("文件目录数:"+dirCount.get());
System.out.println("文件数:"+fileCount.get());
```



#### 拷贝多级目录

```java
String source = "/Users/wh37/Documents/Doc/Netty/netty-demo/logs";
String target = "/Users/wh37/Documents/Doc/Netty/netty-demo/logs-test-copy";

Files.walk(Paths.get(source)).forEach(path -> {
    try {
        String targetName = path.toString().replace(source, target);
        // 是目录
        if (Files.isDirectory(path)) {
            Files.createDirectory(Paths.get(targetName));
        }
        // 是普通文件
        else if (Files.isRegularFile(path)) {
            Files.copy(path, Paths.get(targetName));
        }
    } catch (IOException e) {
        e.printStackTrace();
    }
});
```

