package processer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"passMneumonic/structs"
	"strconv"
	"strings"
	"sync"
)

type byteRay struct {
	sync.RWMutex
	array   []byte
	rowSize int
}

//BigramsRaw ...
var BigramsRaw = "/mnt/2_TB_HD/dataset/ngrams/en.2grams"

//Bigrams is path to bigrams
var Bigrams = "/mnt/2_TB_HD/dataset/ngrams/en.2grams.bytes"

//TrigramsRaw ....
var TrigramsRaw = "/mnt/2_TB_HD/dataset/ngrams/en.3grams"

//Trigrams ...
var Trigrams = "/mnt/2_TB_HD/dataset/ngrams/en.3grams.bytes"

//ProcessNgrams converts the dataset given by Roverto Twitter Corpus into something compact
func ProcessNgrams(path string, rowSize int) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var ray byteRay
	ray.rowSize = rowSize

	count := 0
	scanner := bufio.NewScanner(file)
	var wg *sync.WaitGroup
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go ray.processLine(line, wg)
		count++
		if count%100000 == 0 {
			fmt.Printf("%d, %.1f%%\n", count, float64(count)/float64(19746857)*100)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()

	ioutil.WriteFile(path+".bytes", ray.array, 0777)
}

//CreateMap takes a list of Ngrams and returns a map where the key is the first n-1 words and the value i the last word
func CreateMap(ngrams []structs.Ngram) map[string][]string {
	nmap := make(map[string][]string)
	for _, n := range ngrams {
		key := n.Key()
		if _, exist := nmap[key]; !exist {
			nmap[key] = make([]string, 1)
		}
		nmap[key] = append(nmap[key], n.Words[len(n.Words)-1])
	}
	return nmap
}

//OpenBytes reads the compact dataformat
func OpenBytes(path string) []structs.Ngram {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var ngrams []structs.Ngram

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, " ")
		c, err := strconv.Atoi(words[len(words)-1])
		if err == nil {
			words = words[:len(words)-1]
			n := structs.Ngram{words, c}
			ngrams = append(ngrams, n)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ngrams
}

func (ray *byteRay) processLine(line string, wg *sync.WaitGroup) {
	defer wg.Done()
	parts := strings.Split(line, "\t")
	if len(parts) > 1 {
		words := strings.Split(parts[0], " ")
		b := structs.Ngram{words, 0}
		c := 0
		for i, part := range parts {
			if i > 0 {
				val, err := strconv.Atoi(part)
				if err == nil {
					c += val
				}
			}
		}
		b.Count = c
		bytes, err := b.ToBytes(ray.rowSize)
		if err == nil {
			ray.Lock()
			ray.array = append(ray.array, bytes...)
			ray.Unlock()
		}
	}
}
