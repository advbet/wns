package wns

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

// FTPPullClient is Betradar WNS feed client for feed consumption with FTP-pull
// delivery method.
//
// For FTP-pull client usage examples see demo applications in bin directory.
type FTPPullClient struct {
	username string
	password string
	hostname string
	baseDir  string
}

// NewFTPPull creates a new instance of FTP-pull WNS client. Parameter `baseURL`
// must be a ftp protocol URL and include username, password and path to odds
// documents (usually /wns). Port is optional, if unspecified 21 will be used.
// Example:
//
//   ftp://user:pass@ftp.betradar.com/wns
func NewFTPPull(baseURL string) (*FTPPullClient, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if url.Scheme != "ftp" {
		return nil, errors.New("wns: only ftp URL scheme is supported")
	}

	if strings.IndexByte(url.Host, ':') == -1 {
		url.Host = url.Host + ":21"
	}
	password, _ := url.User.Password()
	return &FTPPullClient{
		username: url.User.Username(),
		password: password,
		hostname: url.Host,
		baseDir:  strings.TrimSuffix(url.Path, "/"),
	}, nil
}

func missingFiles(lastFile string, filenames []string) []string {
	for i, filename := range filenames {
		if filename == lastFile {
			return filenames[i+1:]
		}
	}
	return filenames
}

// Stream starts a goroutine for continuous odds documents delivery. Stream can
// be stopped by stopping `ctx` context. Parameter `lastFile` is used to skip
// restreaming already processed odds documents. Polling will be performed every
// `interval` time duration. Recommended value for `interval` is one minute.
func (c *FTPPullClient) Stream(ctx context.Context, lastFile string, interval time.Duration) <-chan Data {
	ch := make(chan Data)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		lastFile = c.streamPoll(ch, lastFile)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lastFile = c.streamPoll(ch, lastFile)
			}
		}
	}()
	return ch
}

func (c *FTPPullClient) streamPoll(ch chan<- Data, lastFile string) string {
	files, err := c.List()
	if err != nil {
		ch <- Data{Error: err}
		return lastFile
	}
	missing := missingFiles(lastFile, files)
	if len(missing) == 0 {
		return lastFile
	}
	docs, err := c.Get(missing)
	if err != nil {
		ch <- Data{Error: err}
		return lastFile
	}
	for i, doc := range docs {
		ch <- Data{Data: *doc, Filename: missing[i]}
	}
	return missing[len(missing)-1]
}

// Snapshot returns a slice of all feed documents available on the server sorted
// in chronological order, starting with the oldest document. Second returned
// value is newest odds document file name.
//
// If there are no feed documents on the server this function will return nil
// slice and empty string for last document name.
func (c *FTPPullClient) Snapshot() ([]*BetradarBetData, string, error) {
	files, err := c.List()
	if err != nil {
		return nil, "", err
	}
	if len(files) == 0 {
		return nil, "", nil
	}
	docs, err := c.Get(files)
	if err != nil {
		return nil, "", err
	}
	return docs, files[len(files)-1], nil
}

// List returns a list of odds documents available on the FTP server sorted in
// file creation order.
func (c *FTPPullClient) List() ([]string, error) {
	conn, err := ftp.Dial(c.hostname)
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	if err = conn.Login(c.username, c.password); err != nil {
		return nil, err
	}

	items, err := conn.List(c.baseDir)
	if err != nil {
		return nil, err
	}

	// Sort files in chronological order by server create/modify time, as
	// primary key and file name as secondary key.
	sort.Slice(items, func(i, j int) bool {
		if items[i].Time.Equal(items[j].Time) {
			return items[i].Name < items[j].Name
		}
		return items[i].Time.Before(items[j].Time)
	})

	var files []string
	for _, item := range items {
		if item.Type != ftp.EntryTypeFile {
			continue
		}
		files = append(files, item.Name)
	}
	return files, nil
}

// Get retrieves a batch of odds documents from the FTP server.
func (c *FTPPullClient) Get(filenames []string) ([]*BetradarBetData, error) {
	conn, err := ftp.Dial(c.hostname)
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	if err = conn.Login(c.username, c.password); err != nil {
		return nil, err
	}

	docs := make([]*BetradarBetData, 0, len(filenames))
	for _, filename := range filenames {
		doc, err := c.getFile(conn, filename)
		if err != nil {
			return nil, fmt.Errorf("wns: getting %s file: %s", filename, err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// Remove performs a batch deletion of odds documents.
func (c *FTPPullClient) Remove(filenames []string) error {
	conn, err := ftp.Dial(c.hostname)
	if err != nil {
		return err
	}
	defer conn.Quit()

	if err = conn.Login(c.username, c.password); err != nil {
		return err
	}

	for _, filename := range filenames {
		if err := conn.Delete(fmt.Sprintf("%s/%s", c.baseDir, filename)); err != nil {
			return err
		}
	}
	return nil
}

func (c *FTPPullClient) getFile(conn *ftp.ServerConn, filename string) (*BetradarBetData, error) {
	body, err := conn.Retr(fmt.Sprintf("%s/%s", c.baseDir, filename))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var data BetradarBetData
	if err = xml.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
