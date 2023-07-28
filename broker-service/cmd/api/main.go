package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	AuthURL   string `env:"AUTH_URL,required"`
	LoggerURL string `env:"LOGGER_URL,required"`
	MailerURL string `env:"MAILER_URL,required"`
	AmqpURl   string `env:"AMQP_URL,required"`
	WebPort   string `env:"WEB_PORT,required"`
	Rabbit    *amqp.Connection
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

	app.Rabbit = conn

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
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
