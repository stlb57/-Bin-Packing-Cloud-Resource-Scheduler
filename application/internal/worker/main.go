package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error retreiving hostname: %v\n", err)
		return
	}
	fmt.Printf("Booted succesfully. Hostname: %s\n", hostname)
	for {
		fmt.Print("Simulating Video Rendering")
		time.Sleep(time.Second * 10)
	}
}
