package com.moon.sort1;

public class Test {

    public static void main(String[] args) {
        // 1-100的和
//        System.out.println(getSum(100));

        // 求阶乘
        // 5! = 5 * 4 * 3 * 2 * 1
        // 100! = 100 * 99 * 98 * 97 * 96 ... * 2 * 1;

        System.out.println(getChen(3));
    }

    private static int getChen(int number) {
        if (number == 1) {
            return 1;
        }

        return number * getChen(number - 1);

    }

    private static int getSum(int number) {
        if (number == 1) {
            return 1;
        }
        return number + getSum(number - 1);
    }

}
