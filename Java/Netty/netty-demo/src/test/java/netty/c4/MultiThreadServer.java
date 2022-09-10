package netty.c4;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.Iterator;

public class MultiThreadServer {

    public static void main(String[] args) throws IOException {
        Thread.currentThread().setName("boss");

        ServerSocketChannel ssc = ServerSocketChannel.open();
        ssc.configureBlocking(false);
        ssc.bind(new InetSocketAddress(8880));

        Selector boss = Selector.open();
        ssc.register(boss, SelectionKey.OP_ACCEPT, null);

        while (true) {
            boss.select();
            Iterator<SelectionKey> iter = boss.selectedKeys().iterator();
            if (iter.hasNext()) {
                SelectionKey key = iter.next();
                iter.remove();
                if (key.isAcceptable()) {
                    SocketChannel sc = ssc.accept();
                    sc.configureBlocking(false);
                }
            }
        }
    }

    class Worker implements Runnable{
        private Thread thread;
        private Selector worker;
        private String name;

        public Worker(String name) {
            this.name = name;
        }

        // 初始化线程 和 selector
        public void register() throws IOException {
             thread = new Thread(this, name);
             thread.start();
             worker = Selector.open();
        }

        @Override
        public void run() {
            while (true) {
                try {
                    worker.select();
                    Iterator<SelectionKey> iter = worker.selectedKeys().iterator();
                    while (iter.hasNext()) {
                        SelectionKey key = iter.next();
                        iter.remove();
                        if (key.isReadable()) {

                        }
                    }
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        }
    }
}
