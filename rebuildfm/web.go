package rebuildfm

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
	"net/http"
	"reflect"
	"strings"
)

type APIError struct {
	Code    int
	Message string
}

func RunServer(client *elastic.Client, addr string) {
	// https://echo.labstack.com/guide
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Static("/", "public")

	e.GET("/episodes", func(c echo.Context) error {
		castName := c.QueryParam("cast_name")
		title := c.QueryParam("title")
		note := c.QueryParam("note")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")

		result, err := SearchEpisodes(client, castName, title, note)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, result)
	})

	e.Run(standard.New(addr))
}

type SearchResult struct {
	Episodes []*Episode `json:"episodes"`
}

func SearchEpisodes(client *elastic.Client, castName string, title string, note string) (*SearchResult, *APIError) {
	q := elastic.NewBoolQuery()

	if len(castName) > 0 {
		q = q.Must(elastic.NewTermQuery("casts.name", strings.ToLower(castName)))
	}

	if len(title) > 0 {
		q = q.Filter(elastic.NewTermQuery("title", strings.ToLower(title)))
	}

	if len(note) > 0 {
		s := strings.ToLower(note)
		q = q.Filter(elastic.NewMatchQuery("subtitle", s), elastic.NewMatchQuery("description", s))
	}

	s := client.Search().
		Index(IndexName).
		Type(TypeName).
		Query(q).
		Sort("no", false).
		From(0).Size(100).
		Pretty(true)

	result, err := s.Do(context.TODO())

	episodes := make([]*Episode, len(result.Hits.Hits))

	if err != nil {
		return nil, &APIError{Code: 900001, Message: "Search operation failed"}
	}

	var x Episode
	for i, item := range result.Each(reflect.TypeOf(x)) {
		if e, ok := item.(Episode); ok {
			episodes[i] = &e
		}
	}

	return &SearchResult{Episodes: episodes}, nil
}
