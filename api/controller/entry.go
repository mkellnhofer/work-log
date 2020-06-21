package controller

import (
	"net/http"

	"kellnhofer.com/work-log/service"
)

// EntryController handles requests for entry endpoints.
type EntryController struct {
	eServ *service.EntryService
}

// NewEntryController create a new enty controller.
func NewEntryController(es *service.EntryService) *EntryController {
	return &EntryController{es}
}

// --- Endpoints ---

// GetEntriesHandler returns a handler for "GET /entries".
func (c *EntryController) GetEntriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateEntryHandler returns a handler for "POST /entries".
func (c *EntryController) CreateEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryHandler returns a handler for "GET /entries/{id}".
func (c *EntryController) GetEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateEntryHandler returns a handler for "PUT /entries/{id}".
func (c *EntryController) UpdateEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteEntryHandler returns a handler for "DELETE /entries/{id}".
func (c *EntryController) DeleteEntryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryTypesHandler returns a handler for "GET /entry_types".
func (c *EntryController) GetEntryTypesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetEntryActivitiesHandler returns a handler for "GET /entry_activites".
func (c *EntryController) GetEntryActivitiesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateEntryActivityHandler returns a handler for "POST /entry_activites".
func (c *EntryController) CreateEntryActivityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateEntryActivityHandler returns a handler for "PUT /entry_activites/{id}".
func (c *EntryController) UpdateEntryActivityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteEntryActivityHandler returns a handler for "DELETE /entry_activites/{id}".
func (c *EntryController) DeleteEntryActivityHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
