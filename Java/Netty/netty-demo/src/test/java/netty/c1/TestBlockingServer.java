package netty.c1;

import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.ArrayList;
import java.util.List;

import static netty.c2.ByteBufferUtil.debugRead;

@Slf4j // 单线程 阻塞模式
public class TestBlockingServer {

    public static void main(String[] args) throws IOException {
        // 0. ByteBuffer
        ByteBuffer buffer = ByteBuffer.allocate(16);
        // 1. 创建服务器
        ServerSocketChannel ssc = ServerSocketChannel.open();
        // 2. 绑定监听端口
        ssc.bind(new InetSocketAddress(8880));
        // 3. 连接集合
        List<SocketChannel> channels = new ArrayList<>();

        while (true) {
            // 4. accept 建立与客户端连接，SocketChannel 用来与客户端直接通信
            log.debug("connecting...");
            SocketChannel sc = ssc.accept(); // 阻塞方法，线程停止运行
            log.debug("connected... {}", sc);
            channels.add(sc);
            for (SocketChannel channel : channels) {
                // 5. 接收客户端发送的消息
                log.debug("before read... {}", channel);
                channel.read(buffer); // 阻塞方法，线程停止运行; 当通道中没有数据可读时，会阻塞线程
                buffer.flip(); // 切换为读模式
                debugRead(buffer);
                buffer.clear(); // 切换为写模式
                log.debug("after read...{}", channel);
            }
        }
    }
}
