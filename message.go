package main

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	ErrMsgConn         = errors.New("Message: failed to connect to RabbitMQ")
	ErrMsgOpenChan     = errors.New("Message: failed to open channel")
	ErrMsgDeclareQueue = errors.New("Message: failed to declare queue")
	ErrMsgSetQos       = errors.New("Message: failed to set QOS")
	ErrMsgRegConsumer  = errors.New("Message: failed to register consumer")
	ErrMsgPublish      = errors.New("Message: failed to publish message")
	ErrMsgParse        = errors.New("Message: failed to parse message")
)

type Message interface {
	Server() error
}

type hawkbitMessage struct{}

func NewHawkbitMessage() Message {
	return &hawkbitMessage{}
}

func (m *hawkbitMessage) Server() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return ErrMsgConn
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return ErrMsgOpenChan
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return ErrMsgDeclareQueue
	}

	// FAQ: How to Optimize the RabbitMQ Prefetch Count, LOVISA JOHANSSON, CloudAMQP team
	// https://www.cloudamqp.com/blog/how-to-optimize-the-rabbitmq-prefetch-count.html
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return ErrMsgSetQos
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    //args
	)
	if err != nil {
		return ErrMsgRegConsumer
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		for d := range msgs {
			body, _ := parseBody(d.Body)
			// if err != nil {
			// 	return ErrMsgParse
			// }
			err = ch.PublishWithContext(ctx,
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // madatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          body,
				})
			if err != nil {
				return ErrMsgPublish
			}

			d.Ack(false)
		}
	}
}

func parseBody(req []byte) ([]byte, error) {
	var n Deployment
	if err := json.Unmarshal(req, &n); err != nil {
		return nil, err
	}

	d, err := dp.GetDeployment(n.BID)
	if err != nil {
		return nil, ErrDeploymentNotFound
	}
	if n.Distribution != (distribution{}) {
		d.Distribution = n.Distribution
		err = dp.PutDeploymentDistribution(n.BID, d)
		if err != nil {
			return nil, nil
		}
	}

	response, err := json.Marshal(d.Status)
	if err != nil {
		return nil, nil
	}

	return response, nil
}
