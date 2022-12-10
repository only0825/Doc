package com.moon.lord3;

import java.util.*;

public class PokerGame {

    // 准备牌
    static ArrayList<String> list = new ArrayList<>();
    static HashMap<String, Integer> hm = new HashMap<>();

    static {
        String[] color = {"♦", "♠", "♥", "♣"};
        String[] number = {"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"};
        for (String c : color) {
            for (String n : number) {
                list.add(c + n);
            }
        }

        list.add(" 大王");
        list.add(" 小王");

        // 指定牌的价值，不用指定3-10的牌，如果存在，则获取价值，如果不存在，则本身的数字就是价值
        hm.put("J", 11);
        hm.put("Q", 12);
        hm.put("K", 13);
        hm.put("A", 14);
        hm.put("2", 15);
        hm.put("大王", 50);
        hm.put("小王", 100);
    }

    public PokerGame() {
        // 洗牌
        Collections.shuffle(list);
        // 发牌
        ArrayList<String> lord = new ArrayList<>();
        ArrayList<String> player1 = new ArrayList<>();
        ArrayList<String> player2 = new ArrayList<>();
        ArrayList<String> player3 = new ArrayList<>();

        // 遍历牌盒得到每一张牌
        for (int i = 0; i < list.size(); i++) {
            String poker = list.get(i);
            if (i <= 2) {
                lord.add(poker);
                continue;
            }
            if (i % 3 == 0) {
                player1.add(poker);
            } else if (i % 3 == 1) {
                player2.add(poker);
            } else {
                player3.add(poker);
            }
        }

        order(lord);
        order(player1);
        order(player2);
        order(player3);
        lookPorker("底牌", lord);
        lookPorker("Ben", player1);
        lookPorker("Fly", player2);
        lookPorker("Moon", player3);
    }

    public void order(ArrayList<String> list) {
        list.sort((o1, o2) -> {
            String color1 = o1.substring(0, 1);
            int value1 = getValue(o1);
            String color2 = o2.substring(0, 1);
            int value2 = getValue(o2);
            int i = value1 - value2;
            return i == 0 ? color1.compareTo(color2) : i;
        });
    }

    public int getValue(String poker) {
         String key = poker.substring(1);
        if (hm.containsKey(key)) {
            return hm.get(key);
        } else {
            return Integer.parseInt(key);
        }
    }

    public void lookPorker(String name, ArrayList<String> list) {
        System.out.print(name + ": ");
        for (String poker : list) {
            System.out.print(poker + " ");
        }
        System.out.println();
    }

}
