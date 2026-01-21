package onlyoffice

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/royalrick/go-onlyoffice/models"
)

type HistoryVersion struct {
	Version     string          `json:"version"`
	Key         string          `json:"key"`
	Created     time.Time       `json:"created"`
	User        *models.User    `json:"user"`
	ChangesData []models.Change `json:"changes"`
}

func (c *Client) CreateHistory(callback models.Callback, storagePath string) error {
	historyDir := filepath.Join(storagePath, ".history", callback.Key)
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return err
	}

	changesFile := filepath.Join(historyDir, "changes.json")
	data, err := json.MarshalIndent(callback.History, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(changesFile, data, 0644); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetHistory(filename, storagePath string) ([]HistoryVersion, error) {
	historyDir := filepath.Join(storagePath, ".history")
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return []HistoryVersion{}, nil
	}

	files, err := os.ReadDir(historyDir)
	if err != nil {
		return nil, err
	}

	var versions []HistoryVersion
	for _, file := range files {
		if file.IsDir() {
			versionFile := filepath.Join(historyDir, file.Name(), "changes.json")
			if data, err := os.ReadFile(versionFile); err == nil {
				var history models.History
				if err := json.Unmarshal(data, &history); err == nil {
					created, _ := time.Parse("2006-01-02 15:04:05", history.Created)
					version := HistoryVersion{
						Version:     file.Name(),
						Key:         history.Key,
						Created:     created,
						User:        history.User,
						ChangesData: history.Changes,
					}
					versions = append(versions, version)
				}
			}
		}
	}

	return versions, nil
}

func (c *Client) CountVersion(storagePath string) int {
	historyDir := filepath.Join(storagePath, ".history")
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		return 0
	}

	files, err := os.ReadDir(historyDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, file := range files {
		if file.IsDir() {
			count++
		}
	}
	return count
}
