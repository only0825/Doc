package com.moon.exception;

import java.util.Scanner;

public class Pra1 {

    public static void main(String[] args) {

        Scanner sc = new Scanner(System.in);
        GirlFriend girl = new GirlFriend();
        while (true) {
            try {
                System.out.println("请输入女朋友的姓名");
                String name = sc.nextLine();
                girl.setName(name);
                System.out.println("请输入女朋友的年龄");
                int age = Integer.parseInt(sc.nextLine());
                girl.setAge(age);
                break;
            } catch (NumberFormatException | NameFormatException | AgeOutOfBoundsException e) {
                e.printStackTrace();
            }
        }

        System.out.println(girl);
    }
}
