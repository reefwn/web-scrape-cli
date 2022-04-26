package utils

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetDownloadPath() string {
	home, _ := os.UserHomeDir()
	download := filepath.Join(home, "Downloads")

	return download
}

func CreateFolderIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func DownloadFile(url, folderName, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}

	splitUrl := strings.Split(url, ".")
	fileExt := splitUrl[len(splitUrl) - 1]

	// create folder
	downloadFolder := filepath.Join(GetDownloadPath(), folderName)
	err = CreateFolderIfNotExists(downloadFolder)
	if err != nil {
		return err
	}

	// create an empty file
	fileFullName := fileName + "." + fileExt
	file, err := os.Create(filepath.Join(downloadFolder, fileFullName))
	if err != nil {
		return err
	}
	defer file.Close()

	// write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}