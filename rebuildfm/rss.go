package rebuildfm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	elastic "gopkg.in/olivere/elastic.v5"
)

func PollFeed(client *elastic.Client, uri string, timeout int, cr xmlx.CharsetFunc) {
	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		if err := feed.Fetch(uri, cr); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s\n", uri, err)
			return
		}

		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))
	}
}

func chanHandler(feed *rss.Feed, newchannels []*rss.Channel) {
	fmt.Printf("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {

	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	l := len(newitems)
	episodes := make([]*Episode, l)

	for i, item := range newitems {
		fmt.Printf("%s\n", item.Title)
		fmt.Printf("%s\n", item.Links[0].Href)
		fmt.Printf("%s\n", item.Description)

		itunes := item.Extensions["http://www.itunes.com/dtds/podcast-1.0.dtd"]
		subtitle := itunes["subtitle"][0].Value

		// duration := itunes["duration"][0].Value
		// fmt.Printf("duration = %s\n", duration)

		contributors := item.Extensions["http://www.w3.org/2005/Atom"]["contributor"]

		casts := make([]*Cast, len(contributors))
		for j, c := range contributors {
			name := c.Childrens["name"][0].Value
			uri := c.Childrens["uri"][0].Value
			casts[j] = &Cast{Name: name, Uri: uri}
		}

		no := l - i

		episode := &Episode{
			No:          no,
			Title:       item.Title,
			Link:        item.Links[0].Href,
			Description: item.Description,
			Subtitle:    subtitle,
			Casts:       casts,
		}
		episodes[i] = episode
	}
	AddEpisodes(client, episodes)

	fmt.Printf("%d new item(s) in %s\n", len(episodes), feed.Url)
}

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "iso-8859-1" {
		return r, nil
	}
	return nil, errors.New("Unsupported character set encoding: " + charset)
}
