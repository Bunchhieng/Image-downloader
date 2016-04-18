// This program make HTTP GET request to Google custom API.
// Then download all images that contain in the JSON response.
// TODO: Make the download image run concurrently i.e go func downloadImage(url string){}
package main

import (
	"net/http"
	"log"
	"encoding/json"
	"os"
	"io"
	"strings"
	"strconv"
	"fmt"
	"time"
)

// Simplify JSON struct to get images in `items` array
type Query struct {
	Items []Item
}

// Right now we only want Title, Image Link and the file extension
type Item struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Mime string `json:"mime"`
}

func main() {
	start := time.Now()
	// Link to the Google custom search API
	URL := "https://www.googleapis.com/customsearch/v1?cx=001106611627702700888%3Aaonktv-oz_w&q=bells%20palsy%20mouth&exactTerms=palsy&fileType=png&imgColorType=color&imgType=face&searchType=image&key=AIzaSyAYqQ4IxUHnF7rfvzSvnczxQ-u93AbkC8k"

	res, err := http.Get(URL)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode == 200 {
		decoder := json.NewDecoder(res.Body)
		var i Query
		err := decoder.Decode(&i)
		if err != nil {
			log.Fatal(err)
		}
		j := 0
		dir, err := os.Getwd()
		for _, i := range i.Items {
			res, err := http.Get(i.Link)
			if err != nil {
				log.Fatal(err)
			}
			os.Mkdir(dir + "/img/", 0777)
			file, err := os.Create(dir + "/img/" + strconv.Itoa(j) + "." + strings.Split(i.Mime, "/")[1])
			j++
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(file, res.Body)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
			res.Body.Close()
		}
		elapsed := time.Since(start).Seconds()
		fmt.Printf("Run time: %f secs\n", elapsed)
		fmt.Println("Sucess! Images located at " + dir + "/img/")
	}
}
