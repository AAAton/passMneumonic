package structs

import (
	"errors"
	"strconv"
)

//Ngram is a struct for containing an N-gram
type Ngram struct {
	Words []string
	Count int
}

//ToString returns a string ending in a newline
func (n Ngram) ToString() string {
	var s string
	for _, word := range n.Words {
		s += word + " "
	}
	return s + strconv.Itoa(n.Count) + "\n"
}

//ToBytes returns ToString in byte form, with a max cap
func (n Ngram) ToBytes(maxLen int) ([]byte, error) {
	s := []byte(n.ToString())
	if len(s) > maxLen {
		return nil, errors.New("Too long")
	}
	byteRay := []byte(n.ToString())
	return byteRay, nil
}

//Key generates a key from the first n-1 words
func (n Ngram) Key() string {
	var key string
	for i, w := range n.Words {
		if i < len(n.Words)-1 {
			key += w + " "
		}
	}
	key = key[:len(key)-1]
	return key
}
