package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/advbet/wns"

	"github.com/sirupsen/logrus"
)

func main() {
	var baseURL string

	flag.StringVar(&baseURL, "url", "ftp://username:password@ftp.betradar.com/wns", "Base URL for FTP-pull delivery method")
	flag.Parse()

	c, err := wns.NewFTPPull(baseURL)
	if err != nil {
		logrus.WithError(err).Fatal("creating FTP-pull client")
	}

	docs, lastFile, err := c.Snapshot()
	if err != nil {
		logrus.WithError(err).Fatal("getting state snapshot")
	}

	for _, d := range docs {
		bs, err := json.Marshal(d)
		if err != nil {
			logrus.WithError(err).Fatal("marshaling wns document to JSON")
		}
		fmt.Printf("%s\n", bs)
	}
	fmt.Printf("\nLast filename: %s\n", lastFile)
}
