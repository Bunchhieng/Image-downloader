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
	"fmt"
	"time"
)

// Simplify JSON struct to get images in `items` array
type Query struct {
	Items []Item
}

// Right now we only want Title, Image Link and the file extension
type Item struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	DisplayLink string `json:"displayLink"`
	Mime        string `json:"mime"`
}

func main() {
	count := 0
	start := time.Now()
	// Link to the Google custom search API
	// Each request return 10 images.
	// To get to the next page, we need to change startindex=i + 10
	URL := "https://www.googleapis.com/customsearch/v1?start=%d&cx=001106611627702700888%%3Aaonktv-oz_w&q=bells%%20palsy%%20mouth&exactTerms=palsy&fileType=png&imgColorType=color&imgType=face&searchType=image&key=AIzaSyAYqQ4IxUHnF7rfvzSvnczxQ-u93AbkC8k"
	for v := 1; v < 100; v += 10 {
		val := fmt.Sprintf(URL, v)
		res, err := http.Get(val)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		if res.StatusCode == 200 {
			decoder := json.NewDecoder(res.Body)
			var i Query
			err := decoder.Decode(&i)
			if err != nil {
				log.Fatal(err)
			}
			dir, err := os.Getwd()
			for _, i := range i.Items {
				imgRes := downloadImg(i.Link)
				os.Mkdir(dir + "/img/", 0777)
				file, err := os.Create(dir + "/img/" + i.DisplayLink + "." + strings.Split(i.Mime, "/")[1])
				count++
				if err != nil {
					log.Fatal(err)
				}
				_, err = io.Copy(file, imgRes.Body)
				if err != nil {
					log.Fatal(err)
				}
				file.Close()
			}
		}
	}
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Run time: %f secs\n", elapsed)
	fmt.Printf("Sucess! %d images downloaded", count)
}

func downloadImg(url string) (res *http.Response) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
