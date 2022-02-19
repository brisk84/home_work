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
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}
	// defer conn.Close()

	r.ch, err = conn.Channel()
	if err != nil {
		r.conn.Close()
		return err
	}
	// defer ch.Close()

	r.q, err = r.ch.QueueDeclare("hello", false, false, false, false, nil)
	if err != nil {
		r.ch.Close()
		r.conn.Close()
		return err
	}
	// r.ch.ExchangeDeclare()
	return nil
}

func (r *Rabbit) Send(msg []byte) error {
	err := r.ch.Publish("", r.q.Name, false, false, amqp.Publishing{
		// ContentType: "text/plain",
		ContentType: "application/json",
		Body:        msg,
	})
	return err
}

func (r *Rabbit) Get() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume("hello", "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return msgs, err
	// forever := make(chan bool)

	// go func() {
	// for d := range msgs {
	// 	log.Printf("Received a message: %s\n", d.Body)
	// }
	// }()

	// fmt.Println(2)
	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// <-forever
	// return []byte{}, nil
}
