两张表的建表SQL:
```sql
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `student`;
CREATE TABLE `student`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `stuno` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '学号',
  `stuname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '姓名',
  `sex` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '性别',
  `classno` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '班级编号',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
INSERT INTO `student` VALUES (1, '20181101', '梅长苏', '男', 'Class001');
INSERT INTO `student` VALUES (2, '20181102', '萧景炎', '男', 'Class001');
INSERT INTO `student` VALUES (3, '20181103', '宫羽', '女', 'Class001');
INSERT INTO `student` VALUES (4, '20181104', '霓凰', '女', 'Class003');
INSERT INTO `student` VALUES (5, '20181105', '瓜娃子', '男', 'Class003');
SET FOREIGN_KEY_CHECKS = 1;
```
```sql
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS `classes`;
CREATE TABLE `classes`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '班级id',
  `classno` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '班级编号',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '班级名称',
  `major` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '专业',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
INSERT INTO `classes` VALUES (1, 'Class001', '一班', '计算机');
INSERT INTO `classes` VALUES (2, 'Class002', '二班', '计算机');
INSERT INTO `classes` VALUES (3, 'Class003', '三班', '会计');
INSERT INTO `classes` VALUES (4, 'Class004', '四班', '会计');
SET FOREIGN_KEY_CHECKS = 1;
```



classes表和student表是一对多的关系，外键是classno

实体类：

```java
public class Student {
    private Integer id;
    private String stuno;
    private String stuname;
    private String sex;
    private String classno;
    private Classes classes; // 一个学生对应一个班级
    
    ......(get和set等方法)
```

```java
public class Classes {
    private Integer id;
    private String classno;
    private String name;
    private String major;
    private List<Student> students;  // 一个班级有多个学生
    ......(get和set等方法)
```

student的Mapper文件：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper
        PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="studentDetail">
    <select id="selectByClassno" parameterType="String"
            resultType="com.imooc.mybatis.entity.Student">
        select * from student where classno = #{value}
    </select>

    <resultMap id="sqlStudent" type="com.imooc.mybatis.entity.Student">
        <id column="id" property="id"/>
        <result column="classno" property="classno"/>
        <association property="classes" select="classes.selectByClassno" column="classno"/>
    </resultMap>
    <select id="selectManyToOne" parameterType="String" resultMap="sqlStudent">
        select * from student limit 0,5
    </select>
</mapper>
```

classes的Mapper文件：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper
        PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="classes">
    <resultMap id="sqlClasses" type="com.imooc.mybatis.entity.Classes">
        <id column="classno" property="classno"/>
        <collection property="students" select="studentDetail.selectByClassno" column="classno"/>
    </resultMap>
    <select id="selectOneToMany" parameterType="String" resultMap="sqlClasses">
        select * from classes where name = #{value}
    </select>

    <select id="selectByClassno" parameterType="String" resultType="com.imooc.mybatis.entity.Classes">
        select * from classes where classno = #{value}
    </select>
</mapper>
```

别忘了再mybatis-config.xml中添加映射

```xml
<mappers>
    <mapper resource="mappers/classes.xml"/>
    <mapper resource="mappers/student.xml"/>
</mappers>
```

一对多测试：

```java
@Test
public void classOneToMany() throws Exception {
    SqlSession session = null;
    try {
        session = MyBatisUtils.openSession();
        List<Classes> list = session.selectList("classes.selectOneToMany", "一班");
        for(Classes classes:list) {
            List<Student> student = classes.getStudents();
            for(Student tmp:student){
                System.out.println(tmp.getId() + ":" + tmp.getStuno() + ":" + tmp.getStuname() + ":" + tmp.getSex() + ":" + tmp.getClassno());
            }
        }
    } catch (Exception e) {
        throw e;
    } finally {
        MyBatisUtils.closeSession(session);
    }
}
```

多对一测试：

```java
@Test
public void studentManyToOne() throws Exception {
    SqlSession session = null;
    try {
        session = MyBatisUtils.openSession();
        List<Student> list = session.selectList("studentDetail.selectManyToOne");
        for (int i = 0; i < list.size(); i++) {
            Student stu = list.get(i);
            Classes cla = stu.getClasses();
            System.out.println(cla.getId() + ":" + cla.getClassno() + ":" + cla.getName() + ":" + cla.getMajor());
        }
    } catch (Exception e) {
        throw e;
    } finally {
        MyBatisUtils.closeSession(session);
    }
}
```