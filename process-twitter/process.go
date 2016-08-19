package process

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type bigram struct {
	first  string
	second string
	count  int
}

func processBigrams(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var ray []byte
	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) > 1 {
			words := strings.Split(parts[0], " ")
			b := bigram{words[0], words[1], 0}
			c := 0
			for i, part := range parts {
				if i > 0 {
					val, err := strconv.Atoi(part)
					if err == nil {
						c += val
					}
				}
			}
			b.count = c
			bytes, err := b.toBytes()
			if err == nil {
				ray = append(ray, bytes...)
			}
		}
		count++
		if count%100000 == 0 {
			fmt.Printf("%d, %.1f%%\n", count, float64(count)/float64(19746857)*100)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(path+".bigrams", ray, 0777)
}

func (b bigram) toString() string {
	return b.first + " " + b.second + " " + strconv.Itoa(b.count) + "\n"
}

func (b bigram) toBytes() ([]byte, error) {
	maxLen := 64
	byteRay := make([]byte, maxLen, maxLen)
	s := []byte(b.toString())
	if len(s) > maxLen {
		return nil, errors.New("Too long")
	}
	byteRay = append(byteRay, []byte(b.toString())...)
	return byteRay, nil
}
