package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/beeploop/aes-encrypt/frontend/cli"
)

func main() {
	headless := flag.Bool("headless", false, "run without gui")
	help := flag.Bool("help", false, "display help")
	key := flag.String("key", "", "32 character AES key. You can pass in a txt file")
	input := flag.String("input", "", "Input file")
	output := flag.String("output", "", "Output location with filename, will use default filename if not specified.")
	decrypt := flag.Bool("decrypt", false, "run in decrypt mode")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	switch *headless {
	case true:
		fmt.Println("Running in headless mode")

		app, err := cli.New()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Set input file
		if err := app.SetInputFile(*input); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Set key
		if err := app.SetKey(*key); err != nil {
			fmt.Println("Error setting key: ", err)
			os.Exit(1)
		}

		// Set output location
		if err := app.SetOutput(*output); err != nil {
			fmt.Println("Error setting output: ", err)
			os.Exit(1)
		}

		// Set mode
		if *decrypt {
			app.SetMode(cli.DECRYPT)
		} else {
			app.SetMode(cli.ENCRYPT)
		}

		// Run
		if err := app.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case false:
		fmt.Println("GUI mode currently not supported")
		os.Exit(1)
	}
}
