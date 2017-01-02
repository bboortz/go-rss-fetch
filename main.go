package main

/*
This is a minimal sample application, demonstrating how to set up an RSS feed
for regular polling of new channels/items.
Build & run with:
 $ go run example.go
*/

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/nu7hatch/gouuid"
	"github.com/op/go-logging"
	"rsslib"
)

var cacheTimeout int = 5
var log = logging.MustGetLogger("rss-fetch")

func main() {
	// This sets up a new feed and polls it for new channels/items.
	// Invoking it with 'go PollFeed(...)' to have the polling performed in a
	// separate goroutine, so we can poll mutiple feeds.
	//go PollFeed("http://blog.case.edu/news/feed.atom", 5, nil)

	// Poll with a custom charset reader. This is to avoid the following error:
	// ... xml: encoding "ISO-8859-1" declared but Decoder.CharsetReader is nil.
	//PollFeed("https://status.rackspace.com/index/rss", 5, charsetReader)

	go PollFeed("https://www.heise.de/newsticker/heise-top-atom.xml", cacheTimeout, nil)
	go PollFeed("http://www.spiegel.de/schlagzeilen/tops/index.rss", cacheTimeout, nil)
	go PollFeed("http://www.faz.net/rss/aktuell/", cacheTimeout, nil)
	PollFeed("http://www.welt.de/?service=Rss", cacheTimeout, nil)
}

func PollFeed(uri string, timeout int, cr xmlx.CharsetFunc) {
	handlers := &MyHandlers{}
	feed := rss.NewWithHandlers(timeout, true, handlers, handlers)
	//	feed := rss.New(timeout, true, chanHandler, itemHandler)

	for {
		log.Infof("%s\tfeed processing: %s\n", log.Module, uri)
		if err := feed.Fetch(uri, cr); err != nil {
			log.Errorf("%s\tERROR: %s - %s\n", log.Module, uri, err)
			return
		}

		<-time.After(time.Duration(60 * time.Second))
	}
}

type MyHandlers struct{}

func (m *MyHandlers) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	log.Infof("%s\tnew channels in %s: %d\n", log.Module, feed.Url, len(newchannels))
}

func (m *MyHandlers) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	log.Infof("%s\tnew items in %s: %d\n", log.Module, feed.Url, len(newitems))

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
			log.Errorf("%s\tERROR: %s - %s\n", log.Module, feed.Url, err)
			return
		}
		rssitem.Uuid = u5.String()

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

	}

}
