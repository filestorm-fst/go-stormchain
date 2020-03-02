# go-stormchain
go-stormchain is a blockchain technology that implements a multichain ecosystem by supporting multiple blockchains including public, private and federated chains. It also supports multiple consensuses.

## 公有链技术

由全球代码爱好者组成，共同开发和管理的，支持多种共识的多链架构区块链技术。信息公开，技术开源，可以为应用提供专属区块链，支持高TPS。自带区块链存储功能，使用方便。本技术还在持续开发中。

## 联盟链技术

由中国工程师自主研发基于PBFT共识的区块链技术，支持国家监督和管控，多节点备份，数据存证不可篡改。以国密局标准算法加密原始数据，提取文件Hash上链传输，保护原文件及用户信息安全。

## 用户指南

下载运行程序 storm （选择正确的操作系统）

* Linux: [storm-linux](https://github.com/filestorm-fst/go-stormchain/blob/master/storm-linux)
* Mac Os: [storm-os](https://github.com/filestorm-fst/go-stormchain/blob/master/storm-os)

和创世文件 
* [storm_chain.json](https://github.com/filestorm-fst/go-stormchain/blob/master/storm_chain.json)

将 storm 和 storm_chain.json 复制到需要安装的节点服务器上。在这个指南中，我们将使用三台装有 Ubuntu 操作系统的服务器。（Storm区块链至少需要两个节点才能运行。）在服务器上，可以建立新文件夹 storm_node 作为节点运行根目录。将两个文件存到这个目录下。(如果要在同一个服务器上跑多个节点，可以生成多个文件夹如 storm_node1, storm_node2，然后执行下面的步骤。）

### 第一步 建立账号

在区块链上，账号是用户身份的代表。每个节点程序都必须通过一个账号来运行。这个账号文件必须存在节点服务器上，密码由指定管理员保管。所以第一步，我们在三台服务器上各建一个账号。

``````````````````````````````````
$ cd storm_node
$ ./storm --datadir data account new
INFO [02-06|14:13:49.711] Maximum peer count                       FST=50 LES=0 total=50
Your new account is locked with a password. Please give a password. Do not forget this password.
Password:
Repeat password:

Your new key was generated

Public address of the key:   0x084149366635cD8E727d1eD50c363107b1b2a565
```````````````````````````````````
`0x084149366635cD8E727d1eD50c363107b1b2a565`是账号的公钥，可以公开。

使用 `--datadir data` 的目的是把账号建立在节点运行根目录下的 data 文件夹中，便于查找和管理。运行完后，将会看到 data  中有一个新的名为 keystore 的文件夹。钥匙文件就在这个文件夹内。

我们可以将密码写到一个密码文件中，方便节点启动。这个文件在节点启动后需要手动删除。

``````````````````````````````````
$ echo 'password' > password.txt
``````````````````````````````````

在其他两个节点上也做同样的操作，并将密码保管好。三个公钥收集起来，等下需要使用。


### 第二步 修改创世文件

利用stormchain技术，FileStorm与MOAC共同打造了一条应用链。可以通过在MOAC主网上发布一个智能合约生成创世文件。具体介绍可以看[这里](https://github.com/filestorm-fst/go-stormchain/tree/master/solidity)。

开发者可以对上面的链接获取的创世文件做如下修改。

```````````````````````
{
  "config": {
    "chainId": 20090103,
    "pbft": {
      "period": 7,
      "epoch": 36000
    }
  },
  "nonce": "0x0",
  "timestamp": "0x5de22b51",
  "extraData": "0x00000000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111222222222222222222222222222222222222222233333333333333333333333333333333333333330000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "gasLimit": "0x47b760",
  "difficulty": "0x1",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "coinbase": "0x0000000000000000000000000000000000000000",
  "alloc":
  {
    "1111111111111111111111111111111111111111": {
      "balance": "50000000000000000000000000"
    }
  },
  "number": "0x0",
  "gasUsed": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
}
```````````````````````````````

这是缺省的创世文件。其中需要解释和修改的配置为：

* chainId： Storm区块链是一个多链生态，每个区块链都有自己独立的ID，20090103是主链的ID，（第一个区块链Bitcoin的诞生日为01/03/2009），联盟链可以采用任何整数进行测试。MOAC应用链通过发合约注册到基础链可得到认证的chainId。

* pbft：这是风暴联盟链使用的共识，未来Storm区块链将会支持更多的共识机制。

* period：出块时间，缺省设定是7秒钟。可以根据实际情况修改。

* extraData：创始区块验证人账号。Storm区块链启动的时候，需要指定创始区块验证人账号，最少两个，最多无限个。这里给出的例子是三个验证人，分别用1...111, 2...222, 3...333来表示。我们可以用前面生成的三个账号去替代 extraData 值里的 1...111, 2...222, 3...333。在这里记得要把账号前的 0x 去掉。如果有更多创始区块验证人，可以把账号加到3...333后面。也可以减少，但是最少要有两个账号。

* 所以 extraData 的格式是 0x[64个0][账号1][账号2][账号3]...[账号n][130个0]

* 注意：这只是创始区块验证人的账号设置。区块链跑起来以后，我们还可以通过让验证人投票的方式增加或减少验证人。

* alloc：如果是有通证的区块链，可以在这里给一些特殊账号或者管理员账号一些通证。

* 可以用真实账号（如前面为三个节点生成的账号）替代`1111111111111111111111111111111111111111`。

* `"balance": "50000000000000000000000000"` 是为账号设定原生通证的数量。此处设置数量为5千万，包括小数点后18个0。

* 如果需要给多个账号通证，可以在这里加
``````````````````````````
  "alloc":
  {
    "1111111111111111111111111111111111111111": {
      "balance": "50000000000000000000000000"
    },
    "2222222222222222222222222222222222222222": {
      "balance": "20000000000000000000000000"
    }
  }
``````````````````````````

注意：要在前一个}后面加个逗号。

修改完以后将文件保存好。然后复制到其他两个节点上。一条区块链上的所有节点的创世文件必须一模一样。所以，以后需要建新节点，就到这三个节点这里拿创世文件。

### 第三步 初始化节点

在三台节点服务器上运行下面的指令。节点初始化就完成了。区块链的数据库内容就被存到了 data 中。

``````````````````````````````````
$ ./storm --datadir data init storm_chain.json
``````````````````````````````````

节点初始化成功你会看到这样一条信息
``````````````````````````````````
INFO [02-08|12:50:17.448] Successfully wrote genesis state         database=lightchaindata hash=33b98a…e3f240
``````````````````````````````````

### 第四步 启动节点

``````````````````````````````````
$ ./storm --datadir data --networkid 20090103 --gasprice '0' --syncmode 'full' --port 30317 --unlock '0x084149366635cD8E727d1eD50c363107b1b2a565' --password password.txt --mine
``````````````````````````````````

* --datadir data 是指定区块数据保存在当前目录下的 data 目录中。

* --networkid 20090103 这个networkid必须和创世文件中的chainId相同。

* --gasprice '0' 这是设置交易费。无通证的联盟链要把交易费用设为0。若是有通证区块链可以不填，用创世区块中的gasLimit。

* --syncmode 'full' 是指同步方式为全同步。创始节点必须用全同步方式启动。后面加入的节点可以不写，或者采用其他同步发方式。

* --port 30317 这是 Storm 区块链的通讯端口。缺省值为30303。这是一个可选项，如果在一台服务器上跑多个节点，就必须修改这个值。

* --unlock '0x084149366635cD8E727d1eD50c363107b1b2a565' 这里的数字记得要改成前面为这台服务器生成的管理员账号。解锁后，这个节点就可以出块。

* --password password.txt 使用存在密码文件中的密码。

* --mine 启动挖矿（验证出块）

运行成功后，打开一个新窗口，将密码文件删除。

``````````````````````````````````
$ rm password.txt
``````````````````````````````````

这时候，我们会看到所有节点都停在第二个块
``````````````````````````````````
INFO [02-06|13:01:02.505] Commit new mining work                   block=1 sealhash=56a809…a06463 uncles=0 txs=0 gas=0 fees=0 elapsed=158.258µs
INFO [02-06|13:01:02.634] Sealed a new block                       block=1 sealhash=56a809…a06463 hash=17d48d…cb74e4 elapsed=128.207ms
INFO [02-06|13:01:02.634] Mined a potential block                  block=1 hash=17d48d…cb74e4
INFO [02-06|13:01:02.635] Commit new mining work                   block=2 sealhash=72548f…ba20b0 uncles=0 txs=0 gas=0 fees=0 elapsed=809.66µs
INFO [02-06|13:01:02.635] 
``````````````````````````````````
这是因为三个节点还没有互相连上，所以，下一步，就是要做节点的连接。


### 第五步 节点网络连接

选择一个节点做为连接主节点，在服务器上新开一个窗口，做如下操作，就可以进入节点控制面板。

``````````````````````````````````
$ cd storm_node
$ ./storm attach data/storm.ipc
``````````````````````````````````

然后在面板中做如下操作，就可以得到节点的网络ID
``````````````````````````````````
> admin.nodeInfo.enode
"enode://0e3a9317c4c9e8910e5d34f627ab798ac0814ac732cb433db4d73382a39c1cb2e8fd35fc53ad406f5883f86a45bc9fb083503ed357ed1a07b5c0a078815c534f@[::]:30317"
``````````````````````````````````

在另一个节点上打开节点控制面板
``````````````````````````````````
$ cd storm_node
$ ./storm attach data/storm.ipc
``````````````````````````````````
然后做如下操作，就可以将两个节点相连
``````````````````````````````````
> admin.addPeer("enode://0e3a9317c4c9e8910e5d34f627ab798ac0814ac732cb433db4d73382a39c1cb2e8fd35fc53ad406f5883f86a45bc9fb083503ed357ed1a07b5c0a078815c534f@[::]:30317")
``````````````````````````````````
记得在[::]填入正确的节点IP。如果是内网必须使用内网IP。

在第三个节点上也做同样操作，就可以将三个节点连起来。然后三个节点就开始轮流出块。


### 第六步 加新节点（同步节点）

启动节点
``````````````````````````````````
$ ./storm --datadir data  --networkid 20090103 --syncmode 'full' --port 30329
``````````````````````````````````
networkid和port可能根据实际情况修改。

在另一个窗口上打开节点控制面板
``````````````````````````````````
$ cd storm_node
$ ./storm attach data/storm.ipc
``````````````````````````````````
然后连接主节点，从而与整个网络相连
``````````````````````````````````
> admin.addPeer("enode://0e3a9317c4c9e8910e5d34f627ab798ac0814ac732cb433db4d73382a39c1cb2e8fd35fc53ad406f5883f86a45bc9fb083503ed357ed1a07b5c0a078815c534f@[::]:30317")
``````````````````````````````````
这时候，新节点还不能验证出块，只能接受区块信息，也就是同步节点。

如果需要在stormchain上做开发，建议使用同步节点。如何在stormchain上做开发，可以查看[wiki](https://github.com/filestorm-fst/go-stormchain/wiki/Development)。

如果同步节点想要成为出块节点，就必须遵守所用共识机制。本指南采用的是pbft共识，增加新节点，必须通过验证节点投票，得到超过50%的赞同票就可以成为验证节点。下面，我们就来看看验证节点的投票功能。


### 第七步 验证节点投票

联盟链上的每一验证节点都可以代表一个独立的利益体。增加和减少验证节点都需要超过半数的验证节点投票决定。验证节点有一组投票功能，可以在控制面板中输入 pbft 来查看

``````````````````````````````````
> pbft
{
  votes: {},
  blockStatus: function(),
  dropVote: function(),
  getBlockStatusByHash: function(),
  getBlockStatusByNumber: function(),
  getValidatorsByHash: function(),
  getValidatorsByNumber: function(),
  vote: function()
}
``````````````````````````````````

  vote: 给一个新节点投赞成或者反对票。

  dropVote: 撤销给一个节点的投票。

  votes: 显示节点所有的投票。

  blockStatus: 显示区块状态，包括出块节点信息。

  getBlockStatusByHash: 通过区块哈希查找区块状态。

  getBlockStatusByNumber: 通过区块数查找区块状态。

  getValidatorsByHash: 通过区块哈希查找验证人信息。

  getValidatorsByNumber: 通过区块数查找验证人信息。

给一个节点投赞成票的指令是
``````````````````````````````````
> pbft.vote('0x53e5c08cb895599e7cfa5da58a783a56e9f140db', true)
``````````````````````````````````
投反对票是
``````````````````````````````````
> pbft.vote('0x53e5c08cb895599e7cfa5da58a783a56e9f140db', false)
``````````````````````````````````
投完票后可以撤销
``````````````````````````````````
> pbft.dropVote('0x53e5c08cb895599e7cfa5da58a783a56e9f140db')
``````````````````````````````````
如果在投票生效前撤销，相当于弃权，且不能再投票。如果投票已生效，撤销将没有意义。

查看在某个区块，如第300个块时的验证人列表，可以使用
``````````````````````````````````
> pbft.getValidatorsByNumber(300)
``````````````````````````````````
如果联盟成员达成共识，决定同时创建一条新链，可以通过第一到第六步生成节点创建新链。但是一旦链建成后，就必须用投票的方式增加新的验证节点。确保区块链的安全性。


### 后台执行节点

前面介绍的方法是通过服务器终端界面操作启动节点。终端界面关闭或者掉线，节点程序就会停止。所以，我们要让节点程序在服务器后台操作。命令如下

``````````````````````````````````
$ nohup ./storm --datadir data  --networkid 20090103 --syncmode 'full' --port 30317 --unlock '0x084149366635cD8E727d1eD50c363107b1b2a565' --password password.txt --mine > storm.out 2>&1 &
``````````````````````````````````

这样，我们就把程序在后台跑起来，并且把日志输入到一个叫 storm.out 的文件中去了。可以通过下面的指令来看日志
``````````````````````````````````
$ tail -100f storm.out
``````````````````````````````````

如果要中止程序，可以通过下面指令
``````````````````````````````````
$ ps -ef | grep storm
``````````````````````````````````
找到正在运行的 storm 进程，记住第二个数字（在这里是78115）
``````````````````````````````````
501 78115 78036   0  1:37PM ttys014    0:00.06 ./storm --datadir data --networkid 20090109 --gasprice 0 --syncmode full --port 30343 --unlock 0xfe8c8e5Cf019f0DAF4d16f14eC88d5971b16cdD2 --password password.txt --mine
``````````````````````````````````
然后执行
``````````````````````````````````
kill -9 78115
``````````````````````````````````
这样节点就停止了。


### 节点升级

节点升级，只需要中止节点程序，替换节点程序，然后重新启动即可。如果同时管理几个节点，建议一次升级一个节点。等节点出块后再更新其他节点。节点重新启动可能需要最长10分钟来连接回区块链网络，如果想快速连回，可使用第六步中 admin.addPeer 的方式快速建立连接。

### 节点日志解释

````````
INFO [03-01|14:09:57.005] Commit new mining work                   block=8944 sealhash=d0f9c1…0999be uncles=0 txs=0 gas=0 fees=0 elapsed=799.162µs
INFO [03-01|14:10:02.002] Imported new blocks                      block=8944 hash=f193f4…8d0170 blks=1 txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:10:02.002] Commit new mining work                   block=8945 sealhash=05508e…d47330 uncles=0 txs=0 gas=0 fees=0 elapsed=232.068µs
INFO [03-01|14:10:07.002] Imported new blocks                      block=8945 hash=ca9823…981c56 blks=1 txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:10:07.002] Commit new mining work                   block=8946 sealhash=490317…94c810 uncles=0 txs=0 gas=0 fees=0 elapsed=235.869µs
INFO [03-01|14:10:12.005] Sealed a new block                       block=8946 sealhash=490317…94c810 hash=daff82…69ecd7 elapsed=5.002s
INFO [03-01|14:10:12.005] Mined a potential block                  block=8946 hash=daff82…69ecd7
INFO [03-01|14:10:12.006] Commit new mining work                   block=8947 sealhash=8a7231…20c7bd uncles=0 txs=0 gas=0 fees=0 elapsed=888.573µs
INFO [03-01|14:10:12.006] Waiting for mining work. 

INFO [03-01|14:10:02.002] Reached canonical chain                  block=8937 hash=9c2b14…75ef51

INFO [03-01|15:10:02.012] Throw away a block                       block=8947 hash=8a7231…20c7bd

````````
这是一个出块节点上常见的日志，我们来解释一下。

* Commit new mining work - 把本地交易进行打包生成一个本地区块。
* Imported new blocks - 把通过共识选择的别的节点生成的新区块导入。
* Sealed a new block - 被共识选中成为下一个出块节点，把本地区块封装。
* Mined a potential block - 把自己封装的本地区块导入。
* Waiting for mining work - 出块节点要等待一段时间才能继续出块。
* Reached canonical chain - 本节点前些时候出的一个区块终极确认不可再改。（一般等6个区块）
* Throw away a block - 本节点前些时候出的一个区块无效被别的节点替代（常见于网络不畅时收到出块节点太慢）

其他参数
* block 区块高度
* hash 区块哈希
* sealhash 打包哈希值。PBFT每组投票都需要一个ID，sealhash就可以理解成投票ID。
* blks 区块数量。（导入区块的时候，如果网络不通畅，有时候会导入多个ID。)
* txs 区块里的交易数量
* gas 燃料费用
* fee 总费用
* elapsed 生成区块的时间
* mgas 燃料费用(单位 gwei)
* dirty 是否超时
* uncles 叔块。PBFT共识没有叔块

同步节点上常常会看到这样的日志
````````
INFO [03-01|14:42:34.539] Imported new blocks                      block=194366 hash=9f4c6d…e6d34e blks=1  txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:42:58.498] Imported new blocks                      block=194367 hash=c6be54…7580be blks=1  txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:43:06.749] Synchronisation slowing down             peer=99c0270b228226d2 msg="retrieved hash chain is invalid"
INFO [03-01|14:43:09.538] Imported new blocks                      block=194371 hash=aaae9b…2d1dbc blks=1  txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:43:27.821] Reimported chain segment                 block=104371 elapsed=109.200ms blks=7  hash=421faa…b31354
INFO [03-01|14:43:33.864] Imported new blocks                      block=194372 hash=1fbe72…7420de blks=1  txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:43:41.231] Synchronisation slowing down             peer=99c0270b228226d2 msg="retrieved hash chain is invalid"
INFO [03-01|14:43:44.539] Imported new blocks                      block=194376 hash=eec99b…381eb0 blks=1  txs=0 mgas=0.000 dirty=0.00B
INFO [03-01|14:44:17.783] Imported new blocks                      block=194377 hash=fe3c77…f280fd blks=1  txs=0 mgas=0.000 dirty=0.00B
````````

* Imported new blocks 正常导入区块
* Synchronisation slowing down 同步节点网络不畅就会出现这个信息。
* Reimported chain segment 如果前面导入区块有漏掉的，就会从漏掉的区块一直往前一直拿到当前块。
