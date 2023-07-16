package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Consumer struct {
	conn      *amqp.Connection
	queueName string
	LoggerURL string
}

func NewConsumer(conn *amqp.Connection, url string) (Consumer, error) {
	consumer := Consumer{conn: conn, LoggerURL: url}
	if err := consumer.setup(); err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, t := range topics {
		err = ch.QueueBind(
			queue.Name,
			t,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go c.handlePayload(payload)
		}
	}()

	log.Printf("Waiting for message [Exchange, Topic] [logs_topic, %s]\n", queue.Name)

	<-forever
	return nil

}

func (c *Consumer) handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		if err := c.logEvent(payload); err != nil {
			log.Printf("%e", fmt.Errorf("%w", err))
		}

	case "auth":
		panic("handle authenticate not implemented")

	default:
		if err := c.logEvent(payload); err != nil {
			log.Printf("%e", fmt.Errorf("%w", err))
		}
	}
}

func (c *Consumer) logEvent(payload Payload) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = payload.Name
	entry.Data = payload.Data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := c.LoggerURL
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted || err != nil {
		return err
	}

	return nil
}
