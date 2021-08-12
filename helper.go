package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

func chanErr(ccc chan int) {
	if ccc != nil {
		ccc <- 1
	}
}

func getFileContentType(buffer []byte) string {
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType
}

func fetchRemoteImage(filepath string, url string) error {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(request)
	// resp, err := http.Get(url)
	log.Info(resp.StatusCode)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode == 200 {
		_ = os.MkdirAll(path.Dir(filepath), 0755)
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		return err
	} else {
		return errors.New("resp retrun not 200")
	}
}
