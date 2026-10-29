# IPClues Go SDK

Official Go client for the IPClues IP intelligence API.

- Module: `github.com/ipclues/go-sdk`
- Go: 1.22+

## Installation

```bash
go get github.com/ipclues/go-sdk
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ipclues "github.com/ipclues/go-sdk"
)

func main() {
	client := ipclues.NewClient("YOUR_API_KEY")

	response, err := client.Lookup(context.Background(), "1.1.1.1", &ipclues.LookupOptions{
		Databases: []ipclues.Database{
			ipclues.DatabaseGeoI3Country,
			ipclues.DatabaseNetI3Identity,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// HTTP errors are envelopes too: inspect response.Metadata.Status.
	fmt.Printf("%#v\n", response.Data)
	fmt.Println(response.Metadata.Status, response.Metadata.StatusText)
}
```

## License

MIT
