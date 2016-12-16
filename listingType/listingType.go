package listingType

import (
    "bufio"
    "os"
    "log"
    "encoding/json"
)

type Data struct {
    Title        string `json:"title"`        // description of productType for sale
    Manufacturer string `json:"manufacturer"` // who manufactures the productType for sale
    Currency     string `json:"currency"`     // currency code, e.g. USD, CAD, GBP, etc.
    Price        string `json:"price"`        // price, e.g. 19.99, 100.00
}

func Read(filename string, c chan *Data) {
    defer close(c)
    var scanner *bufio.Scanner
    if f, err := os.Open(filename); err != nil {
        log.Println(err)
        os.Exit(1)
    } else {
        defer f.Close()
        scanner = bufio.NewScanner(f)
    }

    for scanner.Scan() {
        listing := Data{}
        if err := json.Unmarshal(scanner.Bytes(), &listing); err != nil {
            log.Panic(err)
        }
        c <- &listing
    }
}
