package com.moon.function;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.function.Function;

public class Pra2 {

    public static void main(String[] args) {

        ArrayList<Student> list = new ArrayList<>();
        list.add(new Student("张三", 23));
        list.add(new Student("李四", 24));
        list.add(new Student("王武", 25));

        String[] arr = list.stream().map(Student::getName).toArray(String[]::new);

        System.out.println(Arrays.toString(arr));
    }
}
