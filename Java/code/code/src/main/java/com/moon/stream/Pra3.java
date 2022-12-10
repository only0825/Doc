package com.moon.stream;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Predicate;
import java.util.stream.Collectors;
import java.util.stream.Stream;

public class Pra3 {

    public static void main(String[] args) {

        ArrayList<String> list1 = new ArrayList<>();
        ArrayList<String> list2 = new ArrayList<>();

        Collections.addAll(list1, "张三,23", "李四4,24", "李大贵,25", "王武,26", "万福贵,33", "王健林,44");
        Collections.addAll(list2, "杨红,18", "杨芳,19", "小玉,20", "王熙凤,24", "林黛玉,21", "杨翠花,33");

        Stream<String> stream1 = list1.stream().limit(2).filter(s -> s.split(",")[0].length() == 3);
        Stream<String> stream2 = list2.stream().skip(1).filter(s -> s.split(",")[0].startsWith("杨"));

        List<Actor> list = Stream.concat(stream1, stream2)
                .map(s -> new Actor(s.split(",")[0], Integer.parseInt(s.split(",")[1])))
                .collect(Collectors.toList());

        System.out.println(list);
    }
}
