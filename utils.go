package main

import "errors"

func min(values []int) (int, error) {
	if len(values) == 0 {
		return 0, errors.New("cannot detect a minimum value in an empty slice")
	}

	minValue := values[0]

	for _, v := range values {
		if v < minValue {
			minValue = v
		}
	}

	return minValue, nil
}
