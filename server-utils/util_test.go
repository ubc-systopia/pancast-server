package server_utils

import (
	"log"
	"math"
	"testing"
)

func TestSampleLaplacianDistribution(t *testing.T) {
	granularity := 1000
	samples := 50000
	count := int64(500)

	mean := 0
	sensitivity := float64(1344)
	epsilon := 0.1
	delta := 0.02

	distribution := make(map[int]int)

	for i := 0; i < samples; i++ {
		junkCount := SampleLaplacianDistribution(int64(mean), sensitivity, epsilon, delta)
		index := int(math.Floor(float64(junkCount + count) / float64(granularity)) * float64(granularity))
		_, found := distribution[index]
		if !found {
			distribution[index] = 1
		} else {
			distribution[index] += 1
		}
	}
	log.Println(distribution)
}
