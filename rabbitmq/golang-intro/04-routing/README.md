---
title: "Rabbitmq | 04 - Routing"
date: 2019-06-02T17:05:35+08:00
lastmod: 2019-06-02 17:05:35
draft: false
keywords: [mq,rabbitmq,go]
description: ""
tags: [mq,rabbitmq,golang]
categories: [tech]
author: "jayer"
---

<!-- 摘要 -->

这一小节继续前面的内容，给日志系统添加新特性：只订阅一部分消息。

<!--more-->

# Bindings

上一小节我们创建了 Exchange 和 Queue 的一个绑定，代码描述如下：

```go
err = ch.QueueBind(
    q.Name,
    "",
    "logs", // exchange
    false,
    nil,
)
```

> Binding: 描述的就是 Exchange 和 Queue 之间的关系，也可以看作：队列对 Exchange 传来的消息感兴趣。

Bindings 也可以接收一个参数：`routing_key`；这里，为了避免和 `Channel.Publish` 的参数冲突，将其称作 `binding_key`, 代码描述如下：

```go
err = ch.QueueBind(
    q.Name,
    "black", // routing key
    "logs",
    false,
    nil,
)
```

> Note: `binding_key` 的含义取决于 exchange 的类型（回顾一下，有：fanout、direct、topic、header）

# Direct exchange

前面，我们的日志系统会将所有的消息广播给所有的消费者，这里我们将其进行扩展：即依据日志的严重程度来允许过滤消息。例如，我们只将严重的错误写入磁盘，而告警和展示的日志信息则不用写入磁盘，以此来节省存储空间。

上一小节使用了 `fanout` exchange, 其扩展性不够，只能进行无意识的广播，即将发送给 exchange 的消息路由到与该 exchange 绑定的所有 queues 中去

这里，我们使用 `direct` exchange, direct 背后的算法逻辑：队列的 binding_key 和 消息的 routing_key 完全匹配时，消息才会被路由到队列中去。

![](https://res.cloudinary.com/zher-files/image/upload/v1559485554/blog/images/direct-exchange.png)

上图中，`direct` Exchange ‘X’ 绑定了两个队列，第一个队列的 `binding_key` 是 `orange`, 第二个队列有两个 `binding_key`, 分别是 `black`、`green`

按照这样设计，消息路由规则如下：

- 发送到 Exchange 的消息，如果其 binding_key 为 orange，则会被路由到 Q1
- 发送到 Exchange 的消息，如果其 binding_key 为 black 或者 green 时，则会被路由到 Q2
- 其他的 binding_key 的消息则会被直接丢弃掉

# Multiple bindings

![](https://res.cloudinary.com/zher-files/image/upload/v1559486171/blog/images/direct-exchange-multiple.png)

同样的 `binding_key` 绑定在多个队列上完全可行，上图实列中：发送到 Exchange X 的消息，如果 `binding_key` 为 black, 则都会传送到 Q1 和 Q2 队列中去。

# See Also

> Thanks to the authors 🙂

* [Routing](https://www.rabbitmq.com/tutorials/tutorial-four-go.html)

# Content

[#1]:https://www.rabbitmq.com/tutorials/tutorial-four-go.html
- [01-hello world](../01-hello-world)
- [02-work-queues](../02-work-queues)
- [03-publish/subscribe](../03-publish-subscribe)
- [04-routing](../04-routing)
- [05-topics](../05-topics)
- [06-rpc](../06-rpc)
