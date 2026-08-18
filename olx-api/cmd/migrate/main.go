package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println(os.Args[1])
	if len(os.Args) < 2 {
		log.Fatal("usage: make migrate <up | down>")
	}
	switch os.Args[1] {
	case "up":
		log.Printf("up called")
	case "down":
		log.Printf("down")
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
