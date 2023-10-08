package utils

import "errors"

func FailIfZeroValue[T comparable](data []T) error {
	var zeroValue T
	for _, v := range data {
		if v == zeroValue {
			return errors.New("one of the required values is empty")
		}
	}
	return nil
}
