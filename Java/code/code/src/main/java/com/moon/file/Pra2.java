package com.moon.file;

import java.io.File;
import java.io.FileFilter;
import java.util.Arrays;

public class Pra2 {

    public static void main(String[] args) {

        File file = new File("./aaa");
        File[] arr = file.listFiles(pathname -> pathname.isFile() && pathname.getName().endsWith(".avi"));
        System.out.println(Arrays.toString(arr));
    }
}
