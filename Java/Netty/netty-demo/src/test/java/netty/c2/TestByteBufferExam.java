package netty.c2;

import java.nio.ByteBuffer;

import static netty.c2.ByteBufferUtil.debugAll;

public class TestByteBufferExam {

    // 解决粘包和半包问题，下面只是原始解法，要会，这样才能知道netty帮我们做了哪些事情
    public static void main(String[] args) {

        /*
            网络上有多条数据发送给服务端，数据之间使用 \n 进行分隔
            但由于某种原因这些数据在接收时，被进行了重新组合，例如原始数据有3条为
            - Hello,world\n
            - I’m Nyima\n
            - How are you?\n
            变成了下面的两个 byteBuffer (粘包，半包)
            - Hello,world\nI’m Nyima\nHo
            - w are you?\n
            现在要求你编写程序，将错误的数据恢复成原始按 \n 分割的数据
         */
        // 模拟粘包+半包
        ByteBuffer source = ByteBuffer.allocate(32);
        source.put("Hello,world\nI'm zhangsan\nHo".getBytes());
        split(source);
        source.put("w are you ?\n".getBytes());
        split(source);
    }

    public static void split(ByteBuffer source) {
        // 切换为读模式
        source.flip();
        for (int i = 0; i < source.limit(); i++) {
            // 找到一条完整消息
            if (source.get(i) == '\n') { // get(i)不会移动position
                int length = i + 1 - source.position(); // 当前位置 - 读指针位置
                // 把这条完整消息存入新的 ByteBuffers
                ByteBuffer target = ByteBuffer.allocate(length);
                // 从source 读，向 target 写
                for (int j = 0; j < length; j++) {
                    target.put(source.get()); // get 方法会让 position 读指针向后走
                }
                debugAll(target);
            }
        }
        // 切换为写模式，但是缓冲区可能未读完，这里需要使用compact
        source.compact();
    }
}
