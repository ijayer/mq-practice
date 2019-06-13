Message Queue Practices
========================
`AuthorBy: zhe`     
`CreateAt: 20190328`        
`ModifyAt: 20190401`  

<!-- 摘要 -->

> [Message Queue Compares Ref to Here][#1]
> 
> - [RabbitMQ](https://www.rabbitmq.com/)
> - [ActiveMQ](http://activemq.apache.org/)
> - [ZeroMQ](http://zeromq.org/)
> - [Redis](https://redis.io/)
> - [RocketMQ](https://rocketmq.apache.org/)
> - [Kafaka](https://kafka.apache.org/)

<!--more-->

# [RabbitMQ](./rabbitmq)

> RabbitMQ Practices (using the Go RabbitMQ client)

- [01-hello world](rabbitmq/golang-intro/01-hello-world)
    
    send.go: 实现生产者代码，发送一条 `hello` 然后退出; recv.go: 实现消费者代码，从队列中持续取数据

- [02-work-queues](rabbitmq/golang-intro/02-work-queues)

    实现一个用来在多个 Workers 之间分发 `耗时任务` 的工作队列
    
- [03-publish/subscribe](rabbitmq/golang-intro/03-publish-subscribe)

    通过实现一个日志系统来学习 RabbitMQ ` p/b ` 的简单用法：该日志系统由2个程序构成，第一个用来发送日志消息，第二个则用来接收并将其输出到终端窗口。

- [04-routing](rabbitmq/golang-intro/04-routing)

    给日志系统添加新特性：只订阅一部分消息

- [TODO][05-topics](rabbitmq/golang-intro/05-topics)

- [TODO][06-rpc](rabbitmq/golang-intro/06-rpc)

# [TODO][Kafaka](./kafaka)

> Kafaka Practices

[#1]:https://zhezh09.github.io/post/tech/mq/20190327-message-queue-compare/