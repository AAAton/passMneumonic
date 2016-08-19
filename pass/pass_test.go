package pass

import (
	"fmt"
	"testing"
)

func BenchmarkPass(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPass(16)
	}
}

func TestCreateRuneSet(t *testing.T) {
	set := []rune(passcharset)
	if len(set) != 76 {
		t.Error("runeset is not the right length")
	}
}

func TestPass(t *testing.T) {
	counter := make(map[string]int64)
	for i := 0; i < 1000000; i++ {
		s := NewPass(1)
		counter[s]++
	}
	var max, min int64
	var maxKey, minKey string
	for key, val := range counter {
		if val > max {
			max = val
			maxKey = key
		}
		if min == 0 || val < min {
			min = val
			minKey = key
		}
	}

	ratio := float64(min) / float64(max)
	if ratio < .9 {
		for i, c := range []rune(passcharset) {
			fmt.Printf("%d %s %d\n", i, string(c), counter[string(c)])
		}
		fmt.Printf("Max: %s %d, Min: %s %d\n", maxKey, max, minKey, min)
		fmt.Printf("Ratio between min and max: %.2f\n", ratio)
		t.Error("Not very uniform...")
	}
}
