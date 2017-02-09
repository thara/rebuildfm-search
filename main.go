package main

import (
	"flag"
	"fmt"
	rebuildfm "github.com/tomochikahara/rebuildfm-search/rebuildfm"
	elastic "gopkg.in/olivere/elastic.v5"
	"log"
	"os"
)

func main() {
	commands := map[string]command{
		"aggregate": aggregateCmd(),
		"runserver": runserverCmd(),
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.Usage = func() {
		fmt.Println("Usage: rebuildfm-search <command> [command options]")
		for name, cmd := range commands {
			fmt.Printf("\n%s command:\n", name)
			cmd.fs.PrintDefaults()
		}

	}

	fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	if cmd, ok := commands[args[0]]; !ok {
		log.Fatal("Unknown command: %s", args[0])
	} else if err := cmd.fn(args[1:]); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type command struct {
	fs *flag.FlagSet
	fn func(args []string) error
}

func aggregateCmd() command {
	fs := flag.NewFlagSet("rebuildfm-search aggregate", flag.ExitOnError)

	var elasticUrl string
	fs.StringVar(&elasticUrl, "elastic-url", "http://localhost:9200", "ElasticSearch URL")

	return command{fs, func(args []string) error {
		fs.Parse(args)

		fmt.Printf("\nelasticsearch URL: %s\n", elasticUrl)

		// Create a client
		// See https://github.com/olivere/elastic/wiki/Connection-Problems#how-to-figure-out-connection-problems
		client, err := elastic.NewClient(
			elastic.SetURL(elasticUrl), elastic.SetSniff(false))
		if err != nil {
			return err
		}

		rebuildfm.SetupIndex(client)
		rebuildfm.PollFeed(client, "http://feeds.rebuild.fm/rebuildfm", 5, nil)

		return nil
	}}
}

type runserverOpts struct {
	elasticUrl string
	addr       string
	siteUrl    string
}

func runserverCmd() command {
	fs := flag.NewFlagSet("rebuildfm-search runserver", flag.ExitOnError)
	opts := &runserverOpts{}

	fs.StringVar(&opts.elasticUrl, "elastic-url", "http://localhost:9200", "ElasticSearch URL")
	fs.StringVar(&opts.addr, "addr", ":8080", "Listen Address and Port")
	fs.StringVar(&opts.siteUrl, "siteUrl", "http://127.0.0.1:8080", "API base url")

	return command{fs, func(args []string) error {
		fs.Parse(args)

		// Create a client
		// See https://github.com/olivere/elastic/wiki/Connection-Problems#how-to-figure-out-connection-problems
		client, err := elastic.NewClient(
			elastic.SetURL(opts.elasticUrl), elastic.SetSniff(false))
		if err != nil {
			return err
		}

		fmt.Printf("Start running Web API server at %s ..\n", opts.addr)

		apiBaseUrl := opts.siteUrl + "/_api"
		rebuildfm.RunServer(client, opts.addr, apiBaseUrl)

		return nil
	}}
}
