package nio.thread;

import lombok.extern.slf4j.Slf4j;

import java.util.concurrent.Callable;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.FutureTask;

@Slf4j(topic = "c.TestTread")
public class Test1 {
    public static void main(String[] args) throws ExecutionException, InterruptedException {
        // 构造方法的参数是给线程指定名字，，推荐给线程起个名字
        Thread t1 = new Thread("t1") {
            @Override
            public void run() {  // run 方法内实现了要执行的任务
                log.debug("hello");
            }
        };
        t1.start();
        log.debug("running");

        // 创建任务对象
        Runnable task2 = new Runnable() {
            @Override
            public void run() {
                log.debug("hello");
            }
        };
        // 参数1 是任务对象； 参数2 是线程名字，推荐给线程起个名字
        Thread t2 = new Thread(task2, "ts");
        t2.start();


        // 实现多线程的第三种方法可以返回数据
        FutureTask futureTask = new FutureTask<>(new Callable<Integer>() {
            @Override
            public Integer call() throws Exception {
                log.debug("多线程任务");
                Thread.sleep(100);
                return 100;
            }
        });
        // 主线程阻塞，同步等待 task 执行完毕的结果
        new Thread(futureTask, "我的名字").start();
        log.debug("主线程");
        log.debug("{}", futureTask.get());
    }
}
