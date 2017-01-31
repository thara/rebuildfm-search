package rebuildfm

import (
	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Episode struct {
	No          int     `json:"no"`
	Title       string  `json:"title"`
	Link        string  `json:"link"`
	Description string  `json:"description"`
	Subtitle    string  `json:"subtitle"`
	Casts       []*Cast `json:"casts"`
}

type Cast struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

func SetupIndex(client *elastic.Client) {
	ctx := context.Background()
	exists, err := client.IndexExists(IndexName).Do(ctx)
	if err != nil {
		panic(err)
	}

	if !exists {
		// Create an index
		_, err := client.CreateIndex(IndexName).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
	}
}

func AddEpisodes(client *elastic.Client, episodes []*Episode) {
	service := elastic.NewBulkService(client)

	for _, e := range episodes {
		r := elastic.NewBulkIndexRequest().
			Index(IndexName).
			Type(TypeName).
			Doc(e)
		service.Add(r)
	}

	service.Do(context.TODO())
}
