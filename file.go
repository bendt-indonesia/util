package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetImageExt(fullPath string) string {
	ext := strings.ToLower(string(fullPath[len(fullPath)-4:]))

	switch ext {
	case ".jpg":
		return "jpg"
	case "jpeg":
		return "jpeg"
	case ".png":
		return "png"
	case ".svg":
		return "svg"
	case ".gif":
		return "gif"
	case "heic":
		return "heic"
	case "webp":
		return "webp"
	case "avif":
		return "avif"
	case "apng":
		return "apng"
	default:
		return "jpg"
	}
}

func CheckDirExists(folderPath string) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.MkdirAll(folderPath, os.ModePerm)
	}
}

func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// For other errors, assume the file doesn't exist
	return false
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func RemoveDir(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}
	return nil
}

func FmtFile(absolutePath string) {
	if strings.HasSuffix(absolutePath, ".go") {
		formatter := "/usr/local/go/bin/gofmt -w " + absolutePath
		cmd := exec.Command("bash", "-c", formatter)
		_, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Execute Failed: " + formatter)
			panic("Unable to gofmt file " + absolutePath)
		}
	}
	return
}

func ReadFile(filePath string) (string, error) {
	// read in the contents of the localfile.data
	data, err := ioutil.ReadFile(filePath)
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		return "", err
	}

	// if it was successful in reading the file then
	// print out the contents as a string
	return string(data), nil
}

func WriteFile(folderPath string, fileName string, content string) error {
	CheckDirExists(folderPath)

	d1 := []byte(content)
	absoluteFilePath := folderPath + fileName
	err := os.WriteFile(absoluteFilePath, d1, 0644)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}
