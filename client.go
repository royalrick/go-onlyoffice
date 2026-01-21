package onlyoffice

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/royalrick/go-onlyoffice/models"
)

type Config struct {
	DocumentServerURL string
	JWTSecret         string
	JWTEnabled        bool
	HTTPClient        *http.Client
}

type Client struct {
	config *Config
	http   *http.Client
}

func NewClient(cfg *Config) (*Client, error) {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	return &Client{
		config: cfg,
		http:   cfg.HTTPClient,
	}, nil
}

func (c *Client) GenerateFileHash(filename string) (string, error) {
	h := sha1.New()
	_, err := h.Write([]byte(filename + time.Now().String()))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (c *Client) CreateToken(claims jwt.Claims) (string, error) {
	if !c.config.JWTEnabled || c.config.JWTSecret == "" {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.config.JWTSecret))
}

func (c *Client) ParseToken(tokenString string) (jwt.MapClaims, error) {
	if !c.config.JWTEnabled {
		return nil, nil
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (c *Client) BuildEditorConfig(params models.EditorParams, fileURL string) (*models.Config, error) {
	if params.Filename == "" {
		return nil, errors.New("filename is required")
	}

	ext := strings.TrimPrefix(strings.ToLower(params.Filename), ".")
	if idx := strings.LastIndex(ext, "."); idx > 0 {
		ext = ext[idx+1:]
	}

	fileKey, err := c.GenerateFileHash(params.Filename)
	if err != nil {
		return nil, err
	}

	cfg := &models.Config{
		Type:         "desktop",
		DocumentType: c.getDocumentType(ext),
		Document: models.Document{
			FileType: ext,
			Key:      fileKey,
			Title:    params.Filename,
			Url:      fileURL,
			Info: models.MetaInfo{
				Author:  params.UserId,
				Created: time.Now().Format("2006-01-02 15:04:05"),
			},
			Permissions: models.Permissions{
				Chat:      true,
				Download:  params.CanDownload,
				Edit:      params.CanEdit,
				FillForms: true,
				Print:     true,
			},
			ReferenceData: models.ReferenceData{
				FileKey: fileKey,
			},
		},
		EditorConfig: models.EditorConfig{
			User: models.UserInfo{
				Id:    params.UserId,
				Name:  params.UserName,
				Email: params.UserEmail,
			},
			CallbackUrl: params.CallbackUrl,
			Lang:        params.Language,
			Mode:        params.Mode,
			Customization: models.Customization{
				About:    true,
				Feedback: true,
			},
		},
	}

	if c.config.JWTEnabled {
		claims := jwt.MapClaims{
			"document": map[string]any{
				"key":      fileKey,
				"url":      fileURL,
				"fileType": ext,
			},
			"editorConfig": map[string]any{
				"user": map[string]any{
					"id": params.UserId,
				},
			},
			"exp": time.Now().Add(5 * time.Minute).Unix(),
		}

		token, err := c.CreateToken(claims)
		if err != nil {
			return nil, err
		}
		cfg.Token = token
	}

	return cfg, nil
}

func (c *Client) getDocumentType(ext string) string {
	switch strings.ToLower(ext) {
	case "docx", "doc", "odt", "rtf", "txt", "html", "htm", "mht", "pdf":
		return "text"
	case "xlsx", "xls", "ods", "csv":
		return "spreadsheet"
	case "pptx", "ppt", "odp":
		return "presentation"
	default:
		return "word"
	}
}
