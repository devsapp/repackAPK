
> 注：当前项目为 Serverless Devs 应用，由于应用中会存在需要初始化才可运行的变量（例如应用部署地区、函数名等等），所以**不推荐**直接 Clone 本仓库到本地进行部署或直接复制 s.yaml 使用，**强烈推荐**通过 `s init ${模版名称}` 的方法或应用中心进行初始化，详情可参考[部署 & 体验](#部署--体验) 。

# start-repack-apk-v3 帮助文档
<p align="center" class="flex justify-center">
    <a href="https://www.serverless-devs.com" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk-v3&type=packageType">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-repack-apk-v3" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk-v3&type=packageVersion">
  </a>
  <a href="http://www.devsapp.cn/details.html?name=start-repack-apk-v3" class="ml-1">
    <img src="http://editor.devsapp.cn/icon?package=start-repack-apk-v3&type=packageDownload">
  </a>
</p>

<description>

基于 CDN + Custom 运行时实现 apk 实时打渠道包

</description>

<codeUrl>

- [:smiley_cat: 代码](https://github.com/devsapp/repackAPK/tree/V3/src)

</codeUrl>
<preview>



</preview>


## 前期准备

使用该项目，您需要有开通以下服务并拥有对应权限：

<service>



| 服务/业务 |  权限  | 相关文档 |
| --- |  --- | --- |
| 函数计算 |  AliyunFCFullAccess | [帮助文档](https://help.aliyun.com/product/2508973.html) [计费文档](https://help.aliyun.com/document_detail/2512928.html) |
| 硬盘挂载 |  AliyunFCServerlessDevsRolePolicy | [帮助文档](https://help.aliyun.com/zh/nas) [计费文档](https://help.aliyun.com/zh/nas/product-overview/billing) |
| 专有网络 |  AliyunFCServerlessDevsRolePolicy | [帮助文档](https://help.aliyun.com/zh/vpc) [计费文档](https://help.aliyun.com/zh/vpc/product-overview/billing) |

</service>

<remark>



</remark>

<disclaimers>



</disclaimers>

## 部署 & 体验

<appcenter>
   
- :fire: 通过 [Serverless 应用中心](https://fcnext.console.aliyun.com/applications/create?template=start-repack-apk-v3) ，
  [![Deploy with Severless Devs](https://img.alicdn.com/imgextra/i1/O1CN01w5RFbX1v45s8TIXPz_!!6000000006118-55-tps-95-28.svg)](https://fcnext.console.aliyun.com/applications/create?template=start-repack-apk-v3) 该应用。
   
</appcenter>
<deploy>
    
- 通过 [Serverless Devs Cli](https://www.serverless-devs.com/serverless-devs/install) 进行部署：
  - [安装 Serverless Devs Cli 开发者工具](https://www.serverless-devs.com/serverless-devs/install) ，并进行[授权信息配置](https://docs.serverless-devs.com/fc/config) ；
  - 初始化项目：`s init start-repack-apk-v3 -d start-repack-apk-v3`
  - 进入项目，并进行项目部署：`cd start-repack-apk-v3 && s deploy -y`
   
</deploy>

## 案例介绍

<appdetail id="flushContent">

基于本案例， 您可以快捷部署生成一个弹性高可用的 “Serverless 实现实时 apk 渠道分包“ 服务。


随着移动游戏市场的多元化，游戏开发商和发行商需针对不同渠道分发定制化的游戏版本，针对此类需求，该方案适用于大规模、多渠道的游戏分发场景，特别是在面临频繁且多变的下载请求时， 用户能够通过系统实时地获取并下载包含定制渠道号的游戏 APK 包，同时享受断点续传的下载体验，确保即使在连接不稳定的情况下也能顺利完成游戏的获取，该方案已被多家游戏公司采纳,。

![](https://img.alicdn.com/imgextra/i2/O1CN019seP901UxWBt9D8h7_!!6000000002584-2-tps-2120-668.png)


如上图所示，游戏 APK 包需要根据实时请求中的的参数获取指定的渠道号，并将渠道号写入 APK 文件固定位置， 如果每天有大量且不同渠道的下载请求， 能**实时**让用户**断点下载**指定渠道的 apk 游戏包。

</appdetail>

## 使用流程

<usedetail id="flushContent">

### CDN 配置

本应用部署成功后， 您会获取一个 访问域名的 url， 比如为 `https://get-apk-apk-repack-evbilghzjb.cn-hangzhou.fcapp.run`

之后登录 [CDN 控制台](https://cdn.console.aliyun.com/) 完成配置：

#### 1. 添加域名

比如您有一个名为 `functioncompute.com` 的域名, 如下图所示， 我添加了 `apk-cdn.functioncompute.com`, 源站的域名为前面应用部署的访问域名 url(_注意是 host，不用填写前面的 https://_), 比如本示例为 `get-apk-apk-repack-evbilghzjb.cn-hangzhou.fcapp.run`

> 其中前缀 apk-cdn 可以随便， 由您这边自己想最后暴露出去的 url 决定

![](https://img.alicdn.com/imgextra/i2/O1CN01KX6FhL1sjp9I1US8M_!!6000000005803-2-tps-1372-840.png)

#### 2. 域名管理

##### 2.1. 根据控制台引导， 完成域名的 CNAME 解析

![](https://img.alicdn.com/imgextra/i4/O1CN01tmlyC222ln0TTrFt1_!!6000000007161-2-tps-956-1372.png)

![](https://img.alicdn.com/imgextra/i1/O1CN01htbiOc1DZNsDqDC9C_!!6000000000230-2-tps-2348-670.png)

![](https://img.alicdn.com/imgextra/i4/O1CN01vKUcG21RGWBEd8eBT_!!6000000002084-2-tps-2586-244.png)

##### 2.2. 完成管理配置, 主要完成回源配置的域名和开启 Range 回源强制

![](https://img.alicdn.com/imgextra/i4/O1CN01d9cRsx23rZckwYqmF_!!6000000007309-2-tps-2646-716.png)

> 域名应用部署成功后返回的访问域名 url 的 host, 比如本示例为 `get-apk-apk-repack-evbilghzjb.cn-hangzhou.fcapp.run`

![](https://img.alicdn.com/imgextra/i3/O1CN01W8rPnG1R1rVDcK7TN_!!6000000002052-2-tps-2612-854.png)

#### 3. 使用浏览器断点下载指定渠道 apk 包

比如:

`http://apk-cdn.functioncompute.com/foo?src=fc-imm-demo/test-apk/qq.apk&cid=uc`

`http://apk-cdn.functioncompute.com/foo?src=fc-imm-demo/test-apk/qq.apk&cid=xiaomi`

其中

- `apk-cdn.functioncompute.com` 表示 cdn 对外的域名
- `src=fc-imm-demo/test-apk/qq.apk` 表示处理的母包， 其中 fc-imm-demo 为 bucket(和函数在同一个 region), test-apk/qq.apk 为 object
- `cid=xiaomi` 表示渠道为 xiaomi, 这个可以自定义

### Tips

- 用户在自己程序中获取渠道信息， 只需要读取 apk 包中 `assets/dap.properties` 文件中的内容即可

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

### 二次开发

#### 本地调试

1. 将测试证书放置在如下位置

```bash
CertPEM_PATH = "/tmp/cert/test-cert.pem"
PrivateKeyPEM_PATH = "/tmp/cert/test-priv.pem"
```

2. 编译， 生成的二进制可执行文件名字为 repack

3. Run Local

```bash
$ RUN_LOCAL=true OSS_ENDPOINT=http://oss-cn-qingdao.aliyuncs.com SOURCE_OBJECT=test/test_pack.apk CHANNEL_ID=xiaomi ACCESS_KEY_ID=xxx ACCESS_KEY_SECRET=yyy  ./repack
```

> 注意将相关 ENV 设置您自己的值即可

####  打包原理

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

- origin apk 存储在 OSS 中
- repack 过程中不需要将 origin apk 下载到本地磁盘
- 新的 apk 实时分段回传给 CDN

这个方案使用很少的磁盘空间并且非常高效。

**方案原理图**

![](https://img.alicdn.com/imgextra/i4/O1CN01ARFir41xyXwDIpAng_!!6000000006512-2-tps-711-463.png)

</usedetail>

## 注意事项

<matters id="flushContent">
</matters>


<devgroup>


## 开发者社区

您如果有关于错误的反馈或者未来的期待，您可以在 [Serverless Devs repo Issues](https://github.com/serverless-devs/serverless-devs/issues) 中进行反馈和交流。如果您想要加入我们的讨论组或者了解 FC 组件的最新动态，您可以通过以下渠道进行：

<p align="center">  

| <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407298906_20211028074819117230.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407044136_20211028074404326599.png" width="130px" > | <img src="https://serverless-article-picture.oss-cn-hangzhou.aliyuncs.com/1635407252200_20211028074732517533.png" width="130px" > |
| --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| <center>微信公众号：`serverless`</center>                                                                                         | <center>微信小助手：`xiaojiangwh`</center>                                                                                        | <center>钉钉交流群：`33947367`</center>                                                                                           |
</p>
</devgroup>
