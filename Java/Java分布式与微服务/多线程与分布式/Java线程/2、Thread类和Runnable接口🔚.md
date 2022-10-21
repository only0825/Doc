[TOC]

### 创建多线程的三种方式：

1、继承 Thread 类，重写run()方法，run()方法代表线程要执行的任务。
2、实现 Runnable 接口，重写 run()方法，run()方法代表线程要执行的任务。
3、实现 callable 接口，重写 call()方法，call()作为线程的执行体，具有返回值，并且可以对异常进行声明和抛出；使用start()方法来启动线程

#### 第一种方法：Thread类

Thread是一个线程类，位于java.lang包下
**构造方法：**

* Thread()：创建一个线程对象
* Thread(String name)：创建一个具有指定名称的线程对象
* Thread(Rummable target)：创建一个基于Runnable接口实现类的线程对象
* Thread(Runnable target,String name)：创建一个基于Runnable接口实现类，并且具有指定名称的线程对象。

**Thread类的常用方法：**

* public void run()：线程相关的代码写在该方法中，一般需要重写。
* public void start()：启动线程的方法
* public static void sleep(long m)：线程休眠m毫秒的方法
* public void join()：优先执行调用join()方法的线程。

#### 第二种方法：Runnable接口

* 只有一个方法run();
* Runnable是Java中用以实现线程的接口
* 任何实现线程功能的类都必须实现该接口

#### 第三种方法：callable 接口

1. 创建Callable 接口的实现类，并实现 call()方法，该 call()方法将作为线程执行体，并且有返回值。
2. 创建Callable 实现类的实例，使用 FutureTask 类来包装Callable 对象，该 FutureTask 对象封装了该Callable 对象的 call()方法的返回值。
3. 使用FutureTask 对象作为 Thread 对象的target 创建并启动新线程。
4. 调用FutureTask 对象的 get()方法来获得子线程执行结束后的返回值。 