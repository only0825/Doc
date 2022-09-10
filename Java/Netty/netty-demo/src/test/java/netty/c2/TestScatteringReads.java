package netty.c2;

import java.io.IOException;
import java.io.RandomAccessFile;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;

import static netty.c2.ByteBufferUtil.debugAll;

public class TestScatteringReads {

    /*
     需求：读取words.txt中的onetwothree，并输出one、two、three
     读取的两种思路
        1. 将数据存入一个ByteBuffer，之后再用别的方法拆分成三组，涉及到数据重新的分割复制
        2. 读取时一次读到3个ByteBuffer (分散读取)
     */
    public static void main(String[] args) {
        // 分散读取
        try (FileChannel channel = new RandomAccessFile("words.txt", "r").getChannel()) {
            ByteBuffer b1 = ByteBuffer.allocate(3);
            ByteBuffer b2 = ByteBuffer.allocate(3);
            ByteBuffer b3 = ByteBuffer.allocate(5);
            channel.read(new ByteBuffer[]{b1, b2, b3});
            b1.flip();
            b2.flip();
            b3.flip();
            debugAll(b1);
            debugAll(b2);
            debugAll(b3);
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
