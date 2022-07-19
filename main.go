package main

import (
	"fmt"
	"os"

	"github.com/68696c6c/capricorn_rnd/cmd"
)

func main() {
	if err := cmd.Capricorn.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
