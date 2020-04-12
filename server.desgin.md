## micro framework (命名：aimgo)
> define framework nameing rules

### 全局变量声明
> 例如：
>
> **aimgo.web.v1.passport** = ${BU}.${stype}.${sver}.${module}
> aimgo 业务线
> passport 模块
> http2/web/srv 服务类型
> v1/v2/v3  服务版本

- ${proto_path} = protobuf 文件路径
- **${host}** = map{${BU}.${stype} => "srv.uqudu.com",...,} 域名
- **${sname}** = ${BU}.${module}.${stype}.${sver}  服务名
    - **${BU}**  业务线 
    - **${module}** 模块
    - **${stype}** 服务类型
    - **${sver}** 服务大版本 v1、v2、v3
    - **${MICRO_SERVER_VERSION}** 版本详细版本v99.99.99 = ${sver}.${month}.${day},最长9位
    - **${endpoint}** 端点
        - **${ip}**  ip使用内网IP{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "100.64.0.0/10", "fd00::/8"}
        - **${port}** 端口使用自动生成的
    - **${args}** 请求参数
- golang(gin路由变量与全局映射)   
    - ${group}  = ${module}/${sver}
    - ${endpoint} = /

### server.path
> ${BU}/srv/${module}/${stype}/handler.${sver}.go
> ${BU}/srv/${module}/${stype}/main.go

### proto.path(服务名称与proto路径一定要匹配)

> 例如：aimgo/passport/web/v1
> ${proto_path}/${BU}/${module}/${stype}/${sver} 

### 服务版本与网关模板匹配规范
> ${sver} 是大版本，一般接口路径体现
> MICRO_SERVER_VERSION 是版本的细节 = ${sver}.xx.xx
> 建议服务${web,http2}每次版本升级: 
- 参考网站：
- 小版本接口只可以增加
- 大版本才可以自由的修改

> 网关模板匹配规则（范围不允许交叉重叠）：
- 版本一般v1.3.2.2.2.2.2 会转成float64位进行比对。比对来讲更简单，但是float64小数位超过15位精度就丢失，请规范的使用版本
-   > >= < <= = 浮点数数学比对工具
-  ~ v1.3  匹配 [v1.3.00 ,v1.33.9999999] 闭区间
- v1.3,v1.5 匹配 [v1.3.00,v1.5.9999999] 闭区间
- * 如果版本没有匹配上，保底的。

### 匹配规则
> ${A} = micro.registry.list 注册中心服务列表
> ${B} = apisix.registry.list 网关的配置中心
> **->**  = 左向右同步符号
> **<-**  = 右向左同步符号
> ${KEY} = 自动配置标记

#### ${ID} = {service,route,proto}.ID 唯一ID生成规则如下：
- {$service}.ID, ${route}.ID  根据网关自动生成
- ${proto}.ID = md5(${sname}.${service})

#### 任务每次重启服务初始化
- ${A}->${B}(拉取注册中心服务列表，全量同步到apisix网关)
- ${A}<-${B}(从网关拉取数据与注册中心)
- ${B}根据Desc匹配${KEY}(表示自动配置),diff(${A},${B}),获取${B}多出的service列表
- 删除${B}.Service.ID对应的路由列表

#### 任务常驻执行逻辑
- ${A}全量同步到->${B}
- watcher(${A}),同步->${B}。操作{create,update,delete}

#### 定时clear功能
- 每隔10m进行处理
- 清理多余的proto列表
- 清理多余的service列表

### api gateway (命名：aimapisix)

#### 优化级的，防止匹配出现错误
- grpc > web > 静态服务(每隔Ns拉取api服务) > 控制台手动配置

#### web模块
> ${host}/${sver}/${module}/${endpoint}?${args}
> http://web.uqudu.com/v1/passport/say/hello?name

#### grpc模块
> ${host}/${sver}/${module}/${endpoint}?${args}
#### websock模块

## monitor tool (命名：aimmonitor)
> search.prefix: [aimgo.\*.web.\*,aimgo.\*.http2.\*]

- services 

  ```json
  {{- $len := sub (len .Nodes) 1 }}
  {
    "desc": "{{ .Desc }}",
    "plugins": {},
    "upstream": {
      "type": "roundrobin",
      "nodes": {
        {{- range $key,$val := .Nodes }}
        "{{ $val.Address }}": 1{{- if lt $key $len }},{{- end }}
        {{- end }}
      }
    }
  }

  ```

- grpc_routes

  ```json
  {{- $methods := split .Endpoint.Metadata.method "," }}
  {{- $lenMethod := sub (len $methods) 1 }}
  {{- $lenHosts := sub (len .Svariable.ModuleVeh.Hosts) 1 }}
  {
    "desc": "{{ .Desc }}",
  "priority": 10,
    "methods": [
      {{- range $k,$v := $methods }}
        "{{ $v }}"{{- if lt $k $lenMethod }},{{- end }}
      {{- end }}
    ],
    "uris": [
      "/{{ .Svariable.Sver }}/{{ .Svariable.Module }}/{{ trimPrefix .Endpoint.Metadata.path "/" }}"
    ],
    "hosts": [
      {{- range $key,$val := .Svariable.ModuleVeh.Hosts }}
          "{{ $val }}"{{- if lt $key $lenHosts }},{{- end }}
      {{- end }}
    ],
    "plugins": {
      "grpc-transcode": {
          "proto_id": "{{ .ProtoID }}",
          "service": "{{ .Svariable.Sname }}.{{ .Endpoint.Metadata.grpc_service }}",
          "method": "{{ .Endpoint.Metadata.grpc_method }}",
          "pb_option":["int64_as_string"]
      },
      "gheader-plugin" : {}
    },
    "service_protocol": "grpc",
    "service_id": "{{ .ServiceID }}"
  }

  
  ```
## 目录规划 
apisix   ---> aimapisix网关
cmd      ---> aimtool(工具)
core     ---> aimgosmile(go语言框架)
micro    ---> aimmicro(micro工具集)

global   --->  aimgo.global 全局常量
config   ---> aimgo.configs业务配置文件目录
proto    ---> aimgo.protos文件目录
srv      ---> aimgo.examples业务模块

  


