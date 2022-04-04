package style

import (
	"math"
)

var n []int

func init() {
	n = make([]int, 0)
	for i, v := range []int{47, 68, 40, 40, 40, 21} {
		for j := 0; j < v; j++ {
			n = append(n, i)
		}
	}
}

// rgbToXterm returns nearest 256 color for given rgb
// source: https://stackoverflow.com/a/62219320
func rgbToXterm(r, g, b int) int {
	mx := math.Max(math.Max(float64(r), float64(g)), float64(b))
	mn := math.Min(math.Min(float64(r), float64(g)), float64(b))

	if (mx-mn)*(mx+mn) <= 6250 {
		c := 24 - (252-((r+g+b)%3))%10
		if 0 <= c && c <= 23 {
			return 232 + c
		}
	}
	return 16 + 36*n[r] + 6*n[g] + n[b]
}
