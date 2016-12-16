package cli

import (
        "flag"
        "fmt"
        "os"
)

func ParseCommandLine() (listings, products, output string) {
	listingsFilename := flag.String("listings", "", "path to listings file")
	productsFilename := flag.String("products", "", "path to products file")
	outputFilename := flag.String("output", "", "path to output file")
	flag.Parse()

	if *listingsFilename == "" {
		fmt.Fprintln(os.Stderr, "Please provide path to listings file.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *productsFilename == "" {
		fmt.Fprintln(os.Stderr, "Please provide path to products file.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outputFilename == "" {
		fmt.Fprintln(os.Stderr, "Please provide path to output file.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	listings = *listingsFilename
	products = *productsFilename
	output = *outputFilename
	return
}
