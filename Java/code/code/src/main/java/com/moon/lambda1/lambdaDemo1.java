package com.moon.lambda1;

public class lambdaDemo1 {

    public static void main(String[] args) {



        method(() -> System.out.println("游泳～～～"));
    }

    public static void method(Swim s) {
        s.swim();
    }
}


interface Swim {
    public abstract void swim();
}