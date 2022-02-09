package base

import "strconv"

func ParsePrice(str string) (int64, error) {
	studentPrice, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return int64(studentPrice * 100), err
}
