package handlers

/*
	Concept:
	Share a channel which is created via a HTTP request
	Iterate over a list of YouTube channels
	Scrape the data of each video, for each channel, as described in the Video struct
	Pipe *all* videos into the channel
*/

import (
	"errors"
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
	length     string
	views      string
}

// Channel - Struct used to define a YouTube channel
type Channel struct {
	name string
	url  string
}

// Response - Concurrency channel return type
type Response struct {
	videos []Video
	err    error
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

func scrapeVideos(channel Channel, c chan Response) {
	var output []Video // Output

	// Request the HTML page.
	res, err := http.Get(channel.url)
	if err != nil {
		c <- Response{videos: nil, err: err}
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		c <- Response{videos: nil, err: errors.New("source not available")}
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		c <- Response{videos: nil, err: err}
	}

	doc.Find(".yt-lockup-dismissable").Each(func(i int, s *goquery.Selection) {
		// Thumbnail
		thumbnail, _ := s.Find("img").Attr("src")
		// Title
		titleAndURLAnchor := s.Find("a.yt-uix-sessionlink")
		title := titleAndURLAnchor.Text()
		// URL
		param, _ := titleAndURLAnchor.Attr("href")
		url := "https://youtube.com" + param
		// Upload date
		uploadDateAndViewsList := s.Find("ul.yt-lockup-meta-info").Children()
		uploadDate := uploadDateAndViewsList.First()
		// Views
		views := uploadDate.Next().Text()
		// Length
		length := s.Find("span.accessible-description").Text()

		// Form output
		output = append(output, Video{
			thumbnail:  thumbnail,
			title:      title,
			url:        url,
			uploadDate: uploadDate.Text(),
			views:      views,
			length:     length})
	})

	c <- Response{videos: output, err: nil}
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	c := make(chan Response)

	for _, channel := range channels {
		go scrapeVideos(channel, c)
	}

	for i := 0; i < len(channels); i++ {
		res := <-c
		if res.err != nil {
			log.Fatal(res.err)
		}
		for i, video := range res.videos {
			fmt.Fprintf(w, "<p>%d. %s</p>", i, video.title)
		}
	}
}
