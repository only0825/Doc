package netty.c2;

import lombok.extern.slf4j.Slf4j;

import java.io.FileInputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;

@Slf4j // Slf4j这个是日志记录器
public class TestByteBuffer {

    /*
        ByteBuffer 正确使用姿势
        1. 向 buffer 写入数据，例如调用 channel.read(buffer)
        2. 调用 flip() 切换至**读模式**
        3. 从 buffer 读取数据，例如调用 buffer.get()
        4. 调用clear() 或 compact() 切换至**写模式**
        5. 重复 1~4 步骤
     */
    public static void main(String[] args) {
        // FileChannel
        // 1. 输入输出流   2. RandomAccessFile
        try (FileChannel channel = new FileInputStream("data.txt").getChannel()) {
            // 准备缓冲区
            ByteBuffer buffer = ByteBuffer.allocate(10); // 缓冲区要占用内存空间 不能无限大 所以要分多次读取
            while (true) {
                // 从 channel 读取数据， 向 buffer 写入
                int len = channel.read(buffer);
                log.debug("读取到的字节数 {}", len);
                if (len == -1) { // 没有内容了
                    break;
                }
                // 打印 buffer 的内容
                buffer.flip(); // 切换至读模式
                while (buffer.hasRemaining()) { // 是否还有剩余未读数据
                    byte b = buffer.get();
                    log.debug("实际字节 {}", (char) b);
                }
                buffer.clear(); // 切换为写模式
            }
        } catch (IOException ignored) {
        }
    }
}
