package monitor

func getAverageInt64(input []int64) float32 {
	length := len(input)
	if length == 0 {
		return 0
	}
	sum := int64(0)
	for _, v := range input {
		sum += v
	}
	return float32(sum) / float32(length)
}
