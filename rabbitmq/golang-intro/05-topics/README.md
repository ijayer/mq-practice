---
title: "Rabbitmq | 05 - Topics"
date: 2019-06-13T22:48:56+08:00
lastmod: 2019-06-22 22:48:56
draft: false
keywords: [mq,rabbitmq,go]
description: ""
tags: [mq,rabbitmq,golang]
categories: [tech]
author: "jayer"
---

<!-- 摘要 -->

前面，分别使用了 `fanout` 和 `direct` 类型实现了简易的日志生产、路由和消费，虽然 `direct` 可以按照 `Binding_Key` 绑定关系实现日志过滤，但其仍有局限性，即不能基于多个标准进行消息路由。例如，我不仅想按照日志的严重程度来收取日志，也想把不同日志源的日志信息也收集过来，就好比 `syslog`，既基于日志的严重程度(error/warn/info) 来路由消息，也基于不同模块(auth/cron/kern...) 来路由。

如果我们只想监听来自模块 `cron` 的严重消息，而忽略掉所有来自 `kern` 模块的日志，应该怎么做呢？

> 接下来，就要学习另一个 Exchange Type：`Topic` 
> 
<!--more-->

## Topic Exchange

在 Topic 模式下对 `routing_key` 是有要求的：

- 首先，消息发送给 Topic Exchange 时，其 routing_key 不能为任意的字符，必须是由 `.` 隔开的单词列表
- 其次，组成 routing_key 的单词可以任意，但通常使用与所传输消息相关的词语
- 最后，routing_key 的长度限制是 255 字节，只要不超出阈值，怎么定义都行

再来看 `binding_key`: binding_key 命名规范和 routing_key 一致

在 RabbitMQ 的数据流转结构图中，Topic 的 Exchange 之后的逻辑实际上和 `direct` 模式是类似的：一个设定有 routing_key 的消息，将会和所有绑定到 Exchange 的队列做一次匹配，匹配项就是 binding_key，匹配成功的消息将会发送到该队列中去。

那么消息过来后是如何做匹配的呢？ 这就要来了解下关于 `binding_key` 的两项特别重要要求：

- *（星号）：用来替代一个单词
- #（井号）：用来替代零个或多个单词

我们通过例子来理解一下：如下图：

![](https://www.rabbitmq.com/img/tutorials/python-five.png)

在例子中，所有发送的消息都是来描述动物的。所发送消息的路由键(routing_key)由三个点(.)分单词组成，第一个单词描述：速度，第二个单词描述：颜色，第三个单词描述：种类。如下：`<speed>.<colour>.<species>`

接下来，我们创建绑定：

- Q1 和 binding_key `*.orange.*` 绑定
- Q2 和 binding_key `*.*.rabbit` 和 `lazy.#` 绑定

上述两条绑定规则可以描述如下：

- 所发送的消息如果是描述橙色动物的，那将会被 Q1 接收
- Q2呢，则会接收所有与 rabbit 和 lazy 的动物有关的消息

来，看看被设置为以下 `routing_key` 的消息将如何路由

- `routing_key = quick.orange.rabbit` 的消息将会被传递到 Q1 & Q2
- `routing_key = lazy.orange.elephant` 的消息也会被路由到 Q1 & Q2
- `routing_key = quick.orange.fox` 的消息将会被路由到 Q1
- `routing_key = lazy.brown.fox` 的消息将会被路由到 Q2
- `routing_key = lazy.pink.rabbit` 的消息只会被路由到 Q2 一次（尽管它和两个队列匹配）
- `routing_key = quick.brown.fox` 的消息将不会被传递到任何队列中去

虽然有了命名规范，但就是有人不按规范来做啊！比如，我就将 `routing_key` 设置为了 `orange` | `quick.orange.male.rabbit`，如果这个消息被发送出去，将会发生什么呢？毫无疑问，这些消息没法匹配到任务队列，且消息将被丢失掉。

另一个例子，`routing_key = lazy.orange.male.rabbit` 尽管由四个单词组成，也会匹配到最后一个绑定，并将消息路由到 Q2

现在，关于 `Topic Exchange` 我们应该知道：

- Topic Exchange 能力很强，可以实现像其他类型的 Exchange 一样的特性
- 当一个队列被绑定到带有 `#` 的 `binding_key` 时，这个队列将会接收到所有的消息，而不用管 `routing_key` 是如何设置的。就像在 `faout` 模式的 Exchange 一样
- 当 `binding_key` 没有用到 `*` 和 `#` 时，Topic Exchange 就起和 `direct` 类型的 Exchange 一样的作用

来，接着上一节的内容开始写代码吧：这里将日志消息的 `routing_key` 设置为由两个单词组成：`<faclity>.<severity>`

完整代码：[Golang Intro: 05 - Topics][#3]

# See Also

> Thanks to the authors 🙂

* [Topics][#1]
  
[#1]:https://www.rabbitmq.com/tutorials/tutorial-five-go.html
[#3]:https://github.com/ijayer/mq-practice/tree/master/rabbitmq/golang-intro/05-topics

# Content

- [01-hello world](../01-hello-world)
- [02-work-queues](../02-work-queues)
- [03-publish/subscribe](../03-publish-subscribe)
- [04-routing](../04-routing)
- [05-topics](../05-topics)
- [06-rpc](../06-rpc)
