package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/kkdai/youtube"
	"github.com/rylio/ytdl"
)

func main() {
	url := os.Args[1]

	// Channels
	chUrls := make(chan string)

	go crawl(url, chUrls)

	for {
		select {
		case temp := <-chUrls:
			go downloadVid(temp)

		case <-time.After(5 * time.Second):
			return
		}
	}
}

func crawl(url string, chUrls chan string) {

	file, err := ioutil.ReadFile(url)
	if err != nil {
		fmt.Printf("couln't open %s, bc %v", url, err)
	}

	clasicPrefix := "https://www.youtube.com/watch?v="

	fileString := string(file)

	lastVidID := ""

	for {
		foundAt := strings.Index(fileString, clasicPrefix)
		if foundAt == -1 {
			return
		}

		fileString = fileString[foundAt+len(clasicPrefix):]

		videoID := fileString[:11]

		if videoID != lastVidID {
			chUrls <- videoID
		}

		lastVidID = videoID
		fileString = fileString[len(clasicPrefix):]

	}
}

func downloadVid(vidID string) {

	vid, err := ytdl.GetVideoInfoFromID(vidID)
	if err != nil {
		fmt.Printf("failed bc %v", err)
	}

	title := strings.Replace(vid.Title, " ", "", -1)
	author := strings.Replace(vid.Author, " ", "", -1)
	nameNoSpaces := title + author

	file, err := os.Create(nameNoSpaces + ".mp3")
	if err != nil {
		fmt.Printf("i couln't open file, because im a dumbfuck and %v\n\n", err)
	}
	defer file.Close()

	y := NewYoutube(true)
	y.DecodeURL("https://www.youtube.com/watch?v=" + vidID)
	currFile, err := filepath.Abs(file.Name())
	if err != nil {
		fmt.Printf("I cant get the absolute path, retard i know, but i have a reason: %v\n\n", err)
	}
	y.StartDownload(currFile)

}
