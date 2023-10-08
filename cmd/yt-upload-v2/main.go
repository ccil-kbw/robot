package main

import (
	"fmt"
	"github.com/ccil-kbw/robot/pkg/yt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	root = "/home/seraf/Documents/MEGAKBW"
)

type File struct {
	Name     string
	FullPath string
	BasePath string
	Uploaded bool
}

func (f *File) flagUploaded() {

	u, err := os.Create(fmt.Sprintf("%s/.%s-uploaded", f.BasePath, f.Name))
	if err != nil {
		fmt.Printf("error flagging file %s as uploaded", f.Name)
	}
	defer func(u *os.File) {
		err := u.Close()
		if err != nil {
			fmt.Println("error closing file, that's life I guess, /shrug")
		}
	}(u)

	_, _ = u.WriteString(time.Now().Format(time.RFC3339))
	_ = u.Sync()
}

type Files struct {
	List []File
	mu   sync.Mutex
}

var (
	files Files
)

func main() {
	go func() {
		for {
			refreshFileList()
			time.Sleep(10 * time.Minute)
		}
	}()

	// Wait until file list is set the first time, 30seconds is overkill but read issues could happen on older disks
	time.Sleep(30 * time.Second)
	for {

		for _, f := range files.List {
			if f.Uploaded {
				fmt.Printf("%s : already uploaded, al hamdulilah\n", f.Name)
				continue
			}

			fmt.Printf("%s : uploading...\n", f.Name)

			err := yt.UploadVideo(yt.Video{
				Title:    f.Name,
				FilePath: f.FullPath,
			})
			if err != nil {
				if strings.Contains(err.Error(), "quotaExceeded") {
					log.Printf("quota exceeded on google youtube api v3, will try %s again later\n", f.Name)
				}
				break
			}

			f.flagUploaded()
		}

		sleep := 2 * time.Hour
		fmt.Printf("next schedule: %s\n", time.Now().Add(sleep))
		time.Sleep(sleep)
	}

}

func refreshFileList() {
	folders, err := os.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	var refreshedFileList []File

	for _, f := range folders {
		fileBasePath := fmt.Sprintf("%s/%s", root, f.Name())
		fs, err := os.ReadDir(fileBasePath)
		if err != nil {
			log.Fatal(err)
		}

		for _, file := range fs {
			fullPath := fmt.Sprintf("%s/%s/%s", root, f.Name(), file.Name())

			if strings.HasSuffix(file.Name(), ".mkv") {
				uploaded := false
				if _, err := os.Stat(fmt.Sprintf("%s/.%s-uploaded", fileBasePath, file.Name())); err == nil {
					uploaded = true
				}
				refreshedFileList = append(refreshedFileList, File{
					Name:     file.Name(),
					BasePath: fileBasePath,
					FullPath: fullPath,
					Uploaded: uploaded,
				})
			}
		}
	}

	files.mu.Lock()
	defer files.mu.Unlock()
	files.List = refreshedFileList
}
