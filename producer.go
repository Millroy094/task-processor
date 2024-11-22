import ("encoding/json", "net/http", "log", "github.com/streadway/amqp")

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendTask(ch *amqp.Channel, queue string, task Task) {
	body, err := json.Marshal(task)
	failOnError(err, "Failed to serialize task")

	// Publish the task
	err = ch.Publish(
		"", 
		queue,  
		false,  
		false,  
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish task")

	log.Printf("Task ID %d sent to queue %s\n", task.ID, queue)
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

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		sendTask(ch, q.Name, task)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Task ID %d queued successfully", task.ID)
	})

	log.Println("API server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}