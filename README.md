# learn-grpc
## Grpc-learning
### 1. mode 介绍grpc四种模式
- 简单模式
    > 简单模式只是使用参数和返回值作为服务器与客户端传递数据的方式，最简单。
- 客户端流模式
    > 从客户端往服务器端发送数据使用的是流，即服务器端的参数为流类型，然而在服务器相应后返还数据给客户端，使用的也是流的send方法。一般在服务器端的代码，需要先recv再send，而客户端与此相反。但是在后面的双向模式中可以使用go的协程协作。
- 服务端流模式
    > 即服务器端返回结果的时候使用的是流模式，即传入的数据是通过参数形式传入的。但是在往客户端发送数据时使用send方法，与客户端返回数据的方式大同小异。
- 双向流模式
    > 客户端如果不适用协程，那么发送必须在接收之前。如果使用协程，发送与接收并没有先后顺序。为了保证协程的同步，可以使用互斥量进行约束。
### 2. healthcheck 健康检查
### 3. authentication 授权与认证
### 4. metadata
    > 1. 可以做链路追踪 2.一些公共的请求参数可以放在metadata
### 5. encryption 加密传输
- ALTS 加密方式
- TLS 加密方式
> https://www.aqniu.com/tools-tech/30315.html
### 6. keepalive
### 7. interceptor 拦截器
### 8. loadbalance 负载均衡
### 9. multiplex 多路复用
### 10. client retry 客户端重试
### 11. errors 错误处理
### 12. deadline 超时机制
### 13. debugging 调试
> 目前grpc提供两种调试方式, 1.日志 2.Channelz
> 日志方式: https://github.com/grpc/grpc-go/blob/master/Documentation/log_levels.md
> To turn on the logs for debugging, run the code with the following environment variable: 
  `GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info`.
>  Channelz方式: https://grpc.io/blog/a-short-introduction-to-channelz/
### 14. name_resolving 域名反解
### 15. profiling 性能检测