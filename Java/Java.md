



## IDEA构建jar包（打包）

配置：

1. 打开File/Project Structure，依次点击Artifacts->Add(+)->JAR->Empty，命名新建的jar包

2. 单击视图中的xxx.jar，下面会出现CREATE NANIFEST，点击并选择对应的工程，再选择Main Class.
   Class path用于让我们的jar包加载相应的目录，目前用不到，就暂时不做处理。

3. 双击Available Elements中相应的类，将其移动到左侧jar包结构中，最后点击ok。

构建：点击导航条中的Build/Build Artifacts/Build。

jar包的测试：找到新生成的out\artifacts\helloWorld\目录下的jar包，右键Show in Explorer，在文件管理器中查看该jar。并在文件管理器地址栏输入cmd进入命令行界面，输入java -jar helloWorld.jar，可以得到该jar包中的主类的执行结果。

![image-20200517035630524](C:\Users\benve\AppData\Roaming\Typora\typora-user-images\image-20200517035630524.png)



## IDEA开发JavaWeb应用

在Java Enterprise中新建一个Web Application应用，选择SDK和Java EE 版本和 Tomcat



#### IDEA中Tomcat热部署

(不用重启服务就可以加载更新类和资源，但只能在调试模式下运行)

Tomcat配置中的On frame deactivation选择Update classes and resources 

#### 修改工程首页访问目录

Tomcat配置中的Deployment中javaweb:war exploded意思是以目录形式发布，下面的Application context修改路径

#### 打包发布war包

1. 打开File/Project Structure，依次点击Artifacts->Add(+)->Web Application Archive->Empty
2. 双击右边的'项目名称' compile output (编译过后的类文件)与 Web facet resources(项目资源) 移到左边，点击OK
3. 点击菜单栏中的Build/Build Artifacts/Build。打包后的war文件项目的在out/artifacts下面
4. 测试一下，将war包放入tomcat/webapps下面，tomcat/bin/startup.bat运行tomcat，打开浏览器访问项目。

