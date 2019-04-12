package main

import (
	"os"
)

func main() {
	f, _ := os.Create("test1")
	f.WriteString("1")
	defer f.WriteString("2")

	f, _ = os.Create("test2")
	f.WriteString("a")

}
