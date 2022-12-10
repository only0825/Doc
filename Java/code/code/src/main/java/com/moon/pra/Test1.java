package com.moon.pra;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Random;

public class Test1 {

    public static void main(String[] args) {

        ArrayList<String> list = new ArrayList<>();
        Collections.addAll(list, "Ben", "Jack", "Ale", "Boss", "Duke", "Zhu", "Gou");

        // 方法一：随机
        Random random = new Random();
        int i = random.nextInt(list.size());
        String name1 = list.get(i);
        System.out.println(name1);

        // 方法二：打乱
        Collections.shuffle(list);
        String name2 = list.get(0);
        System.out.println(name2);

    }
}
