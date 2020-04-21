## x-apisix 项目介绍
* 实现了**go-micro**服务同步到**apisix**守护工具
* [apisix](https://github.com/apache/incubator-apisix) 是基于openresty开发的api网关，支持etcdv2配置中心与yaml本地文件配置
* [go-micro](https://github.com/micro/go-micro) 是一款基于go开发的微服务框架
* [x-apisix](https://github.com/micro-in-cn/x-apisix)是基于etcdv3提供库下保证服务的高可用
* 项目愿景：*apisix+gin+go-micro* 

## 依赖软件
* go-micro/v2
* micro/cli/v2
* etcdv3
* http

## go-micro
* 注册中心
	* etcdv3 - 其实注册中心，需要后期优化代码
* web handler需要增加
	* /stats(服务状态) service.HandleFunc("/stats", webhandler.StatusHandler())
	* /health(健康检查) service.HandleFunc("/health", webhandler.HealthHandler())

## apisix
* [安装文档](https://github.com/apache/incubator-apisix/blob/master/doc/install-dependencies.md)
* [使用文档](https://github.com/apache/incubator-apisix/blob/master/doc/README_CN.md)
* 兼容协议
	* http
	* https
	* websocket
* 注意事项
	* 使用**x-apisix**同步，[apisix-dashboard](https://github.com/apache/incubator-apisix-dashboard)千万别编辑,只可以查看，同步的配置会被UI冲掉(目前是这样)
	* **apisix**模板配置选项，如果初始化，最好设置好格式，否则apisix会报错，由其是key=>val
	* **x-apisix**理论上讲可以支持grpc的，但是由于**apisix**的grpc转http插件并不稳定，所以不建议使用

## TODO
* 兼容apisix v1.2版本测试
* 增加使用apisix插件使用组合
* 兼容go-micro所有的注册中心
* 丰富demo，让大家可以快速使用apisix网关
* 完善守护程序的单元测试
* 兼容grpc转http协议更可靠性
* 快速进入pre-release状态
* 麻烦使用go-micro/v1用户看看是否可以兼容，如果不兼容，是否需要兼容


## demo
* 本项目长期维护,最
* 招一个开源合作者，与我一起长期维护些项目
* [TODO...](https://github.com/aimpsmile/x-apisix-example)

## x-apisix
* 设计文档
	* [服务设计文档](./server.desgin.md)
	* [版本规范](https://semver.org/ "api版本").
	* server_name格式: aimgo.web.v1.passport
	* server_version: v1.22.3

### 执行命令
```shell
# 安装
make install
# 执行
make monitor
```
