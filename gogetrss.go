package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var rssURL string

// https://golang.org/pkg/encoding/xml/#Unmarshaler
// Unmarshaller interface to support decoding time in xml
type atomTime struct {
	time.Time
}

func (c *atomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	*c = atomTime{parse}
	return nil
}

type Feed struct {
	Items []FeedItem `xml:"entry"`
}

type FeedItem struct {
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Updated atomTime `xml:"updated"`
}

func main() {
	log.Print("Started download")
	client := &http.Client{}

	req, err := http.NewRequest("GET", rssURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %s", err.Error())
	}

	req.Header.Set("User-Agent", "Golang_RSS_Bot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	feed := &Feed{}
	err = xml.Unmarshal(body, feed)
	if err != nil {
		log.Fatalf("Unmarshaling xml problem: %s", err.Error())
	}

	for _, item := range feed.Items {
		log.Printf("Date: %s, Title: %s, Link: %s", item.Updated, item.Title, item.Link)
	}
}

func init() {
	flag.StringVar(&rssURL, "rss-url", "", "atom feed to parse")
	flag.Parse()
}
