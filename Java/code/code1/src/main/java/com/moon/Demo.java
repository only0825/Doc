package com.moon;

import java.util.Random;

public class Demo {


    public static void main(String[] args) {

        int[] arr = new int[10];
        int total = 0;
        Random r = new Random();
        for (int i = 0; i < arr.length; i++) {
            arr[i] = r.nextInt(100) + 1;
            total += arr[i];
            System.out.println(arr[i]);
        }

        int avg = total / arr.length;
        int smallAvg = 0;
        for (int i = 0; i < arr.length; i++) {
            if (arr[i] < avg) {
                smallAvg += 1;
            }
        }

        System.out.println(total);
        System.out.println(avg);
        System.out.println(smallAvg);
    }
}
