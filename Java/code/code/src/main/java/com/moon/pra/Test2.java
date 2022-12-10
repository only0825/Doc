package com.moon.pra;

import java.sql.Array;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Random;

public class Test2 {

    public static void main(String[] args) {

        ArrayList<Integer> list = new ArrayList<>();
        Collections.addAll(list, 1, 1, 1, 1, 1, 1, 1);
        Collections.addAll(list, 0, 0, 0);

        Random random = new Random();
        int index = random.nextInt(list.size());
        int number = list.get(index);

        ArrayList<String> boyList = new ArrayList<>();
        ArrayList<String> girlList = new ArrayList<>();
        Collections.addAll(boyList, "大王", "小王", "老马", "光头强", "熊大");
        Collections.addAll(girlList, "小芳", "小微", "小玉");

        if (number == 1) {
            int i = random.nextInt(boyList.size());
            String s = boyList.get(i);
            System.out.println(s);
        } else {
            int i = random.nextInt(girlList.size());
            String s = girlList.get(i);
            System.out.println(s);
        }
    }
}
