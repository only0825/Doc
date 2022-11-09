package com.moon.interface1;

import java.util.Scanner;

public class Test {

    public static void main(String[] args) {

//        Dog dog = new Dog("大黄", 3);
//        dog.swim();
//        dog.eat();
//
//        Rabbit ra = new Rabbit("兔兔", 1);
//        ra.eat();
//
//        Frog f = new Frog("小蛙", 1);
//        f.swim();
//        f.eat();
//
//        Scanner sc = new Scanner(System.in);
//        sc.next();

        Test test = new Test();

        test.eat(
                new Animal() {
                    @Override
                    public void eat() {
                        System.out.println("匿名吃啊");
                    }
                }
        );


    }

    public void eat(Animal a) {
        a.eat();
    }
}
