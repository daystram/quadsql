package utils

func Exp2(x int) int {
	y := 1
	for i := 0; i < x; i++ {
		y *= 2
	}
	return y
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
