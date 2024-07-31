package gotil

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	_qp *amqp.Connection
	_ch *amqp.Channel
}

func NewRabbitMQ(user, pass, host, port string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%v:%v@%v:%v/", user, pass, host, port))
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{_qp: conn}, nil
}

func (r *RabbitMQ) createChannel(name string) error {
	if r._qp == nil {
		return errors.New("Disconnected")
	}

	if r._ch == nil {
		ch, err := r._qp.Channel()
		if err != nil {
			return err
		}
		r._ch = ch
	}

	_, err := r._ch.QueueDeclare(
		name,  // Queue name
		false, // Durable (messages survive broker restarts)
		false, // Delete when unused
		false, // Exclusive (for this connection only)
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) Publish(channel, data string) error {
	if r._qp == nil {
		return errors.New("Disconnected")
	}
	if r._ch == nil {
		r.createChannel(channel)
	}
	return r._ch.Publish(
		"",      // Exchange
		channel, // Routing key
		false,   // Mandatory
		false,   // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
}

func (r *RabbitMQ) Consume(channel string, action func(data string)) error {
	if r._qp == nil {
		return errors.New("Disconnected")
	}
	if r._ch == nil {
		r.createChannel(channel)
	}

	msgs, err := r._ch.Consume(
		channel, // Queue name
		"",      // Consumer
		true,    // Auto-acknowledgement
		false,   // Exclusive
		false,   // No-local
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		action(string(d.Body))
	}
	return nil
}
