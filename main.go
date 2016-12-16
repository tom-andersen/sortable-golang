package main

import (
	"runtime"
	"./cli"
	"./listingType"
    "./productType"
	"./resultType"
	"./productByManufacturerType"
	"fmt"
)

func main() {
	runtime.GOMAXPROCS(4)

	listingsFilename, productsFilename, outputFilename := cli.ParseCommandLine()

	listingChannel := make(chan *listingType.Data, 32)
	go listingType.Read(listingsFilename, listingChannel)

	productChannel := make(chan *productType.Data, 32)
	go productType.Read(productsFilename, productChannel)

	manufacturer := productByManufacturerType.New()
	for product := range productChannel {
		manufacturer.AddProduct(product)
	}

	manufacturer.AddAlias("hp", "hewlett packard")
	manufacturer.AddAlias("konica minolta", "minolta")
	manufacturer.AddAlias("konica minolta", "konica")
	manufacturer.AddAlias("fujifilm", "fuji")
	manufacturer.AddAlias("kodak", "eastman kodak")

	matchedChannel := resultType.Write(outputFilename, manufacturer.ResultChannel)

	cnt := 0
	notfoundCnt := 0
	for listing := range listingChannel {
		cnt++
		if sent := manufacturer.SendListingOnManufacturerChannel(listing); !sent {
				notfoundCnt++
		}
	}
	manufacturer.Close()

	fmt.Println(cnt, "listings read.")
	fmt.Println(notfoundCnt, "listing with unknown manufucturer.")
	fmt.Println(<- matchedChannel, "matches written.")
}



