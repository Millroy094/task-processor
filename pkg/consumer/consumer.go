package consumer

import (
	"github.com/millroy094/task-processor/pkg/common"
)

func main() {
	connection, channel, _ := common.RetrieveRabbitMQQueue()

	defer connection.Close()
	defer channel.Close()
}
