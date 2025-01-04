### 开发须知

1. data/dev 下面的docker-compose 启动依赖项.

执行: 
```shell
docker compose up -d 
```

2. 初始化Mysql 

```shell
# 进入mysql的容器
docker exec -it dev-mysql-1 bash 

mysql -u root -p124x8Xawdasdx1r140xs$
```

初始化数据库命令
```sql
CREATE DATABASE kf CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```


3. 目录结构:
```shell

.
├── Dockerfile      # 构建本项目的 Dockerfile 文件
├── Makefile        # Makefile  
├── cmd             # main 包
├── config          # 配置包
│   ├── config.go
│   └── load.go
├── config.yaml     # 开发用的配置文件
├── data            # 数据目录，用来保存开发、线上、静态资源等文件, 开发中的文件会被 gitignore 忽略
│   ├── dev   # 开发目录
│   ├── prod  # 线上配置目录
│   └── static  # 静态资源目录
├── develop.md        # 开发手册
├── docker-compose.yaml   # 线上启动本项目的 docker compose 文件
├── domain            # 领域层 ( 目前比价空洞，逻辑暂时都写 controller 了)
│   ├── converter # 转换层
│   ├── factory
│   └── repository  # 一般放通用查询逻辑, 写逻辑直接在 controller 执行了 
├── endpoints             # 外部触点，如 mq、定时任务、http、等
│   ├── common      # 通用文件
│   ├── cron        # 定时任务包
│   ├── http        # 接口层
│   └── nsq         # mq消费
├── go.mod          
├── go.sum
├── infrastructure        # 持久层
│   ├── caches      # 缓存
│   ├── mysql       # mysql
│   ├── nsq         # 队列
│   └── redis       # redis 
├── pkg                   # 存储面向过程无状态的通用方法 
│   ├── consumer
│   ├── utils
│   ├── version
│   └── xerrors
└── push-next-tag.sh      # 执行该脚本会打一个tag、并且推到github、执行ci、cd

```
4. 