package com.moon.sort1;

import java.util.Arrays;
import java.util.Comparator;

public class SortDemo1 {

    public static void main(String[] args) {

        Integer[] arr = {3, 44, 38, 5, 47, 15, 36, 26, 27, 2, 46, 4, 19, 50, 48};

        Arrays.sort(arr, Comparator.comparingInt(o -> o));

        System.out.println(Arrays.toString(arr));


    }
}
