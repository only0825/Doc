package thread;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class Test2 {

    public static void main(String[] args) {

        // 创建任务对象
        Runnable task2 = new Runnable() {
            @Override
            public void run() {
                log.debug("hello");
            }
        };

        Thread t2 = new Thread(task2, "ts");
        t2.start();

    }
}
