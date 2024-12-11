package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	var lastValue = ""
	var isSlash = false

	for _, val := range str {
		value := string(val)
		num, err := strconv.Atoi(value)

		if err == nil {
			switch {
			case lastValue == "":
				return "", ErrInvalidString
			case isSlash:
				result.WriteString(lastValue)
				lastValue = value
				isSlash = false
			default:
				result.WriteString(strings.Repeat(lastValue, num))
				lastValue = ""
			}
		} else {
			switch {
			case isSlash && value != `\`:
				return "", ErrInvalidString
			case !isSlash && value == `\`:
				isSlash = true
			case isSlash && value == `\`:
				isSlash = false
				fallthrough
			case lastValue != "":
				result.WriteString(lastValue)
				fallthrough
			default:
				lastValue = value
			}
		}
	}

	result.WriteString(lastValue)

	return result.String(), nil
}
