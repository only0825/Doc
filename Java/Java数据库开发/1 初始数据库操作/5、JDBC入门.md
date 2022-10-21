#### Java连接MySQL

Project Structure->Libraries导入Java 连接 MySQL 需要的驱动包

```java
package com.imooc.test;
import java.sql.*;

public class JDBCTest1 {

    public static final String URL = "jdbc:mysql://localhost:3306/oa";
    public static final String USER = "root";
    public static final String PASSWORD = "root";

    public static void main(String[] args) throws SQLException, ClassNotFoundException {
        //1.加载驱动程序
        Class.forName("com.mysql.jdbc.Driver");
        //2.获得数据库连接
        Connection conn = DriverManager.getConnection(URL, USER, PASSWORD);
        //3.操作数据库，实现增删改查
        Statement stmt = conn.createStatement();
        ResultSet rs = stmt.executeQuery("SELECT * FROM employee");
        //如果有数据，rs.next()返回true
        while (rs.next()) {
            System.out.println(rs.getString("post"));
        }
    }
}
```

插入

```java
package com.imooc.jdbc;

import org.junit.Test;

import java.sql.*;

public class JDBCDemo2 {

    public static final String URL = "jdbc:mysql://localhost:3306/oa";
    public static final String USER = "root";
    public static final String PASSWORD = "root";

    @Test
    /**
     * 保存操作
     */
    public void demo1(){
        Connection conn = null;
        Statement stmt = null;
        try{
            // 注册驱动:
            Class.forName("com.mysql.jdbc.Driver");
            // 获得连接:
            conn = DriverManager.getConnection(URL, USER, PASSWORD);

            // 获得执行SQL语句的对象:
            stmt = conn.createStatement();
            // 编写SQL：
            String sql = "insert into employee values (10005,'000000', '焚风',10002, '其它')";
            // 执行SQL:
            int i = stmt.executeUpdate(sql);
            if(i > 0){
                System.out.println("保存成功！");
            }
        }catch(Exception e){
            e.printStackTrace();
        }finally{
            // 释放资源:
            if(stmt != null){
                try {
                    stmt.close();
                } catch (SQLException e) {
                    e.printStackTrace();
                }
                stmt = null;
            }
            if(conn != null){
                try {
                    conn.close();
                } catch (SQLException e) {
                    e.printStackTrace();
                }
                conn = null;
            }
        }
    }
}
```

#### 封装JDBC工具类

```java
package com.imooc.utils;

import java.io.IOException;
import java.io.InputStream;
import java.sql.*;
import java.util.Properties;

/**
 * JDBC的工具类
 */
public class JDBCUtils {
    private static final String driverClass;
    private static final String url;
    private static final String username;
    private static final String password;

    static{
        // 加载属性文件并解析：
        Properties props = new Properties();
        // 如何获得属性文件的输入流？
        // 通常情况下使用类的加载器的方式进行获取：
        InputStream is = JDBCUtils.class.getClassLoader().getResourceAsStream("jdbc.properties");
        try {
            props.load(is);
        } catch (IOException e) {
            e.printStackTrace();
        }

        driverClass = props.getProperty("driverClass");
        url = props.getProperty("url");
        username = props.getProperty("username");
        password = props.getProperty("password");
    }

    /**
     * 注册驱动的方法
     * @throws ClassNotFoundException
     */
    public static void loadDriver() throws ClassNotFoundException{
        Class.forName(driverClass);
    }

    /**
     * 获得连接的方法:
     * @throws SQLException
     */
    public static Connection getConnection() throws Exception{
        loadDriver();
        Connection conn = DriverManager.getConnection(url, username, password);
        return conn;
    }

    /**
     * 资源释放
     */
    public static void release(Statement stmt, Connection conn){
        if(stmt != null){
            try {
                stmt.close();
            } catch (SQLException e) {
                e.printStackTrace();
            }
            stmt = null;
        }
        if(conn != null){
            try {
                conn.close();
            } catch (SQLException e) {
                e.printStackTrace();
            }
            conn = null;
        }
    }

    public static void release(ResultSet rs, Statement stmt, Connection conn){
        if(rs!= null){
            try {
                rs.close();
            } catch (SQLException e) {
                e.printStackTrace();
            }
            rs = null;
        }
        if(stmt != null){
            try {
                stmt.close();
            } catch (SQLException e) {
                e.printStackTrace();
            }
            stmt = null;
        }
        if(conn != null){
            try {
                conn.close();
            } catch (SQLException e) {
                e.printStackTrace();
            }
            conn = null;
        }
    }
}
```

src下创建jdbc.properties

```properties
driverClass=com.mysql.jdbc.Driver
url=jdbc:mysql:///oa
username=root
password=root
```

测试

```java
package com.imooc.jdbc;
import com.imooc.utils.JDBCUtils;
import org.junit.Test;
import java.sql.Connection;
import java.sql.Statement;

public class JDBCDemo3 {

    @Test
    // 保存记录
    public void demo1() {
        Connection conn = null;
        Statement stmt = null;
        try {
            // 获得连接:
            conn = JDBCUtils.getConnection();
            // 创建执行SQL语句的对象
            stmt = conn.createStatement();
            // 编写SQL:
            String sql = "insert into employee values (10006,'000000', '焚风',10002, '其它')";
            // 执行SQL:
            int num = stmt.executeUpdate(sql);
            if (num > 0) {
                System.out.println("保存成功!");
            }
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            // 释放资源:
            JDBCUtils.release(stmt, conn);
        }
    }
}
```

#### 预防SQL注入

```java
package com.imooc.jdbc;

import com.imooc.utils.JDBCUtils;
import org.junit.Test;
import java.sql.*;

public class JDBCDemo4 {

    @Test
    /**
     * 测试SQL注入漏洞的方法
     */
    public void demo1(){
        // 真正密码是000000
        boolean flag = JDBCDemo4.login2("刘备' or '1=1", "1fsdsdfsdf");
        if(flag){
            System.out.println("登录成功！");
        }else{
            System.out.println("登录失败！");
        }
    }

    /**
     * 避免SQL注入漏洞的方法
     */
    public static boolean login2(String name, String password) {
        Connection conn = null;
        PreparedStatement pstmt = null;
        ResultSet rs = null;
        boolean flag = false;
        try {
            // 获得连接:
            conn = JDBCUtils.getConnection();
            // 编写SQL:
            String sql = "select * from employee where name = ? and password = ?";
            // 预处理SQL:
            pstmt = conn.prepareStatement(sql);
            pstmt.setString(1, name);
            pstmt.setString(2, password);
            // 执行SQL:
            rs = pstmt.executeQuery();
            // 判断结果
            if (rs.next()) {
                flag = true;
            }
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            JDBCUtils.release(rs, pstmt, conn);
        }
        return flag;
    }

    /**
     * 产生SQL注入漏洞的方法
     */
    public static boolean login1(String name,String password){
        Connection conn = null;
        Statement stmt  = null;
        ResultSet rs = null;
        boolean flag = false;
        try{
            conn = JDBCUtils.getConnection();
            // 创建执行SQL语句的对象:
            stmt = conn.createStatement();
            // 编写SQL:
            String sql = "select * from employee where name = '"+name+"' and password = '"+password+"'";
            System.out.println(sql);

            // 执行SQL:
            rs = stmt.executeQuery(sql);
            // 判断结果集中是否有数据。
            if(rs.next()){
                flag = true;
            }
        }catch(Exception e){
            e.printStackTrace();
        }finally{
            JDBCUtils.release(rs, stmt, conn);
        }
        return flag;
    }
}
```