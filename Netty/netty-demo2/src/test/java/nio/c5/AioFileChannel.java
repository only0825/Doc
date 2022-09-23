package nio.c5;


import lombok.extern.slf4j.Slf4j;

import java.io.IOException;
import java.nio.ByteBuffer;
import java.nio.channels.AsynchronousFileChannel;
import java.nio.channels.CompletionHandler;
import java.nio.file.Paths;
import java.nio.file.StandardOpenOption;

import static nio.c2.ByteBufferUtil.debugAll;

@Slf4j
public class AioFileChannel {

    public static void main(String[] args) throws IOException {

        try {
            AsynchronousFileChannel channel = AsynchronousFileChannel.open(Paths.get("data.txt"), StandardOpenOption.READ);
            ByteBuffer buffer = ByteBuffer.allocate(16);
            log.error("read begin...");
            channel.read(buffer, 0, buffer, new CompletionHandler<Integer, ByteBuffer>() {
                @Override
                public void completed(Integer result, ByteBuffer attachment) {
                    log.debug("read completed...{}", attachment);
                    attachment.flip();
                    debugAll(attachment);
                }
                @Override
                public void failed(Throwable exc, ByteBuffer attachment) {
                    exc.printStackTrace();
                }
            });
            log.debug("read end...");
        } catch (IOException e) {
            e.printStackTrace();
        }
        System.in.read();
    }
}
