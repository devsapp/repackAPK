# start-repack-apk 帮助文档

<p align="center" class="flex justify-center">
    <a href="https://www.serverless-devs.com" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk&type=packageType">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-repack-apk" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk&type=packageVersion">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-repack-apk" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk&type=packageDownload">
  </a>
</p>

<description>
</description>

<table>

## 前期准备
使用该项目，推荐您拥有以下的产品权限 / 策略：

| 服务/业务 | 函数计算 |  硬盘挂载 |  VPC |  其它(安全组) |     
| --- |  --- |   --- |   --- |   --- |   
| 权限/策略 | AliyunFCFullAccess |  AliyunNASFullAccess |  AliyunVPCFullAccess |  AliyunECSFullAccess |  

</table>

<codepre id="codepre">

# 代码 & 预览

- [ :smiley_cat:  源代码](https://github.com/devsapp/start-repack-apk/tree/main/src)

</codepre>

<deploy>

## 部署 & 体验

<appcenter>

-  :fire:  通过 [Serverless 应用中心](https://fcnext.console.aliyun.com/applications/create?template=start-repack-apk) ，
[![Deploy with Severless Devs](https://img.alicdn.com/imgextra/i1/O1CN01w5RFbX1v45s8TIXPz_!!6000000006118-55-tps-95-28.svg)](https://fcnext.console.aliyun.com/applications/create?template=start-repack-apk)  该应用。 

</appcenter>

- 通过 [Serverless Devs Cli](https://www.serverless-devs.com/serverless-devs/install) 进行部署：
    - [安装 Serverless Devs Cli 开发者工具](https://www.serverless-devs.com/serverless-devs/install) ，并进行[授权信息配置](https://www.serverless-devs.com/fc/config) ；
    - 初始化项目：`s init start-repack-apk -d start-repack-apk`   
    - 进入项目，并进行项目部署：`cd start-repack-apk && s deploy -y`

</deploy>

<appdetail id="flushContent">

# 应用详情
**Serverless 实现实时 apk 渠道分包**

游戏分发平台的游戏 APK 包需要根据实时请求中的的参数获取指定的渠道号，并将渠道号写入 APK 文件固定位置， 如果每天有大量且不同渠道的下载请求， 能**实时**让用户**断点下载**指定渠道的 apk 游戏包

应用原理图如下：

![](https://img.alicdn.com/imgextra/i2/O1CN019seP901UxWBt9D8h7_!!6000000002584-2-tps-2120-668.png)

本应用主要部署后端的函数，部署成功后， 您会获取一个 url， 比如为 `http://get-apk.apk-repack.1986114430573743.cn-hangzhou.fc.devsapp.net`

之后登录 [CDN 控制台](https://cdn.console.aliyun.com/) 完成配置：

#### 添加域名
![](https://img.alicdn.com/imgextra/i3/O1CN01aDJ7ed1U2j50vyStH_!!6000000002460-2-tps-1291-942.png)

#### 域名管理
1. 根据控制台引导， 完成域名的 CNAME 解析
![](https://img.alicdn.com/imgextra/i2/O1CN01iZLk411kJFDJZ9Pn8_!!6000000004662-2-tps-1535-229.png)

2. 完成管理配置, 主要完成回源配置的域名和开启 Range 回源强制
![](https://img.alicdn.com/imgextra/i2/O1CN01HDjf111JqE79AaKgQ_!!6000000001079-2-tps-1024-327.png)

![](https://img.alicdn.com/imgextra/i1/O1CN01JKdUXF1VhJndHKVPs_!!6000000002684-2-tps-934-417.png)

#### 使用浏览器断点下载指定渠道 apk 包

比如:

`http://xiliu.functioncompute.com/foo?src=fc-imm-demo/test-apk/qq.apk&cid=uc`

`http://xiliu.functioncompute.com/foo?src=fc-imm-demo/test-apk/qq.apk&cid=xiaomi`

其中 
- `xiliu.functioncompute.com` 表示 cdn 对外的域名
- `src=fc-imm-demo/test-apk/qq.apk` 表示处理的母包， 其中  fc-imm-demo 为 bucket(和函数在同一个region), test-apk/qq.apk 为 object
- `cid=xiaomi` 表示渠道为 xiaomi, 这个可以自定义

**Tips**

- 用户在自己程序中获取渠道信息， 只需要读取 apk 包中 `assets/dap.properties` 文件中的内容即可
- 应用生成的 url， 比如上篇幅中介绍的 `http://get-apk.apk-repack.1986114430573743.cn-hangzhou.fc.devsapp.net` 为临时测试域名， 您可以切换成您自己的生产域名
    > [配置自定义域名](https://help.aliyun.com/document_detail/90763.htm)

- 换用自己的证书， 只需要换掉 target/cert 下面的文件即可：
    > jarsigner 将 .keystore 文件作为 RSA 密钥的来源，要将其转换为 golang 可识别的 .pem，我们需要以下几行：
    ```bash
    # key store
    $ keytool -genkey -keystore test.keystore  -alias test -keyalg RSA -validity 10000
    
    # convert to pkcs12 format
    $ keytool -importkeystore -srckeystore test.keystore -destkeystore test.p12 -deststoretype PKCS12
    
    # private key pem
    $ openssl pkcs12 -in test.p12 -nocerts -nodes -out tmp-test-priv.pem
    $ openssl rsa -in tmp-test-priv.pem -out test-priv.pem
    
    # cert pem
    $ openssl pkcs12 -in test.p12 -nokeys -out test-cert.pem
    ```

# 打包原理
对于一个原始的 apk 文件，将一个新文件添加到存档中，然后对 apk 重新签名获取新的 apk 文件。等价于以下命令相同的效果：

```bash
# adds a file to origin.apk and results in new.apk

$ unzip origin.apk -d origin/
$ echo "1234" > /tmp/cpid
$ cp /tmp/cpid origin/
$ rm -rf origin/META-INF
$ cd origin
$ jar -cf new-unsigned.apk *
$ jarsigner -keystore test.keystore -signedjar new.apk new-unsigned.apk 'test'
```

但是这个应用的方案一些区别：
          
- origin apk存储在OSS中
- repack过程中不需要将origin apk下载到本地磁盘
- 新的apk实时分段回传给 CDN

这个方案使用很少的磁盘空间并且非常高效。

**方案原理图**

![](https://img.alicdn.com/imgextra/i4/O1CN01ARFir41xyXwDIpAng_!!6000000006512-2-tps-711-463.png)


</appdetail>

<devgroup>

## 参考

1. The [zip format][zip-format] allows appending to zip files without rewrite the entire file
2. The [great zipmerge][zip-merge] makes appending to zip files as easy as a charm
3. The great design in [great zipmerge][zip-merge] makes using OSS as the storage backend possible
4. The great [OSS][oss] features like multipart/uploadPartCopy/getObjectByRange makes OSS as a perfect storage backend

[zip-format]: https://en.wikipedia.org/wiki/Zip_(file_format)
[zip-merge]: https://github.com/rsc/zipmerge
[oss]: https://www.aliyun.com/product/oss

## 开发者社区

您如果有关于错误的反馈或者未来的期待，您可以在 [Serverless Devs repo Issues](https://github.com/serverless-devs/serverless-devs/issues) 中进行反馈和交流。如果您想要加入我们的讨论组或者了解 FC 组件的最新动态，您可以通过以下渠道进行：

<p align="center">

| <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407298906_20211028074819117230.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407044136_20211028074404326599.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407252200_20211028074732517533.png" width="130px" > |
|--- | --- | --- |
| <center>微信公众号：`serverless`</center> | <center>微信小助手：`xiaojiangwh`</center> | <center>钉钉交流群：`33947367`</center> | 

</p>

</devgroup>