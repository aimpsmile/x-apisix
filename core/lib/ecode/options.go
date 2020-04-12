package ecode

//	OK 通用正确的请求
var OK = New(2002, "ok")

//	SysErr系统错误信息
var SysErr = New(5000, "当前页面维护中，请访问其它页面服务！")

//	系统致命错误
var SysPanicErr = New(5001, "当前页面维护中，请访问其它页面服务！")

// 	hystrix-暂停，处理请求
var SysHystrixBreakErr = New(5011, "当前页面拥堵中，请喝杯茶后进行尝试")

// 	hystrix-超时，处理请求
var SysHystrixTimeoutErr = New(5012, "当前页面拥堵中，请喝杯茶后进行尝试")

// 	hystrix-最大连接，处理请求
var SysHystrixMaxConnectErr = New(5013, "当前页面拥堵中，请喝杯茶后进行尝试")

//	ratelimit限流-需要用户手动重试
var SysRatelimitErr = New(5014, "前方发生塞车,请喝杯茶后进行尝试")

//	json错误处理
var JSONErr = New(5045, "系统运输数据错误,请换个姿势刷新进行尝试")

//	micro srv 错误-并没有返回任何数据
var MicroErr = New(5046, "系统运输数据错误,请换个姿势刷新进行尝试")

//	处理未知错误-专门包go-micro框架返回的错误信息
var UnknownErr = New(5047, "系统未知错误,请换个姿势刷新进行尝试")

//	srv未知错误
var SRVUnknownErr = New(5048, "系统服务错误,请换个姿势刷新进行尝试")

//	api未知道错误
var APIUnknownErr = New(5049, "系统接口错误,请换个姿势刷新进行尝试")

//  http2未知道错误
var HTTP2UnknownErr = New(50450, "系统接口错误,请换个姿势刷新进行尝试")

//	web未知错误
var WEBUnknownErr = New(5051, "系统网页错误,请换个姿势刷新进行尝试")

//  web未知页面
var WEBNotFound = New(5052, "网页不存在，请尝试正确页面")

//  post请求参数未找到
var POSTNotFound = New(5053, "POST参数未找到")

//  get请求参数未找到
var GETNotFound = New(5053, "GET参数未找到")

//  url请求参数未找到
var URLNotFound = New(5053, "请求路径未找到")
