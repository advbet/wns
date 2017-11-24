package main

import (
	"flag"
	"fmt"

	"bitbucket.org/advbet/wns"

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

	files, err := c.List()
	if err != nil {
		logrus.WithError(err).Fatal("listing remote files")
	}

	for _, f := range files {
		fmt.Printf("\t%s\n", f)
	}
	fmt.Printf("Total: %d\n", len(files))
}
