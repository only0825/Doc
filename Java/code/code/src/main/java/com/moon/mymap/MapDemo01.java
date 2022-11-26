package com.moon.mymap;

import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class MapDemo01 {

    public static void main(String[] args) {

        Map<String, String> map = new HashMap<>();
        map.hashCode();
        map.put("A", "艾斯");
        map.put("B", "毕飞");
        map.put("C", "C罗");

        map.forEach((s, s2) -> System.out.println(s + ":" +s2));
    }
}
