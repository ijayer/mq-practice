Publish / Subscribe
===================
`AuthorBy: zhe`     
`CreateAt: 20190519`        
`ModifyAt: 20190602`  

<!-- 摘要 -->

> Publish/Subscribe: Send messages to many consumers at once.

这一小节，通过实现一个日志系统来学习 RabbitMQ ` p/b ` 的简单用法：该日志系统由2个程序构成，第一个用来发送日志消息，第二个则用来接收并将其输出到终端窗口。在开始之前，先了解下 Exchange ！了解它是什么，怎么路由消息等？

<!--more-->

# Exchange

RabbitMQ 中消息传递的核心思想是：生产者永远不会直接向队列发送任何消息。实际上，通常生产者甚至不知道消息是否会被传递到任何队列。 

相反，生产者只能把消息发送给 `exchange`，`exchange` 则负责了两件事：

- 从 `Producer` 接收消息
- 把接收到的消息发送给 `queues`

而 Exchange 怎么处理他收到的消息，则由 `Exchange Type` 决定。

Exchange 可以用 `X` 表示，模型图如下：

![](https://www.rabbitmq.com/img/tutorials/exchanges.png)

## Exchange Type 

> `Exchange Type` 决定了来到 Exchange 的消息该如何分发到 Queues

RabbitMQ 的 Exchange 类型有：`fanout`、`direct`、`topic`、`headers`

- fanout: 将所有发送到 Exchange 的消息路由到所有与该 Exchange 绑定的队列
- direct：将消息路由到那些 BindingKey 和 RoutingKey 完全匹配的队列中去
- topic: 相对于 direct 的严格匹配来说，topic 进行了扩展：即支持模糊匹配
- headers: headers 类型的 Exchange 不依赖于路由键的匹配规则来路由消息，而是根据发送消息的内容中的 `headers` 属性进行匹配

> 后面将在写一篇关于 Exchange Type 的笔记，辅助实例代码进行理解。这里就先做认知学习吧。

## The Default Exchange

默认的 `Exchange` 在声明时用 Empty String `""` 表示，代码描述如下:

```go
err = ch.Publish(
  "",     // exchange
  q.Name, // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing{
    ContentType: "text/plain",
    Body:        []byte(body),
})
```

- 这段代码我们使用了空的(默认) `exchange`。当然，如果该值被指定为具体的值时，消息将被路由到 `routing_key` 表示的队列中去。

# Temproary Queues

> 能够给要使用的队列起名字对我们来说至关重要，因为我们需要将 worker 指向同一个队列中去。其，当你想要在消费者和生产者之间共享 queue 的时候，给队列一个给定的名字也很重要。

接下来，我们看看关于日志系统需要注意的两点：

- 首先，无论我们什么时候连接到 Rabbit，都需要一个新的且为空的队列。为了做到这个，我们可以用随机名来创建队列，或者让 Rabbit Server 为我们选择队列名也是不错的选择

- 其次，一旦我们断开了消费者的连接，都需要将队列自动删除

> 在 amqp 客户端，当我们提供命名为空的队列时，我们会用随即名创建一个非持久化的队列。例如：`amq.gen-JzTY20BRgKO-HjmUJj0wLg`

# Bindings

![](https://res.cloudinary.com/zher-files/image/upload/v1559316042/blog/images/bindings.png)

Bindings: Exchange 和 Queue 之间的绑定关系就叫做 `bindings-绑定`，代码描述如下：

```go
err = ch.QueueBind {
  q.Name,   // queue name
  "",       // routing key
  "logs",   // exchange
  false,
  nil,
}
```

![](https://res.cloudinary.com/zher-files/image/upload/v1559317325/blog/images/python-three-overall.png)

# See Also

- [Publish/Subscribe](https://www.rabbitmq.com/tutorials/tutorial-three-go.html)

[#3]:https://github.com/ijayer/mq-practice/tree/master/rabbitmq/golang-intro/03-publish-subscribe

# Content

- [01-hello world](../01-hello-world)
- [02-work-queues](../02-work-queues)
- [03-publish/subscribe](../03-publish-subscribe)
- [04-routing](../04-routing)
- [05-topics](../05-topics)
- [06-rpc](../06-rpc)
