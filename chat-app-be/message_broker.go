package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

type MessageBroker struct {
	conn             *amqp.Connection
	subscribedUsers  map[string]bool
	subscriptionChan chan string
}

func (mb *MessageBroker) initCon() error {
	if mb.conn != nil {
		return nil
	}

	conn, err := amqp.Dial("amqp://guest:guest@message_broker:5672/")
	if err != nil {
		return err
	}
	mb.conn = conn

	return nil
}

func (mb *MessageBroker) close() error {
	return mb.conn.Close()
}

var MB *MessageBroker

func newMessageBroker() *MessageBroker {
	mb := MessageBroker{}
	mb.subscribedUsers = map[string]bool{}
	mb.subscriptionChan = make(chan string)
	// mb.init()

	return &mb
}

func initMessageBroker() {
	MB = newMessageBroker()
	// err := MB.initCon()
	go messageBrokerHandler()
}

func messageBrokerHandler() {
	for userId := range MB.subscriptionChan {
		if MB.subscribedUsers[userId] {
			continue
		}
		MB.subscribedUsers[userId] = true
		err := MB.initCon()
		if err != nil {
			fmt.Println(err)
			return
		}

		go MB.listenToUserTopic(userId)
	}
}

func (mb *MessageBroker) listenToUserTopic(userId string) error {
	ch, err := mb.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		userExchange(userId), // name
		"fanout",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,               // queue name
		"",                   // routing key
		userExchange(userId), // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	deliveredMessages, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	for message := range deliveredMessages {
		WS_HUB.sendMessageChan <- WsMessageData{
			userId:  userId,
			message: message.Body,
		}
	}

	return nil
}

func (mb *MessageBroker) sendMessages(message []byte, userIds []string) {
	for _, memberId := range userIds {
		err := mb.sendMessageToTopic(message, memberId)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (mb *MessageBroker) sendMessageToTopic(message []byte, userId string) error {
	ch, err := mb.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		userExchange(userId), // name
		"fanout",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		userExchange(userId), // exchange
		"",                   // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}

	return nil
}

func userExchange(userId string) string {
	return "users/" + userId
}
