package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JobRequest struct {
	Weight int `json:"weight"`
}

func HandleJob(ch chan int) http.HandlerFunc {
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

		ch <- jobReq.Weight

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Job with weight %d accepted!\n", jobReq.Weight)
	}
}

func StartServer(ch chan int) {
	http.HandleFunc("/job", HandleJob(ch))
	log.Println("API Server started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
