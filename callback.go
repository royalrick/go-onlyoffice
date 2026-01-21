package onlyoffice

import (
	"encoding/json"
	"errors"

	"github.com/royalrick/go-onlyoffice/models"
)

func (c *Client) ParseCallback(jsonData []byte, tokenHeader string) (*models.Callback, error) {
	var callback models.Callback

	if err := json.Unmarshal(jsonData, &callback); err != nil {
		return nil, err
	}

	if c.config.JWTEnabled {
		if tokenHeader == "" {
			return nil, errors.New("missing token")
		}

		if len(tokenHeader) > 7 && tokenHeader[:7] == "Bearer " {
			tokenHeader = tokenHeader[7:]
		}

		claims, err := c.ParseToken(tokenHeader)
		if err != nil {
			return nil, err
		}

		if claims["status"] != nil {
			if status, ok := claims["status"].(float64); ok {
				callback.Status = int(status)
			}
		}

		if claims["key"] != nil {
			if key, ok := claims["key"].(string); ok {
				callback.Key = key
			}
		}
	}

	callback.Token = tokenHeader

	return &callback, nil
}

func (c *Client) ValidateCallback(callback *models.Callback) error {
	switch callback.Status {
	case 2, 6:
		return nil
	case 7:
		return errors.New("document is corrupted")
	default:
		return errors.New("unknown callback status")
	}
}

func (c *Client) GetDownloadURL(callback *models.Callback) (string, error) {
	if callback.Url == "" {
		return "", errors.New("empty download url")
	}
	return callback.Url, nil
}

func (c *Client) GenerateHistoryKey(filename string) (string, error) {
	return c.GenerateFileHash(filename)
}
