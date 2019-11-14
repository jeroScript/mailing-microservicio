package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func startRabbit() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Fallo conexion con rabbit")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mailing/specific", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			procesarMensaje(&d, ch)
		}
	}()

	log.Printf(" [*] Esperando mensajes. Para salir presione CTRL+C")
	<-forever
}

func procesarMensaje(d *amqp.Delivery, ch *amqp.Channel) {
	var respbody string
	var body callBody
	err := json.Unmarshal(d.Body, &body)
	if err != nil {
		log.Print(err)
		respbody = "Error en payload: " + err.Error() + ". CorrelationID: " + d.CorrelationId
	} else {
		// err = body.sendSpecific(d.Body["template_name"].(string))
		err = body.sendSpecific()
		if err != nil {
			log.Print(err)
			respbody = "Error: no se pudo enviar el/los email. " + err.Error() + ". CorrelationID: " + d.CorrelationId
		} else {
			log.Print("Se envió el/los email correctamente")
			respbody = "Éxito: se realizo la accion send solicitada. CorrelationID: " + d.CorrelationId
		}
	}
	d.Ack(false)

	respq, err := ch.QueueDeclare(
		d.ReplyTo, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",         // exchange
		respq.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          []byte(respbody),
			CorrelationId: d.CorrelationId,
		})
	failOnError(err, "Failed publishing response to queue: "+d.ReplyTo)
}
