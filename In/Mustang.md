* #### 体育直播项目和综合盘

  主要负责

  涉及技术：Go、fasthttp、goqu、Mysql、Redis、ants、beanstalkd

  负责：

  * 系统公告：负责后台管理公告功能与前台的接口
  * 登录日志：登录的日志存TD数据库
  * 广告管理：增删改查、以及用延时队列实现banner自动展示与隐藏
  * 301管理后台： 使用beanstalk延迟队列来做公告的自动过期。APP的启动图、公告、版本的增删改查。APP的接口从Redis拿数据