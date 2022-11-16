package com.moon.search;

public class BlockSearch {

    public static void main(String[] args) {
        // 扩展的分块查找（无规律的数据）
        int[] arr = {27, 22, 30, 40, 36,
                13, 19, 16, 20,
                7, 10,
                43, 50, 48};

        // 创建4个块的对象
        Block b1 = new Block(22, 40, 0, 4);
        Block b2 = new Block(13, 20, 5, 8);
        Block b3 = new Block(7, 10, 9, 10);
        Block b4 = new Block(43, 50, 11, 13);

        // 定义数组用来管理4个块的对象（索引表）
        Block[] blockArr = {b1, b2, b3, b4};

        // 定义一个标亮用来记录要查找的元素
        int number = 16;

        // 调用方法，传递索引表，数组，要查找的元素
        int index = getIndex(blockArr, arr, number);

        System.out.println(index);
    }

    private static int getIndex(Block[] blockArr, int[] arr, int number) {
        int index = -1;
        for (Block b : blockArr) {
            if (number <= b.getMax() && number >= b.getMin()) {
                for (int i = b.getStartIndex(); i < b.getEndIndex(); i++) {
                    if (arr[i] == number) {
                        index = i;
                        return index;
                    }
                }
            }
        }

        return index;
    }


}

class Block {
    int min;
    int max;
    int startIndex;
    int endIndex;

    public Block() {
    }

    public Block(int min, int max, int startIndex, int endIndex) {
        this.min = min;
        this.max = max;
        this.startIndex = startIndex;
        this.endIndex = endIndex;
    }

    /**
     * 获取
     *
     * @return min
     */
    public int getMin() {
        return min;
    }

    /**
     * 设置
     *
     * @param min
     */
    public void setMin(int min) {
        this.min = min;
    }

    /**
     * 获取
     *
     * @return max
     */
    public int getMax() {
        return max;
    }

    /**
     * 设置
     *
     * @param max
     */
    public void setMax(int max) {
        this.max = max;
    }

    /**
     * 获取
     *
     * @return startIndex
     */
    public int getStartIndex() {
        return startIndex;
    }

    /**
     * 设置
     *
     * @param startIndex
     */
    public void setStartIndex(int startIndex) {
        this.startIndex = startIndex;
    }

    /**
     * 获取
     *
     * @return endIndex
     */
    public int getEndIndex() {
        return endIndex;
    }

    /**
     * 设置
     *
     * @param endIndex
     */
    public void setEndIndex(int endIndex) {
        this.endIndex = endIndex;
    }

    public String toString() {
        return "Block{min = " + min + ", max = " + max + ", startIndex = " + startIndex + ", endIndex = " + endIndex + "}";
    }
}
