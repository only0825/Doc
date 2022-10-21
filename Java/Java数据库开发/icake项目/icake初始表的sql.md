```sql
drop database if exists icake;

create database icake;
use icake;


drop table if exists account;

drop table if exists cake;

drop table if exists catalog;

/*==============================================================*/
/* Table: account                                               */
/*==============================================================*/
create table account
(
   id                   int not null auto_increment,
   account              varchar(20),
   password             varchar(20),
   nick_name            varchar(20),
   primary key (id)
);

/*==============================================================*/
/* Table: cake                                                  */
/*==============================================================*/
create table cake
(
   id                   int not null auto_increment,
   title                varchar(20),
   cid                  int,
   image_path           varchar(100),
   price                double,
   taste                varchar(20),
   sweetness            int,
   weight               double,
   size                 int,
   material             varchar(100),
   status               varchar(20),
   primary key (id)
);

/*==============================================================*/
/* Table: catalog                                               */
/*==============================================================*/
create table catalog
(
   id                   int not null auto_increment,
   title                varchar(20),
   pid                  int,
   info                 varchar(100),
   primary key (id)
);

alter table cake add constraint FK_Reference_1 foreign key (cid)
      references catalog (id) on delete restrict on update restrict;

alter table catalog add constraint FK_Reference_2 foreign key (pid)
      references catalog (id) on delete restrict on update restrict;



insert into account(id,account,password,nick_name) values(10000,'admin','admin','管理员');

insert into catalog(id,title) values(10000,'蛋糕');
insert into catalog(id,title,pid) values(11000,'婚礼',10000);
insert into catalog(id,title,pid) values(11100,'西式',11000);
insert into catalog(id,title,pid) values(11101,'多层',11100);
insert into catalog(id,title,pid) values(11102,'花朵',11100);
insert into catalog(id,title,pid) values(11103,'造型',11100);
insert into catalog(id,title,pid) values(11200,'传统',11000);
insert into catalog(id,title,pid) values(11201,'复古',11200);
insert into catalog(id,title,pid) values(12000,'生日',10000);
insert into catalog(id,title,pid) values(12100,'儿童',12000);
insert into catalog(id,title,pid) values(12101,'男孩',12100);
insert into catalog(id,title,pid) values(12102,'女孩',12100);
insert into catalog(id,title,pid) values(12103,'满月',12100);
insert into catalog(id,title,pid) values(12200,'老人',12000);
insert into catalog(id,title,pid) values(12201,'寿糕',12200);
insert into catalog(id,title,pid) values(12300,'成人',12000);
insert into catalog(id,title,pid) values(12301,'夫妻',12300);
insert into catalog(id,title,pid) values(12302,'亲戚',12300);
insert into catalog(id,title,pid) values(12303,'朋友',12300);
insert into catalog(id,title,pid) values(13000,'节日',10000);
insert into catalog(id,title,pid) values(13100,'家庭',13000);
insert into catalog(id,title,pid) values(13101,'母亲节',13100);
insert into catalog(id,title,pid) values(13102,'父亲节',13100);
insert into catalog(id,title,pid) values(13103,'儿童节',13100);
insert into catalog(id,title,pid) values(13104,'情人节',13100);
insert into catalog(id,title,pid) values(13200,'送人',13000);
insert into catalog(id,title,pid) values(13201,'教师节',13200);
insert into catalog(id,title,pid) values(13202,'圣诞节',13200);
insert into catalog(id,title,pid) values(14000,'专用',10000);
insert into catalog(id,title,pid) values(14100,'聚会',14000);
insert into catalog(id,title,pid) values(14101,'商业宴会',14100);
insert into catalog(id,title,pid) values(14102,'同学聚会',14100);
insert into catalog(id,title,pid) values(14103,'单身派对',14100);
insert into catalog(id,title,pid) values(14200,'纪念',14000);
insert into catalog(id,title,pid) values(14201,'组织庆典',14200);
insert into catalog(id,title,pid) values(14202,'活动仪式',14200);

insert into cake values(100,'芒果奶油蛋糕',11201,'/download/images/1537290653487.png',198,'巧克力',3,12,2,'芒果、 百香果酱','特卖');
insert into cake values(110,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(111,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(112,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(113,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(114,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(115,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(116,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(117,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(118,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(119,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(120,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(121,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(122,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(123,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(124,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(125,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(126,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(127,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(128,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(129,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(130,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(131,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(132,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(133,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','');
insert into cake values(134,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','');
insert into cake values(135,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','推荐');
insert into cake values(136,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','推荐');
insert into cake values(137,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','推荐');
insert into cake values(138,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','推荐');
insert into cake values(139,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','推荐');
insert into cake values(140,'浓情巧克力',11201,'/download/images/1537289675245.jpg',298,'奶油',3,8,4,'巧克力','推荐');
insert into cake values(141,'奶香巧克力',11201,'/download/images/1537289838080.png',298,'巧克力',3,8,2,'巧克力','推荐');
insert into cake values(198,'桂圆冰淇淋',11201,'/download/images/1537289865096.png',198,'巧克力',3,12,1,'白兰地酒 、金黄桂圆肉干','推荐');
insert into cake values(199,'爱尔兰咖啡',11201,'/download/images/1537289956018.jpg',268,'布丁',3,12,2,'爱尔兰威士忌 、咖啡慕斯和咖啡坯','推荐');
insert into cake values(200,'深艾尔',11201,'/download/images/1537289970651.jpg',268,'奶油',3,10,1,'帝国世涛、樱桃果酱、 修道院啤酒','推荐');
```

