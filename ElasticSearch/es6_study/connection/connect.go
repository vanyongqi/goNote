package connection

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

func Connect() {
	// TODO
	// Create a client and connect to http://192.168.2.10:9201
	client, err := elastic.NewClient(elastic.SetURL("http://10.10.10.124:9200"))
	if err != nil {
		// Handle error
	} else {
		fmt.Println("Connected to elasticsearch")
	}
	exists, err := client.IndexExists("twitter").Do(context.Background())
	if err != nil {
		// Handle error
	}
	if !exists {
		// Index does not exist yet.
		fmt.Println("index not exist")
	}

	exists1, err := client.IndexExists("fofaee_assets").Do(context.Background())
	if err != nil {
		// Handle error
	}
	if !exists1 {
		// Index does not exist yet.
		fmt.Println("index not exist")
	} else {
		fmt.Println("index exist")
	}

}
