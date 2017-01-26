package main

import (
	"flag"
	"fmt"
	rebuildfm "github.com/tomochikahara/rebuildfm-search/rebuildfm"
	elastic "gopkg.in/olivere/elastic.v5"
	"os"
)

var usage = `Usage rebuildfm <Command>
Commands:
  polling Polling RSS feed
  runserver Run web API server
`

func main() {

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Create a client
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	switch args[0] {
	case "polling":
		fmt.Println("Start polling RSS feed.")

		rebuildfm.SetupIndex(client)
		rebuildfm.PollFeed(client, "http://feeds.rebuild.fm/rebuildfm", 5, nil)

	case "runserver":
		fmt.Println("Start running Web API server..")
		rebuildfm.RunServer(client)

	default:
		flag.Usage()
		os.Exit(1)
	}
}
