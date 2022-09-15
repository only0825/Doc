package nio.c4;

import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.*;
import java.util.Iterator;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.atomic.AtomicInteger;

import static nio.c2.ByteBufferUtil.debugAll;

@Slf4j
public class MultiThreadServer {

    public static void main(String[] args) throws IOException {
        Thread.currentThread().setName("boss");
        ServerSocketChannel ssc = ServerSocketChannel.open(); // 这个Channel可以监听新来的连接
        ssc.configureBlocking(false); // 非阻塞
        ssc.bind(new InetSocketAddress(8880)); // 绑定8880 端口
        Selector boss = Selector.open(); // 创建Selector
        ssc.register(boss, SelectionKey.OP_ACCEPT, null); // 注册并绑定ACCEPT事件

        // 1. 创建固定数量的 worker 并初始化
        Worker[] workers = new Worker[2];
        for (int i = 0; i < workers.length; i++) {
            workers[i] = new Worker("worker-" + i);
        }
        AtomicInteger index = new AtomicInteger();
        while (true) {
            boss.select();
            Iterator<SelectionKey> iter = boss.selectedKeys().iterator();
            if (iter.hasNext()) {
                SelectionKey key = iter.next();
                iter.remove();
                if (key.isAcceptable()) {
                    SocketChannel sc = ssc.accept();
                    sc.configureBlocking(false);
                    log.debug("connected...{}", sc.getRemoteAddress());
                    // 2. 关联 selector
                    log.debug("before register...{}", sc.getRemoteAddress());
                    workers[index.getAndIncrement() % workers.length].register(sc); // boss 调用 初始化 selector ， 启动 worker-0
                    log.debug("after register...{}", sc.getRemoteAddress());
                }
            }
        }
    }

    static class Worker implements Runnable {
        private Thread thread;
        private Selector selector;
        private String name;
        // Volatile可以看做是轻量级的 Synchronized，它只保证了共享变量的可见性。在线程 A 修改被 volatile 修饰的共享变量之后，线程 B 能够读取到正确的值。
        private volatile boolean start = false; // 还未初始化
        private ConcurrentLinkedQueue<Runnable> queue = new ConcurrentLinkedQueue<>();

        public Worker(String name) {
            this.name = name;
        }

        // 问题：如果sc.register(selector, SelectionKey.OP_READ, null) 写在 Boss 线程，
        // 但 selector.select(); 是在work线程，而且他们用的同一个selector，
        // 当 selector.select(); 先执行就会阻塞，从而影响到sc.register
        // 下面有两种解决办法：队列方法 和 wakeup方法

        // 初始化线程 和 selector
        public void register(SocketChannel sc) throws IOException {
            if (!start) {
                selector = Selector.open();
                thread = new Thread(this, name);
                thread.start();
                start = true;
            }

            // 队列方法解决问题
            // 向队列添加了任务，但这个任务并没有立刻执行 boss
            queue.add(()->{
                try {
                    sc.register(selector, SelectionKey.OP_READ, null);
                } catch (ClosedChannelException e) {
                    e.printStackTrace();
                }
            });
            selector.wakeup(); // 唤醒 select 方法

            // wakeup方法解决问题  这个解法有问题 当第三个客户端连接时会阻塞
//            selector.wakeup(); // 唤醒 select 方法 boss
//            sc.register(selector, SelectionKey.OP_READ, null); // boss
        }

        @Override
        public void run() {
            while (true) {
                try {
                    selector.select(); // worker-0 阻塞， wakeup

                    // 队列方法解决问题
                    Runnable task = queue.poll();
                    if (task != null) {
                        task.run(); // 执行了 sc.register(selector, SelectionKey.OP_READ, null);
                    }

                    Iterator<SelectionKey> iter = selector.selectedKeys().iterator();
                    while (iter.hasNext()) {
                        SelectionKey key = iter.next();
                        iter.remove();
                        if (key.isReadable()) {
                            ByteBuffer buffer = ByteBuffer.allocate(16);
                            SocketChannel channel = (SocketChannel) key.channel();
                            log.debug("read...{}", channel.getRemoteAddress());
                            channel.read(buffer);
                            buffer.flip();
                            debugAll(buffer);
                        }
                    }
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
    }
}
