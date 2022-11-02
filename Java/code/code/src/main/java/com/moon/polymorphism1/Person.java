package com.moon.polymorphism1;

public class Person {

    int age;
    String name;

    public Person() {
    }

    public Person(int age, String name) {
        this.age = age;
        this.name = name;
    }

    public int getAge() {
        return age;
    }

    public void setAge(int age) {
        this.age = age;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public void keepPet(Animal a, String something) {
        if (a instanceof Dog) {
            System.out.println("年龄为" + this.age + "的" + this.name + "养了一只" + a.color + "颜色的" + a.age + "岁的🐶");
            a.eat(something);
            Dog d = (Dog) a;
            d.lookHome();
        } else if (a instanceof Cat) {
            System.out.println("年龄为" + this.age + "的" + this.name + "养了一只" + a.color + "颜色的" + a.age + "岁的🐱");
            a.eat(something);
            Cat d = (Cat) a;
            d.catchMouse();
        }
    }

}
