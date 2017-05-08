package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"bitbucket.org/advbet/wns"
)

func main() {
	var name string
	var key string
	var url string
	var delete bool

	flag.StringVar(&name, "name", "", "Bookmaker name/ID")
	flag.StringVar(&key, "key", "", "Feed access key")
	flag.StringVar(&url, "url", "https://www.betradar.com/betradar/getXmlFeed.php", "Feed HTTP pull base URL")
	flag.BoolVar(&delete, "delete", false, "Delete feed document after transfer")
	flag.Parse()

	client := wns.HTTPPullClient{
		Username: name,
		Key:      key,
		URL:      url,
		Interval: 30 * time.Second,
		HTTPClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}

	data, err := client.Get(context.TODO(), delete)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(data)
}
