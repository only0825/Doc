package com.moon.stream;

import java.util.ArrayList;
import java.util.Map;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collectors;

public class Pra2 {

    public static void main(String[] args) {

        ArrayList<String> list = new ArrayList<>();
        list.add("zhangsan, 23");
        list.add("lisi, 24");
        list.add("wangwu, 25");
        list.add("ben, 28");

        Map<String, Integer> map = list.stream()
                .filter(s -> Integer.parseInt(s.split(",")[1].trim()) >= 24)
                .collect(Collectors.toMap(s -> s.split(",")[0], s -> Integer.valueOf(s.split(",")[1].trim())));

        map.forEach((s, integer) -> System.out.println(s + " : " + integer));
    }
}
