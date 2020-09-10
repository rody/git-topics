package main

import (
	"fmt"
	"os"

	"github.com/rody/find-commits/cmd/topics"
)

func main() {

	if err := topics.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
