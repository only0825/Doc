[TOC]

## Docker是什么？

* Docker是一个用来装程序及其环境的容器，属于Linux容器的一种封装，提供简单易用的容器使用接口。它是目前最流行的Linux容器解决方案。

- IT 软件中所说的 “Docker” ，是指容器化技术，用于支持创建和使用 [Linux® 容器](https://www.redhat.com/zh/topics/containers)。

#### 什么是 Linux 容器？

Linux 容器技术能够让您对应用及其整个运行时环境（包括全部所需文件）一起进行打包或隔离。从而让您可以在不同环境（如开发、测试和生产等环境）之间轻松迁移应用，同时还可保留应用的全部功能。

## Docker的用途

**提供统一的环境**

* 在本地测试其他人软件，或者把软件和环境一起打包，给其他人进行测试

**提供快速拓展，弹性伸缩的云服务**

* 比如双十一时候业务量会比平时多好几倍，为了应对提前扩容是个好办法，但是平时又不需要用到，就会造成浪费，一种方法是可以节前把机器运过来，拓展，开机，过了大型活动后下线，这种就比较麻烦。用了docker后，可以只进行简单配置，10台服务器可以变成上千台，节日时就能快速进行拓展，后续又快速的把能力进行转移下线，节省资源。有了标准化程序之后，不管在什么机器上，直接把它下载过来并启动，而无需担心有什么问题。

**防止其他用户的进程把服务器资源占用过多**

* 使用docker可以把不同内容进行隔离，对方发生的错误只会在固定范围内产生影响，不会对其他模块产生牵连



## Docker的特点

### docker特点

##### 标准化

1. 运输方式（只需要通过命令即可把程序和环境从一个机器运到另一个机器上）

2. 存储方式（程序和环境的存储，不需要关心程序存在什么地方，只需要关心使用或停止它的时候执行对应的命令即可）

3. API接口（不需要tomcat，rabbitmq等应用等自己的命令了，使用同样命令即可控制所有应用）

   

**灵活**：即使是最复杂的应用也可以集装箱化
轻量级：容器利用并共享主机内核

**便携式**：可以在本地构建，部署到云，并在任何地方运行
   docker带来的好处：
   1，开发团队得到的好处：1）可以完全控制环境，之前这些操作是依赖运维工程师的，现在就不需要依赖他们，更加灵活，也降低了风险。



## Docker的组成、架构、重要概念

![image-20210719021908745](/Users/sofia/Library/Application Support/typora-user-images/image-20210719021908745.png)

![image-20210719021928754](/Users/sofia/Library/Application Support/typora-user-images/image-20210719021928754.png)

**image镜像**：本质就是文件，包括应用程序文件，运行环境等存储：联合文件系统，UnionFS，是分层文件系统，把不同文件目录挂到一个虚拟的文件系统下，具体分几层可以自己控制。每一层构建好后，下一层不再变化，上面一层是以下面一层作为基础的，最终多层组成了文件系统。

![image-20210719022218903](/Users/sofia/Library/Application Support/typora-user-images/image-20210719022218903.png)

**Container容器**：

* 镜像类似于Java中的类，而容器就是实例
* 容器的这一层是可以修改的，而镜像是不可以修改的
* 同一个镜像可以生成多个容器独立运行，而他们之间没有任何的干扰

**Repositories仓库**：

1. 相当于一个中转站，类似百度网盘，把文件上传上去，如果要使用就自己去下载
2. 地址：- 国外：hub.docker.com- 国内：https://c.163yun.com/hub#/m/home/
3. 也分为公有和私有可以自己搭建一个镜像中心

**client和deamon**：

* client：客户端，提供给用户一个终端，用户输入docker提供的命令来管理本地或远程的服务器

- deamon：服务端守护进程，接收client发送的命令并执行相应的操作





## Docker的安装

### Windows和Mac下安装Docker

https://www.docker.com/ 官网下载即可

### 在Cent OS安装Docker

1. **先有一个Cent OS 7.6系统** 这个很重要，不同版本安装的时候是不一样的。
   查看CentOS版本 cat /etc/redhat-release 

2. **用root账户登录进去**

3. **配置国内yum源** 
   wget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
   yum clean all
   yum makecache

4. **卸载旧版本**
   较旧的Docker版本称为docker或docker-engine。如果已安装这些程序，请卸载它们以及相关的依赖项。
   yum remove docker \
   docker-client \
   docker-client-latest \
   docker-common \
   docker-latest \
   docker-latest-logrotate \
   docker-logrotate \
   docker-engine
   如果yum报告未安装这些软件包，也没问题。

5. **更新yum**
   yum check-update
   yum update

6. **安装所需的软件包**
   yum install -y yum-utils \
   device-mapper-persistent-data \
   lvm2

7. **使用以下命令来设置稳定的存储库**
   sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo

8. **查看docker版本**
   yum list docker-ce --showduplicates | sort -r

9. **安装指定的版本**
   yum install docker-ce-18.09.0 docker-ce-cli-18.09.0 containerd.io

10. **Docker 是服务器----客户端架构。命令行运行docker命令的时候，需要本机有 Docker 服务。用下面的命令启动**
    systemctl start docker

11. **安装完成后，运行下面的命令，验证是否安装成功**
    docker version	或者	docker info



## 第一个Docker容器

**下载镜像**

* docker pull [OPTIONS] NAME[:TAG]  拉取镜像

* docker images[OPTIONS] [REPOSITORY[:TAG]]  查看当前左右的镜像


**运行镜像**

* docker run [OPTIONS] IMAGE [COMMAND] [ARG...]



## 运行Nginx镜像，并访问到Docker容器内部

**docker images** 查看镜像

**docker pull xx** 拉取镜像

**docker run xx** 运行镜像

**docker ps** 查看当前容器列表

**docker run -d xxx 让docker**容器在后台运行

 **docker exec  -it [OPTIONS] CONTAINER COMMAND [ARG...**] 访问容器内部  

it：把容器内部终端映射到当前终端

例如：**docker exec -it 45 bash**    	 bash是启动终端 45是CONTAINER ID

