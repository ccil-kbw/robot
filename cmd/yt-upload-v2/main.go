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
	root = "/data"
)

type File struct {
	Name     string
	Size     int64
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
		log.Printf("Error reading directory %s: %v", root, err)
		return
	}

	var refreshedFileList []File

	for _, file := range folders {
		fullPath := fmt.Sprintf("%s/%s", root, file.Name())

		if strings.HasSuffix(file.Name(), ".mp4") {
			uploaded := false
			if _, err := os.Stat(fmt.Sprintf("%s/.%s-uploaded", root, file.Name())); err == nil {
				uploaded = true
			}

			fileInfo, err := file.Info()
			if err != nil {
				continue
			}

			// If size is smaller than 5 MegaByte, skip for now (is 0 if still not uploaded; skipping bad videos too)
			if fileInfo.Size() < 5<<20 {
				fmt.Printf("%s video too small, don't think should upload (%d bytes)\n", file.Name(), fileInfo.Size())
				continue
			}

			refreshedFileList = append(refreshedFileList, File{
				Name:     file.Name(),
				Size:     fileInfo.Size(),
				BasePath: root,
				FullPath: fullPath,
				Uploaded: uploaded,
			})
		}
	}

	files.mu.Lock()
	defer files.mu.Unlock()
	files.List = refreshedFileList
}
