package rsslib

import ()

type RssItem struct {
	Id          uint64 `json:"Id"`
	Uuid        string `json:"Uuid"`
	Channel     string `json:"Channel"`
	Title       string `json:"Title"`
	Link        string `json:"Link"`
	Description string `json:"Description"`
	Thumbnail   string `json:"Thumbnail"`
	PublishDate string `json:"PublishDate"`
	UpdateDate  string `json:"UpdateDate"`
}

type RssItems []RssItem

func (s RssItems) Len() uint64 {
	return uint64(len(s))
}

func (s1 RssItem) Compare(s2 RssItem) bool {
	if s1.Id != s2.Id {
		return false
	} else if s1.Uuid != s2.Uuid {
		return false
	} else if s1.Channel != s2.Channel {
		return false
	} else if s1.Title != s2.Title {
		return false
	} else if s1.Link != s2.Link {
		return false
	} else if s1.Thumbnail != s2.Thumbnail {
		return false
	} else if s1.PublishDate != s2.PublishDate {
		return false
	} else if s1.UpdateDate != s2.UpdateDate {
		return false
	}

	return true
}

func (s1 RssItem) CompareWithoutPublishDate(s2 RssItem) bool {
	if s1.Id != s2.Id {
		return false
	} else if s1.Uuid != s2.Uuid {
		return false
	} else if s1.Channel != s2.Channel {
		return false
	} else if s1.Title != s2.Title {
		return false
	} else if s1.Link != s2.Link {
		return false
	} else if s1.Thumbnail != s2.Thumbnail {
		return false
	} else if s1.UpdateDate != s2.UpdateDate {
		return false
	}

	return true
}

func (s1 RssItem) CompareWithoutDate(s2 RssItem) bool {
	if s1.Id != s2.Id {
		return false
	} else if s1.Uuid != s2.Uuid {
		return false
	} else if s1.Channel != s2.Channel {
		return false
	} else if s1.Title != s2.Title {
		return false
	} else if s1.Link != s2.Link {
		return false
	} else if s1.Thumbnail != s2.Thumbnail {
		return false
	}

	return true
}

func (s1 RssItem) Diff(s2 RssItem) {
	if s1.Id != s2.Id {
		log.Infof("Id not equal. Expected: %d - Received: %d", s1.Id, s2.Id)
	} else if s1.Uuid != s2.Uuid {
		log.Infof("Uuid not equal. Expected: %s - Received: %s", s1.Uuid, s2.Uuid)
	} else if s1.Channel != s2.Channel {
		log.Infof("Channel not equal. Expected: %s - Received: %s", s1.Channel, s2.Channel)
	} else if s1.Title != s2.Title {
		log.Infof("Title not equal. Expected: %s - Received: %s", s1.Title, s2.Title)
	} else if s1.Link != s2.Link {
		log.Infof("Link not equal. Expected: %s - Received: %s", s1.Link, s2.Link)
	} else if s1.Thumbnail != s2.Thumbnail {
		log.Infof("Thumbnail not equal. Expected: %s - Received: %s", s1.Thumbnail, s2.Thumbnail)
	} else if s1.PublishDate != s2.PublishDate {
		log.Infof("PublishDate not equal. Expected: %s - Received: %s", s1.PublishDate, s2.PublishDate)
	} else if s1.UpdateDate != s2.UpdateDate {
		log.Infof("UpdateDate not equal. Expected: %s - Received: %s", s1.UpdateDate, s2.UpdateDate)
	}
}
