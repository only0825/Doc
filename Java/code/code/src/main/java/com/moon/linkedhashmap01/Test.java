package com.moon.linkedhashmap01;

import java.util.LinkedHashMap;

public class Test {

    public static void main(String[] args) {

        LinkedHashMap<String, Integer> lhm = new LinkedHashMap<>();

        lhm.put("wang" , 11);
        lhm.put("xiao" , 12);
        lhm.put("che" , 14);

        System.out.println(lhm);
    }
}
