## x-apisix 项目介绍
* 实现了**go-micro**服务同步到**apisix**开源网关工具
* [apisix](https://github.com/apache/incubator-apisix,"网关") 是基于openresty开发的api网关，支持etcdv2配置中心与yaml本地文件配置
* [go-micro](https://github.com/micro/go-micro,"微服务") 是一款基于go开发的微服务框架
* [x-apisix](https://github.com/micro-in-cn/x-apisix,"监控服务")是基于etcdv3提供库下保证服务的高可用
* 项目愿景：
** apisix+gin+go-micro

## 依赖软件
* go-micro/v2
* micro/cli/v2
* etcdv3
* http

## go-micro
* 注册中心
** etcdv3
** 其实注册中心，需要后期优化代码
* handler需要增加
** /stats  服务状态
** /health 健康检查

## apisix
* [安装文档](https://github.com/apache/incubator-apisix/blob/master/doc/install-dependencies.md,"安装")
* [使用文档](https://github.com/apache/incubator-apisix/blob/master/doc/README_CN.md,"操作文档")
* 兼容协议
** http
** https
** websocket
* 注意事项
** 使用**x-apisix**，使用apisix-dashboard不要编辑,只可以查看，会冲掉配置
** apisix模板有一些配置选项，如果初始化，最好设置好格式，否则对面会报错，由其是map操作
** 理论上**x-apisix**是支持grpc的，但是由于**apisix**的grpc转http插件并不稳定，所以不建议使用

## demo
* TODO...

## x-apisix
* 设计文档
** [服务设计文档](./server.desgin.md)
** [版本规范](https://semver.org/ "api版本").
** server_name格式 e.g. : aimgo.web.v1.passport
** server_version e.g. : v1.22.3

### 执行命令
```shell
# 安装
make install
# 执行
make monitor
```




