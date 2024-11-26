package consumer

import (
	"github.com/millroy094/task-processor/pkg/common"
)

func processor(id int,messages <-chan amp.Delivery) {
	for message := range messages {
		var task common.Task
		err:= json.Unmarshal(message.Body, &task)

		if err != nil {
			log.Printf("Worker %d: Failed to deserialize task: %v\n", id, err)
			continue
		}

		fmt.Printf("Worker %d processing task ID %d of type %s\n", id, task.ID, task.Type)

		message.Ack(false)
	}
}

func main() {

	connection, channel, queue := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()

	messages, err: channel.Consume(queue.Name, "", false, false, false, false, nil)

	common.FailOnError(err, "Failed to register a consumer")

	for i := 1; i <= 3; i++ {
		go worker(i, messages)
	}
}
