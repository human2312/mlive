package util

// @Time : 2020年3月3日16:17:24
// @Author : Lemyhello
// @Desc: rabbitmq

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"log"
	"mlive/library/graceful"
	"mlive/library/snowflakeId"
	"mlive/library/wg"
	"mlive/service/admin"
	"net"
	"runtime/debug"
	"strconv"
	"time"
)

var (
	mqPool Pool
	isRun bool
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewMQ() {
	//factory 创建连接的方法
	factory := func() (interface{}, error) { return amqp.Dial(viper.GetString("rabbitmq.mqurl")) }
	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }
	//创建一个连接池： 初始化2，最大连接5，空闲连接数是4
	poolConfig := &Config{
		InitialCap: 2,
		MaxIdle:    100,
		MaxCap:     300,
		Factory:    factory,
		Close:      close,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 60 * 60 * 24 * 30 * time.Second,
	}
	mqPool, _ = NewChannelPool(poolConfig)
	//从连接池中取得一个连接
	//v, err := p.Get()
	//do something
	//conn :=v.(*amqp.Connection)
	//将连接放回连接池中
	//p.Put(v)
	//释放连接池中的所有连接
	//p.Release()
	//查看当前连接中的数量
	current := mqPool.Len()
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"len=", current)
	return
}

//创建RabbitMQ结构体实例
func NewRabbitMQ(p Pool) *RabbitMQ {
	rabbitmq := &RabbitMQ{}
	var err error
	//从连接池中取得一个连接
	v, err := p.Get()
	//将连接放回连接池中
	defer p.Put(v)
	//do something
	rabbitmq.conn =v.(*amqp.Connection)
	//rabbitmq.conn, err = amqp.Dial(viper.GetString("rabbitmq.mqurl"))
	rabbitmq.failOnErr(err, "创建连接错误!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败!")
	return rabbitmq
}

//断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}

//生产
func (r *RabbitMQ) Publish(queueName string,message string) (err error){
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"rabbitMq连接池数量",mqPool.Len())
	r = NewRabbitMQ(mqPool)
	defer r.channel.Close()
	exchange := viper.GetString("rabbitmq.exchange")
	kind := viper.GetString("rabbitmq.kind")
	err = r.channel.ExchangeDeclare(exchange, kind, true, false, false,false, nil)
	//保证队列存在,消息队列能发送到队列中
	_, err = r.channel.QueueDeclare(
		queueName,
		//是否持久化
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"QueueDeclare:", err)
		return
	}
	err = r.channel.QueueBind(queueName, queueName, exchange, true, nil)
	//2.发送消息到队列中
	err = r.channel.Publish(
		exchange,
		queueName,
		//如果为true,根据exchange类型和routekey规则,如果无法找到符合条件的队列那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息队列到队列后发现队列上没有绑定消费者,则会把消息发还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	messageId := snowflakeId.GetSnowflakeId()
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"Publish:", err)
		go logDb(strconv.FormatInt(messageId,10),queueName,"push",message,0)
		return
	}
	//log to mysql
	go logDb(strconv.FormatInt(messageId,10),queueName,"Publish",message,1)
	return nil
}

//消费
func (r *RabbitMQ) Consume(queueName string,f func([] byte) (bool)) (err error) {
	r = NewRabbitMQ(mqPool)
	defer r.channel.Close()
	exchange := viper.GetString("rabbitmq.exchange")
	kind := viper.GetString("rabbitmq.kind")
	err = r.channel.ExchangeDeclare(exchange, kind, true, false, false,false, nil)
	//1.申请队列,如果队列不存在会自动创建,如果存在则跳过创建
	//保证队列存在,消息队列能发送到队列中
	_, err = r.channel.QueueDeclare(
		queueName,
		//是否持久化
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"QueueDeclare:", err)
		return
	}
	err = r.channel.QueueBind(queueName, queueName, exchange, true, nil)
	// 获取消费通道
	r.channel.Qos(1, 0, true) // 确保rabbitmq会一个一个发消息
	//2.接受消息
	msgs, err := r.channel.Consume(
		queueName,
		//用来区分多个消费者
		queueName,
		//是否自动应答
		false,
		//是否具有排他性
		false,
		//如果设置为true,表示不能将同一个connection中发送消息传递给这个connection中的消费者
		false,
		//队列消费是否阻塞
		false,
		nil)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"Consume:", err)
		return
	}

	//平滑重启
	isRun = false
	go func() {
		grace := graceful.GetChan()
		for g := range grace {
			if g == true {
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"处理结束消费")
				r.channel.Cancel(queueName,false)
				for {
					if isRun == false {
						fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"事务处理结束")
						wg.Done()
					} else {
						fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"事务还没结束")
					}
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()

	forever := make(chan bool)
	res := false
	//3.启用协程处理消息
	go func() {
		defer func() {
			if err := recover(); err != nil {
				SendDingDing(queueName,"队列",err,"【堆栈信息】", string(debug.Stack()))
				fmt.Println(queueName,"队列",err,"【堆栈信息】", string(debug.Stack()))
				panic(recover())
			}
		}()

		//实现我们要处理的逻辑函数
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),"queue name" ,queueName ,"get messages,To exit press CTRL+C\n")
		for d := range msgs {
			isRun = true
			//log to mysql
			messageId := snowflakeId.GetSnowflakeId()
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"接收到消息messageid:",messageId,"内容:",string(d.Body))
			go logDb(strconv.FormatInt(messageId,10),queueName,"Consume",string(d.Body),0)
			res = f(d.Body)
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"),messageId,"消息处理结果:", res)
			if res != true {
				// 业务消费处理失败
				err = d.Ack(true)
				if err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"接收到",messageId,"信息确认失败",string(d.Body))
				}
			} else {
				// 业务消费处理成功
				err = d.Ack(false)
				//d.Nack(true,true)
				if err != nil {
					fmt.Println(time.Now().Format("2006-01-02 15:04:05"),queueName,"接收到",messageId,"信息确认失败",string(d.Body))
				}
			}
			go func() {
				if res == true {
					logRes(strconv.FormatInt(messageId,10),1)
				} else {
					logRes(strconv.FormatInt(messageId,10),0)
				}
			}()
			isRun = false
		}
	}()
	<-forever
	return
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s,%s", message, err))
	}
}

//log 日志记录
func logDb(messageId string,queue string , types string ,message string,result int)  {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err,"【堆栈信息】", string(debug.Stack()))
		}
	}()
	admin.MqAdd(messageId,queue,types,message,result)
}

//logRes 结果更新
func logRes(messageId string,result int)  {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err,"【堆栈信息】", string(debug.Stack()))
		}
	}()
	admin.MqUpdate(messageId,result)
}
