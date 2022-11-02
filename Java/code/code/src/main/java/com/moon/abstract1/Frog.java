package com.moon.abstract1;

import com.moon.abstract1.Animal;

public class Frog extends Animal {

    public Frog() {
    }

    public Frog(String name, int age) {
        super(name, age);
    }

    @Override
    public void eat() {
        System.out.println("青蛙吃在虫子");
    }
}
