package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"sort"
	"time"
)

func ffd(ch chan int, c int) int {
	n := len(ch)
	if n == 0 {
		return 0
	}
	weight := make([]int, 0)
	for range len(ch) {
		var num int
		num = <-ch
		weight = append(weight, num)

	}
	sort.Slice(weight, func(i, j int) bool {
		return weight[i] > weight[j]
	})

	buckets := []int{c}

	for i := 0; i < n; i++ {
		flag := false

		for j := 0; j < len(buckets); j++ {
			if buckets[j] >= weight[i] {
				flag = true
				buckets[j] -= weight[i]
				break
			}
		}

		if !flag {
			buckets = append(buckets, c-weight[i])
		}
	}

	return len(buckets)
}

func triggerTerraform(serverCount int) {
	varFlag := fmt.Sprintf("-var=server_count=%d", serverCount)
	fmt.Println("Calculated servers:", serverCount)
	fmt.Println("Triggering terraform...")

	cmd := exec.Command("terraform", "plan", varFlag)
	output, err := cmd.CombinedOutput()

	fmt.Println(string(output))

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ch := make(chan int, 100)
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
	for range ticker.C {
		res := ffd(ch, c)
		triggerTerraform(res)
	}
}
