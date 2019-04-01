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
    
- [TODO][03-publish/subscribe](rabbitmq/golang-intro/03-publish-subscribe)

- [TODO][04-routing](rabbitmq/golang-intro/04-routing)

- [TODO][05-topics](rabbitmq/golang-intro/05-topics)

- [TODO][06-rpc](rabbitmq/golang-intro/06-rpc)

# [TODO][Kafaka](./kafaka)

> Kafaka Practices

[#1]:https://zhezh09.github.io/post/tech/mq/20190327-message-queue-compare/