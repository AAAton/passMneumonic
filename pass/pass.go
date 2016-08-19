package pass

import (
	"crypto/rand"
	"math/big"
	"unicode/utf8"
)

const passcharset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,?!#'-;:@$€%"`

var runeset []rune

func createRuneSet() []rune {
	var set []rune
	b := []byte(passcharset)

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		set = append(set, r)
		b = b[size:]
	}
	return set
}

func NewPass(length int) string {
	if len(runeset) == 0 {
		runeset = createRuneSet()
	}
	var pass string
	for i := 0; i < length; i++ {
		index, err := randomInt(len(runeset))
		if err != nil {
			i--
			continue
		}
		runeValue := runeset[index]
		pass += string(runeValue)
	}
	return pass
}

func randomInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return -1, err
	}
	n := nBig.Int64()
	return int(n), nil
}
