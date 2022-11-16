package com.moon.sort1;

// 选择排序
public class SelectDemo1 {

    public static void main(String[] args) {

        int[] arr = {3, 44, 38, 5, 47, 15, 36, 26, 27, 2, 46, 4, 19, 50, 48};

        // 1.找到无序的那一组数组是从哪个索引开始的
        int startIndex = -1;
        for (int i = 0; i < arr.length; i++) {
            if (arr[i] > arr[i + 1]) {
                startIndex = i + 1;
                break;
            }
        }

        // 2.遍历从startIndex开始到最后一个元素，依次得到无序的那一组数据中的每个元素
        for (int i = startIndex; i < arr.length; i++) {
            // 问题，如何把遍历到的数据，插入到前面有序的这一组当中？
            // 记录当前要插入数据的索引
            int j = i;
            // 如果j大于0 且 当前元素小于上一个元素就交换
            while (j > 0 && arr[j] < arr[j - 1]) {
                int temp = arr[j];
                arr[j] = arr[j - 1];
                arr[j - 1] = temp;
                j--;
            }
        }

        for (int j : arr) {
            System.out.println(j);
        }
    }
}
