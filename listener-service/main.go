package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"listener-service/event"

	"github.com/caarlos0/env/v6"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	AmqpURl   string `env:"AMQP_URL"`
	LoggerURL string `env:"LOGGER_URL"`
}

func main() {
	app := Config{}
	if err := env.Parse(&app); err != nil {
		log.Fatal(err.Error())
	}

	conn, err := connect(&app)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Start listening for and consuming RabbitMQ messages")

	consumer, err := event.NewConsumer(conn, app.LoggerURL)
	if err != nil {
		log.Panic(err)
	}

	if err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"}); err != nil {
		log.Printf("%e", fmt.Errorf("%w", err))
	}
}

func connect(app *Config) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial(app.AmqpURl)
		if err != nil {
			log.Println("RabbitMQ not yet ready..")
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			log.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("Backing of on ", backOff.Seconds())
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
