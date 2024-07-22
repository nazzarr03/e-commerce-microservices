package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/nazzarr03/product-service/models"
	"github.com/streadway/amqp"
)

var (
	RabbitMQConn *amqp.Connection
)

func ConnectRabbitMQ() {
	var err error
	RabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err)
	}
}

func ProductConsumer() {
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"product_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			productID, err := strconv.Atoi(string(d.Body))
			if err != nil {
				log.Printf("Error converting product ID: %s", err)
				continue
			}

			var product models.Product
			Db.Where("id = ?", productID).First(&product)
			if product.ProductID == 0 {
				log.Println("Product not found")
				continue
			}
			log.Println("Product found")

			response := map[string]interface{}{
				"price": product.Price,
				"stock": product.Stock,
			}

			responseBody, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshalling response: %s", err)
				continue
			}

			err = ch.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          responseBody,
				},
			)
			if err != nil {
				log.Printf("Error publishing message: %s", err)
			}
		}
	}()

	log.Println("ProductConsumer: Waiting for messages")
	<-forever
}
