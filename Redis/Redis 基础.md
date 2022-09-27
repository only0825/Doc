

https://www.oz6.cn/articles/58

# 一、初识 Redis

## 1.认识NoSQL

### 1.1 什么是NoSQL


- NoSQL最常见的解释是"`non-relational`"， 也有人说它是"***Not Only SQL***"
- NoSQL仅仅是一个概念，泛指**非关系型的数据库**
- 区别于关系数据库，它们不保证关系数据的 ACID 特性
- NoSQL是一项全新的数据库革命性运动，提倡运用非关系型的数据存储，相对于铺天盖地的关系型数据库运用，这一概念无疑是一种全新的思维的注入
- 常见的NoSQL数据库有：`Redis`、`MemCache`、`MongoDB`等

### 1.2 NoSQL与SQL的差异

|          |                             SQL                              | NoSQL                                                        |
| :------: | :----------------------------------------------------------: | ------------------------------------------------------------ |
| 数据结构 |                            结构化                            | 非结构化                                                     |
| 数据关联 |                            关联的                            | 无关联的                                                     |
| 查询方式 |                           SQL查询                            | 非SQL                                                        |
| 事务特性 |                             ACID                             | BASE                                                         |
| 存储方式 |                             磁盘                             | 内存                                                         |
|  扩展性  |                             垂直                             | 水平                                                         |
| 使用场景 | 1）数据结构固定 <br />2）相关业务对数据安全性、一致性要求较高 | 1）数据结构不固定 <br />2）对一致性、安全性要求不高 <br />3）对性能要求 |



## 2.认识Redis

> Redis诞生于2009年全称是Remote Dictionary Server，远程词典服务器，是一个基于内存的键值型NoSQL数据库。

**Redis的特征：**

- 键值（`key-value`）型，value支持多种不同数据结构，功能丰富
- 单线程，每个命令具备原子性
- 低延迟，速度快（基于内存、IO多路复用、良好的编码）。
- 支持数据持久化（防断电）
- 支持主从集群、分片集群
- 支持多语言客户端



## 3.安装Redis

### 3.1 前置准备

> 本次安装Redis是基于Linux系统下安装的，因此需要一台Linux服务器或者虚拟机。
>
> Ps：由于提供的CentOS操作系统为mini版，因此需要自行配置网络，不会配置的请联系我，如果您使用的是自己购买的服务器，请提前开放`6379`端口，避免后续出现的莫名其妙的错误！

- **虚拟机**：[VMware16](https://pan.baidu.com/s/1Zn13h9G7MtSgz-xdkQFeJg?pwd=1234)
- **操作系统**：[CentOS-7-x86_64-Minimal-1708](https://pan.baidu.com/s/1SiYip29cYqiNBqjGGV0JgA?pwd=1234)
- **Redis**：[redis-6.2.6.tar](https://pan.baidu.com/s/1hsoEz1NTCDCCWZmaiZrIgg?pwd=1234)
- **xShell及xFtp**：https://www.xshell.com/zh/free-for-home-school/

### 3.2 安装Redis依赖

------

> Redis是基于C语言编写的，因此首先需要安装Redis所需要的gcc依赖

```sh
yum install -y gcc tcl
```

### 3.3 正式安装Redis

- **将`redis-6.2.6.tar`上传至`/usr/local/src`目录**

- **在xShell中`cd`到`/usr/local/src`目录执行以下命令进行解压操作**

  ```sh
  tar -xzf redis-6.2.6.tar.gz
  sh
  ```

- **解压成功后依次执行以下命令**

  ```sh
  cd redis-6.2.6
  make
  make install
  sh
  ```

- **安装成功后打开/usr/local/bin目录（该目录为Redis默认的安装目录）**

## 4.启动Redis

> Redis的启动方式有很多种，例如：**前台启动**、**后台启动**、**开机自启**

### 4.1 前台启动（不推荐）

------

> **这种启动属于前台启动，会阻塞整个会话窗口，窗口关闭或者按下`CTRL + C`则Redis停止。不推荐使用。**

- **安装完成后，在任意目录输入`redis-server`命令即可启动Redis**

  ```
  redis-server
  ```

### 4.2 后台启动（不推荐）

------

> **如果要让Redis以后台方式启动，则必须修改Redis配置文件，配置文件所在目录就是之前我们解压的安装包下**

- **因为我们要修改配置文件，因此我们需要先将原文件备份一份**

  ```sh
  cd /usr/local/src/redis-6.2.6
  ```

  ```sh
  cp redis.conf redis.conf.bck
  ```

- **然后修改`redis.conf`文件中的一些配置**

  ```sh
  # 允许访问的地址，默认是127.0.0.1，会导致只能在本地访问。修改为0.0.0.0则可以在任意IP访问，生产环境不要设置为0.0.0.0
  bind 0.0.0.0
  # 守护进程，修改为yes后即可后台运行
  daemonize yes 
  # 密码，设置后访问Redis必须输入密码
  requirepass 1325
  ```

- **Redis其他常用配置**

  ```sh
  # 监听的端口
  port 6379
  # 工作目录，默认是当前目录，也就是运行redis-server时的命令，日志、持久化等文件会保存在这个目录
  dir .
  # 数据库数量，设置为1，代表只使用1个库，默认有16个库，编号0~15
  databases 1
  # 设置redis能够使用的最大内存
  maxmemory 512mb
  # 日志文件，默认为空，不记录日志，可以指定日志文件名
  logfile "redis.log"
  ```

- **启动Redis**

  ```sh
  # 进入redis安装目录 
  cd /usr/local/src/redis-6.2.6
  # 启动
  redis-server redis.conf
  ```

- **停止Redis服务**

  ```sh
  # 通过kill命令直接杀死进程
  kill -9 redis进程id
  ```

  ```sh
  # 利用redis-cli来执行 shutdown 命令，即可停止 Redis 服务，
  # 因为之前配置了密码，因此需要通过 -a 来指定密码
  redis-cli -a 132537 shutdown
  ```

### 4.3 开机自启（推荐）

------

> **我们也可以通过配置来实现开机自启**

- **首先，新建一个系统服务文件**

  ```sh
  vi /etc/systemd/system/redis.service
  ```

- **将以下命令粘贴进去**

  ```sh
  [Unit]
  Description=redis-server
  After=network.target
  
  [Service]
  Type=forking
  ExecStart=/usr/local/bin/redis-server /usr/local/src/redis-6.2.6/redis.conf
  PrivateTmp=true
  
  [Install]
  WantedBy=multi-user.target
  ```

- **然后重载系统服务**

  ```sh
  systemctl daemon-reload
  ```

- **现在，我们可以用下面这组命令来操作redis了**

  ```sh
  # 启动
  systemctl start redis
  # 停止
  systemctl stop redis
  # 重启
  systemctl restart redis
  # 查看状态
  systemctl status redis
  ```

- **执行下面的命令，可以让redis开机自启**

  ```sh
  systemctl enable redis
  ```



## 5. Redis客户端

> 安装完成Redis，我们就可以操作Redis，实现数据的CRUD了。这需要用到Redis客户端，包括：

- **命令行客户端**
- **图形化桌面客户端**
- **编程客户端**

### 1.命令行客户端

- **Redis安装完成后就自带了命令行客户端：`redis-cli`，使用方式如下：**

  ```sh
  redis-cli [options] [commonds]
  ```

- **其中常见的`options`有：**

  - `-h 127.0.0.1`：指定要连接的redis节点的IP地址，默认是127.0.0.1
  - `-p 6379`：指定要连接的redis节点的端口，默认是6379
  - `-a 132537`：指定redis的访问密码

- **其中的`commonds`就是Redis的操作命令，例如：**

  - `ping`：与redis服务端做心跳测试，服务端正常会返回`pong`
  - 不指定commond时，会进入`redis-cli`的交互控制台：
  - 

### 2.图形化客户端

> Windows 版: https://github.com/lework/RedisDesktopManager-Windows
>
> Mac 版：https://macwk.com/soft/redis-desktop-manager

- **安装图形化客户端**

  ```
  安装步骤过于简单不再演示
  ```



# 二、Redis常见命令

> 我们可以通过Redis的中文文档：http://www.redis.cn/commands.html，来学习各种命令。
>
> 也可以通过菜鸟教程官网来学习：https://www.runoob.com/redis/redis-keys.html

## 1.Redis数据结构介绍

> **Redis是一个key-value的数据库，key一般是String类型，不过value的类型多种多样**

<img src="./imgs/Redis value type.png" alt="AAC规格" style="zoom:100%;" />

## 2.通用命令

> **通用指令是部分数据类型的，都可以使用的指令，常见的有如下表格所示**

|  指令  |                        描述                        |
| :----: | :------------------------------------------------: |
|  KEYS  | 查看符合模板的所有key，不建议在生产环境设备上使用  |
|  DEL   |                 删除一个指定的key                  |
| EXISTS |                  判断key是否存在                   |
| EXPIRE | 给一个key设置有效期，有效期到期时该key会被自动删除 |
|  TTL   |              查看一个KEY的剩余有效期               |

**可以通过`help [command] `可以查看一个命令的具体用法！**（redis-cli 中 help @set 就可以看到相关的命令 tab切换）

## 3.String类型

> **String类型，也就是字符串类型，是Redis中最简单的存储类型。**

其value是字符串，不过根据字符串的格式不同，又可以分为3类：

- `string`：普通字符串
- `int`：整数类型，可以做自增、自减操作
- `float`：浮点类型，可以做自增、自减操作

> 不管是哪种格式，底层都是字节数组形式存储，只不过是编码方式不同。字符串类型的最大空间不能超过**512m**.

|  KEY  |    VALUE    |
| :---: | :---------: |
|  msg  | hello world |
|  num  |     10      |
| score |    92.5     |

> **String的常见命令有如下表格所示**

|    命令     |                             描述                             |
| :---------: | :----------------------------------------------------------: |
|     SET     |         添加或者修改已经存在的一个String类型的键值对         |
|     GET     |                 根据key获取String类型的value                 |
|    MSET     |                批量添加多个String类型的键值对                |
|    MGET     |             根据多个key获取多个String类型的value             |
|    INCR     |                     让一个整型的key自增1                     |
|   INCRBY    | 让一个整型的key自增并指定步长，例如：incrby num 2 让num值自增2 |
| INCRBYFLOAT |              让一个浮点类型的数字自增并指定步长              |
|    SETNX    | 添加一个String类型的键值对，前提是这个key不存在，否则不执行  |
|  **SETEX**  |          添加一个String类型的键值对，并且指定有效期          |

> **Redis的key允许有多个单词形成层级结构，多个单词之间用” ：“隔开，格式如下：**

```
项目名:业务名:类型:id
tex
```

这个格式并非固定，也可以根据自己的需求来删除或添加词条。

例如我们的项目名称叫 `heima`，有`user`和`product`两种不同类型的数据，我们可以这样定义key：

- **user**相关的key：`heima:user:1`
- **product**相关的key：`heima:product:1`

如果Value是一个Java对象，例如一个User对象，则可以将对象序列化为JSON字符串后存储

|       KEY       |                   VALUE                   |
| :-------------: | :---------------------------------------: |
|  heima:user:1   |    {“id”:1, “name”: “Jack”, “age”: 21}    |
| heima:product:1 | {“id”:1, “name”: “小米11”, “price”: 4999} |