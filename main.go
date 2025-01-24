package main

import (
	"fmt"
	"os"
	"osekpahesh/internal/flags"
	"osekpahesh/internal/osek"
)

func main() {
	fl := flags.New()
	fl.Parse()

	myOsek := osek.New()
	if err := myOsek.LoadConfiguration(fl.GetPathToConfig()); err != nil {
		fmt.Printf("Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	if fl.IsPrintTransactions() {
		myOsek.PrintTransactions()
	}

}
