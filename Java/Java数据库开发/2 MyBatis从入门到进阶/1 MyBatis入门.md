[TOC]



## SSM开发框架

Spring 职责是对系统中的各个对象进行有效的管理，是框架的框架，其它框架都要基于Spring这个框架进行底层开发

Spring MVC 替代Servelt，更有效的进行Web层面的开发

MyBatis 简化数据库的交互

简言之:

```
Spring提供了底层的对象管理，
Spring MVC提供了Web层面的交互
Mybatis 提供了数据库的增删改查的便捷操作
```



## 什么是MyBatis

*MyBatis是优秀的持久层框架 --将内存中的数据保存在数据库中
*MyBatis使用XML将SQL与程序解耦，便于维护
*MyBatis学习简单，执行高效，是JDBC的延伸

## MyBatis开发流程

* 引入MyBatis依赖
* 创建核心配置文件
* 创建实体（Entity）类
* 创建Mapper映射文件
* 初始化SessionFactory
* 利用SqlSession对象操作数据



## Mybatis环境配置

Maven中的

* Mybatis采用xml文件配置数据库环境信息
* Mybatis环境配置标签<environment>
* environment配置包含数据库驱动，URL，用户名和密码

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE configuration
        PUBLIC "-//mybatis.org//DTD Config 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-config.dtd">
<configuration>
    <environments default="dev">
        <environment id="dev">
            <transactionManager type="JDBC"/>
            <dataSource type="POOLED">
                <property name="driver" value="com.mysql.jdbc.Driver"/>
                <property name="url" value="jdbc:mysql://localhost:3306/babytun?useUnicode=true&amp;characterEncoding=UTF-8"/>
                <property name="username" value="root"/>
                <property name="password" value="root"/>
            </dataSource>
        </environment>
        <environment id="prd">
            <!-- 采用JDBC方式对数据库事务进行commit/rollback -->
            <transactionManager type="JDBC"/>
            <!--采用连接池方式管理数据库连接-->
            <dataSource type="POOLED">
                <property name="driver" value="com.mysql.jdbc.Driver"/>
                <property name="url" value="jdbc:mysql://192.168.1.155:3306/babytun?useUnicode=true&amp;characterEncoding=UTF-8"/>
                <property name="username" value="root"/>
                <property name="password" value="root"/>
            </dataSource>
        </environment>
    </environments>
</configuration>
```



## 从 XML 中构建 SqlSessionFactory

SqlSessionFactorpy  sql会话工厂对象，是MyBatis的核心对象
用于初始化MyBatis，创建SqlSession对象
应保证SqlSessionFactory在应用中的全局唯一性

SqlSession是MyBatis操作数据库的核心对象
SqlSession使用JDBC方式与数据库交互
SqlSession对象提供了数据表CRUD对应方法

SqlSessionFactory 的实例可以通过 SqlSessionFactoryBuilder 获得
而 SqlSessionFactoryBuilder 则可以从 XML 配置文件或一个预先配置的 Configuration 实例来构建出 SqlSessionFactory 实例

```java
@Test
public void testSqlSessionFactory() throws IOException {
    //利用Reader加载classpath下的mybatis-config.xml核心配置文件
    Reader reader = Resources.getResourceAsReader("mybatis-config.xml");
	//InputStream config = Resources.getResourceAsStream("mybatis-config.xml");
    
    //初始化SqlSessionFactory对象,同时解析mybatis-config.xml文件
    SqlSessionFactory sqlSessionFactory = new SqlSessionFactoryBuilder().build(reader);
    System.out.println("SessionFactory加载成功");
    SqlSession sqlSession = null;
    try {
        //创建SqlSession对象,SqlSession是JDBC的扩展类,用于与数据库交互
        sqlSession = sqlSessionFactory.openSession();
        //创建数据库连接(测试用)
        Connection connection = sqlSession.getConnection();
        System.out.println(connection);
    }catch (Exception e){
        e.printStackTrace();
    }finally {
        if(sqlSession != null){
            //如果type="POOLED",代表使用连接池,close则是将连接回收到连接池中
            //如果type="UNPOOLED",代表直连,close则会调用Connection.close()方法关闭连接
            sqlSession.close();
        }
    }
}
```



## MyBatis数据查询

**1、创建实体类**
在main/java下创建com.i.mybatis.entity包，entity包下创建Goods商品实体类，将要查询的数据表中的字段对应的在实体类中增加一系列的私有属性及getter/setter方法，属性采用驼峰命名。

**2、在main/resources下创建新的子目录mappers，代表映射器，里面存放的都是xml文件**。创建goods.xml文件来说明实体类和数据表的对应关系（和哪个数据表对应，属性和字段怎么对应）。
a、通过增加不同的命名空间namespace来区分不同的表SQL语句
b、id为别名，相当于sql名称，因此namespace的设置就很有必要，不然分不清是哪个id；同一个namespace下id要唯一，不同的namespace可以重名
c、resultType代表返回的对象是什么，为实体类的完整路径，在SQL语句执行完后会自动的将得到的每一条记录包装成对应的实体类的对象；
com.imooc.mybatis.entity.Goods

**3、让mybatis认识goods.xml：**
在mybatis-config.xml中添加mappers标签
mybatis在初始化的时候才指定这个goods.xml的存在
<mappers><mapper resource="mapper/goods.xml”></mappers>

**4、创建SqlSession对象**
调用其selectList方法
goods 是namespace
selectAll 是id

**5、开启驼峰命名映射**
在mybatis-config.xml中添加：
setting设置项，开启驼峰命名与字段名的转换。
<setting name="mapUnderscoreToCamelCase" value="true"></setting>



## SQL传参

```java
/**
* 传递单个SQL参数
* @throws Exception
*/
@Test
public void testSelectById() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        Goods goods = session.selectOne("goods.selectById" , 1603);
        System.out.println(goods.getTitle());
    }catch (Exception e){
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}

/**
 * 传递多个SQL参数
 * @throws Exception
 */
@Test
public void testSelectByPriceRange() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        Map param = new HashMap();
        param.put("min",100);
        param.put("max" , 500);
        param.put("limt" , 10);
        List<Goods> list = session.selectList("goods.selectByPriceRange", param);
        for(Goods g:list){
            System.out.println(g.getTitle() + ":" + g.getCurrentPrice());

        }
    }catch (Exception e){
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}
```

```xml
<!-- 单参数传递,使用parameterType指定参数的数据类型即可,SQL中#{value}提取参数-->
<select id="selectById" parameterType="Integer" resultType="com.imooc.mybatis.entity.Goods">
    select * from t_goods where goods_id = #{value}
</select>

<!-- 多参数传递时,使用parameterType指定Map接口,SQL中#{key}提取参数 -->
<select id="selectByPriceRange" parameterType="java.util.Map" resultType="com.imooc.mybatis.entity.Goods">
    select * from t_goods
    where
      current_price between  #{min} and #{max}
    order by current_price
    limit 0,#{limt}
</select>
```



## 获取多表关联查询结果

```java
/**
 * 利用Map接收关联查询结果
 * @throws Exception
 */
@Test
public void testSelectGoodsMap() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        List<Map> list = session.selectList("goods.selectGoodsMap");
        for(Map map : list){
            System.out.println(map);
        }
    }catch (Exception e){
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}
```

```xml
<!-- 利用LinkedHashMap保存多表关联结果
    MyBatis会将每一条记录包装为LinkedHashMap对象
    key是字段名  value是字段对应的值 , 字段类型根据表结构进行自动判断
    优点: 易于扩展,易于使用
    缺点: 太过灵活,无法进行编译时检查
 -->
<select id="selectGoodsMap" resultType="java.util.LinkedHashMap" flushCache="true">
    select g.* , c.category_name,'1' as test from t_goods g , t_category c
    where g.category_id = c.category_id
</select>
```



## ResultMap结果映射

DTO -> Data Transfer Object--数据传输对象

com.imooc.mybatis.entity下有Goods和Category两个实体类

com.imooc.mybatis.dto下有GoodDTO类：

```java
public class GoodsDTO {
    private Goods goods = new Goods();
    private Category category = new Category();
    private String test;

    public Goods getGoods() {
        return goods;
    }

    public void setGoods(Goods goods) {
        this.goods = goods;
    }

    public Category getCategory() {
        return category;
    }

    public void setCategory(Category category) {
        this.category = category;
    }

    public String getTest() {
        return test;
    }

    public void setTest(String test) {
        this.test = test;
    }
}
```

映射的GoodsDTO这个类，这个类中又实例化了 Goods和Category这两个实体类，所以可以利用GoodsDTO这个类中实例化的对象存多表关联查询出来的数据

```xml
<!--结果映射-->
<resultMap id="rmGoods" type="com.imooc.mybatis.dto.GoodsDTO">
    <!--设置主键字段与属性映射-->
    <id property="goods.goodsId" column="goods_id"></id>
    <!--设置非主键字段与属性映射-->
    <result property="goods.title" column="title"></result>
    <result property="goods.originalCost" column="original_cost"></result>
    <result property="goods.currentPrice" column="current_price"></result>
    <result property="goods.discount" column="discount"></result>
    <result property="goods.isFreeDelivery" column="is_free_delivery"></result>
    <result property="goods.categoryId" column="category_id"></result>
    <result property="category.categoryId" column="category_id"></result>
    <result property="category.categoryName" column="category_name"></result>
    <result property="category.parentId" column="parent_id"></result>
    <result property="category.categoryLevel" column="category_level"></result>
    <result property="category.categoryOrder" column="category_order"></result>

    <result property="test" column="test"/>
</resultMap>
<select id="selectGoodsDTO" resultMap="rmGoods">
    select g.* , c.*,'1' as test from t_goods g , t_category c
    where g.category_id = c.category_id
</select>
```



## MyBatis数据插入操作

```xml
<!--flushCache="true"在sql执行后强制清空缓存-->
<insert id="insert" parameterType="com.imooc.mybatis.entity.Goods" flushCache="true">
    INSERT INTO t_goods(title, sub_title, original_cost, current_price, discount, is_free_delivery, category_id)
    VALUES (#{title} , #{subTitle} , #{originalCost}, #{currentPrice}, #{discount}, #{isFreeDelivery}, #{categoryId})
  <selectKey resultType="Integer" keyProperty="goodsId" order="AFTER">
      select last_insert_id()
  </selectKey>
</insert>
```

```java
/**
 * 新增数据
 * @throws Exception
 */
@Test
public void testInsert() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        Goods goods = new Goods();
        goods.setTitle("测试商品");
        goods.setSubTitle("测试子标题");
        goods.setOriginalCost(200f);
        goods.setCurrentPrice(100f);
        goods.setDiscount(0.5f);
        goods.setIsFreeDelivery(1);
        goods.setCategoryId(43);
        //insert()方法返回值代表本次成功插入的记录总数
        int num = session.insert("goods.insert", goods);
        session.commit();//提交事务数据
        System.out.println(goods.getGoodsId());
    }catch (Exception e){
        if(session != null){
            session.rollback();//回滚事务
        }
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}
```

```java
/**
 * MyBatisUtils工具类,创建全局唯一的SqlSessionFactory对象
 */
public class MyBatisUtils {
    //利用static(静态)属于类不属于对象,且全局唯一
    private static SqlSessionFactory sqlSessionFactory = null;
    //利用静态块在初始化类时实例化sqlSessionFactory
    static {
        Reader reader = null;
        try {
            reader = Resources.getResourceAsReader("mybatis-config.xml");
            sqlSessionFactory = new SqlSessionFactoryBuilder().build(reader);
        } catch (IOException e) {
            e.printStackTrace();
            //初始化错误时,通过抛出异常ExceptionInInitializerError通知调用者
            throw new ExceptionInInitializerError(e);
        }
    }

    /**
     * openSession 创建一个新的SqlSession对象
     * @return SqlSession对象
     */
    public static SqlSession openSession(){
        //默认SqlSession对自动提交事务数据(commit)
        //设置false代表关闭自动提交,改为手动提交事务数据
        return sqlSessionFactory.openSession(false);
    }

    /**
     * 释放一个有效的SqlSession对象
     * @param session 准备释放SqlSession对象
     */
    public static void closeSession(SqlSession session){
        if(session != null){
            session.close();
        }
    }
}
```



## SelectKey和useGeneratedKeys的区别

```xml
<selectKey resultType="Integer" keyProperty="goodsId" order="AFTER">
select last_insert_id()
</selectKey>
```

useGeneratedKeys是在insert标签里边添加标签的，useGeneratedKeys=”true”//这里默认为false，
keyProperty=”goodsId” keyColumn=”goods_id”>
INSERT INTO SQL语句
SelectKey适用于所有数据库，但需要根据不同的数据库编写对应的获得最后改变主键值得查询语句
userGenerateKeys只支持“自增主键”的数据库，但使用简单，会根据不同的数据库驱动自动编写查询语句

--在Oracle中selectKey的用法

```xml
<insert id=”insert” parameterType=”com.imooc.mybatis.entity.Message”>
INSERT INTO SQL语句
<selectKey resultType=”Integer” order=”BEFORE” keyProperty=”id”>
SELECT seq_message_nextval as id from dual
</selectKey>
</insert>
```



## 更新与删除操作

```xml
<update id="update" parameterType="com.imooc.mybatis.entity.Goods">
    UPDATE t_goods
    SET
      title = #{title} ,
      sub_title = #{subTitle} ,
      original_cost = #{originalCost} ,
      current_price = #{currentPrice} ,
      discount = #{discount} ,
      is_free_delivery = #{isFreeDelivery} ,
      category_id = #{categoryId}
    WHERE
      goods_id = #{goodsId}
</update>
<!--delete from t_goods where goods_id in (1920,1921)-->
<delete id="delete" parameterType="Integer">
    delete from t_goods where goods_id = #{value}
</delete>
```

```java
/**
 * 更新数据
 * @throws Exception
 */
@Test
public void testUpdate() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        Goods goods = session.selectOne("goods.selectById", 739);
        goods.setTitle("更新测试商品");
        int num = session.update("goods.update" , goods);
        session.commit();//提交事务数据
    }catch (Exception e){
        if(session != null){
            session.rollback();//回滚事务
        }
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}

/**
 * 删除数据
 * @throws Exception
 */
@Test
public void testDelete() throws Exception {
    SqlSession session = null;
    try{
        session = MyBatisUtils.openSession();
        int num = session.delete("goods.delete" , 739);
        session.commit();//提交事务数据
    }catch (Exception e){
        if(session != null){
            session.rollback();//回滚事务
        }
        throw e;
    }finally {
        MyBatisUtils.closeSession(session);
    }
}
```





## MyBatis两种传值方式

${}原文传值(文本替换)，未经任何处理对SQL文本替换
#{}预编译传值，使用预编译传值可以**预防SQL注入**

${}原文传值用处在哪：
像排序order by title desc 是需要原文传值的



## MyBatis工作流程

![image-20200525014609830](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200525014609830.png)