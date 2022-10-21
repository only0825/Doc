[TOC]



## MyBatis日志管理

![image-20200526205750778](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200526205750778.png)



引入logback包

```xml

<dependency>
<groupId>ch.qos.logback</groupId>
<artifactId>logback-classic</artifactId>
<version>1.2.3</version>
</dependency>
```



logback允许对日志进行自定义:
在resources目录下新增logback.xml(强制要求这个名称)

```xml

<?xml version="1.0" encoding="UTF-8"?>
<configuration>
   <appender name="console" class="ch.qos.logback.core.ConsoleAppender">
       <encoder>
           <pattern>[%thread] %d{HH:mm:ss.SSS} %-5level %logger{36} - %msg%n</pattern>
       </encoder>
   </appender>
    <!--
        日志输出级别(优先级高到低):
        error: 错误 - 系统的故障日志
        warn: 警告 - 存在风险或使用不当的日志
        info: 一般性消息
        debug: 程序内部用于调试信息
        trace: 程序运行的跟踪信息
     -->
    <root level="debug">
        <appender-ref ref="console"/>
    </root>
</configuration>
```



## MyBatis 动态SQL

```xml
<select id="dynamicSQL" parameterType="java.util.Map" resultType="com.imooc.mybatis.entity.Goods">
    select * from t_goods
    <where>
      <if test="categoryId != null">
          and category_id = #{categoryId}
      </if>
      <if test="currentPrice != null">
          and current_price &lt; #{currentPrice}
      </if>
    </where>
</select>
```



## MyBatis二级缓存

MyBatis缓存：用于数据优化、提高程序执行效率的有效方式

为什么使用缓存？
二次查询相同的元素，得出相同的记录，若两次都需从硬盘中获取，读取的速度很慢，这次需要开启缓存。

一级缓存默认开启，缓存范围SQLSession会话，
二级缓存手动开启，属于范围Mapper Namespace

每一个SqlSession都共用一个二级缓存
写操作commit提交时对该namespace里面的缓存强制清空
在查询标签中配置useCache=false 可以不使用缓存
在查询标签中配置flushCache=true 代表对执行完某一条sql后对该namespace下所有缓存进行强制清空

cache标签开启二级缓存,在对应的namespace空间里边添加此标签： 

```xml


<!--开启了二级缓存
    eviction是缓存的清除策略,当缓存对象数量达到上限后,自动触发对应算法对缓存对象清除
        1.LRU – 最近最少使用的:移除最长时间不被使用的对象。
        O1 O2 O3 O4 .. O512
        14 99 83 1     893
        2.FIFO – 先进先出:按对象进入缓存的顺序来移除它们。
        3.SOFT – 软引用:移除基于垃圾回收器状态和软引用规则的对象。
        4.WEAK – 弱引用:更积极地移除基于垃圾收集器状态和弱引用规则的对象。

    flushInterval 代表间隔多长时间自动清空缓存,单位毫秒,600000毫秒 = 10分钟
    size 缓存存储上限,用于保存对象或集合(1个集合算1个对象)的数量上限
    readOnly 设置为true ,代表返回只读缓存,每次从缓存取出的是缓存对象本身.这种执行效率较高
             设置为false , 代表每次取出的是缓存对象的"副本",每一次取出的对象都是不同的,这种安全性较高
-->
<cache eviction="LRU" flushInterval="600000" size="512" readOnly="true"/>
```

Cache Hit Ratio [goods]:0.5   是缓存命中率
第一次没使用缓存，第二次使用了缓存，所有是0.5
缓存命中率越高，缓存的使用效率就越高，对程序的优化的效果就越好



## MyBastis的对象关联查询(OneToMany和ManyToOne)

多表关联查询：两个表通过主外键在一条SQL中完成所有数据的提取
多表级联查询：通过一个对象来获取与它关联的另外一个对象，执行的SQL语句分为多条
确定对象之间的关系是 双向的
双向的一对多 应该变成多对多 需要单独抽象出一张中间表来
一对多的表设计，在多的这方需要持有一的这方的主键作为外键

双向一对多的关系 ---- 多对多

具体看对象关联查询demo

![image-20200526220047110](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200526220047110.png)

## 分页

pom.xml引入包：

```xml
<dependency>
    <groupId>com.github.pagehelper</groupId>
    <artifactId>pagehelper</artifactId>
    <version>5.1.10</version>
</dependency>
<dependency>
    <groupId>com.github.jsqlparser</groupId>
    <artifactId>jsqlparser</artifactId>
    <version>2.0</version>
</dependency>
```

mybatis-config.xml中启用Pagehelper分页插件

```xml
<!--启用Pagehelper分页插件-->
<plugins>
    <plugin interceptor="com.github.pagehelper.PageInterceptor">
        <!--设置数据库类型-->
        <property name="helperDialect" value="mysql"/>
        <!--分页合理化-->
        <property name="reasonable" value="true"/>
    </plugin>
</plugins>
```

mapper中（`&lt;` 是<符号）:

```xml
<select id="selectPage" resultType="com.imooc.mybatis.entity.Goods">
    select * from t_goods where current_price &lt; 1000    
</select>
```

测试:

```java
@Test
/**
 * PageHelper分页查询
 */
public void testSelectPage() throws Exception {
    SqlSession session = null;
    try {
        session = MyBatisUtils.openSession();
        /*startPage方法会自动将下一次查询进行分页*/
        PageHelper.startPage(2,10);
        Page<Goods> page = (Page) session.selectList("goods.selectPage");
        System.out.println("总页数:" + page.getPages());
        System.out.println("总记录数:" + page.getTotal());
        System.out.println("开始行号:" + page.getStartRow());
        System.out.println("结束行号:" + page.getEndRow());
        System.out.println("当前页码:" + page.getPageNum());
        List<Goods> data = page.getResult();//当前页数据
        for (Goods g : data) {
            System.out.println(g.getTitle());
        }
        System.out.println("");
    } catch (Exception e) {
        throw e;
    } finally {
        MyBatisUtils.closeSession(session);
    }
}
```



## mybatis和c3p0整合

1.pom中增加c3p0的依赖

```xml
<dependency>
    <groupId>com.mchange</groupId>
    <artifactId>c3p0</artifactId>
    <version>0.9.5.2</version>
</dependency>
```

2.新增一个C3P0DataSourceFactory

```java
/**
* c3p0与mybatis兼容使用的数据源工厂类
*/
public class C3P0DataSourceFactory extends UnpooledDataSourceFactory {
    public C3P0DataSourceFactory(){
    	this.dataSource = new ComboPooledDataSource();
    }
}
```

3.修改mybatis-config.xml,type改成指向新增的类：C3P0DataSourceFactory

```xml
<environment id="c3p0">
    <!-- 采用JDBC方式对数据库事务进行commit/rollback -->
    <transactionManager type="JDBC"></transactionManager>
    <!--采用连接池方式管理数据库连接-->
    <dataSource type="com.imooc.mybatis.datasource.C3P0DataSourceFactory">
        <property name="driverClass" value="com.mysql.jdbc.Driver"/>
        <property name="jdbcUrl" value="jdbc:mysql://localhost:3306/babytun?useUnicode=true&amp;characterEncoding=UTF-8"/>
        <property name="user" value="root"/>
        <property name="password" value="12345678"/>
        <property name="initialPoolSize" value="5"/>
        <property name="maxPoolSize" value="20"/>
        <property name="minPoolSize" value="5"/>
    </dataSource>
</environment>
```



## MyBatis批处理

```xml
<!--INSERT INTO table-->
<!--VALUES ("a" , "a1" , "a2"),("b" , "b1" , "b2"),(....)-->
<insert id="batchInsert" parameterType="java.util.List">
    INSERT INTO t_goods(title, sub_title, original_cost, current_price, discount, is_free_delivery, category_id)
    VALUES
    <foreach collection="list" item="item" index="index" separator=",">
        (#{item.title},#{item.subTitle}, #{item.originalCost}, #{item.currentPrice}, #{item.discount}, #{item.isFreeDelivery}, #{item.categoryId})
    </foreach>
</insert>
<!--in (1901,1902)-->
<delete id="batchDelete" parameterType="java.util.List">
    DELETE FROM t_goods WHERE goods_id in
    <foreach collection="list" item="item" index="index" open="(" close=")" separator=",">
        #{item}
    </foreach>
</delete>
```

```java
    /**
     * 批量插入测试
     * @throws Exception
     */
    @Test
    public void testBatchInsert() throws Exception {
        SqlSession session = null;
        try {
            long st = new Date().getTime();
            session = MyBatisUtils.openSession();
            List list = new ArrayList();
            for (int i = 0; i < 10000; i++) {
                Goods goods = new Goods();
                goods.setTitle("测试商品");
                goods.setSubTitle("测试子标题");
                goods.setOriginalCost(200f);
                goods.setCurrentPrice(100f);
                goods.setDiscount(0.5f);
                goods.setIsFreeDelivery(1);
                goods.setCategoryId(43);
                //insert()方法返回值代表本次成功插入的记录总数

                list.add(goods);
            }
            session.insert("goods.batchInsert", list);
            session.commit();//提交事务数据
            long et = new Date().getTime();
            System.out.println("执行时间:" + (et - st) + "毫秒");
//            System.out.println(goods.getGoodsId());
        } catch (Exception e) {
            if (session != null) {
                session.rollback();//回滚事务
            }
            throw e;
        } finally {
            MyBatisUtils.closeSession(session);
        }
    }
```



## MyBatis注解

![image-20200527233958205](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200527233958205.png)

mybatis-config.xml中

```xml
<mappers>
    <!--<mapper class="com.imooc.mybatis.dao.GoodsDAO"/> 二选一-->
    <package name="com.imooc.mybatis.dao"/>
</mappers>
```