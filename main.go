package main

import (
	"encoding/json"
	"fmt"
	amqp "github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
)
const MQURL = "amqp://admin:123456@192.168.0.115:5672"
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// Key
	Key string
	// 连接信息
	Mqurl string
}

func NewRabbitmq(queueName, exchagne, key string) *RabbitMQ {
	mq := &RabbitMQ{
		QueueName: queueName,
		Exchange: exchagne,
		Key: key,
		Mqurl: MQURL,
		}
	var err error
	mq.conn, err = amqp.Dial(mq.Mqurl)
	mq.failOnErr(err, "create conn fail .")

	mq.channel, err = mq.conn.Channel()

	mq.failOnErr(err, "create channel fail .")

	return mq
}
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

func (r *RabbitMQ) Destory() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func (r *RabbitMQ) Publish(message string) {
	_,err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
		)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		},
		)
}

func (r *RabbitMQ) Consume()  {
	_,err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		false,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协和处理消息
	go func() {
		for d := range msgs {
			// 实现我们要实现的逻辑函数
			log.Printf("Received a message: %s", d.Body)
			msgObj := message{}
			json.Unmarshal(d.Body, &msgObj)
			if msgObj.Action == "sync_point" {
				handleSyncPoint(msgObj)
			} else if msgObj.Action == "sync_member" {
				handleSyncMember(msgObj)
			}
			d.Ack(false)
		}
	}()
	log.Printf("[*] Waiting for message, To exit press CTRL+C")
	<-forever
}

type message struct {
	Action string `json:"action"`
	MemberId string `json:"memberId"`
	Point int `json:"point"`
	Tel string `json:"tel"`
	Name string `json:"name"`
	Sex int `json:"sex"`
}

func handleSyncMember(msg message)  {
	log.Print(msg.Action)
}
func handleSyncPoint(msg message)  {
	log.Print(msg.Action)
}
func getAccessToken() string {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	content := string(body)
	log.Println(content)

	return content
}

func main() {
	fmt.Println("gogogo")

	//content := getAccessToken()
	//log.Printf(content)

	mq := NewRabbitmq("q1","ex1","k1")
	mq.Consume()
}