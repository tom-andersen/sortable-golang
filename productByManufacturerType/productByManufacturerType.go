package productByManufacturerType

import (
    "../productType"
    "../listingType"
    "../resultType"
    "strings"
    btree "github.com/google/btree"
    "sync"
)

// The data structure wraps a tree implementation for looking up manufacturers.
// Each manufacturer has a list of associated products.
// Every time a new manufacturer is added, a dedicated processor is started to handle listings.
// The processors emit their final results on resultChannels
type Data struct {
    tree              *btree.BTree             // Google BTree implementation providing a sorted/searchable data structure.
    productChannels   []chan *listingType.Data // To allow concurrent processing, listings of manufactuer are sent over channel.
    ResultChannel     chan *resultType.Data    // All the processors of listings emit their final result on these channels.
    processorsRunning sync.WaitGroup           // Count of processors actively running.
}

func New() Data {
    return Data{
        tree: btree.New(32),
        ResultChannel: make(chan *resultType.Data, 32),
    }
}

// Node in tree
type ProductItem struct {
    Key      string                 // Lowercase Manufacturer Name
    Products *[]productType.Data    // List of products associated with Manufacturer
    Listings chan *listingType.Data // Channel for listings that match manufacturer name.
}

// Node comparison function required by tree.
func (a ProductItem) Less(b btree.Item) bool {
    return a.Key < b.(ProductItem).Key
}

func (x *Data) newListingChannel() (c chan *listingType.Data) {
    c = make(chan *listingType.Data, 32)
    x.productChannels = append(x.productChannels, c)
    x.processorsRunning.Add(1)
    return c
}

// Add product to existing manufacturer node.
// If manufacturer node does not exist:
// 1. Create node.
// 2. Create listing channel.
// 3. Start manufacturer processor
func (x *Data) AddProduct(product *productType.Data) {
    newItem := ProductItem{
        Key: strings.ToLower(product.Manufacturer),
    }
    if existingItem := x.tree.Get(newItem); existingItem == nil {
        newItem.Products = &[]productType.Data{*product }
        newItem.Listings = x.newListingChannel()
        x.tree.ReplaceOrInsert(newItem)
        go x.processor(newItem.Listings, newItem.Products)
    } else {
        products := existingItem.(ProductItem).Products
        for _, v := range *products {
            if v.ProductName == product.ProductName || contains(v.Model, product.Model) && contains(v.Family, product.Family) {
                // Skip duplicate
                // It was found that sometimes products are listed twice.
                return
            }
        }
        *products = append(*products, *product)
    }
}

// Add alias, such that single manufactuer can be found under differend names.
// Reuse listing channel and processor.
func (x *Data) AddAlias(canonical string, alias string) {
    productItem := x.tree.Get(ProductItem{Key: canonical}).(ProductItem)
    x.tree.ReplaceOrInsert(ProductItem{
        Key: alias,
        Products: productItem.Products,
        Listings: productItem.Listings,
    })
}

// Send listing on manufacturer channel.
func (x *Data) SendListingOnManufacturerChannel(listing *listingType.Data) (sent bool) {
    sent = false
    // Find manufacturer by only matching length of name.
    // This compensates for listings that add extra text to manufacturer name.
    item := ProductItem{Key: strings.ToLower(listing.Manufacturer) }
    x.tree.DescendLessOrEqual(item, func(i btree.Item) bool {
        product := i.(ProductItem)
        if (strings.HasPrefix(item.Key, product.Key)) {
            product.Listings <- listing
            sent = true
        }
        return false
    })
    return
}

// When no more listings will be added, we can close the channels to signal we are done.
func (x *Data) Close() {
    for _, c := range x.productChannels {
        close(c)
    }
    x.processorsRunning.Wait()
    close(x.ResultChannel)
}

// Custom matching of strings, with additional rules:
// 1. Comparison is done in lowercase.
// 2. We expect match to be entire word - no partial matches.
// 3. Hyphens are removed, concatentating words.
// 4. Hyphens are also replaced by space, separating words in string.
func contains(s string, substr string) bool {
    if len(substr) == 0 {
        return true
    }
    s = strings.ToLower(s)
    split := strings.Split(s, "-")
    s = strings.Join(split, "") + " " + strings.Join(split, " ")

    substr = strings.ToLower(substr)
    substr = strings.Join(strings.Split(substr, "-"), "")

    index := strings.Index(s, substr)
    if index == -1 {
        return false
    }
    if index + len(substr) < len(s) {
        char := s[index + len(substr)]
        if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
            return false
        }
    }
    if index > 0 {
        char := s[index - 1]
        if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
            return false
        }
    }
    return true
}

// Given a list of products, attempt match on Model and Family.
// Duplicates matches are ambiguous, so we fail to match.
func findMatch(products *[]productType.Data, listing *listingType.Data) (match *productType.Data) {
    title := strings.ToLower(listing.Title)
    for _, product := range *products {
        if contains(title, product.Model) && contains(title, product.Family) {
            if match == nil {
                copy := product
                match = &copy
            } else {
                // Ambiguous - matches more than one productType.
                return nil
            }
        }
    }
    return
}

// Processor that listens for listings, and emits results when listings channel is closed.
func (x *Data) processor(listingChannel chan *listingType.Data, products *[]productType.Data) {
    accumulator := make(map[string]*resultType.Data)
    for listing := range listingChannel {
        if match := findMatch(products, listing); match != nil {
            result := accumulator[match.ProductName]
            if result == nil {
                accumulator[match.ProductName] = &resultType.Data{
                    ProductName: match.ProductName,
                    Listings: []listingType.Data{*listing },
                }
            } else {
                result.Listings = append(result.Listings, *listing)
            }
        }
    }
    for _, v := range accumulator {
        x.ResultChannel <- v
    }
    x.processorsRunning.Done()
}

