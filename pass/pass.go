package pass

import (
	"crypto/rand"
	"math/big"
)

// const passcharset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,?!#'-;:@$€%"`
const passcharset = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,?!#`

//NewPass generates a password from passcharset with the given length
func NewPass(length int) string {
	runeset := []rune(passcharset)
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
