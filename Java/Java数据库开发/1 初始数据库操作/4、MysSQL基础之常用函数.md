## MyISAM存储引擎

1、默认MyISAM的表会在磁盘中产生3个文件：.frm表结构文件、.MYD数据文件、.MYI索引文件，

2、可以在建表的时候指定数据文件和索引文件的存储位置，只有MyISAM支持，方法：
DATA DIRECORY[=]数据保存的绝对路径，或INDEX DIRECORY[=]索引文件保存的绝对路径

3、MyISAM单表最大支持的数据量2的64次方条记录，每个表最多可以建立64个索引

4、如果是复合索引，每个复合索引最多包含16个列，索引值最大长度是1000B

5、MyISAM引擎的存储格式：
①定长（FIXED静态）是指字段中不包含VARCHAR/TEXT/BLOB等
②动态（DYNAMIC）只要字段中包含了VARCHAR/TEXT/BLOB等
③压缩（COMPRESSED）myisampack创建
在创建数据表时，可以在ENGINE参数后添加 ROW_FORMAT=FIXED/DYNAMIC进行格式化

查看引擎存储格式: SHOW TABLE STATUS LIKE 'tbl_name'\G



## InnoDB存储引擎

1、设计遵循ACID模型，支持事务，具有从服务崩溃中恢复的能力，能够最大限度保护用户的数据
原子性（Atomiocity）、一致性（Consistency）、隔离性（Isolation）、持久性（Durability）；

2、支持行级锁，可以提升多用户并发时的读写性能；

3、支持外键，保证数据的一致性和完整性；

4、拥有自己独立的缓冲池，常用的数据和索引都存在缓存中；

5、对于INSERT、UPDATE、DELETE操作，使用一种change buffering的机制来自动优化，还可以提供一致性的读，并且还能够缓存变更的数据，减少磁盘I/O，提高性能

6、创建InnoDB表之后会产生两个文件：①.frm表结构文件②.ibd，数据和索引存储空间中

7、所有的表都需要创建主键，最好配合上AUTO_INCRAMENT，也可以放到经常查询的列作为主键