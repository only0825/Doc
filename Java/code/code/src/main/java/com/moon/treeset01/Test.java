package com.moon.treeset01;

import java.util.Comparator;
import java.util.HashSet;
import java.util.Set;
import java.util.TreeSet;
import java.util.function.Consumer;

public class Test {

    public static void main(String[] args) {

//        TreeSet<Student> set = new TreeSet<>();
//        Student s6 = new Student("张6", 15);
//        Student s3 = new Student("张三", 11);
//        Student s4 = new Student("张4", 12);
//        Student s5 = new Student("张5", 14);
//        Student s7 = new Student("张7", 5);
//        set.add(s3);
//        set.add(s4);
//        set.add(s5);
//        set.add(s6);
//        set.add(s7);
//        System.out.println(set);
        HashSet<String> hashSet = new HashSet<>();

        TreeSet<String> ts = new TreeSet<>((o1, o2) -> {
            int i = o1.length() - o2.length();
            i = i == 0 ? o1.compareTo(o2) : i;
            return i;
        });
        ts.add("c");
        ts.add("ab");
        ts.add("df");
        ts.add("qwer");

        System.out.println(ts);
    }
}
