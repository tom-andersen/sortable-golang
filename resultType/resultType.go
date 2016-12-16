package resultType

import (
    "../listingType"
    "bufio"
    "os"
    "log"
    "encoding/json"
)

type Data struct {
    ProductName string    `json:"product_name"`
    Listings    []listingType.Data `json:"listings"`
}

func Write(filename string, c chan *Data) (matchedChannel chan int) {
	matchedChannel = make(chan int)
    go func(matchedChannel chan int) {
        defer close(matchedChannel)
        matched := 0
        var writer *bufio.Writer
        if f, err := os.Create(filename); err != nil {
            log.Println(err)
            os.Exit(1)
        } else {
            defer f.Close()
            writer = bufio.NewWriter(f)
        }

        for result := range c {
            matched += len(result.Listings)
            if bytes, err := json.Marshal(result); err == nil {
                if _, err := writer.Write(bytes); err != nil {
                    log.Fatal(err)
                }
                if _, err := writer.WriteRune('\n'); err != nil {
                    log.Fatal(err)
                }
            } else {
                log.Fatal(err)
            }
        }
        writer.Flush()
        matchedChannel <- matched
    }(matchedChannel)
    return
}

