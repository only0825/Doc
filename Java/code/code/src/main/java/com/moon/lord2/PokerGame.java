package com.moon.lord2;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.TreeSet;

public class PokerGame {

    // 牌子 Map
    // 把牌跟序号产生对应关系，不需要按照序号排序，所以用HashMap
    static HashMap<Integer, String> hm = new HashMap<>();
    static ArrayList<Integer> list = new ArrayList<>();

    static {
        String[] color = {"♦️", "♣️", "♥️", "♠️"};
        String[] number = {"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"};

        // 序号
        int serialNumber = 1;

        for (String n : number) {
            for (String c : color) {
                hm.put(serialNumber, c + n);
                list.add(serialNumber);
                serialNumber++;
            }
        }

        hm.put(serialNumber, "小王");
        list.add(serialNumber);
        serialNumber++;
        hm.put(serialNumber, "大王");
        list.add(serialNumber);
    }

    public PokerGame() {
        // 洗牌
        Collections.shuffle(list);
        // 发牌
        TreeSet<Integer> lord = new TreeSet<>();
        TreeSet<Integer> player1 = new TreeSet<>();
        TreeSet<Integer> player2 = new TreeSet<>();
        TreeSet<Integer> player3 = new TreeSet<>();

        // 遍历牌盒得到每一张牌
        for (int i = 0; i < list.size(); i++) {
            int serialNumber = list.get(i);
            if (i <= 2) {
                lord.add(serialNumber);
                continue;
            }
            if (i % 3 == 0) {
                player1.add(serialNumber);
            } else if (i % 3 == 1) {
                player2.add(serialNumber);
            } else {
                player3.add(serialNumber);
            }
        }

        lookPorker("底牌", lord);
        lookPorker("Ben", player1);
        lookPorker("Fly", player2);
        lookPorker("Moon", player3);
    }

    public void lookPorker(String name, TreeSet<Integer> ts) {
        System.out.print(name + ": ");
        for (int serialNumber : ts) {
            String poker = hm.get(serialNumber);
            System.out.print(poker + " ");
        }
        System.out.println();
    }
}
