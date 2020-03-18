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

fstorm节点需要使用三个端口4001，5001和8080，如果要对其进行修改，可做如下操作。
```
vi ~/.ipfs/config
```
在Addresses下面找到三个端口进行修改，然后保存。
```
  "Addresses": {
    "Swarm": [
      "/ip4/127.0.0.1/tcp/4001",
      "/ip6/::/tcp/4001"
    ],
    "Announce": [],
    "NoAnnounce": [],
    "API": "/ip4/127.0.0.1/tcp/5001",
    "Gateway": "/ip4/127.0.0.1/tcp/8080"
  },
```

为确保fstorm节点正常运行，必须保证防火墙设置中将这三个端口对外开放。

另外，fstorm节点可设置为对固定IP开放，或者对全网开放。

* 如果采用服务器/客户端的架构来实现应用和对fstorm节点的调用，可以将下面两处设置[server ip]改成服务器地址。
```
    "API": "/ip4/127.0.0.1/tcp/5001",
    "Gateway": "/ip4/127.0.0.1/tcp/8080"
```
```    
    "API": "/ip4/[server ip]/tcp/5001",
    "Gateway": "/ip4/[server ip]/tcp/8080"    
```
* 如果采用 DAPP 远程调用架构，则应将[server ip]改成 0.0.0.0

## 创建私有网络

删除公网引导节点bootstrap，把自带的bootstrap信息去掉。不与公共 IPFS 网络连接。

``````
ipfs bootstrap rm --all
``````
将来有详细介绍如何添加私有网络节点，或者与公共网络连接。

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
