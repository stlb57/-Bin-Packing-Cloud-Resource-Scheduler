package main

import (
	"autoscaler/internal/api"
	"autoscaler/pkg/provisioner"
	"autoscaler/pkg/scheduler"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// ch := make(chan int, 100)
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, os.Interrupt)
	// rand.Seed(time.Now().UnixNano())
	// go func() {

	// 	for {

	// 		rand_time := rand.Intn(4) + 1
	// 		time.Sleep(time.Duration(rand_time) * time.Second)
	// 		rand_weight := rand.Intn(8) + 1
	// 		ch <- rand_weight
	// 	}
	// }()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ. Did you start the Docker container?", err)
	}
	defer conn.Close()

	mq, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", err)
	}
	defer mq.Close()

	q, err := mq.QueueDeclare(
		"render_jobs",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue", err)
	}

	go api.StartServer(mq, q.Name)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	c := 10
mainLoop:
	for {
		select {
		case <-ticker.C:
			// res := scheduler.Ffd(ch, c)
			// provisioner.TriggerTerraform(res)

			var batchWeights []int

			for {
				msg, ok, err := mq.Get(q.Name, true)
				if err != nil || !ok {
					break
				}

				weightVal, _ := strconv.Atoi(string(msg.Body))
				batchWeights = append(batchWeights, weightVal)
			}

			if len(batchWeights) > 0 {
				fmt.Printf("📦 Pulled %d jobs from queue. Calculating load...\n", len(batchWeights))
				res := scheduler.Ffd(batchWeights, c)
				provisioner.TriggerTerraform(res)
			} else {
				fmt.Println("💤 No jobs in queue. Waiting for next tick.")
			}

		case <-sigChan:
			fmt.Println("Graceful sHUTDOWN")
			cmd := exec.Command("terraform", "destroy", "-auto-approve")
			_, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}
			break mainLoop
		}

	}
	fmt.Println("Daemon stopped safely")
}
