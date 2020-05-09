package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
)

// Video - Struct used to define a video
type Video struct {
	Title      string `json:"Title"`
	Thumbnail  string `json:"Thumbnail"`
	UploadDate string `json:"UploadDate"`
	URL        string `json:"URL"`
	Length     string `json:"Length"`
	Views      string `json:"Views"`
	Author     string `json:"Author"`
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

func shuffle(videos []Video) {
	// Create a truly random RNG
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// Swap indices
	for i := range videos {
		newPos := r.Intn(len(videos) - 1)
		videos[i], videos[newPos] = videos[newPos], videos[i] // Swap syntax
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
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
		titleText := strings.ReplaceAll(titleAndURLAnchor.Text(), "\n", "")
		title := strings.Trim(titleText, " ")
		// URL
		param, _ := titleAndURLAnchor.Attr("href")
		url := "https://youtube.com" + param
		// Upload date
		uploadDateAndViewsList := s.Find("ul.yt-lockup-meta-info").Children()
		views := uploadDateAndViewsList.First()
		// Views
		uploadDate := views.Next().Text()
		// Length
		length := s.Find("span.accessible-description").Text()

		// Form output
		output = append(output, Video{
			Thumbnail:  thumbnail,
			Title:      title,
			URL:        url,
			UploadDate: uploadDate,
			Views:      views.Text(),
			Length:     strings.Replace(length, " - ", "", 1),
			Author:     channel.name})
	})

	c <- Response{videos: output, err: nil}
}

// GetDocumentaries gets all docs from the listed channels
func GetDocumentaries(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	c := make(chan Response)
	defer close(c)

	videos := []Video{}

	for _, channel := range channels {
		go scrapeVideos(channel, c)
	}

	for i := 0; i < len(channels); i++ {
		res := <-c
		if res.err != nil {
			json.NewEncoder(w).Encode(res.err)
		} else {
			for _, video := range res.videos {
				videos = append(videos, video)
			}
		}
	}

	shuffle(videos)

	json.NewEncoder(w).Encode(videos)
}

// Router
func createRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/all", GetDocumentaries)

	return router
}

// Start the web server
func Start(port int) {
	router := createRouter()

	// Server init
	p := ":" + strconv.Itoa(port)
	fmt.Printf("Server running!\nhttp://localhost%s/\n", p)
	http.ListenAndServe(p, router)
}

func main() {
	var port int
	p, exists := os.LookupEnv("PORT")
	if !exists {
		port = 8000
	} else {
		parsed, _ := strconv.Atoi(p)
		port = parsed
	}
	Start(port)
}
