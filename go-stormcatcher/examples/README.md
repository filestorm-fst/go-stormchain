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
(async function f() {
  const files = [{
    path: '/tmp/myfile2.txt',
    // content: 'ABC'
  }];

  for await (const result of fStormHttpClient.add(files)) {
    console.log(result);
  }
}());
``````  

读文件
``````
const FStormHttpClient = require('../src/index');
const config = require('./config');
const path = require('path');
const fs = require('fs');
const FileType = require('file-type');

(async function f() {
  const fStormHttpClient = new FStormHttpClient({
    IPFSProvider: config.IPFSProvider,
    chainProvider: config.chainProvider,
    chainAccount: {
      privateKey: '4099a9cdefc9c23466bb5552e90f07292bc8c367957d389feacb737b8530dad9'
    },
  });
  // 001 文本(在线输出)
  /*
  let files = await fStormHttpClient.get('QmfDmsHTywy6L9Ne5RXsj5YumDedfBLMvCvmaxjBoe6w4d');
  for await (let res of files) {
    console.log((await res.content.next()).value.toString());
  }
  */

  let files = fStormHttpClient.get('QmNWM9kwgnsskedBMNm7U2T9MXg2Cd4KkwRhMJNuD5UFXt');

  for await (const file of files) {
    let bufferArray = [];
    for await (const chunk of file.content) {
      bufferArray.push(chunk._bufs[0]);
    }

    // 合并 buffer
    let buffer = Buffer.concat(bufferArray);
    console.log(buffer);
    // 获取
    let fileType = await FileType.fromBuffer(buffer);
    // 后缀
    let ext = fileType.ext;
    // mime
    let mimeType = fileType.mime;
    console.log(mimeType);
    // 保存文件
    fs.writeFileSync(path.resolve(__dirname, './save-file.' + ext), buffer);

  }


}());
``````
