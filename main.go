package main

import (
	"fmt"
	"log"
	"os"

	"tv-shows-manager/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tv-shows-manager <csv_file>")
		os.Exit(1)
	}

	csvFile := os.Args[1]

	if err := ui.Run(csvFile); err != nil {
		log.Fatal(err)
	}
}