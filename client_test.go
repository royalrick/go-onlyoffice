package onlyoffice_test

import (
	"testing"

	"github.com/royalrick/go-onlyoffice"
	"github.com/royalrick/go-onlyoffice/models"
)

func TestNewClient(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
		JWTEnabled:        false,
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client == nil {
		t.Fatal("Client should not be nil")
	}
}

func TestBuildEditorConfig(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
		JWTEnabled:        false,
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	params := models.EditorParams{
		Filename:    "test.docx",
		Mode:        "edit",
		Language:    "en",
		UserId:      "user1",
		UserName:    "Test User",
		UserEmail:   "test@example.com",
		CallbackUrl: "https://example.com/callback",
		CanEdit:     true,
		CanDownload: true,
	}

	cfg, err := client.BuildEditorConfig(params, "https://example.com/storage/test.docx")
	if err != nil {
		t.Fatalf("Failed to build editor config: %v", err)
	}

	if cfg.Document.Title != "test.docx" {
		t.Errorf("Expected title 'test.docx', got '%s'", cfg.Document.Title)
	}

	if cfg.EditorConfig.User.Id != "user1" {
		t.Errorf("Expected user ID 'user1', got '%s'", cfg.EditorConfig.User.Id)
	}
}

func TestCanConvert(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		ext      string
		expected bool
	}{
		{"docx", true},
		{"xlsx", true},
		{"pptx", true},
		{"pdf", true},
		{"txt", true},
		{"xyz", false},
	}

	for _, tt := range tests {
		if got := client.CanConvert(tt.ext); got != tt.expected {
			t.Errorf("CanConvert(%s) = %v, expected %v", tt.ext, got, tt.expected)
		}
	}
}

func TestGetInternalExtension(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"doc", "docx"},
		{"xls", "xlsx"},
		{"ppt", "pptx"},
		{"docx", "docx"},
		{"xlsx", "xlsx"},
	}

	for _, tt := range tests {
		if got := client.GetInternalExtension(tt.input); got != tt.expected {
			t.Errorf("GetInternalExtension(%s) = %v, expected %v", tt.input, got, tt.expected)
		}
	}
}

func TestParseCallback(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
		JWTEnabled:        false,
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	jsonData := []byte(`{
		"status": 2,
		"key": "test-key",
		"url": "https://example.com/file.docx",
		"changesurl": "https://example.com/changes.zip"
	}`)

	callback, err := client.ParseCallback(jsonData, "")
	if err != nil {
		t.Fatalf("Failed to parse callback: %v", err)
	}

	if callback.Status != 2 {
		t.Errorf("Expected status 2, got %d", callback.Status)
	}

	if callback.Key != "test-key" {
		t.Errorf("Expected key 'test-key', got '%s'", callback.Key)
	}
}

func TestValidateCallback(t *testing.T) {
	config := &onlyoffice.Config{
		DocumentServerURL: "https://example.com",
	}

	client, err := onlyoffice.NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"Edited", 2, false},
		{"MustSave", 6, false},
		{"Corrupted", 7, true},
		{"Unknown", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			callback := &models.Callback{Status: tt.status}
			err := client.ValidateCallback(callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCallback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
