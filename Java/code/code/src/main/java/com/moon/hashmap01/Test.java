package com.moon.hashmap01;

import java.util.*;
import java.util.function.BiConsumer;

public class Test {

    public static void main(String[] args) {

//        Map<Student, String> students = new HashMap<>();
//
//        students.put(new Student("张三", 11), "四川");
//        students.put(new Student("李四", 13), "广西");
//        students.put(new Student("王武", 16), "新疆");
//
//        students.forEach((student, s) -> System.out.println(student.getName() + " 年龄" + student.getAge() + " 籍贯:" + s));

//        Map<String, String> map = new HashMap<>();
//        String[] attractions = {"A", "B", "C", "D"};
//
//        for (int i = 0; i < 80; i++) {
//            int index = (int) (Math.random() * attractions.length);
//            String rand = attractions[index];
//            map.put("学生" + i, rand);
//        }
//
//        Map<String, Integer> map2 = new HashMap<>();
//        map2.put("A", 0);
//        map2.put("B", 0);
//        map2.put("C", 0);
//        map2.put("D", 0);
//
//        map.forEach((name, s) -> map2.forEach((key, count) -> {
//            if (s.equals(key)) {
//                map2.put(key, count + 1);
//            }
//        }));
//
//        map2.forEach((s, integer) -> System.out.println(s + " " + integer));

        Map<String, Integer> map = new HashMap<>();
        String[] arr = {"A", "B", "C", "D"};
        for (int i = 0; i < 80; i++) {
            int index = (int) (Math.random() * arr.length);
            String rand = arr[index];
            boolean b = map.containsKey(rand);
            if (b) {
                int count = map.get(rand);
                map.put(rand, count + 1);
            } else {
                map.put(rand, 1);
            }
        }

        map.forEach(new BiConsumer<String, Integer>() {
            @Override
            public void accept(String s, Integer integer) {
                System.out.println(s + ":" + integer);
            }
        });

        ArrayList<Integer> list = new ArrayList<>();
        Collections.addAll(list, 1, 2, 3, 4, 5);
        Collections.fill(list, 100);
        System.out.println(list);
    }
}
