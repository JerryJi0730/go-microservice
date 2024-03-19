package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	//try to connect to rabbitmq
	rabbitCoon, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitCoon.Close()

	//start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")
	//create consumer
	consumer, err := event.NewConsumer(rabbitCoon)
	if err != nil {
		panic(err)
	}
	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.Warning", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//don't continue until rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backingoff....")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
