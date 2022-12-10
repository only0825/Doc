package com.moon.pra;

import java.util.ArrayList;
import java.util.Collections;
import java.util.Random;

public class Test3 {
    public static void main(String[] args) {

        ArrayList<String> list1 = new ArrayList<>();
        Collections.addAll(list1, "大王", "小王", "老马", "光头强", "熊大");
        ArrayList<String> list2 = new ArrayList<>();


        for (int i = 1; i < 10; i++) {
            System.out.println("=========" + i + "===========");
            int count = list1.size();
            Random r = new Random();
            for (int j = 0; j < count; j++) {
                int index = r.nextInt(list1.size());
                String name = list1.remove(index);
                list2.add(name);
                System.out.println(name);
            }
            list1.addAll(list2);
            list2.clear();
        }

    }
}
