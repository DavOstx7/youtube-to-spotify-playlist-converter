package utils

func SliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func CollectFromChannel[T any](channel <-chan T) []T {
	var s []T
	for v := range channel {
		s = append(s, v)
	}
	return s
}
