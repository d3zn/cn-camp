# cn-camp

# module 1，module 2 go语言基础
官方文档：https://go.dev/doc/  
如何写出良好的go代码：https://go.dev/doc/effective_go  
重点：  
内存分配：TCMalloc  
内存回收：三色标记，GMP   
并发编程  
性能分析：pprof

# module 3 docker核心技术
传统分层架构与微服务的区别
docker：基于linux内核的cgroup，namespace Union FS等技术，进行操作系统层面的资源隔离  

相关命令  
```shell
#启动
docker run
  -it exec "command"  交互
  -d 后台运行
  -p 端口映射
  -v 磁盘挂载
 
# 停止
docker stop

# 启动已终止容器
docker start

# 查看容器进程
docker ps

# 查看容器细节
docker inspect <containerid>

# 拷贝文件到容器内
docker cp file1 <containerid>:/file-to-path

#打包容器镜像
docker build -f <Dockerfile> -t <tagname> . #根据dockerfile打镜像
docker push <tagname> 将容器镜像推送至hub
```
dockerfile参考文档：https://docs.docker.com/engine/reference/builder/  

# module 4, module 5, module 6 module 7 K8S原理
官方文档介绍：https://kubernetes.io/zh/docs/concepts/overview/what-is-kubernetes/
## 命令式  
我要你做什么，怎么做，严格按照我说的做  
## 声明式  
我需要什么，而不是告诉你应该怎样做  
- 直接声明：直接告诉你我需要什么  
- 间接声明：我把我的需求放在特定的地方，请在方便的时候拿出来处理  
  声明式是幂等的、面向对象的  

k8s将一切资源抽象成对象，核心对象包括
- Node：计算节点抽象，用来描述计算节点的资源抽象、健康状态等，可以是虚拟机或者是物理机（集群管理）
- Namespace：资源隔离的基本单位，是一个逻辑概念， 可以简单理解为文件系统中的目录结构
- Pod：用来描述应用实例的基本单位，调度的基本单元。包括镜像地址、资源需求等，k8s中最核心的对象，也是打通应用和基础架构的核心（作业管理）
- Service：将一组应用实例（pod）看做一个整体对外发布为服务，本质上是负载均衡和域名服务的声明（服务发现）

## k8s相关组件
###主节点 Master Node

#### kube-apiserver  整个k8s的前端  
是个REST API，接收请求，分发请求，并将数据存入到etcd中  
认证（Authentication），鉴权（Authorization），准入（Admission）  
整个控制平面（Control Plane）中唯一代用用户可访问API及用户交互的组件，API服务器会暴露一个REST kubernetes API 并使用JSON格式的清单文件（mainfest files）  
控制平面其他组件只和API server交互，由API server协调控制平面各组件动作，完成整个资源的调度  
接收并分发来自etcd的消息事件  

#### etcd  整个集群数据存储  
etcd 是兼具一致性和高可用性的键值数据库，可以作为保存 Kubernetes 所有集群数据的后台数据库。  
watch 模式：get对象时可以实时查看状态，保持长连接，当对象状态发生变化，etcd会主动推送（事件）给API server，可以理解为一个消息中间件  

#### kube-controller-manager
他本身管理了很多控制器，比如deployment控制器， replicaSet控制器，NodeLifeCycle控制器等等，每个控制器管理不同的对象  
从逻辑上讲，每个控制器都是一个单独的进程， 但是为了降低复杂性，它们都被编译到同一个可执行文件，并在一个进程中运行。  
这些控制器包括:  
- 节点控制器（Node Controller）: 负责在节点出现故障时进行通知和响应
- 任务控制器（Job controller）: 监测代表一次性任务的 Job 对象，然后创建 Pods 来运行这些任务直至完成
- 端点控制器（Endpoints Controller）: 填充端点(Endpoints)对象(即加入 Service 与 Pod)
- 服务帐户和令牌控制器（Service Account & Token Controllers）: 为新的命名空间创建默认帐户和 API 访问令牌

#### kube-scheduler
控制平面组件，负责监视新创建的、未指定运行节点（node）的 Pods，选择节点让 Pod 在上面运行。  
调度决策考虑的因素包括单个 Pod 和 Pod 集合的资源需求、硬件/软件/策略约束、亲和性和反亲和性规范、数据位置、工作负载间的干扰和最后时限。  

#### 一个基本的工作协同流程
1. 用户发起创建一个pod请求，请求到API server
2. API server转发请求到scheduler
3. Scheduler 看pod需要多少资源，看当前集群的资源利用状况（利用状况通过kubelet 上报），上报的状况会存放在etcd中，scheduler有能力知道目前集群状态。
4. Scheduler 根据调度策略把pod和node信息绑定，告诉API server ，API server 把pod里的nodeName属性更新到etcd中
5. kubelet 获取etcd中的 pod node绑定信息，如果当前node的pod信息与我相关，则进行pod的创建


### 工作节点 Worker Node
#### kubelet
- 负责调度到对应Node的Pod的声明周期管理，执行Job并将Pod状态报告给主节点的渠道（类似于easyops agent）
- 通过容器运行时（contrainer runtime，拉取镜像、启动和停止容器）运行容器
- 定期探活
  不会管理非k8s创建的容器，比如手工docker run  
  自带cAdvisor组件可以收集容器的健康状态或者资源用量  

#### kube-proxy
- 负责Node的网络，在主机上维护网络规则并执行连接转发，
- 实现集群内部或外部的网络会话与Pod进行网络通信
- 对正在服务的pod进行负载均衡

所有的控制平面组件都在kube-system命名空间下  

文章参考：  
https://www.cnblogs.com/wwchihiro/p/9261607.html   k8s架构  
https://kubernetes.io/zh/docs/concepts/overview/components/  k8s官方文档  

### etcd
etcd是CoreOS基于Raft协议开发的分布式key-value数据库，可用于服务发现，共享配置，一致性保障（数据库选主，分布式锁等）  

提供存储以及获取数据的接口，它通过协议保证 Etcd 集群中的多个节点数据的强一致性。用于存储元信息以及共享配置。  
提供监听机制，客户端可以监听某个key或者某些key的变更（v2和v3的机制不同，参看后面文章）。用于监听和推送变更。  
提供key的过期以及续约机制，客户端通过定时刷新来实现续约（v2和v3的实现机制也不一样）。用于集群监控以及服务注册发现。  
提供原子的CAS（Compare-and-Swap）和 CAD（Compare-and-Delete）支持（v2通过接口参数实现，v3通过批量事务实现）。用于分布式锁以及leader选举。  
摘自：https://blog.csdn.net/zl1zl2zl3/article/details/79627412

容器内访问etcd的数据  
```shell
# 获取以/开头的key列表
etcdctl --endpoints https://localhost:2379 \--cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt get --keys-only --prefix /

# 获取指定key，get的数据是以pb的形式返回，grpc协议
etcdctl --endpoints https://localhost:2379 \--cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt get --prefix /registry/services/specs/default/kubernetes

# 监听对象变化
etcdctl --endpoints https://localhost:2379 \--cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key \
--cacert /etc/kubernetes/pki/etcd/ca.crt watch --prefix /registry/services/specs/default/kubernetes
```

## API Server（重要，在module6中详细介绍）

kube-apiserver 是k8s最重要的核心组件之一，主要功能：   
提供集群管理的REST API接口  
- 认证 Authentication
- 授权 Authorization
- 准入 Admission （Mutating &Valiating） 类似于参数是否合法
  API server 提供etcd数据缓存以减少集群对etcd的访问  
  模块之间数据通信的枢纽，其他模块通过API server查询或者修改，只有APIserver 直接访问etcd

## Controller manager
整个集群的大脑，保证集群联动的关键  
确保k8s遵循声明式系统规范，确保系统的真实状态（actual state）与用户定义的期望状态（desired state）一致  
多个控制器的组合，每个Controller都是一个control loop，负责侦听管控对象，当对象发生变更时完成配置  
配置失败后通常会触发自动重试，整个集群会在控制器的重试机制下确保最终一致性（Eventual Consistency）  

### Informer 监听对象
处理对象的事件，添加、删除、更新  
事件处理器在事件发生后，将事件处理函数key放入队列中  
worker进程消费key，执行函数做配置  
Lister提供缓存接口，从client-cache拿到对象缓存状态，不需要去api server取  

### 控制器工作原理
每个controller都有自己的监听事件，监听到事件后，做相应的动作  
当kubectl create -f xxx.yaml 会请求到apiserver  
deployment-controller会监听到api-server的创建请求，deployment-controller有自己的意义，代表要根据指定的模板创建ReplicaSet（副本集）  
deployment-controller会发起创建副本集的请求到api-server  
同样的，replicaSet-controller也在监听来自api-server的事件请求，当创建完副本集之后，会发起创建pod的请求  
同样的，scheduler会监听来自api-server的pod的事件请求，scheduler创建pod并选择节点进行绑定，所以说scheduler是个特殊的controller  
kubelet的职责是真正调起pod的角色，通过CRI创建pod，通过CNI创建容器网络，通过CSI挂载容器存储  
pod的创建之初的nodeName为空，由scheduler计算出最佳的node后 更新上去  

### 常用命令
```shell
# 查看资源详细信息
kubectl get pod -o yaml -w
- o # 以什么样的格式输出： yaml 、json、wide
- w # 动态查看

# 查看资源详细信息和相关event
kubectl describe pod  <podname>

# 进入容器
kubectl exec -it <podname> bash # 前提是容器的镜像有bash

# 查看pod的标准输出和错误输出（stdout，stderr）
kubectl logs <podname>
```
## k8s资源对象定义
每个API对象都有四大类属性  
- TypeMeta  
  - Group  
  - Kind  
  - Version  
- MetaData  
  - Label
  - Annotation
  - Finalizer
  - ResourceVersion
- Spec
- Status

## k8s核心资源
### Node
pod运行的真正主机，可以是物理机或者是虚拟机

### Namespace
一组资源和对象的抽象集合

### Pod
一组紧密关联的容器集合，k8s调度基本单位

### ConfigMap
保存非机密性的数据，键值对。  
使用时，Pod可以用于环境变量，命令行参数，或者存储卷中的配置文件  
将环境配置信息和容器镜像解耦，便于应用配置修改

### Secret
保存和传递密码、密钥、认证凭证等敏感信息的对象

### Service
应用服务的抽象， 通过labels为应用提供负载均衡和服务发现，匹配labels的Pod IP和端口类别组成endpoints，由kube-proxy负责将服务IP负载均衡到这些endpoints上

### Replica Set
Pod是单个应用实例的抽象，要构建高可用应用，需要构建多个同样的副本，提供同一个服务，由此抽象出ReplicaSet  
允许用户定义Pod副本数
每个Pod作为无状态成员进行管理
保证用户期望数量的pod正常运行，挂了某个pod会被自动拉起

### Deployment
表示对k8s进群的一次更新操作
滚动升级一个服务，实际是创建一个新的RS，然后逐渐把老的RS的pod替换为新的RS，这样一个复合操作用RS不好描述，所以用一个更加通用的Deployment描述

### Stateful Set
管理有状态应用

### Job
控制批处理型任务的API对象

### Daemon Set
后台支撑服务集

## Raft协议
参考文档：http://thesecretlivesofdata.com/raft/
