package utils

// Source: https://learngolang.net/tutorials/how-to-convert-integer-to-english-words-in-golang/
// https://github.com/divan/num2words/blob/master/num2words.go You can also check this package

import (
	"math"
)

func pow(i int, p int) int {
	return int(math.Pow(1000, float64(p)))
}

func Num2words(n int) string {
	to19 := []string{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve",
		",Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}

	tens := []string{"Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
	if n == 0 {
		return ""
	}
	if n < 20 {
		return to19[n-1]
	}
	if n < 100 {
		return tens[n/10-2] + " " + Num2words(n%10)
	}
	if n < 1000 {
		return to19[n/100-1] + " Hundred " + Num2words(n%100)
	}

	for idx, w := range []string{"Thousand", "Million", "Billion"} {
		p := idx + 1
		if n < pow(1000, (p+1)) {
			return Num2words(n/pow(1000, p)) + " " + w + " " + Num2words(n%pow(1000, p))
		}
	}

	return "error"
}
