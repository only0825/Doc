package com.moon.function;

import java.util.ArrayList;
import java.util.Collections;

public class Demo2 {

    public static void main(String[] args) {

        ArrayList<String> list = new ArrayList<>();
        Collections.addAll(list, "张无忌", "周芷若", "赵敏", "张强", "张三丰");
        StringOperation so = new StringOperation();
        list.stream().filter(so::stringJudge)
                .forEach(s -> System.out.println(s));

    }
}
