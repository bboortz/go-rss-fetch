package main

/*
This is a minimal sample application, demonstrating how to set up an RSS feed
for regular polling of new channels/items.
Build & run with:
 $ go run example.go
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/nu7hatch/gouuid"
	"rsslib"
)

func main() {
	// This sets up a new feed and polls it for new channels/items.
	// Invoking it with 'go PollFeed(...)' to have the polling performed in a
	// separate goroutine, so we can poll mutiple feeds.
	//go PollFeed("http://blog.case.edu/news/feed.atom", 5, nil)

	// Poll with a custom charset reader. This is to avoid the following error:
	// ... xml: encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil.
	//PollFeed("https://status.rackspace.com/index/rss", 5, charsetReader)

	go PollFeed("https://www.heise.de/newsticker/heise-top-atom.xml", 5, nil)
	go PollFeed("http://www.spiegel.de/schlagzeilen/tops/index.rss", 5, nil)
	go PollFeed("http://www.faz.net/rss/aktuell/", 5, nil)
	PollFeed("http://www.welt.de/?service=Rss", 5, nil)
}

func PollFeed(uri string, timeout int, cr xmlx.CharsetFunc) {
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
	//	spew.Dump(newchannels)
}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	fmt.Printf("%d new item(s) in %s\n", len(newitems), feed.Url)
	//	spew.Dump(newitems[0])
	var val *rss.Item

	for _, val = range newitems {
		var uuidString string
		var rssitem rsslib.RssItem = rsslib.RssItem{}
		rssitem.Channel = ch.Title
		rssitem.Title = val.Title
		rssitem.Link = val.Links[0].Href
		rssitem.Description = val.Description
		rssitem.PublishDate = val.PubDate
		rssitem.UpdateDate = val.Updated
		if val.Enclosures != nil {
			rssitem.Thumbnail = val.Enclosures[0].Url
		}
		if val.Guid != nil {
			uuidString = *val.Guid
		} else if val.Id != "" {
			uuidString = val.Id
		} else {
			panic("Cannot generate UUID")
		}

		u5, err := uuid.NewV5(uuid.NamespaceURL, []byte(uuidString))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		rssitem.Uuid = u5.String()

		spew.Dump(rssitem)

		requestJson, _ := json.Marshal(rssitem)
		requestBody := string(requestJson)
		req, err := http.NewRequest("PUT", "http://localhost:9090/item", strings.NewReader(requestBody))
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		//		spew.Dump(resp)
		//		fmt.Println(err)

		//		spew.Dump(rssitem)
	}

}

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "iso-8859-1" {
		return r, nil
	}
	return nil, errors.New("Unsupported character set encoding: " + charset)
}
