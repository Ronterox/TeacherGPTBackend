package main

import (
	"log"
)

func main() {
	completion, err := gpt("What is the meaning of life?")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(completion)
}
