package isbn

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func Valid(n string) bool {
	str := strings.ReplaceAll(n, "-", "")

	if len(str) == 10 {
		return valid10(str)
	} else if len(str) == 13 {
		return valid13(str)
	}
	return false
}

func valid10(n string) bool {
	var sum int
	buf := make([]byte, 1)
	for i, r := range n[:9] {
		_ = utf8.EncodeRune(buf, r)
		val, err := strconv.Atoi(string(buf))
		if err != nil {
			return false
		}
		sum += val * (10 - i)
	}

	cs := (11 - sum%11) % 11
	if cs == 10 {
		return n[9] == 'x' || n[9] == 'X'
	}
	_ = utf8.EncodeRune(buf, rune(n[9]))
	val, err := strconv.Atoi(string(buf))
	if err != nil {
		return false
	}
	return cs == val
}

func valid13(n string) bool {
	var sum, alt int
	buf := make([]byte, 1)
	for i, r := range n {
		if i%2 == 0 {
			alt = 1
		} else {
			alt = 3
		}
		_ = utf8.EncodeRune(buf, r)
		val, err := strconv.Atoi(string(buf))
		if err != nil {
			return false
		}
		if i == 12 { // checksum check
			return (10 - sum%10) == val
		} else {
			sum += val * alt
		}
	}
	return false // never reached
}
