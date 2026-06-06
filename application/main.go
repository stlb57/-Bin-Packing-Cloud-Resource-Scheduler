package main

import (
	"autoscaler/pkg/provisioner"
	"autoscaler/pkg/scheduler"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	ch := make(chan int, 100)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	rand.Seed(time.Now().UnixNano())
	go func() {

		for {

			rand_time := rand.Intn(4) + 1
			time.Sleep(time.Duration(rand_time) * time.Second)
			rand_weight := rand.Intn(8) + 1
			ch <- rand_weight
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	c := 10
mainLoop:
	for {
		select {
		case <-ticker.C:
			res := scheduler.Ffd(ch, c)
			provisioner.TriggerTerraform(res)

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
