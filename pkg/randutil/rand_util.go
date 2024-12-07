package randutil

import "math/rand"

const (
	AlphabetNumeric           = "0123456789"
	AlphabetLowerAlpha        = "abcdefghijklmnopqrstuvwxyz"
	AlphabetLowerAlphaNumeric = AlphabetNumeric + AlphabetLowerAlpha
)

// LowerAlphaNumeric 生成指定长度的包含随机小写字母和数字的字符串
func LowerAlphaNumeric(length int) string {
	return RandString(AlphabetLowerAlphaNumeric, length)
}

// RandString 从 alphabet 中随机挑选字符生成长度为 length 的随机字符串
func RandString(alphabet string, length int) string {
	ret := make([]byte, length)
	alphabetLen := uint64(len(alphabet))
	i := 0
	for {
		randInt := rand.Uint64()
		for randInt > 0 {
			ret[i] = alphabet[randInt%alphabetLen]
			i++
			if i >= length {
				return string(ret)
			}
			randInt /= alphabetLen
		}
	}
}
