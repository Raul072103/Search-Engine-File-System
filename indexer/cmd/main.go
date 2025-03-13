package main

import (
	"fmt"
)

func main() {
	var s = ""
	fmt.Printf(s[:min(len(s), 200)])
}
