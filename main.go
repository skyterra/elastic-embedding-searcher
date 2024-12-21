package main

import (
	"github.com/skyterra/elastic-embedding-searcher/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
