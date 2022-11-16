package com.moon.sort1;

// 插入排序
public class InsertDemo1 {

    public static void main(String[] args) {

        int[] arr = {2, 4, 5, 3, 1};

        // 外循环：几轮
        // i:表示这一轮中，我拿着哪个索引上的数据跟后面的数据进行比较并交换
        for (int i = 0; i < arr.length-1; i++) {
            // 内循环：每一轮我要干什么事情？
            // 拿着i跟i后面的数据进行比较交换
            for (int j = i + 1; j < arr.length; j++) {
                if (arr[i] > arr[j]) {
                    int temp = arr[i];
                    arr[i] = arr[j];
                    arr[j] = temp;
                }
            }
        }

        for (int a : arr) {
            System.out.println(a);
        }
    }
}
