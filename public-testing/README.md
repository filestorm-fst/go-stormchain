# Stormchain Public Testing

## Download storm application

* Linux: [storm-linux](https://github.com/filestorm-fst/go-stormchain/blob/master/storm-linux)
* Mac Os: [storm-os](https://github.com/filestorm-fst/go-stormchain/blob/master/storm-os)

We have two testnet. Stormchain is a public chain and 风暴联盟链 is a federated chain.

### Stormchain Public Testing
Explorer [http://explorer.filestorm.info](http://explorer.filestorm.info)

Connecting Node: 
enode://c9aff5fa6fd978a5935590b8f98212ab0d9ad45f337f0477da7fa20a9cc5d534a1ec3f2a731bb12ccf42caaddc4f7f01cccba54f61848b84c9282fe2a091cfaa@47.115.27.232:30314

Genesis File: [generated_storm.json](generated_storm.json)

Connect

`````````````````````````````
$ mkedir storm_node 
$ cd storm_node
$ ./storm --datadir data account new
INFO [02-06|14:13:49.711] Maximum peer count                       FST=50 LES=0 total=50
Your new account is locked with a password. Please give a password. Do not forget this password.
Password:
Repeat password:

$ ./storm --datadir data init genesis_storm.json
$ ./storm --datadir data  --networkid 20200215 --syncmode 'full' --port 30314 --bootnodes "enode://c9aff5fa6fd978a5935590b8f98212ab0d9ad45f337f0477da7fa20a9cc5d534a1ec3f2a731bb12ccf42caaddc4f7f01cccba54f61848b84c9282fe2a091cfaa@47.115.27.232:30314"

`````````````````````````````

Open a console on a new window 
``````````````````````````````````
$ cd storm_node
$ ./storm attach data/storm.ipc
``````````````````````````````````


### 风暴联盟链

浏览器：[http://federated.filestorm.info](http://federated.filestorm.info)

连接节点：
enode://117c56e1ea11802bedc1a25145272ed01eabde8e8c741fc42ee268e73324442016a95e17738d18dbbe9e06821a396e877159d7fe1d4df82ae547ca129dba9e8d@47.115.0.166:30411

创世文件: [generated_federated.json](generated_federated.json)

连接

`````````````````````````````
$ mkedir storm_node 
$ cd storm_node
$ ./storm --datadir data account new
INFO [02-06|14:13:49.711] Maximum peer count                       FST=50 LES=0 total=50
Your new account is locked with a password. Please give a password. Do not forget this password.
Password:
Repeat password:

$ ./storm --datadir data init genesis_federated.json
$ ./storm --datadir data  --networkid 20200207 --syncmode 'full' --port 30311 --bootnodes "enode://117c56e1ea11802bedc1a25145272ed01eabde8e8c741fc42ee268e73324442016a95e17738d18dbbe9e06821a396e877159d7fe1d4df82ae547ca129dba9e8d@47.115.0.166:30411“

`````````````````````````````

在新窗口打开控制面板
``````````````````````````````````
$ cd storm_node
$ ./storm attach data/storm.ipc
``````````````````````````````````