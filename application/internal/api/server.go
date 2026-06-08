package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type JobRequest struct {
	Weight int `json:"weight"`
}

type messagePublisher interface {
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg any) error
}

// func HandleJob(ch chan int) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPost {
// 			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		var jobReq JobRequest
// 		if err := json.NewDecoder(r.Body).Decode(&jobReq); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}

// 		ch <- jobReq.Weight

// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprintf(w, "Job with weight %d accepted!\n", jobReq.Weight)
// 	}
// }

func HandleJob(mq *amqp.Channel, queueName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var jobReq JobRequest
		if err := json.NewDecoder(r.Body).Decode(&jobReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		body := strconv.Itoa(jobReq.Weight)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := mq.PublishWithContext(ctx,
			"",
			queueName,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		if err != nil {
			http.Error(w, "Failed to publish job", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Job with weight %d safely stored in RabbitMQ!\n", jobReq.Weight)
	}
}

// func StartServer(ch chan int) {
// 	http.HandleFunc("/job", HandleJob(ch))
// 	log.Println("API Server started on http://localhost:8080")
// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatal(err)
// 	}
// }

func StartServer(mq *amqp.Channel, queueName string) {
	http.HandleFunc("/job", HandleJob(mq, queueName))
	log.Println("API Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
