package util

// Distinct 函数用于对输入的切片进行去重，支持整型、浮点型和字符串类型
func Distinct[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string](s []T) []T {
	if len(s) == 0 {
		return s
	}
	seen := make(map[T]struct{})
	result := make([]T, 0, len(s))
	for _, item := range s {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
