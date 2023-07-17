package core

import (
	"github.com/ananaslegend/short-link/utils"
	"strings"
)

const alphabet = "1qazQAZ2wsxWSX3edcEDC4rfvRFV5tgbTGB6yhnYHN7ujmUJM8ikK9oLPp"

var alphabetLen = uint32(len(alphabet))

func MakeShorter(id uint32) string {
	var (
		digits  []uint32
		num     = id
		builder strings.Builder
	)

	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}

	utils.Reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}

	return builder.String()
}
