package com.moon.studentsystem;

import java.util.ArrayList;
import java.util.Scanner;

public class StudentSystem {

    public static void main(String[] args) {

        ArrayList<Student> list = new ArrayList<>();
        loop:
        while (true) {
            System.out.println("----------还原来到MOON学生管理系统-JAVA版");
            System.out.println("1:添加学生");
            System.out.println("2:删除学生");
            System.out.println("3:修改学生");
            System.out.println("4:查询学生");
            System.out.println("5:退出");
            System.out.println("请输入您的选择");
            Scanner sc = new Scanner(System.in);
            String choose = sc.next();
            switch (choose) {
                case "1":
                    addStudent(list);
                    break;
                case "2":
                    delStudent(list);
                    break;
                case "3":
                    updateStudent(list);
                    break;
                case "4":
                    queryStudent(list);
                    break;
                case "5":
                    System.out.println("退出");
                    break loop;
                //System.exit(0); // 停止虚拟机运行
                default:
                    System.out.println("没有这个选项");
            }
        }
    }

    // 添加学生
    public static void addStudent(ArrayList<Student> list) {

        Scanner sc = new Scanner(System.in);
        int id;
        while (true) {
            System.out.println("请输入学生ID");
            id = sc.nextInt();
            boolean flag = contains(list, id);
            if (flag) {
                System.out.println("ID已存在，请重新输入");
            } else {
                break;
            }
        }

        System.out.println("请输入学生姓名");
        String name = sc.next();

        System.out.println("请输入学生年龄");
        int age = sc.nextInt();

        System.out.println("请输入学生家庭住址");
        String address = sc.next();

        Student student = new Student(id, name, age, address);
        list.add(student);
        System.out.println("添加成功！");
    }

    // 删除学生
    public static void delStudent(ArrayList<Student> list) {
        Scanner sc = new Scanner(System.in);
        System.out.println("请输入要删除的学生ID");
        int id = sc.nextInt();
        int index = getIndex(list, id);
        if (index >= 0) {
            list.remove(index);
            System.out.println("id为" + id + "的学生删除成功");
        } else {
            System.out.println("id不存在，删除失败");
        }
    }

    // 修改学生
    public static void updateStudent(ArrayList<Student> list) {

        Scanner sc = new Scanner(System.in);
        System.out.println("请输入要修改的学生ID");
        int id = sc.nextInt();
        int index = getIndex(list, id);
        if (index >= 0) {
            Student stu = list.get(index);
            System.out.println("请输入要修改的信息（姓名、年龄、家庭住址");
            String data = sc.next();
            switch (data) {
                case "姓名":
                    System.out.println("请输入新的姓名");
                    String newName = sc.next();
                    stu.setName(newName);
                    System.out.println("修改成功");
                    break;
                case "年龄":
                    System.out.println("请输入新的年龄");
                    int newAge = sc.nextInt();
                    stu.setAge(newAge);
                    System.out.println("修改成功");
                    break;
                case "地址":
                    System.out.println("请输入新的地址");
                    String newAddress = sc.next();
                    stu.setAddress(newAddress);
                    System.out.println("修改成功");
                    break;
                default:
                    System.out.println("修改信息错误");
            }
        } else {
            System.out.println("id不存在， 修改失败");
        }

    }

    // 查询学生
    public static void queryStudent(ArrayList<Student> list) {

        if (list.size() > 0) {
            System.out.println("id\t\t姓名\t年龄\t家庭住址");
            for (int i = 0; i < list.size(); i++) {
                Student stu = list.get(i);
                System.out.println(stu.getId() + "\t" + stu.getName() + "\t" + stu.getAge() + "\t" + stu.getAddress());
            }
        } else {
            System.out.println("当前无学生信息，请添加后查询");
        }
    }

    // 判断ID在集合中是否存在
    public static boolean contains(ArrayList<Student> list, int id) {
        return getIndex(list, id) >= 0;
    }

    // 通过ID获取索引
    public static int getIndex(ArrayList<Student> list, int id) {
        for (int i = 0; i < list.size(); i++) {
            Student stu = list.get(i);
            if (stu.getId() == id) {
                return i;
            }
        }
        return -1;
    }
}
