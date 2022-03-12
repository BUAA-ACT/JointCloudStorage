# Joint Cloud Storage 云际存储

![GitHub repo size](https://img.shields.io/github/repo-size/BUAA-ACT/JointCloudStorage)
![GitHub last commit](https://img.shields.io/github/last-commit/BUAA-ACT/JointCloudStorage)

# 介绍

Joint Cloud Storage 是一个分布式云际存储系统。云存储服务提供者根据云际标准，将自身提供的存储服务、交易服务等服务标准化，使得多个云服务访问方式统一，实现互操作性。在可互操作的基础上，云存储服务通过统一定义的接口互相协作，形成云际存储系统。此时每个云存储服务不仅仅只提供自身的存储能力，而是代表着整个云际存储服务的能力，即每个云存储服务都能作为云际存储服务的入口。

# 系统架构

Joint Cloud Storage 由多家云存储服务商的多个数据中心共同构成，每个数据中心内有控制台（portal）、后台管理模块（httpserver）、数据传输模块（transporter）、调度器（scheduler）和存储服务（storage service）五个模块。在这个模型中，数据中心内的五个模块都是云服务商自行管理的，云服务商之间先对存储能力共享的机制达成一致后，自主决定是否加入云际协作系统中。如此一来，云际协作的参与主体就是现有的各家云厂商，云际系统也就实现了完全的去中心化。五个模块中我们自行实现了控制台、后台管理模块、数据传输模块和调度器四个模块，用来提供用户访问接口和云际协作接口，而存储服务则使用各家云厂商的对象存储服务。

![系统架构.png](https://i.loli.net/2021/04/20/hFfOqJGTKYCb7lR.png)

# 使用方法

## 运行环境

**软件**

`docker`  `docker-compose` `nginx`

**安全组**

| 端口 | 用途        |
| ---- | ----------- |
| 80   | nginx       |
| 8082 | scheduler   |
| 8083 | transporter |
| 2181 | zookeeper   |

**zookeeper**

```bash
docker run --name zk -d -p 2181:2181 zookeeper
```

**nginx**

宿主机上配置nginx并绑定80和443端口，将发往`*.jointcloudstorage.cn`的请求转发至`portal`容器所对应的8080端口。

## 所需文件

以下配置文件的模板位于项目的`configs`目录下，需要复制到所有服务器的`~/jcspan`中

| 文件名                         | 用途                |
| ------------------------------ | ------------------- |
| docker-compose.yml             | docker-compose脚本  |
| nginx.conf                     | nginx配置文件       |
| httpserver.properties          | httpserver配置文件  |
| transporter_config.json.docker | transporter配置文件 |
| mongo-init.js                  | mongo初始化         |

之后执行 `sed -i "s/<this-cloud-id>/$(hostname)/g"` 命令，将模板中的 `<this-cloud-id>` 置换为主机名，如 `aliyun-hangzhou`。

## 启停命令

**前台启动**

```bash
docker-compose up
```

**后台启动**

```
docker-compose up -d
```

**查看日志**

```
docker-compose logs -f [SERVICE...]
```

**停止**

```
docker-compose down
```
