package main

import (
	"flag"
	"fmt"
	"passMneumonic/pass"
)

func main() {
	passLength := flag.Int("l", 10, "supply the length of the required password")
	flag.Parse()
	pass := pass.NewPass(*passLength)
	fmt.Println(pass)
}
