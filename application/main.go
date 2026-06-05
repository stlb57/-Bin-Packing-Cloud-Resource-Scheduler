package main

import (
	"fmt"
	"log"
	"os/exec"
	"sort"
)

func ffd(weight []int, c int, n int) int {
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
	var weight []int = []int{9, 8, 2, 2, 2, 2}
	c := 10
	n := len(weight)

	res := ffd(weight, c, n)

	triggerTerraform(res)
}
