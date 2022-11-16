package com.moon.test;

import java.sql.Array;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;

public class Test1 {

    public static void main(String[] args) {

        GrilFriend g1 = new GrilFriend("xiaoshishi", 18, 1.67);
        GrilFriend g2 = new GrilFriend("xiaodandan", 19, 1.72);
        GrilFriend g3 = new GrilFriend("xiaohuihui", 19, 1.78);
        GrilFriend g4 = new GrilFriend("abc", 19, 1.78);
        GrilFriend[] arr = {g1, g2, g3, g4};

        // lambda 表达式
        Arrays.sort(arr, (o1, o2) -> {
            // 按age大小排序，age一样按身高，身高一样按姓名的字母
            double temp = o1.getAge() - o2.getAge();
            temp = temp == 0 ? o1.getHeight() - o2.getHeight() : temp;
            temp = temp == 0 ? o1.getName().compareTo(o2.getName()) : temp;

            if (temp > 0) {
                return 1;
            } else if (temp < 0) {
                return -1;
            } else {
                return 0;
            }
        });

//        Arrays.sort(arr, new Comparator<GrilFriend>() {
//            @Override
//            public int compare(GrilFriend o1, GrilFriend o2) {
//                // 按age大小排序，age一样按身高，身高一样按姓名的字母
//                double temp = o1.getAge() - o2.getAge();
//                temp = temp == 0 ? o1.getHeight() - o2.getHeight() : temp;
//                temp = temp == 0 ? o1.getName().compareTo(o2.getName()) : temp;
//
//                if (temp > 0) {
//                    return 1;
//                } else if (temp < 0) {
//                    return -1;
//                } else {
//                    return 0;
//                }
//            }
//        });

        System.out.println(Arrays.toString(arr));
    }
}
