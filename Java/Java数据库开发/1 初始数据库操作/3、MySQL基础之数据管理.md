

## 表结构相关操作

添加列：ALTER TABLE test ADD password VARCHAR<20> FIRST

删除列：ALTER TABLE table_name DROP COLUMN column_name

修改字段类型、字段属性：ALTER TABLE tb_name MODIFY 字段类型[字段属性][FIRST|AFTER 字段名称]
MODIFY只能修改字段的类型和字段属性，CHANGE可修改字段类型，字段属性和字段的名称

添加主键：ALTER TABLE tb_name ADD PRIMARY KEY(字段名称)
删除主键：ALTER TABLE tb_name DROP PRIMARY KEY;
添加唯一：ALTER TABLE tb_name ADD UNIQUE KEY|INDEX [index_name](字段名称)；
删除唯一：ALTER TABLE tb_name DROP INDEX index_name；



## 复合索引

复合索引也叫组合索引；

使用:  

创建复合索引 ：

```sql
CREATE INDEX columnId ON table1(col1,col2,col3) ;
```

查询语句： 

```sql
select * from table1 where col1= A and col2= B and col3 = C
```

这时候查询优化器,不在扫描表了,而是直接的从索引中拿数据,因为索引中有这些数据,这叫覆盖式查询,这样的查询速度非常快; 

查询使用时,最好将条件顺序按找索引的顺序,这样效率最高;



**何时是用复合索引 ?**

根据where条件建索引是极其重要的一个原则;   
注意不要过多用索引,否则对表更新的效率有很大的影响,因为在操作表的时候要化大量时间花在创建索引中