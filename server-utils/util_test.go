package server_utils

import (
	"testing"
)

func TestSampleLaplacianDistribution(t *testing.T) {
	mean := 0
	sensitivity := float64(4032)
	epsilon := 0.1
	delta := 0.02
	for i := 0; i < 1000; i++ {
		_ = SampleLaplacianDistribution(int64(mean), sensitivity, epsilon, delta)
	}
}
