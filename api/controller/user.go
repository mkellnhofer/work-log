package controller

import (
	"net/http"

	"kellnhofer.com/work-log/api/model"
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

// --- Parameters ---

// swagger:parameters getUser
type GetUserParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters createUser
type CreateUserParameters struct {
	// in: body
	// required: true
	Body model.CreateUserData
}

// swagger:parameters updateUser
type UpdateUserParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
	// in: body
	// required: true
	Body model.UpdateUserData
}

// swagger:parameters deleteUser
type DeleteUserParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters updateUserPassword
type UpdateUserPasswordParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
	// in: body
	// required: true
	Body model.UpdateUserPassword
}

// swagger:parameters getUserRoles
type GetUserRolesParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// swagger:parameters updateUserRoles
type UpdateUserRolesParameters struct {
	// The ID of the user.
	//
	// in: path
	// required: true
	Id int `json:"id"`
	// in: body
	// required: true
	Body model.UpdateUserRoles
}

// --- Responses ---

// The list of users.
// swagger:response GetUsersResponse
type GetUsersResponse struct {
	// in: body
	Body model.UserDataList
}

// The user.
// swagger:response GetUserResponse
type GetUserResponse struct {
	// in: body
	Body model.UserData
}

// The created user.
// swagger:response CreateUserResponse
type CreateUserResponse struct {
	// in: body
	Body model.UserData
}

// The updated user.
// swagger:response UpdateUserResponse
type UpdateUserResponse struct {
	// in: body
	Body model.UserData
}

// The list of user roles.
// swagger:response GetUserRolesResponse
type GetUserRolesResponse struct {
	// in: body
	Body model.UserRoles
}

// The list of updated user roles.
// swagger:response UpdateUserRolesResponse
type UpdateUserRolesResponse struct {
	// in: body
	Body model.UserRoles
}

// --- Endpoints ---

// GetCurrentUserHandler returns a handler for "GET /user".
func (c *UserController) GetCurrentUserHandler() http.HandlerFunc {
	// swagger:operation GET /user user getCurrentUser
	//
	// Get the current user.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetUserResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetCurrentUserRolesHandler returns a handler for "GET /user/roles".
func (c *UserController) GetCurrentUserRolesHandler() http.HandlerFunc {
	// swagger:operation GET /user/roles user getCurrentUserRoles
	//
	// Get roles of the current user.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetUserRolesResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUsersHandler returns a handler for "GET /users".
func (c *UserController) GetUsersHandler() http.HandlerFunc {
	// swagger:operation GET /users users listUsers
	//
	// Lists all users.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetUsersResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// CreateUserHandler returns a handler for "POST /users".
func (c *UserController) CreateUserHandler() http.HandlerFunc {
	// swagger:operation POST /users users createUser
	//
	// Create a user.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/CreateUserResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUserHandler returns a handler for "GET /users/{id}".
func (c *UserController) GetUserHandler() http.HandlerFunc {
	// swagger:operation GET /users/{id} users getUser
	//
	// Get a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetUserResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserHandler returns a handler for "PUT /users/{id}".
func (c *UserController) UpdateUserHandler() http.HandlerFunc {
	// swagger:operation PUT /users/{id} users updateUser
	//
	// Update a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/UpdateUserResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   '409':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// DeleteUserHandler returns a handler for "DELETE /users/{id}".
func (c *UserController) DeleteUserHandler() http.HandlerFunc {
	// swagger:operation DELETE /users/{id} users deleteUser
	//
	// Delete a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//     description: No content.
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserPasswordHandler returns a handler for "PUT /users/{id}/password".
func (c *UserController) UpdateUserPasswordHandler() http.HandlerFunc {
	// swagger:operation PUT /users/{id}/password users updateUserPassword
	//
	// Update the password of a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//     description: No content.
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// GetUserRolesHandler returns a handler for "GET /users/{id}/roles".
func (c *UserController) GetUserRolesHandler() http.HandlerFunc {
	// swagger:operation GET /users/{id}/roles users getUserRoles
	//
	// Get roles of a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetUserRolesResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}

// UpdateUserRolesHandler returns a handler for "PUT /users/{id}/roles".
func (c *UserController) UpdateUserRolesHandler() http.HandlerFunc {
	// swagger:operation PUT /users/{id}/roles users updateUserRoles
	//
	// Update roles of a user by its ID.
	//
	// ---
	//
	// security:
	// - basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/UpdateUserRolesResponse"
	//   '400':
	//     "$ref": "#/responses/ErrorResponse"
	//   '401':
	//     "$ref": "#/responses/ErrorResponse"
	//   '403':
	//     "$ref": "#/responses/ErrorResponse"
	//   '404':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
