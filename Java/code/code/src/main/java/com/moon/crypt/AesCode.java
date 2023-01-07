package com.moon.crypt;

import org.apache.commons.codec.digest.DigestUtils;

import java.net.URLEncoder;
import java.security.GeneralSecurityException;
import java.util.Arrays;
import java.util.Base64;
import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

public class AesCode {

    public static byte[] encrypt(String key, byte[] origData) throws GeneralSecurityException {
        byte[] keyBytes = getKeyBytes(key);
        byte[] buf = new byte[16];
        System.arraycopy(keyBytes, 0, buf, 0, Math.max(keyBytes.length, buf.length));
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS7Padding");
        cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(buf, "AES"), new IvParameterSpec(keyBytes));
        return cipher.doFinal(origData);

    }

    public static byte[] decrypt(String key, byte[] crypted) throws GeneralSecurityException {
        byte[] keyBytes = getKeyBytes(key);
        byte[] buf = new byte[16];
        System.arraycopy(keyBytes, 0, buf, 0, Math.max(keyBytes.length, buf.length));
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS7Padding");
        cipher.init(Cipher.DECRYPT_MODE, new SecretKeySpec(buf, "AES"), new IvParameterSpec(keyBytes));
        return cipher.doFinal(crypted);
    }

    private static byte[] getKeyBytes(String key) {
        byte[] bytes = key.getBytes();
        return bytes.length == 16 ? bytes : Arrays.copyOf(bytes, 16);
    }

    public static String encrypt(String key, String val) {
        try {
            byte[] origData = val.getBytes();
            byte[] crafted = encrypt(key, origData);
            return Base64.getEncoder().encodeToString(crafted);
        }catch (Exception e){
            return "";
        }
    }

    public static String decrypt(String key, String val) throws GeneralSecurityException {
        byte[] crypted = Base64.getDecoder().decode(val);
        byte[] origData = decrypt(key, crypted);
        return new String(origData);
    }


    public static void main(String[] args) throws Exception {
        // 密钥
        String key = DigestUtils.md5Hex("ef14b146e989b922dcb5a00a19f907c5").substring(0, 16);

        // 加密
//        String val = "hello,ase";
//        String ret = encrypt(key, val);
//        System.out.println(ret);

        // 解密
        String text = "oDJZrCJy4/RZIPJC4hZ0heeLBGMNlZmQ+mTGrRdyMaxIiWA1Qw1xFHAtofda8FU4UbgoKpSOHxQK78fAwjfYbVPCbSwptgvX+QI0c/kGLqk5+fEQTZLNNDgGOL7gy6ez";
        String rr = decrypt(key, text);
        System.out.println(rr);
    }


}
