package com.moon.function;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;

public class Pra1 {

    public static void main(String[] args) {

        ArrayList<String> list = new ArrayList<>();
        Collections.addAll(list, "张三,23", "李四,24", "王武,25");

        Student[] arr = list.stream().map(Student::new).toArray(Student[]::new);

        System.out.println(Arrays.toString(arr));

    }
}
