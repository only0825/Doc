## Maven介绍

* Maven是项目管理工具，对软件项目提供构建和依赖管理
* Maven是Apache下的Java开源工具
* Maven是Java项目提供了统一的管理方式，已成为业界标准


Maven核心特性：
* 项目设置遵循统一的规则，保证不同开发环境的兼容性
* 强大的依赖管理，项目依赖组件自动下载、自动更新
* 可扩展的插件机制，使用简单，功能丰富



## Maven的安装：

官网：maven.apache.org

目录结构：

bin是可执行文件夹，里面的mvn.cmd是核心执行文件
boot是引导
conf 是配置
lib 相关的依赖包

下载完毕后然后进行环境变量的配置，path中添加maven的bin目录的安装地址即可，命令行执行mvn -v测试



**Eclipse整合Maven**

Window->Preferences->Maven->Installations中其实已经整合好了。也可以自己添加



## Maven的坐标：

1、GroupId：机构或者团体的英文，采用“逆向域名”形式书写
2、ArtifactId：项目名称，说明其用途，例如：cms、oa等
3、Version：版本号，一般采用“版本+单词”形式，例如：1.0.0.RELEASE

## maven项目标准结构：

最常用的是main
pom.xml是最重要的配置文件

![image-20200518224552184](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200518224552184.png)



## Eclipse创建Maven工程

New->Other-Maven->Maven Project

勾选上Create a simple project，下一步填写Group Id , Artifact Id, Version, 和Packaging

## Maven依赖管理

1. Maven利用dependency（依赖）自动下载、管理第三方jar；

2. 在pom.xml文件中配置项目依赖的第三方组件

3. maven自动将依赖从远程仓库下载至本地仓库，并在工程中引用

如：

```xml
<dependencies>
    <dependency>
        <groupId>mysql</groupId>
        <artifactId>mysql-connector-java</artifactId>
        <version>5.1.47</version>
    </dependency>
</dependencies>
```

https://search.maven.org Maven中央仓库检索网站



## 本地仓库与中央仓库

本地仓库位置:
eclipse：Window->Preferences->Maven->User Settings中
IDEA：Settings中的Maven查看

```xml
本地仓库与中央仓库 在pom.xml中通过repositories标签添加仓库，使用maven私服来下载jar包比从国外网站直接下载速度明显提高，我们常用阿里云的Maven私服： maven.aliyun.com
<repositories>
<repository>
<--创建私服的地址-->
<id>aliyun</id>
<name>aliyun</name> <url>https://maven.aliyun.com/repository/public</url>
</repository>
</repositories>

优先从阿里云私服下载，如果私服没有，在从maven官网下载
```

## maven打包流程

通过Plugins（插件）技术实现
输出jar包的插件：maven-assembly-paugin

1.根节点上方添加插件配置

```xml
<build>
    <!-- 打jar包的配置插件 -->
    <plugins>
        <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-assembly-plugin</artifactId>
        <version>2.5.5</version>
        <!-- 配置 -->
        <configuration>
            <!-- 指定入口类 即带有main方法的类-->
            <archive>
                <manifest>
                	<mainClass>com.imooc.maven.PinyinTestor</mainClass>
                </manifest>
            </archive>
            <!-- 额外参数 -->
            <descriptorRefs>
                <!-- all in one 在打包时将所有引用的jar合并到输出的jar文件中 -->
                <descriptorRef>jar-with-dependencies</descriptorRef>
            </descriptorRefs>
        </configuration>
        </plugin>
    </plugins>
</build>
```

2.创建执行命令
在上方运行按钮处点击下拉菜单，-->Run Configurations...
-->Maven Build 右键 --> New configurgtion
--> Name:给该运行命令设置名字
--> Base directory:选择要导出的项目
--> Goals:assembly:assembly(前者为插件的名字，后者为装配的意思)
--> Apply



## Eclipse中Maven构建web工程

1. 配置好Tomcat
2. new --> maven project（创建一个标准的maven工程）
3. 将默认的jre修改为8
4. 项目右键--> p --> java Compiler 将JDK修改为对应版本
5. 在src/main下创建目录webapp（用来保存网页文件）
6. 项目右键--> p - -> Project Facets --> 点击链接 --> 勾选Dynamic Web Module选项 版本自选
7. Runtimes --> 勾选Tomcat
8. 点击下方链接 --> 将Content directory 改为 src/main/webapp --> 勾选创建web.xml
9. ok --> apply



## Web应用打包（Maven工程）

packaging代表输出格式，只有jar和war两种格式，如果没写，默认输出jar包
在version标签下
```<packaging>war</packaging>```



```xml
  <build>
    <finalName>maven-web</finalName>
  	<plugins>
  		<plugin>
			<groupId>org.apache.maven.plugins</groupId>
			<artifactId>maven.war.plugis</artifactId>
			<version>3.2.2</version>
  		</plugin>
  	</plugins>
  </build>
```

1. 播放按钮->Run Configurations->Maven build右键创建新的配置

2. 输入Name:名称，base directory指向当前的工程，Goals输入package（是maven的打包命令)

3. Apply and Run





## Maven常用命令

mvn archetype:generate 创建Maven工程结构
mvn compile 编译源代码
mvn test 执行测试用例
mvn clean 清除产生的项目
mvn package 项目打包
mvn install 安装(发布)至本地仓库



## 修改本地仓库存放路径

1.在maven安装目录的conf目录下有个settings.xml文件
2.找到localRepository标签--> 取消注释
3.将其值修改到指定目录即可，例：<localRepository>d:/maven-repo</localRepository>
4.回到eclipse，window --> p --> Maven -- Installations 确保其指向的是自己安装的maven版本
5.User setting 将User Settings 的地址修改为 D:\apache-maven-3.6.0\conf\setting.xml -->点击update settings --> Reindex --Apply

小技巧： 可以将原先的包全部cpoy到新路径，这样就不用重新下载了。