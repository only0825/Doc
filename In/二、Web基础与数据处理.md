[TOC]

## 请说明Servlet执行流程

**请说明Servlet执行流程？**
客户端向tomcat服务器发送http请求，包含servlet映射地址和要传递的参数--》tomcat解析每一个web.xml文件，找到与之匹配的url和对应的servlet name--》根据servlet name找到对应的servlet，并对这个servlet进行实例化和初始化--》tomcat执行servlet实例中的方法--》方法运行后把程序执行结果通过响应返回浏览器--》浏览器接收到这段代码后进行解释。

![img](http://img.mukewang.com/climg/603a632e000107e119201080.jpg)

**Servlet实例化时机**
servlet默认情况下是第一次访问的时候实例化的，也可以通过web.xml配置loadon-start-up，使其在服务器启动时候实例化。

**servlet在并发环境下是如何处理的？**
servlet是基于单例多线程来处理并发情况。利用多线程技术提供web服务。

servlet**多线程处理的情况下，如何解决线程安全问题？**
所有的线程，都共享一个servlet实例。所以我们在使用servlet时，不允许在servlet内创建，存在状态的变量和对象。因为这样会在并发访问时产生无法预期的结果。



对于serlvet来说，全局在tomcat中，有且只有一个唯一的对象。tomcat不会创建servlet的多个实例。





## Servlet生命周期是什么？

1，装载--java应用程序启动的时候，tomcat会扫描web.xml文件，得知当前有哪些servlet.(装载时并不会实例化Servlet)（创建时java层面的对象创建）

2，创建--当 url 第一次访问servlet地址的时候进行创建。同时执行构造函数，创建对象。

3，初始化--servlet在创建对象以后，马上执行init()初始化函数，对servlet进行初始化。（初始化，是servlet自身专门用于初始化servlet执行资源的方法）

4，提供服务--service()方法。servlce()方法--对于发来的请求（无论是post/get）,一律使用servlet方法接收处理。如果将请求细化，service()方法下还可以细化为doGet()/doPost()方法。doGet()--只处理get请求，doPost()--只处理post请求。

5，销毁--在web应用重启或关闭时使用destory()方法将servlet的资源彻底销毁。

```
    装载-web.xml    
    创建-构造函数    
    初始化-init()    
    提供服务-service()    
    销毁-destory()
```

## 请阐述HTTP请求的结构：

```
1，HTTP请求包含三部分：请求行、请求头、请求体

请求行--通常在HTTP的第一行，说明了发送的方式（get/post），发送的地址、和HTTP的版本号。

请求头--说明了，从浏览器到服务器发送到辅助信息。
Accept-Language：zh-CN  说明浏览器优先使用中文
User-Agent：代表了用户的使用环境（判断用户使用的是手机还会电脑进行的访问，然后根据浏览器的规格不停进行不同的展现）
Content-Type：说明了提交的表单的格式

请求体--由浏览器向服务器发送的真实数据，请求题中，数据使用键值对的形式“键”和“值”之间使用“=”连接。多个键值对之间使用“&”进行分隔。（请求体，只有在post请求中才会存在，get请求中是没有请求体这一项的）请求体会被附加在url 后面发送到服务器。


响应--有服务器返回给浏览器的结果
HTTP响应包含3部分内容：响应行、响应头、响应体。

响应行--通常在响应的第一行，包含http版本、状态码、状态码的英文描述
200--表示访问成功
404--表示资源未找到
500--代表的是服务器的内部错误。

响应头--表述了返回数据的一些辅助信息，使用了哪种web服务器、
Service--表示使用了哪种web服务器、
Content-Type：表示数据返回给浏览器以后，浏览器采用什么样的方式进行处理呢。（text/html--表示把返回的数据解释成html进行显示）
Date--响应数据产生的时间

响应体--服务器向浏览器返回的真实数据，（html片段、二进制的内容、xml）
```



## 请求转发与响应重定向之间的区别

javaweb中有两种资源跳转的方式：

```
1，请求转发--是服务器跳转，只会产生一次请求。会将请求原封不动的转发给下一个请求。(服务器跳转)
语法：request.getRequestDispatcher().forward();
```

![img](http://img.mukewang.com/climg/60380a550001319619201080.jpg)

```
2，响应重定向--是浏览器端的跳转，会产生两次请求。（浏览器客户端跳转）
语法：response.sendRedirect();
```

![img](http://img.mukewang.com/climg/60380b200001c9fb19201080.jpg)



## 请阐述Session的原理

Session的原理：

session--又被称为用户会话，与客户端浏览器窗口绑定的，且存储在服务器内部的用户数据。

session 的工作原理是客户端登录完成之后，服务器会创建对应的 session，session 创建完之后，会把 session 的 id 发送给客户端，客户端再存储到浏览器的cookie中。只要浏览器没关闭，这个cookie是一直存在的。这样客户端每次访问服务器时，都会带着 sessionid，服务器拿到 sessionid 之后，在内存找到与之对应的 session 这样就可以正常工作了。

## JSP九大内置对象是什么

- 输出输入**对象**:request对象、response对象、out对象
- 通信控制**对象**:pageContext对象、session对象、application对象
- Servlet**对象**:page对象、config对象
- 错误处理**对象**:exception对象

![http://img.mukewang.com/climg/60851249094cb42708580486.jpg](http://img.mukewang.com/climg/60851249094cb42708580486.jpg)

## Statement和PreparedStatement的区别

* PreparedStatement是预编译的SQL语句，效率高于Statement
* PreparedStatement支持?操作符(参数化操作)，相对于Statement更灵活
* PreparedStatement可以防止SQL注入，安全性高于Statement

```
Statement 是 Java 执行数据库操作的一个重要方法，用于在已经建立数据库连接的基础上，向数据库发送要执行的SQL语句。Statement对象，用于执行不带参数的简单SQL语句。
与此差不多的还有PreparedStatement，PreparedStatement 继承了Statement，一般如果已经是稍有水平开发者，应该以PreparedStatement代替Statement。
```





## 请说明JDBC使用步骤：

JavaDataBatesConnecter

1. 加载JDBC驱动

2. 创建于数据库的链接（Connection）

3. 创建命令（PreparedStatement/Statement）

4. 对于查询的处理结果，要使用ResultSet对象进行接收，并通过遍历将结果进行处理。

5. 关闭数据库连接

```java
        String driverName = "com.mysql.jdbc.Driver";
        String URL = "jdbc:mysql://127.0.0.1:3306/scott";
        String sql = "SELECT * FROM emp";
        String username = "root";
        String password = "root";
        Connection conn = null;
        try {
            //1.加载JDBC驱动
            Class.forName(driverName);
            //2.建立连接
            conn = DriverManager.getConnection(URL, username, password);
            //3.创建命令（Statement)
            Statement ps = conn.createStatement();
            //4.处理结果
            ResultSet rs = ps.executeQuery(sql);
            while (rs.next()) {
                System.out.println(rs.getString("ename") + "," +
                        rs.getString("job") + "," +
                        rs.getFloat("sal"));
            }
        } catch (ClassNotFoundException e) {
            e.printStackTrace();
        } catch (SQLException e) {
            e.printStackTrace();
        }
        //5.关闭连接
        finally {
            if (conn != null) {
                try {
                    conn.close();
                } catch (SQLException e) {
                    e.printStackTrace();
                }
            }
        }
```



## SQL编程训练

员工表emp

![image-20220331023538236](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20220331023538236.png)

部门表dept

![image-20220331023544213](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20220331023544213.png)

```sql
/*1.按部门编号升序、工资倒序的排列员工信息 考察多条件排序*/
SELECT * FROM emp ORDER BY deptno ASC , sal DESC
/*2.列出deptno=30的部门名称及员工 考察关联查询的使用 */
SELECT dept.dname , emp.* FROM emp , dept WHERE emp.deptno = dept.deptno and dept.deptno = 30
/*3.列出每个部门最高、最低及平均工资 考察的是对于分组的使用*/
SELECT deptno , max(sal) , min(sal) , avg(sal) , count(*)  FROM emp GROUP BY deptno
```

```sql
/*1.列出市场部(SALES)及研发部(RESEARCH)的员工 考察多条件的or的运用*/
SELECT * FROM dept d , emp e where d.deptno = e.deptno and
(d.dname = 'SALES' or d.dname = 'RESEARCH')
/*2.列出人数超过3人的部门 考察having关键字的使用*/
SELECT d.dname , count(*) FROM dept d , emp e WHERE d.deptno = e.deptno GROUP BY d.dname HAVING count(*) > 3
/*3.计算MILLER年薪比SMITH高多少 考察子查询*/
SELECT a.a_sal , b.b_sal , a.a_sal - b.b_sal FROM
(SELECT sal * 12 a_sal FROM emp WHERE ename = 'MILLER') a,
(SELECT sal * 12 b_sal FROM emp WHERE ename = 'SMITH') b
```

```sql
/*1.列出直接向King汇报的员工 考察自查询*/
SELECT * FROM emp WHERE mgr = (SELECT empno FROM emp WHERE ename = 'King')

SELECT e.* FROM emp e , (SELECT empno FROM emp WHERE ename = 'King') k WHERE e.mgr = k.empno

/*2.列出公司所有员工的工龄，并倒序排列  考察Mysql函数的运用*/
/*SELECT DATE_FORMAT(NOW(),"%Y/%m/%d")*/
SELECT * FROM (
   SELECT emp.* , DATE_FORMAT(NOW(),"%Y") -   DATE_FORMAT(hiredate,"%Y") yg FROM emp
) d ORDER BY d.yg DESC

/*3.计算管理者与基层员工平均薪资差额  SQL的综合运用能力*/
SELECT a_avg - b_avg FROM
(SELECT avg(sal) a_avg FROM emp where job = 'MANAGER' OR job = 'PRESIDENT') a ,
(SELECT avg(sal) b_avg FROM emp where job in ('CLERK' , 'SALESMAN' , 'ANALYST')) b
```
