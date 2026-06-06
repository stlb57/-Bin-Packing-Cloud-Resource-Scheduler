package scheduler

import (
	"sort"
)

func Ffd(ch chan int, c int) int {
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
