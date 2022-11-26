package com.moon;

import java.lang.annotation.Native;
import java.lang.reflect.Field;
import java.util.*;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class Demo {
    public static final int   MAX_VALUE = 0x7fffffff;
    private static final int MAX_ARRAY_SIZE = MAX_VALUE - 8;


    public static void main(String[] args) {

    }

    static int getCapacity(List al) throws Exception {
        Field field = ArrayList.class.getDeclaredField("elementData");
        field.setAccessible(true);
        return ((Object[]) field.get(al)).length;
    }

}
