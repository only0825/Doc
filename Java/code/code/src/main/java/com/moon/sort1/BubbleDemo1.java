package com.moon.sort1;

// 冒泡排序
public class BubbleDemo1 {

    public static void main(String[] args) {

        int[] arr = {2, 4, 5, 3, 1};

        // 表示我要执行多少轮
        for (int i = 0; i < arr.length - 1; i++) {
            for (int j = 0; j < arr.length -1 - i; j++) {
                // i 依次表示数组中的每一个索引：0 1 2 3 4
                if (arr[j] > arr[j + 1]) {
                    int temp = arr[j];
                    arr[j] = arr[j + 1];
                    arr[j + 1] = temp;
                }
            }
        }
        for (int a : arr) {
            System.out.println(a);
        }
    }
}
