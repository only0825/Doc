## 个⼈信息

* **Jerome**/男/27/四川/大专
* ⼯作年限：4年+
* 期望职位：Java
* 期望薪资：35k

技能清单

1. 开发语⾔：GO、PHP、Java、Python、JS
2. 框架：熟练掌握SSM、SpringBoot、Laravel、TP
3. 前端技术：熟练掌握、CSS、JS、Jquery、Bootstrap、layui、AdminLTE
4. 能写规范的接口文档，良好的代码习惯（写注释）。
5. 操作系统：熟练掌握Linux 基本命令，能在该系统下搭建环境，如LAMP，LNMP；nginx负载均衡、反向代理
6. 缓存和队列：熟练使⽤redis缓存，熟悉beanstalkd。 
7. 常用工具：Git、GitLab、GitHub、Postman、Docker。
8. 英语日常交流，能阅读英文文档。



## ⼯作经历

##### CG集团 **2020年6⽉-2022年8月**

描述： 负责P3越南站部分功能开发与迭代、维护。

**盈城 2019年2⽉-2020年5月**

描述：负责泛⽬录系统、⼴告系统、词库系统的维护以及更新迭代。

**2018年5⽉-2019年1月**

**国王大厦**

职责：负责游戏后台管理的开发和维护、对接第三方支付。



## 项⽬经验

#### P3越南的域名管理系统

项目描述：域名和APP相关的配置，我负责APP相关的功能

涉及技术：Go、fasthttp、goqu、Mysql、Redis、ants、beanstalkd

技术要点：

1. APP的启动图、公告、版本的增删改查。APP的接口从Redis拿数据
2. 使用beanstalk延迟队列来做公告的自动过期。

#### P3越南站

负责模块：系统公告、登录日志、广告管理

描述：

* 系统公告：负责后台管理公告功能与前台的接口
* 登录日志：登录的日志存TD数据库
* 广告管理：增删改查、以及用延时队列实现banner自动展示与隐藏

#### 课程查询项目

项目描述：查询课程项目

涉及技术：Spring Cloud Netflix Eureka、Feign、Ribbon、Zuul、Hystrix、Maven

技术要点：

1. 服务注册中心使用 Spring Cloud Netflix Eureka
2. 主要有课程价格服务与课程列表服务，使用Feign实现服务间调用
3. 利用Ribbon实现负载均衡
4. 利用Hystrix实现断路器，当某一个服务器发生故障时，返回默认消息，将其隔离出去
5. 网关使用的Zuul

#### Java电商网站

项目描述：

涉及技术：Spring+SpringMVC+Mybatis、Mybaitis-generator、MybaitisX、Mybatis-pagehelper、Logback、Guava、Session、simditor

技术要点：

* 解决横向越权、纵向越权问题。构建高复用服务响应对象
* 回答问题修改密码功能中，UUID生成Token再使用Guava缓存实现token本地缓存
* 分类模块用的时无限层级的树状数据结构
* 产品列表使用PageHelper插件做分页
* 富文本用的simditor
* 使用BigDecimal的String构造器防止运算中丢失精度
* 分类表设计为无限层级的树状数据结构，使用递归获取子节点

#### ⼴告配置后台

项目描述：广告相关的配置

涉及技术：phpcms框架、AdminLTE 3、mysql

主要功能：

1. 登录注册以及权限管理。
2. 关键词分类对应盘⼝、关键词对应盘⼝、盘⼝信息、热⻔下载APP信息、统计代码、⼴告标语、站点模板、等配置的增删改查。

责任描述：负责前端与服务端开发，以及开发完后的维护、新增功能、迭代功能。

#### 数据统计系统

项目描述：统计每日新增泛目录数据的数据和蜘蛛、检查原站友链、百度Site与SEO综合查询、子域名、蜘蛛IP来源数据统计、蜘蛛IP来源统计图表，方便对泛目录进行监控和维护。

涉及技术：laravel、lay-ui 、mysql、redis、Echarts、phpspider

技术要点：

1. 泛目录数据的数据统计从不同的表中查询到需要展示的数据，再采用Redis存储每个泛目录数据，只用计划任务每日执行。
2. 检查原站友链使用计划任务执行检查方法，用redis存储需要检查的网站地址，检查结果存到表中，最后页面展示。
3. 使用layui进行的table.render进行表格渲染
4. 采用phpspider爬取原站在SEO综合查询的结果。
5. 使用echarts展示蜘蛛IP来源数据

责任描述：前后端开发以及后续调整和维护。

#### 对接支付接口

对接第三方或第四方的支付接口

#### 词库后台系统

项目描述：关键词排名记录功能，词库后台主为泛⽬录和站群提供关键词，将需要扩展的词⽤Python去爬取，再插⼊到词库表中。

涉及技术：phpcms框架 、mysql

关键词排名记录功能描述：

1. 主要功能为爬⾍将关键词表中的词拿去百度搜索，搜索的结果做排序并统计关键词数量为前100的站点以及历史记录
2. 由于数据较⼤采⽤了分表功能，当数据达到⼀定数量时就新增⼀张表。

#### 棋牌游戏

开发描述：团队发开发

开发周期：2个⽉

开发⼯具：Easyswoole、Egret、禅道、GitLab、Mysql、Redis

责任描述：负责服务端开发，消息推送，缓存处理。

#### 游戏管理后台

开发描述：团队发开发

开发周期：2个⽉

开发⼯具：Laravel、Vue

责任描述：负责接⼝开发

项⽬描述：历经三个版本1、前端easyui，后端Laravel 2、laravel-admin开发 3、element-ui和laravel。

