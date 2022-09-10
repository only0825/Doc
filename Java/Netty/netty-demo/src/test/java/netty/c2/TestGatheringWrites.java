package netty.c2;

import java.io.IOException;
import java.io.RandomAccessFile;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;
import java.nio.charset.StandardCharsets;

public class TestGatheringWrites {

    /*
        需求：三个ByteBuffer写入到一个文件中
        思路：
            1. 三个ByteBuffer组合到一个大的ByteBuffer中，涉及到数据到多次拷贝
            2. 三个ByteBuffer组合到一起，以一个整体写入
     */
    public static void main(String[] args) {

        ByteBuffer b1 = StandardCharsets.UTF_8.encode("hello");
        ByteBuffer b2 = StandardCharsets.UTF_8.encode("world");
        ByteBuffer b3 = StandardCharsets.UTF_8.encode("你好");

        try (FileChannel channel = new RandomAccessFile("words2.txt", "rw").getChannel()) {
            channel.write(new ByteBuffer[]{b1, b2, b3});
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
