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

func scrapeVideos(channel Channel) ([]Video, error) {
	output := []Video{}
	// Request the HTML page.
	res, err := http.Get(channel.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".ytd-grid-video-renderer").Each(func(i int, s *goquery.Selection) {
		output = append(output, Video{thumbnail: "test", title: "test", uploadDate: "test", url: "test", views: 123})
	})
	return output, nil
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	// for _, channel := range channels {
	channel := channels[0]
	videos, err := scrapeVideos(channel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s -> %d\n", channel.name, len(videos))
	// }
}
