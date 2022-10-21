数据库：是按照数据结构来组织/存储和管理数据的仓库
常见数据库：Oracle / DB2 / SQL Server / Postgre SQL / MySQL
数据库系统（DBS):数据库/数据库管理系统（DBMS)/应用开发工具/管理员及用户
SQL语言（Structured Query Langyage）结构化管理语言：

* DDL（数据定义语言）/DML（数据操作语言）/DQL（数据查询语言）/DCL（数据控制语言）

## 常用命令

my.ini是mySQL的配置文件
登录： mysql -uroot -p或mysql -uroot -p加上密码
退出：exit或quit或\q
获得版本号： mysql -V 或mysql --version
登陆的同时打开指定的数据库：mysql -uroot -p密码 -D db_name
显示现有的库：show database
查看打开的数据库：select database();
命令行结束符默认是;或\g，如select database();
可以通过help或\h或？加上相关关键字来查看手册
\c可以取消当前命令的执行
常用SQL语句：
select user()得到登陆的用户
select version()得到mysql的版本信息
select now()得到当前的日期时间
select database()得到当前打开的数据库
SQL语句语法规范：
常用MySQL的关键字需要大写，库名、表名、字段等使用小写
SQL语句支持折行操作，拆分时不能把完整的单词拆开
数据库名称、表名称、字段名称不要使用MySQL的保留字，如果必须使用，需要使用反引号''将其括起来

SHOW CREATE TABLE tbl_name;    查看创建表的语句

查看引擎存储格式: SHOW TABLE STATUS LIKE 'tbl_name'\G



## 数据库相关操作：

创建数据库：GREATE DATABASE或SCHEMA db_name；（数据库名称最好有意义，名称不要包含特殊字符或MySQL关键字）
查看当前服务器下的全部数据库：SHOW DATABASES或SCHEMAS;
检测数据库名称是否存在，如果不存在则进行创建：CREATE DATABASE IF NOT EXISTS db_name;
创建数据库的同时制定编码方式：CREATE DATABASE IF NOT EXISTS test3 DEFAULT CHARACTER SET 'UTF-8';
修改指定数据库的编码方式：ALTER DATABASE db_name DEFAULT CHARACTER SET = charset;
打开指定数据库：USE db_name;
得到当前打开的数据库：SELECT DATABASE();
删除指定的数据库：DROP DATABASE db_name;
如果数据库存在则删除：DROP DATABASE IF EXISTS db_name;
查看指定数据库的详细信息：SHOW CREATE DATABASE db_name;
mysql 注释： #注释 或者 --注释
常用SQL语句:
查看上一步的警告信息：SHOW WARNINGS;



## 数据表及数据类型简介

数据表相关操作：
数据表是数据库最重要的组成部分之一，数据保存在数据表中
数据表由行（row）和列（column）来组成
每个数据表中至少有一列，行可以有零行或多行来组成
表名要求唯一，不要包含特殊字符，最好含义明确
创建表：CRETAE TABLE IF NOT EXISTS tbl_name(字段名称 字段类型[完整性约束条件]，字段名称 字段类型[完整性约束条件]，. . . )ENGINE=存储引擎 CHARSET=编码方式；
MySQL中的数据类型：
1、数值型：
整数型
浮点数
定点数
2、字符串类型
3、日期时间类型

