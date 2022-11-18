package com.moon;


import com.sun.org.apache.xerces.internal.impl.dv.util.Base64;

public class Demo {

    public static void main(String[] args) {
        String base_data = "conv_time=1668768912201&client_ip=127.0.0.1&ua=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36&sign=962f27d4963026c8b3118c7f09d3bb9a";
        String encrypt_key = "QeRRNxxlsSlFjgJD";
        String res = encrypt(base_data, encrypt_key);
        System.out.println(res);
    }

    public static String encrypt(String base_data, String key) {
        try {
            if (base_data.isEmpty() || key.isEmpty()) {
                return null;
            }

            char[] infoChar = base_data.toCharArray();
            char[] keyChar = key.toCharArray();

            byte[] resultChar = new byte[infoChar.length];
            for (int i = 0; i < infoChar.length; i++) {
                resultChar[i] = (byte) ((infoChar[i] ^ keyChar[i % keyChar.length]) & 0xFF);
            }
            return Base64.encode(resultChar);
        } catch (Exception e) {
            return null;
        }
    }
}
