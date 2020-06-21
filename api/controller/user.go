package controller

import (
	"net/http"

	"kellnhofer.com/work-log/service"
)

// UserController handles requests for user endpoints.
type UserController struct {
	uServ *service.UserService
}

// NewUserController create a new user controller.
func NewUserController(us *service.UserService) *UserController {
	return &UserController{us}
}

// --- Endpoints ---

// GetCurrentUserHandler returns a handler for "GET /user".
func (c *UserController) GetCurrentUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetCurrentUserRolesHandler returns a handler for "GET /user/roles".
func (c *UserController) GetCurrentUserRolesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUsersHandler returns a handler for "GET /users".
func (c *UserController) GetUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateUserHandler returns a handler for "POST /users".
func (c *UserController) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUserHandler returns a handler for "GET /users/{id}".
func (c *UserController) GetUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserHandler returns a handler for "PUT /users/{id}".
func (c *UserController) UpdateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteUserHandler returns a handler for "DELETE /users/{id}".
func (c *UserController) DeleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserPasswordHandler returns a handler for "PUT /users/{id}/password".
func (c *UserController) UpdateUserPasswordHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUserRolesHandler returns a handler for "GET /users/{id}/roles".
func (c *UserController) GetUserRolesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserRolesHandler returns a handler for "PUT /users/{id}/roles".
func (c *UserController) UpdateUserRolesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
