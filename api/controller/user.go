package controller

import (
	"fmt"
	"net/http"

	"kellnhofer.com/work-log/api/mapper"
	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/api/validator"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
	"kellnhofer.com/work-log/service"
	httputil "kellnhofer.com/work-log/util/http"
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute action
		user, err := c.uServ.GetCurrentUserData(r.Context())
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		httputil.WriteHttpResponse(w, http.StatusOK, au)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute action
		userRoles, err := c.uServ.GetCurrentUserRoles(r.Context())
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aur := mapper.ToRoles(userRoles)
		httputil.WriteHttpResponse(w, http.StatusOK, aur)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute action
		users, err := c.uServ.GetUserDatas(r.Context())
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aus := mapper.ToUserDatas(users)
		httputil.WriteHttpResponse(w, http.StatusOK, aus)
	}
}

// CreateUserHandler returns a handler for "POST /users".
func (c *UserController) CreateUserHandler() http.HandlerFunc {
	// swagger:operation POST /users users createUser
	//
	// Create a user.
	//
	// # Username / password rules
	//
	// __Username:__
	//
	// ⦁ Minimum length: 4
	// ⦁ Maximum length: 100
	// ⦁ Allowed characters: `0-9 a-z A-Z - .`
	//
	// __Password:__
	//
	// ⦁ Minimum length: 8
	// ⦁ Maximum length: 100
	// ⦁ Allowed characters: `0-9 a-z A-Z ! \ # $ % & ' ( ) * + , - . / : ; = ? @ [ \ ] ^ _ { | } ~`
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Read API model from request
		var acu model.CreateUserData
		httputil.ReadHttpBody(r, &acu)

		// Validate model
		if err := validator.ValidateCreateUser(&acu); err != nil {
			panic(err)
		}

		// Convert to logic model
		user := mapper.FromCreateUserData(&acu)

		// Execute action
		if err := c.uServ.CreateUserData(r.Context(), user); err != nil {
			panic(err)
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		httputil.WriteHttpResponse(w, http.StatusOK, au)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from request
		userId := getIdPathVar(r)

		// Execute action
		user, err := c.uServ.GetUserDataByUserId(r.Context(), userId)
		if err != nil {
			panic(err)
		}

		// Check if a user was found
		if user == nil {
			err = e.NewError(e.LogicUserNotFound, fmt.Sprintf("Could not find user %d.", userId))
			log.Debug(err.StackTrace())
			panic(err)
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		httputil.WriteHttpResponse(w, http.StatusOK, au)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		id := getIdPathVar(r)

		// Read API model from request
		var auu model.UpdateUserData
		httputil.ReadHttpBody(r, &auu)

		// Validate model
		if err := validator.ValidateUpdateUser(&auu); err != nil {
			panic(err)
		}

		// Convert to logic model
		user := mapper.FromUpdateUserData(id, &auu)

		// Execute action
		if err := c.uServ.UpdateUserData(r.Context(), user); err != nil {
			panic(err)
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		httputil.WriteHttpResponse(w, http.StatusOK, au)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from request
		userId := getIdPathVar(r)

		// Execute action
		err := c.uServ.DeleteUserById(r.Context(), userId)
		if err != nil {
			panic(err)
		}

		// Write response
		httputil.WriteHttpResponse(w, http.StatusNoContent, nil)
	}
}

// UpdateUserPasswordHandler returns a handler for "PUT /users/{id}/password".
func (c *UserController) UpdateUserPasswordHandler() http.HandlerFunc {
	// swagger:operation PUT /users/{id}/password users updateUserPassword
	//
	// Update the password of a user by its ID.
	//
	// # Username / password rules
	//
	// __Username:__
	//
	// ⦁ Minimum length: 4
	// ⦁ Maximum length: 100
	// ⦁ Allowed characters: `0-9 a-z A-Z - .`
	//
	// __Password:__
	//
	// ⦁ Minimum length: 8
	// ⦁ Maximum length: 100
	// ⦁ Allowed characters: `0-9 a-z A-Z ! \ # $ % & ' ( ) * + , - . / : ; = ? @ [ \ ] ^ _ { | } ~`
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		userId := getIdPathVar(r)

		// Read API model from request
		var auupw model.UpdateUserPassword
		httputil.ReadHttpBody(r, &auupw)

		// Validate password
		if err := validator.ValidateUpdateUserPassword(&auupw); err != nil {
			panic(err)
		}

		// Execute action
		if err := c.uServ.UpdateUserPassword(r.Context(), userId, auupw.Password); err != nil {
			panic(err)
		}

		// Write response
		httputil.WriteHttpResponse(w, http.StatusNoContent, nil)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from request
		userId := getIdPathVar(r)

		// Execute action
		userRoles, err := c.uServ.GetUserRoles(r.Context(), userId)
		if err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aur := mapper.ToRoles(userRoles)
		httputil.WriteHttpResponse(w, http.StatusOK, aur)
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
	//   '412':
	//     "$ref": "#/responses/ErrorResponse"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(w http.ResponseWriter, r *http.Request) {
		// Get ID from request
		userId := getIdPathVar(r)

		// Read API model from request
		var auurs model.UpdateUserRoles
		httputil.ReadHttpBody(r, &auurs)

		// Validate model
		if err := validator.ValidateUpdateUserRoles(&auurs); err != nil {
			panic(err)
		}

		// Convert to logic model
		userRoles := mapper.FromRoles(&auurs)

		// Execute action
		if err := c.uServ.SetUserRoles(r.Context(), userId, userRoles); err != nil {
			panic(err)
		}

		// Convert to API model and write response
		aurs := mapper.ToRoles(userRoles)
		httputil.WriteHttpResponse(w, http.StatusOK, aurs)
	}
}
