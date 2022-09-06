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
* 或者调用 get(int i) 方法获取索引 i 的内容， 它就不会移动读指针



