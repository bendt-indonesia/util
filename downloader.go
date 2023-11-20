package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//Filename = override downloaded file name
func DownloadFile(URL string) (*http.Response, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}
	return response, nil
}

func DownloadAndSaveFile(URL string, savePath string, saveNameWithExt *string) (string, error) {
	CheckDirExists(savePath)

	_, orgFileName, orgFileExt, _ := GetFileNameWithExtFromUrl(URL)
	var fileName string
	if saveNameWithExt == nil {
		fileName = orgFileName + "-" + RandomTimestampStr() + "." + orgFileExt
	} else {
		fileName = *saveNameWithExt
	}
	savePath += fileName

	//Create a empty file
	file, err := os.Create(savePath)
	if err != nil {
		return fileName, err
	}

	resp, err := http.Get(URL)
	if err != nil {
		return fileName, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		defer file.Close()
		return fileName, fmt.Errorf("Received non 200 response code")
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fileName, err
	}

	defer resp.Body.Close()
	defer file.Close()

	return fileName, nil
}

//FullName, FileName, Ext, Error
func GetFileNameWithExtFromUrl(rawUrl string) (string, string, string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", "", "", err
	}
	pos := strings.LastIndex(u.Path, ".")
	if pos == -1 {
		return "image.jpg", "image", ".jpg", nil
	}
	fileExt := u.Path[pos+1 : len(u.Path)]
	sp := strings.Split(u.Path, "/")
	fileNameWithExt := sp[len(sp)-1]
	return fileNameWithExt, fileNameWithExt[0 : len(fileNameWithExt)-len(fileExt)-1], fileExt, nil
}

func IsValidURL(URL string) bool {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return false
	}

	return true

}
