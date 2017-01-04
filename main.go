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
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	//	"fmt"
	rsslib "github.com/bboortz/go-rsslib"
	"github.com/davecgh/go-spew/spew"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/nu7hatch/gouuid"
	"github.com/op/go-logging"
)

var cacheTimeout int = 30
var postItemWorkerInstances int = 10
var log = logging.MustGetLogger("rss-fetch")
var newitemsChan = make(chan rsslib.RssItem, 100)

func main() {
	wg := new(sync.WaitGroup)

	// array of feeds
	feedArr := [...]string{
		"https://blog.fefe.de/rss.xml",
		"https://www.heise.de/newsticker/heise-top-atom.xml",
		"http://www.spiegel.de/schlagzeilen/tops/index.rss",
		"http://www.faz.net/rss/aktuell/",
		"http://www.welt.de/?service=Rss",
		"http://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml",
		"http://rss.golem.de/rss.php?feed=ATOM1.0",
		"http://rss.slashdot.org/Slashdot/slashdotatom",
		"http://newsfeed.zeit.de/index",
		"http://rss.sueddeutsche.de/rss/Digital",
		"http://feeds.reuters.com/reuters/topNews",
		"http://www.tagesschau.de/xml/rss2",
		//"https://news.ycombinator.com/rss",
		//"http://www.infoworld.com/news/index.rss",
	}

	// goroutines for polling feeds
	for _, s := range feedArr {
		wg.Add(1)
		go PollFeed(s, cacheTimeout, charsetReader)
	}

	// goroutines for posting items
	for i := 0; i < postItemWorkerInstances; i++ {
		wg.Add(1)
		go postItemWorker(newitemsChan, wg)
	}

	wg.Wait()
	close(newitemsChan)
}

func postItemWorker(localChan chan rsslib.RssItem, wg *sync.WaitGroup) {
	// Decreasing internal counter for wait-group as soon as goroutine finishes
	defer wg.Done()
	//log.Infof("%s\tpostworker started: %s: %d\n", log.Module, localChan, len(localChan))
	for i := range localChan {
		//	time.Sleep(1 * time.Second)
		//fmt.Printf("Done processing link #%s\n", i.Uuid)
		requestJson, _ := json.Marshal(i)
		requestBody := string(requestJson)
		req, err := http.NewRequest("PUT", "http://go-rss-cache.herokuapp.com/item", strings.NewReader(requestBody))
		//		req, err := http.NewRequest("PUT", "http://localhost:9090/item", strings.NewReader(requestBody))
		if err != nil {
			panic(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	}
	//log.Infof("%s\tpostworker done: %s: %d\n", log.Module, localChan, len(localChan))
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

		<-time.After(time.Duration(10 * 60 * time.Second))
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
			spew.Dump(val)
			panic("Cannot generate UUID: no Id or Guid found.")
		}
		u5, err := uuid.NewV5(uuid.NamespaceURL, []byte(uuidString))
		if err != nil {
			log.Errorf("%s\tERROR: %s - %s\n", log.Module, feed.Url, err)
			return
		}
		rssitem.Uuid = u5.String()

		newitemsChan <- rssitem
	}

}

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	if charset == "ISO-8859-1" || charset == "iso-8859-1" {
		return r, nil
	}
	return nil, errors.New("Unsupported character set encoding: " + charset)
}
