package main

import (
	"log"
	"time"
)

func main() {
	i := 0
	j := 10

	for i < j {
		log.Printf("[%dof%d] Hello, World!", i, j)
		time.Sleep(1 * time.Second)
		i++
	}
}
