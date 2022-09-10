package thread;

import lombok.extern.slf4j.Slf4j;

@Slf4j(topic = "c.TestTread")
public class TestThread {
    public static void main(String[] args) {
        // 构造方法的参数是给线程指定名字，，推荐给线程起个名字
        Thread t1 = new Thread("t1") {
            @Override
            public void run() {  // run 方法内实现了要执行的任务
                log.debug("hello");
            }
        };
        t1.start();
        log.debug("running");
    }
}
