package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
	"ttit/timetable"
)

const (
	fileDateFormat    = "2006-01-02"
	fileExtension     = "json"
	storageFolderName = "ttit-storage"
)

func LoadOrCreateTimetable(date time.Time) (*timetable.Timetable, error) {
	filePath, err := getFileForDate(date)
	if err != nil {
		return nil, err
	}

	if tt, err := parseTimetable(filePath); err == nil {
		return tt, nil
	}

	// file doesn't exists or is corrupted, return empty timetable
	tt := timetable.New()
	return &tt, nil
}

func Save(date time.Time, tt *timetable.Timetable) error {
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

func parseTimetable(path string) (*timetable.Timetable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tt timetable.Timetable

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tt)
	if err != nil {
		return nil, err
	}

	return &tt, nil
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
