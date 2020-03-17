# FileStorm 节点的使用

### FileStorm 节点 Javascript 使用指南

``````
npm install js-fstorm-http-client
``````
使用
``````
const FStormHttpClient = require('js-fstorm-http-client');
``````

实例化一个新的 FStormHttpClient 对象
``````
const fStormHttpClient = new FStormHttpClient({
    IPFSProvider:  'http://localhost:5001',
    chainProvider:  'http://localhost:7545',
    chainAccount: {
      privateKey: '0x532ae81c5aab3ce0f6d04d107f3a5b48d76b1c6338e5f94105d1081309bbdc08'
    },
});
``````
options 必选
* options.IPFSProvider - any 必选 fstorm 存储节点服务端配置，支持 ipfs-http-client 的设置方式
* options.chainProvider - any 必选 storm 同步节点服务端配置，支持 storm3 的设置方式
* chainAccount - object 必选 发送到链交易的地址信息
* chainAccount.privateKey - string 可选 地址私钥（如果提供 privateKey 则无需提供 keystore 和 密码）
* chainAccount.keystore - object | string 可选 keystore（如果未提供 privateKey，则必须提供 keystore 和 密码）
* chainAccount.password - string 可选 密码

``````
 const fStorm = new FStormHttpClient({
    IPFSProvider:  'http://localhost:5001',
    chainProvider:  'http://localhost:8645',
    chainAccount: {
      privateKey: ''
    },
  });
``````
或者
``````
 const fStorm = new FStormHttpClient({
   IPFSProvider: config.IPFSProvider,
   chainProvider: config.chainProvider,
   chainAccount: {
     keystore: {},
     password: ''
   }
  });
``````  

加文件
``````  
const fs = require('fs');

(async function f() {
let files = await fStormHttpClient.add(
    fs.readFileSync('/test.jpg').toString())  
  for await (let res of files) {
    console.log(res);
  }
}())
``````  
