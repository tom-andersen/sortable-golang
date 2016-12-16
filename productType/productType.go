package productType

import (
    "bufio"
    "os"
    "log"
    "encoding/json"
)

type Data struct {
    ProductName   string `json:"product_name"`   // A unique id for the productType
    Manufacturer  string `json:"manufacturer"`   //
    Family        string `json:"family"`         // optional grouping of products
    Model         string `json:"model"`          //
    AnnouncedDate string `json:"announced-date"` // ISO-8601 formatted date string, e.g. 2011-04-28T19:00:00.000-05:00
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
        product := Data{}
        if err := json.Unmarshal(scanner.Bytes(), &product); err != nil {
            log.Panic(err)
        }
        c <- &product
    }
}

