package main

import (
	"github.com/ccil-kbw/robot/pkg/yt"
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
	"time"
)

func main() {
	videos := make(chan yt.Video, 5)

	go func(videos chan yt.Video) {
		for {
			// read from channel videos the metadata
			yt.UploadVideo(<-videos)
			// fmt.Printf("got new video %v", <-videos)

		}
	}(videos)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic("something very bad happened. jk we just couldn't initiate the watcher.")
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			panic("something very bad happened. jk we just couldn't close the watcher and the event channel")
		}
	}(watcher)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
				if event.Has(fsnotify.Create) {
					fp := strings.Split(event.Name, "/")
					videos <- yt.Video{
						Title:    fp[len(fp)-1],
						FilePath: event.Name,
					}
				}
				// this is to avoid rate limiting
				time.Sleep(10 * time.Minute)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	if err = watcher.Add("/home/seraf/Documents/MEGASync/Masjid Khalid Ben Walid/archives/to_be_uploaded"); err != nil {
		panic("couldn't add path to watch list")
	}

	<-make(chan struct{})
}
