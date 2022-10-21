
[TOC]
## ThreadLocal两大使用场景

* 典型场景1：每个线程需要一个独享的对象（通常是工具类，如SimpleDateForate和Random等）
* 典型场景2：每个线程内需要保存全局变了（例如在拦截器中获取用户信息），可以让不同方法直接使用，避免参数传递的麻烦

### 典型场景1 - 每个线程需要一个独享的对象

下面有4段代码，阐述了为什么要使用ThreadLocal

```java
/**
 * 描述：10个线程打印日期
 */
public class ThreadLocalNormalUsage01 {
    public static void main(String[] args) throws InterruptedException {
        // 问题：当有1000个或更多任务时，用for循环太难维护了
        for (int i = 0; i < 10; i++) {
            int finalI = i;
            new Thread(new Runnable() {
                @Override
                public void run() {
                    String date = new ThreadLocalNormalUsage01().date(finalI);
                    System.out.println(date);
                }
            }).start();
            Thread.sleep(100);
        }
    }
    public String date(int seconds) {
        //参数的单位是浩渺，葱1970.1.1 00:00:00 GMT计时
        Date date = new Date(1000L * seconds);
        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
        return dateFormat.format(date);
    }
}

```

```java
/**
 * 描述：1000个打印日期的任务，用线程池来执行
 */
public class ThreadLocalNormalUsage02 {

    public static ExecutorService threadPool = Executors.newFixedThreadPool(10);

    public static void main(String[] args) throws InterruptedException {
        for (int i = 0; i < 1000; i++) {
            int finalI = i;
            threadPool.submit(new Runnable() {
                @Override
                public void run() {
                    //问题： 这里创建和销毁了一千次SimpleDateFormat对象 这是不必要的开销
                    String date = new ThreadLocalNormalUsage02().date(finalI);
                    System.out.println(date);
                }
            });
        }
        threadPool.shutdown();
    }
    public String date(int seconds) {
        //参数的单位是浩渺，葱1970.1.1 00:00:00 GMT计时
        Date date = new Date(1000L * seconds);
        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
        return dateFormat.format(date);
    }
}
```

```java
/**
 * 描述：加锁来解决线程安全问题
 */
public class ThreadLocalNormalUsage04 {

    public static ExecutorService threadPool = Executors.newFixedThreadPool(10);
    static SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");

    public static void main(String[] args) throws InterruptedException {
        // 问题：会出现相同当日期，这是线程安全问题
        for (int i = 0; i < 1000; i++) {
            int finalI = i;
            threadPool.submit(new Runnable() {
                @Override
                public void run() {
                    String date = new ThreadLocalNormalUsage04().date(finalI);
                    System.out.println(date);
                }
            });
        }
        threadPool.shutdown();
    }
    public String date(int seconds) {
        //参数的单位是浩渺，葱1970.1.1 00:00:00 GMT计时
        Date date = new Date(1000L * seconds);
        String s = null;
        // 类锁的方式加锁，同一个实例调用会阻塞
        // 问题：利用synchronized加锁后，线程只能一个个排队，效率低
        // 解决：利用ThreadLocal
        synchronized (ThreadLocalNormalUsage04.class) {
            s = dateFormat.format(date);
        }
        return s;
    }
}
```

```java
/**
 * 描述：利用ThreadLocal，给每个线程分配自动dateFormat对象，保证了线程安全，高效利用内存
 */
public class ThreadLocalNormalUsage05 {
    public static ExecutorService threadPool = Executors.newFixedThreadPool(10);
    public static void main(String[] args) throws InterruptedException {
        for (int i = 0; i < 1000; i++) {
            int finalI = i;
            threadPool.submit(new Runnable() {
                @Override
                public void run() {
                    String date = new ThreadLocalNormalUsage05().date(finalI);
                    System.out.println(date);
                }
            });
        }
        threadPool.shutdown();
    }
    public String date(int seconds) {
        //参数的单位是毫秒，葱1970.1.1 00:00:00 GMT计时
        Date date = new Date(1000L * seconds);
//        SimpleDateFormat dateFormat = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
        SimpleDateFormat dateFormat = ThreadSafeFormatter.dateFormatThreadLocal.get();
        return dateFormat.format(date);

    }
}
class ThreadSafeFormatter {
    public static ThreadLocal<SimpleDateFormat> dateFormatThreadLocal
            = new ThreadLocal<SimpleDateFormat>() {
        @Override
        protected SimpleDateFormat initialValue() {
            return new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
        }
    };
    // Lambda表达式写法
    public static ThreadLocal<SimpleDateFormat> dateFormatThreadLocal2
            = ThreadLocal.withInitial(() -> new SimpleDateFormat("yyyy-MM-dd hh:mm:ss"));

}
```

### 什么是ConcurrentHashMap？

ConcurrentHashMap 是Java集合中map的实现，是HashMap的线程安全版本，性能也比较好。

 ConcurrentHashMap在数据结构上和HashMap的数据结构是一致的，区别在于ConcurrentHashMap是线程安全的，而HashMap不是线程安全的。

 在需要用到HashMap且是多线程的情况下，推荐使用ConcurrentHashMap。



### 典型场景2 - 避免传递参数的麻烦

使用Thread Local就无需synchronized，不会影响性能

```java
package com.moon.thread.threadlocal;

/**
 * 描述： 演示ThreadLocal的用法2 - 避免传递参数的麻烦
 */
public class ThreadLocalNormalUsage06 {
    public static void main(String[] args) {
        new Service1().process();
    }
}

class Service1 {
    public void process() {
        User user = new User("超哥");
        UserContextHolder.holder.set(user);
        new Service2().process();
    }
}
class Service2 {
    public void process() {
        User user = UserContextHolder.holder.get();
        System.out.println("Service2" + user.name);
        new Service3().process();
    }
}
class Service3 {
    public void process() {
        User user = UserContextHolder.holder.get();
        System.out.println("Service3" + user.name);
    }
}

class UserContextHolder {
    public static ThreadLocal<User> holder = new ThreadLocal<>();
}

class User {
    String name;
    public User(String name) {
        this.name = name;
    }
}		
```



**练习demo:** 

ThreadLocal往每个线程中保存String类型数据，再获得保存的数据并打印输出信息。运行效果如下图所示：

```java
package com.moon.thread.threadlocal;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class TokenContextHolder {
    public static ThreadLocal<String> holder = new ThreadLocal<>();
}
class GetToken {
    public void getToken() {
        String str = TokenContextHolder.holder.get();
        System.out.println("GetToken " + str);
    }
}
class TokenUtil {
    public void getTokenStr() {
        String str = TokenContextHolder.holder.get();
        System.out.println("TokenUtil " + str);
        new GetToken().getToken();
    }
}
class Test {
    public static ExecutorService threadPool = Executors.newFixedThreadPool(10);
    public static void main(String[] args) {
        for (int i = 0; i < 10; i++) {
            int num = i;
            threadPool.submit(new Runnable() {
                @Override
                public void run() {
                    TokenContextHolder.holder.set("token" + num);
                    new TokenUtil().getTokenStr();
                }
            });
        }
    }
}
```



## ThreadLocal的作用和主要方法

### ThreadLocal的两个作用和好处

#### ThreadLocal的两个作用

1、让某个需要用到的对象在线程间隔离（每个线程都拥有自己的独立的对象）

2、在任何方法中都可以轻松获取到该对象

#### 根据共享对象的生成时机不同，选择initialValue或set来保存对象

**initialValue**: 

* 在ThreadLocal**第一次get**的时候把对象给初始化处理，对象的初始化时机可以**由我们控制**

**set**：

* 如果需要保存到ThreadLocal里的对象的生成时机**不由我们随意控制**，例如拦截器生成的用户信息
* 用ThreadLocal.set直接放到我们的ThreadLocal中去，以便后续使用。

#### 使用ThreadLocal带来的好处

1. 达到线程安全
2. 不需要加锁，提高效率
3. 更高效的利用内存，节省开销（相比每个任务都新建一个SimpleDateForma，显然用ThreadLocal更好）
4. 免去传参的麻烦（不需要每次都传相同的参数，ThreadLocal使得代码耦合度更低、更优雅）

### ThreadLocal的主要方法

#### **T initialValue()**

1. 该方法会返回当前线程对应的“初始值”，这是一个**延迟加载**的方法，只有在**调用get**的时候，才会触发
2. 当线程**第一次使用get**方法访问变量时，将调用此方法
3. 每个线程最多调用**一次**此方法，但如果已经调用了remove()后，再调用get()，则可以再次调用此方法。
4. 如果不重写本方法，这个方法会返回null。一般使用匿名内部类的方法来**重写initialValue()**方法

#### **void set(T t)**

* 为这个线程设置一个新值

#### **T get()**

* 得到这个线程对应的value。如果是首次调用get()，则会调用initialize来得到这个值。

#### **void remove()**

* 删除对应这个线程的值



#### 方法运用的demo

```java
package com.moon.thread.threadlocal;

import java.util.UUID;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class UUIDContextHolder {
    public static ThreadLocal<String> holder = new ThreadLocal<>();
}

class GetUUID {
    public void getUuid() {
        UUIDContextHolder.holder.get();
    }
}

class UpdateUUID {
    public void update() {
        // 获取数据
        String uuid = UUIDContextHolder.holder.get();
        System.out.println("UpdateUUID更新前，拿到" + uuid);
        // 删除之前数据
        UUIDContextHolder.holder.remove();
        // 重新赋值
        UUIDContextHolder.holder.set("最新");
        // 获取更新后的值
        new GetUUID().getUuid();
    }
}

class PutUUID {
    public void put(int num) {
        String uuid = UUID.randomUUID().toString().replaceAll("-", "") + "---" + num;
        UUIDContextHolder.holder.set(uuid);
        new UpdateUUID().update();
        String newUUID = UUIDContextHolder.holder.get();
        System.out.println("GetUUID拿到" + uuid + newUUID);
    }
}

class Test2 {
    public static void main(String[] args) {
        for (int i = 0; i < 10; i++) {
            int finalI = i;
            new Thread(new Runnable() {
                @Override
                public void run() {
                    new PutUUID().put(finalI);
                }
            }).start();
        }
    }
}
```

## ThreadLocal的原理和注意点

### 图解ThreadLocal原理

要搞清楚Thread的原理必须清楚Thread、ThreadLocal以及ThreadMap之前的关系

如图：

![image-20210801000031327](/Users/sofia/Library/Application Support/typora-user-images/image-20210801000031327.png)

****每个Thread对象中都持有一个ThreadLocalMap成员变量

