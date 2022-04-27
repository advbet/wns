package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/advbet/wns"

	"github.com/sirupsen/logrus"
)

func main() {
	var baseURL string
	var lastFilename string
	var interval time.Duration

	flag.StringVar(&baseURL, "url", "ftp://username:password@ftp.betradar.com/wns", "Base URL for FTP-pull delivery method")
	flag.StringVar(&lastFilename, "last", "", "Name of last dowloaded odds documents")
	flag.DurationVar(&interval, "interval", time.Minute, "Recheck interval for detecting new odds documents")
	flag.Parse()

	c, err := wns.NewFTPPull(baseURL)
	if err != nil {
		logrus.WithError(err).Fatal("creating FTP-pull client")
	}

	stream := c.Stream(context.TODO(), lastFilename, interval)
	for msg := range stream {
		if msg.Error != nil {
			logrus.WithError(msg.Error).Error("stream error")
			continue
		}
		fmt.Printf("==== %s ====\n", msg.Filename)
		fmt.Println(msg.Data)
	}
}
