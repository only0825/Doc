## 个⼈信息

* **Gray**/男/27/四川/大专
* ⼯作年限：4年+
* 期望职位：PHP
* 期望薪资：40k

技能清单

1. 开发语⾔：PHP、GO、Java、Python、JS
2. 框架：TP6、Laravel、SSM、SpringBoot、Gin、Zinx
3. 前端技术：熟练掌握、CSS、JS、Jquery、Bootstrap、layui、AdminLTE
4. 能写规范的接口文档，良好的代码习惯（写注释, 代码合理分快,return空行等）。
5. 操作系统：熟练掌握Linux 基本命令，能在该系统下搭建环境，如LAMP，LNMP；nginx负载均衡、反向代理
6. 缓存和队列：熟练使⽤redis缓存和队列。 
7. 常用工具：Git、GitLab、GitHub、Navicat、DataGrip、PhpStorm、Postman、Xmind、Typora。
8. 英语日常交流，能阅读英文文档。



## ⼯作经历

##### CG集团 **2020年6⽉-至今**

描述： 负责体育项目和越南站部分功能开发与迭代、维护。

**盈城 2019年2⽉-2020年5月**

描述：负责泛⽬录系统、⼴告系统、词库系统的维护以及更新迭代。

**2018年5⽉-2019年1月**

**国王大厦**

职责：负责游戏后台管理的开发和维护、对接第三方支付。



## 项⽬经验

#### 海豹体育直播

职责：体育直播项目，独自负责项目前后台的开发、迭代、上线、维护。

技术：PHP7.2、TP6、workman-gateway、Mysql、Redis、Go、gorilla

功能：

* 前台赛事模块，查询功能有日期、分类(篮球足球等)、联赛、比赛状态条件，预约功能，积分榜

* 前台开播模块是用的OBS推到腾讯云后，前端再通过拉流地址直播。
* 聊天室，用的workman的gateway实现，发言、广播、禁言功能
* 后台主要是用户管理、直播间管理、赛事管理、数据统计，会统计直播的在线人数、历史人数以及APP渠道统计

* 对接第三方：赛事数据、微信登陆、广告平台上报、腾讯云拉流鉴权
* 即时比分推送服务，能及时给前端推送最新的比分和指数数据。go-cron做数据缓存，gorilla做ws推送服务



#### 越南站

负责模块：系统公告、登录日志、广告管理

描述：

* 系统公告：负责后台管理公告功能与前台的接口
* 登录日志：登录的日志存TD数据库
* 广告管理：增删改查、以及用延时队列实现banner自动展示与隐藏



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
