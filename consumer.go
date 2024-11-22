import ("encoding/json", "net/http", "log", "github.com/streadway/amqp")

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task-queue",
		false,
		false,
		false,
		false,
		nil, 
	)
	failOnError(err, "Failed to declare a queue")
}