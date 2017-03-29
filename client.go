package wns

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

// WNS is a World Number Service client
type WNS struct {
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
func (w *WNS) Clear(ctx context.Context) error {
	url := fmt.Sprintf("%s?bookmakerName=%s&key=%s&xmlFeedName=FileGet&deleteFullQueue=yes", w.URL, w.Username, w.Key)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := w.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

type Data struct {
	Data  BetradarBetData
	Error error
}

// Stream streams all updates to a returned channel. Under the hood it uses
// Get method on WNS with delete set to `true`
func (w *WNS) Stream(ctx context.Context) chan Data {
	ch := make(chan Data)
	interval := w.Interval
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
				d := w.Get(ctx, true)
				if d.Error != nil {
					if wnserror, ok := d.Error.(*WNSError); ok {
						// if there are no new files, we skip streaming the
						// struct with such error
						if wnserror.Type == noNew {
							continue
						}
					}
				}
				ch <- d
			}
		}
	}()
	return ch
}

// Get gets avialable lottery feed document
// if delete flag is set to false, it does not delete the document from source,
// so next Get request will return the same document until it is Get'ed with
// delete flag set to true, or queue is cleared.
func (w *WNS) Get(ctx context.Context, delete bool) Data {
	keyword := "no"
	if delete {
		keyword = "yes"
	}
	url := fmt.Sprintf("%s?bookmakerName=%s&key=%s&xmlFeedName=FileGet&deleteAfterTransfer=%s", w.URL, w.Username, w.Key, keyword)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Data{Error: err}
	}
	resp, err := w.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return Data{Error: err}
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	if err := checkForErr(tee); err != nil {
		return Data{Error: err}
	}
	var data BetradarBetData
	d := xml.NewDecoder(&buf)
	d.CharsetReader = charsetReader
	err = d.Decode(&data)
	return Data{Data: data, Error: err}
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
			return &WNSError{
				Type: tooFrequent,
				Err:  betradarErrors.Val,
			}
		}
		if betradarErrors.Val == "There are no files ready for transfer at the moment." {
			return &WNSError{
				Type: noNew,
				Err:  betradarErrors.Val,
			}
		}
		return &WNSError{
			Type: unknown,
			Err:  betradarErrors.Val,
		}
	}
	return nil
}

type typeError int

const (
	noNew typeError = iota
	tooFrequent
	unknown
)

// WNSError is an error wrapper to distinguish between known and unknown
// feed errors, so additional logic could be done on errors that are
// actually only informational messages concealed in error format
type WNSError struct {
	Type typeError
	Err  string
}

func (e *WNSError) Error() string {
	var t string
	switch e.Type {
	case noNew:
		t = "No new files available"
	case tooFrequent:
		t = "Too frequent requests"
	}
	return fmt.Sprintf("%s (%s)", t, e.Err)
}

func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	switch charset {
	case "ISO-8859-1":
		return charmap.ISO8859_1.NewDecoder().Reader(input), nil
	default:
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
}
