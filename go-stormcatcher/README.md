# FileStorm 节点的安装

FileStorm底层使用IPFS协议做文件存储，在此基础上实现数据加密，区块链确权，文件冗余复制，自动修复，数据去重，存储奖励等功能。FileStorm所有节点将组成一个私有的IPFS网络。下面我们来介绍一下怎样生成一个FileStorm网络。

本指南使用Ubuntu系统做演示。

使用说明请看[examples](./examples/README.md)

## 安装

### 下载

可以到github [release](https://github.com/filestorm-fst/go-stormchain/releases/) 页下载最新的 fstorm 节点程序。目前仅支持Linux版本。

也可以拿到程序链接后，在服务器上通过下面的指令下载

``````
wget https://github.com/filestorm-fst/go-stormchain/releases/download/1.0/fstorm1.0-linux.tar.gz
``````

### 解压

``````
tar xvfz fstorm1.0-linux.tar.gz
rm fstorm1.0-linux.tar.gz
``````

### 安装

将解压的目录`go-ipfs/ipfs`复制到，然后删除安装文件，IPFS就安装好了。

```
mv go-ipfs/ipfs /usr/local/bin/ipfs
rm -R ./go-ipfs
```

### 初始化

```
ipfs init
```

### 修改端口（可选）

fstorm节点需要使用三个端口4001，5001和8080，如果需要修改，可以用如下指令
```
vi ~/.ipfs/config
```
然后找到如下设置进行修改，然后保存。
```
  "Addresses": {
    "Swarm": [
      "/ip4/0.0.0.0/tcp/4001",
      "/ip6/::/tcp/4001"
    ],
    "Announce": [],
    "NoAnnounce": [],
    "API": "/ip4/127.0.0.1/tcp/5001",
    "Gateway": "/ip4/127.0.0.1/tcp/8080"
  },
```

为确保fstorm节点正常运行，必须保证这三个端口对外开放。

## 创建私有网络，删除公网引导节点bootstrap

### 把自带的bootstrap信息去掉。
``````
ipfs bootstrap rm --all
``````

### 显示节点配置
``````
ipfs config show
``````

### 启动节点

``````
nohup ipfs daemon > fstorm.out 2>&1 &
``````

### 监控节点

``````
tail -f fstorm.out
``````

到此，节点完整完毕。增加节点请看（Coming Soon）。使用 firestorm 存储请看[examples](./examples/README.md)。
