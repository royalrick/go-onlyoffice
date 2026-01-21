package onlyoffice

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/royalrick/go-onlyoffice/models"
)

// CallbackHandlerFunc is a function that handles OnlyOffice callbacks
type CallbackHandlerFunc func(cb *models.Callback) error

// CallbackHandlers defines handlers for different callback statuses
type CallbackHandlers struct {
	OnEditing   CallbackHandlerFunc // status 1: document is being edited
	OnSave      CallbackHandlerFunc // status 2: document is ready for saving
	OnSaveError CallbackHandlerFunc // status 3: document saving error has occurred
	OnClose     CallbackHandlerFunc // status 4: document is closed with no changes
	OnForceSave CallbackHandlerFunc // status 6: document is being forcibly saved
	OnCorrupt   CallbackHandlerFunc // status 7: document is corrupted
}

// callbackHandler implements http.Handler for OnlyOffice callbacks
type callbackHandler struct {
	client   *Client
	handlers CallbackHandlers
}

// CallbackHandler returns an http.Handler that processes OnlyOffice callbacks
func (c *Client) CallbackHandler(handlers CallbackHandlers) http.Handler {
	return &callbackHandler{client: c, handlers: handlers}
}

// ServeHTTP implements the http.Handler interface
func (h *callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Only accept POST requests
	if r.Method != http.MethodPost {
		h.respondError(w, http.StatusMethodNotAllowed)
		return
	}

	// 2. Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.respondError(w, http.StatusBadRequest)
		return
	}

	// 3. Parse callback (includes JWT validation)
	callback, err := h.client.ParseCallback(body, r.Header.Get("Authorization"))
	if err != nil {
		h.respondError(w, http.StatusUnauthorized)
		return
	}

	// 4. Dispatch based on status
	var handler CallbackHandlerFunc
	switch callback.Status {
	case 1:
		handler = h.handlers.OnEditing
	case 2:
		handler = h.handlers.OnSave
	case 3:
		handler = h.handlers.OnSaveError
	case 4:
		handler = h.handlers.OnClose
	case 6:
		handler = h.handlers.OnForceSave
	case 7:
		handler = h.handlers.OnCorrupt
	}

	// 5. Execute handler function
	if handler != nil {
		if err := handler(callback); err != nil {
			h.respondError(w, http.StatusInternalServerError)
			return
		}
	}

	// 6. Return success response
	h.respondOK(w)
}

// respondOK sends a successful response to OnlyOffice
func (h *callbackHandler) respondOK(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"error": 0})
}

// respondError sends an error response to OnlyOffice
func (h *callbackHandler) respondError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]int{"error": 1})
}
