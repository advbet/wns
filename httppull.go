package wns

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

var debug = os.Getenv("WNS_DEBUG") != ""

// HTTPPullClient is Betradar WNS feed client for feed consumption with
// HTTP-pull delivery method.
type HTTPPullClient struct {
	Username   string      // Betradar given Bookmaker Name
	Key        string      // API Key
	URL        string      // URL of feed source
	HTTPClient http.Client // If not specified, http.DefaultClient will be used
	// Interval is a duration of how frequently it should be retried to fetch
	// the feed. As the WNS does a rate limit for 1 request per 10 seconds,
	// it is advisable to not set Interval to less then 10s.
	Interval time.Duration
}

// Clear clears the queue (removes all queued old lottery feed files)
func (c *HTTPPullClient) Clear(ctx context.Context) error {
	url := fmt.Sprintf("%s?bookmakerName=%s&key=%s&xmlFeedName=FileGet&deleteFullQueue=yes", c.URL, c.Username, c.Key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// Data wraps BetradarBetData and error to be returned via streamed channel
type Data struct {
	Data     BetradarBetData
	Filename string
	Error    error
}

// Stream streams all updates to a returned channel. Under the hood it uses
// Get method on WNS with delete set to `true`
func (c *HTTPPullClient) Stream(ctx context.Context) <-chan Data {
	ch := make(chan Data)
	interval := c.Interval
	if interval == 0 {
		interval = 10 * time.Second
	}
	go func() {
		defer close(ch)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				d, err := c.Get(ctx, true)
				if wnserror, ok := err.(*APIError); ok {
					// if there are no new files, we skip streaming the
					// struct with such error
					if wnserror.Type == ErrTypeNoNew {
						continue
					}
				}
				ch <- Data{Data: d, Error: err}
			}
		}
	}()
	return ch
}

// Get gets avialable lottery feed document
// if delete flag is set to false, it does not delete the document from source,
// so next Get request will return the same document until it is Get'ed with
// delete flag set to true, or queue is cleared.
func (c *HTTPPullClient) Get(ctx context.Context, delete bool) (BetradarBetData, error) {
	keyword := "no"
	if delete {
		keyword = "yes"
	}
	url := fmt.Sprintf("%s?bookmakerName=%s&key=%s&xmlFeedName=FileGet&deleteAfterTransfer=%s", c.URL, c.Username, c.Key, keyword)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return BetradarBetData{}, err
	}
	resp, err := c.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return BetradarBetData{}, err
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return BetradarBetData{}, err
	}

	if debug {
		_ = ioutil.WriteFile(fmt.Sprintf("wns-debug-%d.xml", time.Now().Unix()), bs, 0644)
	}

	if err = checkForErr(bytes.NewBuffer(bs)); err != nil {
		return BetradarBetData{}, err
	}
	var data BetradarBetData
	d := xml.NewDecoder(bytes.NewBuffer(bs))
	d.CharsetReader = charsetReader
	if err = d.Decode(&data); err != nil {
		return BetradarBetData{}, err
	}
	return data, nil
}

func checkForErr(r io.Reader) error {
	var betradarErrors struct {
		XMLName xml.Name
		Val     string `xml:",innerxml"`
	}
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	err := d.Decode(&betradarErrors)
	if err != nil {
		return err
	}
	local := betradarErrors.XMLName.Local
	if local == "error" || local == "error-message" {
		if strings.HasPrefix(betradarErrors.Val, "Too frequent download") {
			return &APIError{
				Type: ErrTypeTooFrequent,
				Err:  betradarErrors.Val,
			}
		}
		if betradarErrors.Val == "There are no files ready for transfer at the moment." {
			return &APIError{
				Type: ErrTypeNoNew,
				Err:  betradarErrors.Val,
			}
		}
		return &APIError{
			Type: ErrTypeUnknown,
			Err:  betradarErrors.Val,
		}
	}
	return nil
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "ISO-8859-1":
		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
}
