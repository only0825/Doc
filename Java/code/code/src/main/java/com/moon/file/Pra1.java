package com.moon.file;

import java.io.File;

public class Pra1 {

    public static void main(String[] args) {

        File file = new File("./aaa");
        boolean b = file.mkdir();
        System.out.println(b);
    }
}
