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
	"net/http"
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
	{name: "National Geographic", url: "https://www.youtube.com/user/NationalGeographic/videos"},
	{name: "Journeyman Pictures", url: "https://www.youtube.com/user/journeymanpictures/videos"},
	{name: "DW Documentary", url: "https://www.youtube.com/channel/UCW39zufHfsuGgpLviKh297Q/videos"},
	{name: "RT Documentary", url: "https://www.youtube.com/user/RTDocumentaries/videos"},
	{name: "Free Documentary", url: "https://www.youtube.com/user/FreeDocumentary/videos"},
	{name: "Documentary Tube", url: "https://www.youtube.com/user/TheDocumenteriesTube/videos"},
	{name: "Timeline - World History Documentaries", url: "https://www.youtube.com/channel/UC88lvyJe7aHZmcvzvubDFRg/videos"},
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	// if r.Method == "POST" {
	fmt.Fprintf(w, "Channels: %d", len(channels))
	// }
}
