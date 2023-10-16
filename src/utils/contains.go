package utils

func Contains[T comparable](arr []T, toFind T) bool {
	for _, k := range arr {
		if k == toFind {
			return true
		}
	}
	return false
}
