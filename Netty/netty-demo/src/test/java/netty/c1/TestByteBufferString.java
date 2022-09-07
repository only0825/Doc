package netty.c1;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;

import static netty.c1.ByteBufferUtil.debugAll;

public class TestByteBufferString {

    public static void main(String[] args) {
        // 字符串转为 ByteBuffer 3种方法
        // 1. 原始方法
        ByteBuffer buffer1 = ByteBuffer.allocate(16);
        buffer1.put("hello".getBytes());
        debugAll(buffer1);
        // 2. Charset
        ByteBuffer buffer2 = StandardCharsets.UTF_8.encode("hello"); // 会自动转为
        debugAll(buffer2);
        // 3. wrap
        ByteBuffer buffer3 = ByteBuffer.wrap("hello".getBytes());
        debugAll(buffer3);


        // ByteBuffer 转字符串
        // 2 和 3 都是直接切换到读模式 再转可以不用filp()切换到读模式
        String str1 = StandardCharsets.UTF_8.decode(buffer2).toString();
        System.out.println(str1);
        // 1方法再转就要用filp()转成读模式
        buffer1.flip();
        String str2 = StandardCharsets.UTF_8.decode(buffer1).toString();
        System.out.println(str2);
    }
}
