package main

import (
	"flag"
	"fmt"
	"passMneumonic/pass"
	"passMneumonic/processer"
)

var bimap, trimap map[string][]string

func main() {

	passLength := flag.Int("l", 10, "supply the length of the required password")
	m := flag.Bool("m", false, "Do you want a mneumonic with it?")
	i := flag.Bool("i", false, "Interactive mode")
	flag.Parse()
	if *i {
		interactive()
	}
	password := pass.NewPass(*passLength)

	if !*m {
		fmt.Println(password)
	} else {
		fmt.Println("Creating trigram map")
		trigrams := processer.OpenBytes(processer.Trigrams)
		trimap = processer.CreateMap(trigrams)

		fmt.Println("Creating bigram map")
		bigrams := processer.OpenBytes(processer.Bigrams)
		bimap = processer.CreateMap(bigrams)
		fmt.Println(password)
		val, success := createMneumonic(password)
		if !success {
			fmt.Println("failed...")
		}
		fmt.Println(val)

	}
}

func createMneumonic(password string) (string, bool) {
	runes := []rune(password)
	//Find starting word
	r := runes[0]
	for _, val := range bimap["<<boundary>>"] {
		if startsWith(val, r) {
			if s, exist := findInit(val, runes); exist {
				return s, true
			}
		}

	}
	return "", false
}

func interactive() {

	fmt.Println("Creating trigram map")
	trigrams := processer.OpenBytes(processer.Trigrams)
	trimap = processer.CreateMap(trigrams)

	fmt.Println("Creating bigram map")
	bigrams := processer.OpenBytes(processer.Bigrams)
	bimap = processer.CreateMap(bigrams)

	for {
		var i int
		fmt.Print("Password length:")
		_, err := fmt.Scan(&i)
		if err != nil {
			fmt.Println("not a number")
			continue
		}

		var password, mneumonic string
		done := false
		for !done {
			password = pass.NewPass(i)
			mneumonic, done = createMneumonic(password)
		}
		fmt.Println(mneumonic)
		fmt.Println(password)
	}
}

func startsWith(val string, other rune) bool {
	valRunes := []rune(val)
	return len(valRunes) > 0 && valRunes[0] == other
}

func findInit(key string, runes []rune) (string, bool) {
	ray, success := find([]string{key}, runes, 1)
	if !success {
		return "", success
	}
	s := ""
	for _, v := range ray {
		s += v + " "
	}
	return s, true
}

func find(key []string, runes []rune, i int) ([]string, bool) {
	if i == len(runes) {
		return key, true
	}

	//Start with checking trigrams
	if len(key) > 1 {
		triKey := keyLength(key, 2)
		for _, val := range trimap[triKey] {
			if startsWith(val, runes[i]) {
				fmt.Println("Found a trigram:", key, val)
				if newVal, solves := find(append(key, val), runes, i+1); solves {
					return newVal, true
				}
			}
		}
	}

	bimapKey := key[len(key)-1]
	for _, val := range bimap[bimapKey] {
		if startsWith(val, runes[i]) {
			if newVal, solves := find(append(key, val), runes, i+1); solves {
				return newVal, true
			}
		}
	}
	return []string{}, false
}

func keyLength(keyArray []string, l int) string {
	key := ""
	for _, k := range keyArray[len(keyArray)-l:] {
		key += k + " "
	}
	return key[:len(key)-1]
}
