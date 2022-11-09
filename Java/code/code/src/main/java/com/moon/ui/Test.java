package com.moon.ui;

import javax.swing.*;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.time.LocalDate;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.Scanner;

public class Test {

    public static void main(String[] args) {
        int[] arr = {1, 3, 44, 55, 67, 88, 99, 111, 3333};
        int index = binarySearch(arr, 44);
        System.out.println(index);
    }

    public static int binarySearch(int[] arr, int number) {
        int min = 0;
        int max = arr.length - 1;

        while (true) {
            // 查找数组中间的索引
            int mid = (min + max) / 2;
             if (arr[mid] > number) {
                // 如果该索引的元素大于查找的number  x就为结束索引
                 max = mid;
            } else {
                // 如果该索引的元素小于查找的number  x就为开始索引
                 min = mid;
            }
            if (arr[mid] == number) {
                return mid;
            }
        }


//        System.exit(0);
//        int index;
//        int start;
//        int end;
//        while (true) {
//            int len = arr.length;
//            index = len / 2;
//            if (index == 0) {
//                break;
//            }
//            System.out.println(index);
//            System.out.println(arr[index]);
//            System.out.println("==========");
//            if (arr[index] > number) {
//                arr = Arrays.copyOfRange(arr, index, len);
//            } else {
//                arr = Arrays.copyOfRange(arr, 0, index);
//            }
//        }
    }
}
