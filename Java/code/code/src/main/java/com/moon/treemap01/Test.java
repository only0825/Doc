package com.moon.treemap01;

import java.util.Comparator;
import java.util.StringJoiner;
import java.util.TreeMap;
import java.util.function.BiConsumer;

public class Test {

    public static void main(String[] args) {

        String s = "aababcabcdabcde";

        TreeMap<Character, Integer> tm = new TreeMap<>();
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (tm.containsKey(c)) {
                int count = tm.get(c);
                count++;
                tm.put(c, count);
            } else {
                tm.put(c, 1);
            }
        }

        StringBuilder sb = new StringBuilder();
        tm.forEach((key, value) -> sb.append(key).append("(").append(value).append(")"));
        System.out.println(sb);

        StringJoiner sj = new StringJoiner("", "", "");
        tm.forEach((key, value) -> sj.add(key + "").add("(").add(value + "").add(")"));
        System.out.println(sj);
    }
}
