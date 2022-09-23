package com.moon;

import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.ChannelHandlerContext;
import io.netty.channel.ChannelInboundHandlerAdapter;
import io.netty.channel.ChannelInitializer;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.channel.socket.nio.NioSocketChannel;
import io.netty.handler.codec.string.StringDecoder;

public class Test22 {

    public static void main(String[] args) {
        // 1. 启动器 负责组装 netty 组件，启动服务器
        new ServerBootstrap()
                // 2. 创建 NioEventLoopGroup，可以简单理解为 线程池 + Selector
                .group(new NioEventLoopGroup())
                // 3. 选择 服务器的 ServerSocketChannel 实现
                .channel(NioServerSocketChannel.class) // OIO
                // 4. boss 负责处理连接 worker(child) 负责处理读写 该方法决定了 child 能执行哪些操作（handler）
                .childHandler(
                        // 5. channel 代表和客户端进行数据数据读写的通道 Initializer 初始化，负责添加别的 handler
                        new ChannelInitializer<NioSocketChannel>() {
                            @Override
                            protected void initChannel(NioSocketChannel ch) throws Exception {
                                // 6. 添加具体 handler
                                ch.pipeline().addLast(new StringDecoder()); // 将 ByteBuf 转换为字符串
                                ch.pipeline().addLast(new ChannelInboundHandlerAdapter() { // 自定义 handler
                                    @Override
                                    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
                                        System.out.println(msg);
                                    }
                                });
                            }
                        }
                )
                // 6. 绑定监听端口
                .bind(8880);
    }
}