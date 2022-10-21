## MySQL中的常用数据类型

整数型：
1、TINYINT占1字节： 有符号值（-127~128）和无符号值（0到255）
2、SMALLINT占2字节：有符号值（-32768~32767）和无符号值（0到65535）
3、MEDIUMINT占3字节：有符号值（-8388608~8388607）和无符号值（0到16777215）
4、INT占4字节：有符号值（-2的31次方到2的31次方减一）和无符号值（0到2的32次方-1）
5、BIGINT占8字节：有符号值（-2的63次方到2的63次方减一）和无符号值（0到2的64次方-1）
6、BOOL,BOOLEAN占1字节：等价于TINYINT(1)，0位false，其余为true
浮点数：
FLOAT[(M（该浮点数总长度）,D（小数点后面的长度）)]占4字节、DOUBLE[(M,D)]占8字节、DECIMAL[(M,D)]定点数，内部以字符串的形式存储
字符串
时间

`一个字节有8位二进制数，每一位都有2种选择（1或者0）所以一共有2*2*2*2*2*2*2*2种情况`



## 数据表相关常用操作

```
查看当前数据库下已有数据表：SHOW TABLES；
查看数据库imooc下的数据表：SHOW FULL TABLES FROMIN imooc；
查看指定数据表的详细信息：SHOW CREATE TABLE tbl_name；
查看表结构：DESC tbl_name;或DESCRIBE tbl_name;或SHOW COLUMNS FROM tbl_name；
删除指定数据表：DROP TABLE IF EXISTS tbl_name;
```



## 完整性约束条件：

unsigned无符号，没有负数，从0开始
zerofill零填充，当显示长度不够的时候可以使用前补0的效果填充至指定长度
NOT NULL非空约束，插入值的时候这个字段必须要给值
DEFAULT默认值，如果插入值的时候没有给字段赋值，则使用默认值
PRIMARY KEY主键，标识记录的唯一性，值不能重复，一个表只能有一个主键，自动禁止为空
UNIQUE KEY唯一性，一个表中可以有多个字段是唯一索引，同样的值不能重复，但是NULL除外
AUTO_INCREMENT自动增长，只能用于数值列，而且配合索引所有
FOREIGN KEY外键约束

