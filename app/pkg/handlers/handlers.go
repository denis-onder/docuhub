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
	length     string
	views      string
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
	var output []Video // Output

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
		output = append(output, Video{
			thumbnail:  thumbnail,
			title:      title,
			url:        url,
			uploadDate: uploadDate.Text(),
			views:      views,
			length:     length})
	})

	return output, nil
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	for _, channel := range channels {
		videos, err := scrapeVideos(channel)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "<p>%s -> %d</p>", channel.name, len(videos))
		for i, video := range videos {
			fmt.Fprintf(w, "<p>%d. %s</p>", i, video.title)
		}
		fmt.Fprintf(w, "\n")
	}
}
