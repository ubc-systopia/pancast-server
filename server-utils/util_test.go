package server_utils

import (
	"fmt"
	"testing"
)

func TestSampleLaplacianDistribution(t *testing.T) {
	mean := 0
	sensitivity := float64(4032)
	epsilon := 0.1
	delta := 0.02
	for i := 0; i < 1000; i++ {
		sample := SampleLaplacianDistribution(mean, sensitivity, epsilon, delta)
		fmt.Println(sample)
	}
}
