RabbitMQ Practices
==================
`AuthorBy: zhe`     
`CreateAt: 20190328`        
`ModifyAt: 20190401`  

<!-- 摘要 -->

> RabbitMQ 是一个开源的、使用最广的消息队列。
> 
> - Erlang 开发，对高并发、路由、负载均衡、数据持久化有很好的支持。
> - 支持的协议：AMQP，XMPP, SMTP, STOMP
> - 支持集群部署（私有云或公有云）
> - 可使用 HTTP-API、命令行工具和 UI 界面以便于管理和监控

<!--more-->

---

> - [RabbitMQ API(golang)][#1]
> - [RabbitMQ Tutorials][#2]

# RabbitMQ

> RabbitMQ 消息转发器：可用来接收、存储和转发消息（binary blobs of data ‒ messages）

RabbitMQ 中的几个术语：

- 生产者：只负责发送消息的程序

    ![](https://www.rabbitmq.com/img/tutorials/producer.png)

- 队列：一个很大的消息缓存池，大小取决于宿主机的内存和磁盘容量；多个生产者可同时发消息给一个队列，多个消费者也可以同时从一个队列中取消息

    ![](https://www.rabbitmq.com/img/tutorials/queue.png)

- 消费者：只负责接收消息的程序

    ![](https://www.rabbitmq.com/img/tutorials/consumer.png)
    
> 生产者、队列及消费者通常会运行在不同机器上；而且同一个应用程序即可包含生产者也可包含消费者

# Running RabbitMQ

[Ref: Running RabbitMQ With Management Plugin][#3]

> ? Note: 
>
> - RabbitMQ 会使用容器的主机名称(Hostname) 生成数据存储目录，因此我们需要在运行容器服务时指定一个具体的名称(`--hostname my-rabbit`)以便后续查看数据, 如果不指定则会使用一个随机名称。
> - RabbitMQ server addr for connecting: `amqp://guest:guest@host_ip:5672/`

```
# This will start a RabbitMQ container listening on the default port of 5672

$ docker run -d --hostname my-rabbit -p 5672:5672 --name rabbitmq rabbitmq:3
# - -p: 做端口转发，让宿主机外的机器可以通过 ip:port 方式访问 

$ docker logs rabbitmq
...
              Starting broker...
2019-03-28 13:14:38.489 [info] <0.216.0>
 node           : rabbit@my-rabbit
 home dir       : /var/lib/rabbitmq
 config file(s) : /etc/rabbitmq/rabbitmq.conf
 cookie hash    : rWhct608elv6zt7P4yeo0A==
 log(s)         : <stdout>
 database dir   : /var/lib/rabbitmq/mnesia/rabbit@my-rabbit # 数据存放目录 -> Queue Buffer
 
# 开启带管理插件的 RabbitMQ Server

$ docker run -d --hostname my-rabbit --name rabbitmq-m -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

# Using the Go RabbitMQ client

```$xslt
Design 

> 客户端协议实现说明

Most other broker clients publish to queues, but in AMQP, clients publish
Exchanges(交易中心) instead.  AMQP is programmable, meaning that both the producers and
consumers agree on（达成一致） the configuration of the broker, instead requiring an
operator or system configuration that declares the logical topology in the
broker.  The routing between producers and consumer queues is via Bindings.
These bindings form the logical topology（拓扑） of the broker.

> Hello World Example Routing

                       +---------------------------------+
           Publishing  +                                 + Delivery
Publisher ------------>+--Exchange --> Routes --> Queue--+----------> Consumer
                       +                                 +
                       +---------------------------------+

In this library, a message sent from publisher is called a "Publishing" and a
message received to a consumer is called a "Delivery".  The fields of
Publishings and Deliveries are close（接近） but not exact mappings to the underlying
wire format to maintain stronger types.  Many other libraries will combine
message properties with message headers.  In this library, the message well
known properties are strongly typed fields on the Publishings and Deliveries,
whereas the user defined headers are in the Headers field.

> Publishings 和 Deliveries 类型定义相关说明

The method naming closely matches the protocol's method name with positional
parameters mapping(映射) to named protocol message fields.  The motivation here is to
present a comprehensive(全面的) view over all possible interactions with the server.

Generally, methods that map to protocol methods of the "basic" class will be
elided（消隐） in this interface, and "select" methods of various channel mode selectors
will be elided(隐藏) for example Channel.Confirm and Channel.Tx.

> 接口命名尽量和 AMQP 协议保持了一致 

The library is intentionally（有意的） designed to be synchronous, where responses for
each protocol message are required to be received in an RPC manner.  Some
methods have a noWait parameter like Channel.QueueDeclare, and some methods are
asynchronous like Channel.Publish.  The error values should still be checked for
these methods as they will indicate IO failures like when the underlying
connection closes.

> 同步通信；但也有例外，如 Channel.Publish 属于异步，也需要检查 err 是否是 IO 错误（在底层连接关闭时就会抛出） 

Asynchronous Events

> 异步事件机制相关说明

Clients of this library may be interested in receiving some of the protocol
messages other than Deliveries like basic.ack methods while a channel is in
confirm mode（确认模式）.

The Notify* methods with Connection and Channel receivers model（模拟） the pattern of
asynchronous events like closes due to exceptions, or messages that are sent out
of band from an RPC call like basic.ack or basic.flow.

Any asynchronous events, including Deliveries and Publishings must always have
a receiver until the corresponding chans are closed.  Without asynchronous
receivers, the sychronous methods will block.

> 异步事件必须有接收者，等待 chan 被关闭；同步方法会阻塞

Use Case

It's important as a client to an AMQP topology to ensure the state of the
broker matches your expectations.  For both publish and consume use cases,
make sure you declare the queues, exchanges and bindings you expect to exist
prior to calling Channel.Publish or Channel.Consume.

  // Connections start with amqp.Dial() typically from a command line argument
  // or environment variable.
  connection, err := amqp.Dial(os.Getenv("AMQP_URL"))

  // To cleanly shutdown by flushing kernel buffers, make sure to close and
  // wait for the response.
  defer connection.Close()

  // Most operations happen on a channel.  If any error is returned on a
  // channel, the channel will no longer be valid, throw it away and try with
  // a different channel.  If you use many channels, it's useful for the
  // server to
  channel, err := connection.Channel()

  // Declare your topology here, if it doesn't exist, it will be created, if
  // it existed already and is not what you expect, then that's considered an
  // error.

  // Use your connection on this topology with either Publish or Consume, or
  // inspect your queues with QueueInspect.  It's unwise to mix Publish and
  // Consume to let TCP do its job well.
```

# Practices Content

- [01-hello world](./golang-v/01-hello-world)
- [02-work-queues](./golang-v/02-work-queues)
- [TODO][03-publish/subscribe](./golang-v/03-publish-subscribe)
- [TODO][04-routing](./golang-v/04-routing)
- [TODO][05-topics](./golang-v/05-topics)
- [TODO][06-rpc](./golang-v/06-rpc)

[#1]:https://godoc.org/github.com/streadway/amqp
[#2]:https://www.rabbitmq.com/getstarted.html