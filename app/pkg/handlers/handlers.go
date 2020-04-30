package handlers

/*
	Concept:
	Share a channel which is created via a HTTP request
	Iterate over a list of YouTube channels
	Scrape the data of each video, for each channel, as described in the Video struct
	Pipe *all* videos into the channel
*/

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Video - Struct used to define a video
type Video struct {
	title      string
	thumbnail  string
	uploadDate string
	url        string
	views      int
}

// Channel - Struct used to define a YouTube channel
type Channel struct {
	name string
	url  string
}

var channels []Channel = []Channel{
	{name: "National Geographic", url: "http://www.youtube.com/user/NationalGeographic/videos"},
	{name: "Journeyman Pictures", url: "http://www.youtube.com/user/journeymanpictures/videos"},
	{name: "DW Documentary", url: "http://www.youtube.com/channel/UCW39zufHfsuGgpLviKh297Q/videos"},
	{name: "RT Documentary", url: "http://www.youtube.com/user/RTDocumentaries/videos"},
	{name: "Free Documentary", url: "http://www.youtube.com/user/FreeDocumentary/videos"},
	{name: "Documentary Tube", url: "http://www.youtube.com/user/TheDocumenteriesTube/videos"},
	{name: "Timeline - World History Documentaries", url: "http://www.youtube.com/channel/UC88lvyJe7aHZmcvzvubDFRg/videos"},
}

func scrapeVideos(channel Channel, c chan string, w http.ResponseWriter) {
	// Request the HTML page.
	res, err := http.Get(channel.url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	sel := doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, "/watch") {
			c <- href
		}
		// for i := 0; i < len(s.Nodes); i++ {
		// 	fmt.Fprintf(w, <-c+"\n")
		// }
	})

	fmt.Println(len(sel.Nodes))
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	c := make(chan string)
	for _, channel := range channels {
		go scrapeVideos(channel, c, w)
	}
	for {
		fmt.Fprintf(w, "<p>%s</p>\n", <-c)
	}
}
