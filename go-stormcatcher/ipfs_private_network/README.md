# 生成一个IPFS私有网络

FileStorm底层使用IPFS协议做文件存储，在此基础上实现数据加密，区块链确权，文件冗余复制，自动修复，数据去重，存储奖励等功能。FileStorm所有节点将组成一个私有的IPFS网络。下面我们来介绍一下怎样生成一个私有IPFS网络。

本指南使用Ubuntu系统做演示。

## 第一步 安装 IPFS

### 下载

Go的下载和安装可以看这里 [https://golang.org/doc/install](https://golang.org/doc/install)

下载并解压IPFS。
`````
wget https://dist.ipfs.io/go-ipfs/v0.4.23/go-ipfs_v0.4.23_linux-amd64.tar.gz
tar xvfz go-ipfs_v0.4.23_linux-amd64.tar.gz
`````

然后生成一个GOPATH `/usr/local/go/bin

`````
export PATH=$PATH:/usr/local/go/bin
`````


### 安装

将解压的目录`go-ipfs/ipfs`复制到，然后删除安装文件，IPFS就安装好了。

```
mv go-ipfs/ipfs /usr/local/bin/ipfs
rm go-ipfs_v0.4.18_linux-amd64.tar.gz
rm -R ./go-ipfs
```

## 第二步 初始化节点

通过`IPFS_PATH=~/.ipfs`可以在一个服务器上生成多个IPFS节点。
```
IPFS_PATH=~/.ipfs ipfs init
```
再初始化一个
```
IPFS_PATH=~/.ipfs2 ipfs init
```
第二个节点的端口都需要修改。以后加上。


## 第三步 创建私有网络

安装 git

`````
sudo apt-get install git
`````

下载私有网络钥匙生成工具`ipfs-swarm-key-gen`

`````
go get -u github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen
``````

生成私网钥匙

`````
./go/bin/ipfs-swarm-key-gen > ~/.ipfs/swarm.key
``````

将私网钥匙复制到其他节点上。如果是同服务器上，可以用下面的指令。

``````
cp ~/.ipfs/swarm.key ~/.ipfs2/swarm.key
``````

## 第四步 生成私网引导节点bootstrap

把自带的bootstrap信息去掉。

``````
IPFS_PATH=~/.ipfs ipfs bootstrap rm --all
``````

显示节点配置

``````
IPFS_PATH=~/.ipfs ipfs config show
``````

显示节点peer

``````
IPFS_PATH=~/.ipfs ipfs config show | grep "PeerID"
``````

在`IPFS_PATH=~/.ipfs2`上加上新的bootstrap节点

``````
IPFS_PATH=~/.ipfs2 ipfs bootstrap add /ip4/<ip address of bootnode>/tcp/4001/ipfs/<peer identity hash of bootnode>
``````

实际操作

``````
IPFS_PATH=~/.ipfs2 ipfs bootstrap add /ip4/127.0.0.1/tcp/4000/ipfs/QmQBBwHr8fShxZmNG91rfnurdkbW37Qpez6PaNfhhbbS3T
``````

## 第五步 启动网络

强迫节点只能连接私网

``````
export LIBP2P_FORCE_PNET=1
``````

启动节点

``````
IPFS_PATH=~/.ipfs ipfs daemon &
``````

如果看到这个，表示运行成功。
```
Successfully raised file descriptor limit to 2048.
Swarm is limited to private network of peers with the swarm key
Swarm key fingerprint: 5782836a3e5ef7a10a06a1bd7fe30a3e
Swarm listening on /ip4/10.33.31.85/tcp/4000
Swarm listening on /ip4/127.0.0.1/tcp/4000
Swarm listening on /ip6/::1/tcp/4000
Swarm listening on /p2p-circuit/ipfs/QmQBBwHr8fShxZmNG91rfnurdkbW37Qpez6PaNfhhbbS3T
Swarm announcing /ip4/10.33.31.85/tcp/4000
Swarm announcing /ip4/127.0.0.1/tcp/4000
Swarm announcing /ip6/::1/tcp/4000
API server listening on /ip4/127.0.0.1/tcp/5001
Gateway (readonly) server listening on /ip4/127.0.0.1/tcp/8080
Daemon is ready
````
