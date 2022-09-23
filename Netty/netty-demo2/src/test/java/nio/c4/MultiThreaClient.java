package nio.c4;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.nio.channels.SocketChannel;
import java.nio.charset.Charset;

public class MultiThreaClient {

    public static void main(String[] args) throws IOException {

        SocketChannel sc = SocketChannel.open();
        sc.connect(new InetSocketAddress("localhost", 8880));
        sc.write(Charset.defaultCharset().encode("1234567890abcdef"));
        System.in.read();
    }
}
