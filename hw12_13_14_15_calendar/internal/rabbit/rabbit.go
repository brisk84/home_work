package rabbit

import (
	"github.com/streadway/amqp"
)

type Rabbit struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewRabbit(url string) *Rabbit {
	return &Rabbit{url: url}
}

func (r *Rabbit) Connect() error {
	var err error
	r.conn, err = amqp.Dial(r.url)
	if err != nil {
		return err
	}
	r.ch, err = r.conn.Channel()
	if err != nil {
		return err
	}
	err = r.ch.ExchangeDeclare("CalE", "direct", false, false, false, false, nil)
	if err != nil {
		return err
	}
	r.q, err = r.ch.QueueDeclare("CalQ", false, false, false, false, nil)
	if err != nil {
		return err
	}
	err = r.ch.QueueBind("CalQ", "CalQ", "CalE", false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *Rabbit) Stop() {
	r.ch.Close()
	r.conn.Close()
}

func (r *Rabbit) Send(msg []byte) error {
	err := r.ch.Publish("CalE", r.q.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
	return err
}

func (r *Rabbit) Get() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(r.q.Name, "sender", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, err
}
