package rabbit

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// docker run -d --name some-rabbit -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management

// func TestSendRabbit(t *testing.T) {
// 	url := "amqp://guest:guest@localhost:5672/"
// 	conn, err := amqp.Dial(url)
// 	if err != nil {
// 		fmt.Println(1, err)
// 		return
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
// 	if err != nil {
// 		fmt.Println(2, err)
// 		return
// 	}
// 	body := "Hello world12"
// 	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
// 		ContentType: "text/plain",
// 		Body:        []byte(body),
// 	})
// 	if err != nil {
// 		fmt.Println(3, err)
// 		return
// 	}
// }

// func TestGetRabbit(t *testing.T) {
// 	url := "amqp://guest:guest@localhost:5672/"
// 	conn, err := amqp.Dial(url)
// 	if err != nil {
// 		fmt.Println(1, err)
// 		return
// 	}
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
// 	if err != nil {
// 		fmt.Println(2, err)
// 		return
// 	}

// 	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
// 	if err != nil {
// 		fmt.Println(3, err)
// 		return
// 	}
// 	// forever := make(chan bool)

// 	// go func() {
// 	for d := range msgs {
// 		log.Printf("Received a message: %s\n", d.Body)
// 	}
// 	// }()

// 	fmt.Println(2)
// 	// log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
// 	// <-forever
// }

func TestRabbit(t *testing.T) {
	url := "amqp://guest:guest@localhost:5672/"
	r := NewRabbit(url)
	err := r.Connect()
	require.NoError(t, err)
	err = r.Send([]byte("Test123"))
	require.NoError(t, err)
	ch, err := r.Get()
	require.NoError(t, err)

	for d := range ch {
		log.Printf("Received a message: %s\n", d.Body)
	}
}
