package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	sdl "github.com/NazarDunets/sitt/schedule"
)

const (
	fileDateFormat    = "2006-01-02"
	fileExtension     = "json"
	storageFolderName = "sitt-storage"
)

func LoadOrCreateSchedule(date time.Time) (*sdl.Schedule, error) {
	filePath, err := getFileForDate(date)
	if err != nil {
		return nil, err
	}

	if s, err := parseSchedule(filePath); err == nil {
		return s, nil
	}

	// file doesn't exists or is corrupted, return empty schedule
	s := sdl.New()
	return &s, nil
}

func Save(date time.Time, tt *sdl.Schedule) error {
	filePath, err := getFileForDate(date)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tt)
	return err
}

func parseSchedule(path string) (*sdl.Schedule, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var s sdl.Schedule

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func getFileForDate(date time.Time) (string, error) {
	fileName := fmt.Sprintf("%s.%s", date.Format(fileDateFormat), fileExtension)
	storageFolder, err := getStorageFolder()
	if err != nil {
		return "", err
	}

	return path.Join(storageFolder, fileName), nil
}

func getStorageFolder() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	storageFolder := path.Join(homedir, storageFolderName)
	err = os.MkdirAll(storageFolder, os.ModePerm)
	if err != nil {
		return "", err
	}

	return storageFolder, nil
}
