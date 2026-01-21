package onlyoffice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ConvertOptions struct {
	DocumentURL string
	FromExt     string
	ToExt       string
	DocumentKey string
	Async       bool
	Title       string
}

type ConvertResult struct {
	FileURL  string `json:"fileUrl"`
	FileType string `json:"fileType"`
	Percent  int    `json:"percent"`
	IsEnd    bool   `json:"endConvert"`
	Error    int    `json:"error"`
	Key      string `json:"key"`
}

func (c *Client) ConvertDocument(opts ConvertOptions) (*ConvertResult, error) {
	if opts.FromExt == "" {
		opts.FromExt = getExtension(opts.DocumentURL)
	}

	convertURL := fmt.Sprintf("%s/ConvertService.ashx", c.config.DocumentServerURL)

	payload := map[string]any{
		"url":         opts.DocumentURL,
		"outputtype":  opts.ToExt,
		"filetype":    opts.FromExt,
		"title":       opts.Title,
		"key":         opts.DocumentKey,
		"async":       opts.Async,
		"region":      "en",
		"embedded":    false,
		"canDownload": true,
	}

	if c.config.JWTEnabled {
		token, err := c.CreateToken(jwt.MapClaims(payload))
		if err != nil {
			return nil, err
		}
		payload["token"] = token
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", convertURL, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("conversion failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ConvertResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) CanConvert(ext string) bool {
	convertibleExts := map[string]bool{
		"doc": true, "docx": true, "odt": true, "rtf": true,
		"xls": true, "xlsx": true, "ods": true, "csv": true,
		"ppt": true, "pptx": true, "odp": true,
		"pdf": true, "txt": true, "html": true, "htm": true,
	}
	return convertibleExts[strings.ToLower(ext)]
}

func (c *Client) GetInternalExtension(ext string) string {
	switch strings.ToLower(ext) {
	case "doc", "odt", "rtf":
		return "docx"
	case "xls", "ods", "csv":
		return "xlsx"
	case "ppt", "odp":
		return "pptx"
	default:
		return ext
	}
}

func getExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

func (c *Client) DownloadFile(fileURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
