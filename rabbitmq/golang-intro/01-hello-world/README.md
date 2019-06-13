Hello World
===========
`AuthorBy: zhe`     
`CreateAt: 20190328`        
`ModifyAt: 20190329`  

- send.go: 实现生产者代码，发送一条 `hello` 然后退出.  
- recv.go: 实现消费者代码，从队列中持续取数据

    ```
    p -> [][][][][][] -> c
            queue
    ```
    
    >  hello world data flow diagram
            
## 发送

> send.go：连接到RabbitMQ后，发送一条消息后，退出。

- 连接到 RabbitMQ 服务器：即：建立 Socket 连接，处理协议转换、版本对接以及一些登陆授权问题 For Us.

    ```go
    conn, err := amqp.Dial(url)
    defer conn.Close()
    ```
    
- 打开通道，然后通过这个 ch 来实现我们生产者业务相关的 API

    ```go
	ch, err := conn.Channel()
	defer ch.Close()
    ```

- 创建队列，用于接收、存储生产者消息

    ```go
	q, err := ch. QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
    ```
    
    > Note: 声明队列的操作幂等，即多次创建只会返回一个名为 'name' 的 queue，不存在时则会创建。
    
- 发布消息

    ```go
	err = ch.Publish(exchange, key string, mandatory, immediate bool, msg Publishing) error
    ```
    
    > Note: 消息的内容是字节数组，因此可以按照业务需求进行编码。

## 接收

> recv.go: 连接到 RabbitMQ 服务器，然后从 queue 中不断地读取消息

- 连接到 RabbitMQ 服务器：即：建立 Socket 连接，处理协议转换、版本对接以及一些登陆授权问题 For Us.

    ```go
    conn, err := amqp.Dial(url)
    defer conn.Close()
    ```
    
- 打开通道，然后通过这个 ch 来实现我们消费者业务相关的 API

    ```go
	ch, err := conn.Channel()
	defer ch.Close()
    ```
                   
- 创建队列，用于存储、转发消息给消费者
  
    ```go
  	q, err := ch. QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table) (Queue, error)
    ```
      
    > Note: 这里创建的队列名称必须和生产者程序 `(send.go)` 中的队列名一致, 否则无法完成 `发送/接收` 绑定。
      
- 将消费者程序中的 queue 注册到 RabbitMQ 服务器中，然后准备接收来自 RabbitMQ 转发的消息；消息发送时异步通信，在 Go 中消息将被发送到 `chan amqp.Delivery` 通道中, 因此我们需要通过 `range` 从这个通道来持续取消息

    ```go
    consumeCh, err := ch.Consume(q.Name,"",true, false, false, false, nil)
    for d = range consumeCh {
        // do something with msg
    }
    ```
    
## 运行

- Run Consumer

    ```bash
    # Open terminal A
    $ go run recv.go
    ```
    
- Run Producer

    ```
    # Open terminal B
    $ go run send.go
    ```    

- Result

    ![](https://res.cloudinary.com/zher-files/image/upload/v1553843154/blog/images/rabbitmq-helloworld.png)

> Note：如果运行了带 `管理插件` 的 RabbitMQ Docker Container, 可以访问 `http://host_ip:15672` 查看队列中的详细信息. 如下图：

![](https://res.cloudinary.com/zher-files/image/upload/v1553842308/blog/images/rabbitmq-manager-1.png)
      

# See Also

- [using the Go RabbitMQ client: hello world](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)

[#1]:https://godoc.org/github.com/streadway/amqp
[#2]:https://www.rabbitmq.com/getstarted.html
[#3]:https://hub.docker.com/_/rabbitmq/?tab=description

# Content

- [01-hello world](../01-hello-world)
- [02-work-queues](../02-work-queues)
- [03-publish/subscribe](../03-publish-subscribe)
- [04-routing](../04-routing)
- [05-topics](../05-topics)
- [06-rpc](../06-rpc)
