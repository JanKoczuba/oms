package main

import (
	"context"
	"encoding/json"
	pb "github.com/JanKoczuba/commons/api"
	"github.com/JanKoczuba/commons/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)

				d.Nack(false, false)

				continue
			}

			log.Printf("Payment link created %s", paymentLink)
			d.Ack(false)
		}
	}()

	<-forever
}
