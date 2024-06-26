package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/api/mapper"
	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/api/validator"
	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/service"
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

// swagger:parameters updateCurrentUserPassword
type UpdateCurrentUserPasswordParameters struct {
	// in: body
	// required: true
	Body model.UpdateUserPassword
}

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
func (c *UserController) GetCurrentUserHandler() echo.HandlerFunc {
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
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Execute action
		user, err := c.uServ.GetCurrentUserData(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		return writeResponse(eCtx, http.StatusOK, au)
	}
}

// UpdateCurrentUserPasswordHandler returns a handler for "PUT /user/password".
func (c *UserController) UpdateCurrentUserPasswordHandler() echo.HandlerFunc {
	// swagger:operation PUT /user/password user updateCurrentUserPassword
	//
	// Update the password of the current user.
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-318]: Invalid password"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Read API model from request
		var auupw model.UpdateUserPassword
		if err := readRequestBody(eCtx, &auupw); err != nil {
			return err
		}

		// Validate password
		if err := validator.ValidateUpdateUserPassword(&auupw); err != nil {
			return err
		}

		// Execute action
		if err := c.uServ.UpdateCurrentUserPassword(getContext(eCtx), auupw.Password); err != nil {
			return err
		}

		// Write response
		return writeResponse(eCtx, http.StatusNoContent, nil)
	}
}

// GetCurrentUserRolesHandler returns a handler for "GET /user/roles".
func (c *UserController) GetCurrentUserRolesHandler() echo.HandlerFunc {
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
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Execute action
		userRoles, err := c.uServ.GetCurrentUserRoles(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aur := mapper.ToRoles(userRoles)
		return writeResponse(eCtx, http.StatusOK, aur)
	}
}

// GetUsersHandler returns a handler for "GET /users".
func (c *UserController) GetUsersHandler() echo.HandlerFunc {
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
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-201]: No right to list users"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Execute action
		users, err := c.uServ.GetUserDatas(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aus := mapper.ToUserDatas(users)
		return writeResponse(eCtx, http.StatusOK, aus)
	}
}

// CreateUserHandler returns a handler for "POST /users".
func (c *UserController) CreateUserHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-308]: Null field\n
	//       ⦁ [-309]: Negative number\n
	//       ⦁ [-310]: Negative or zero number\n
	//       ⦁ [-311]: Empty string\n
	//       ⦁ [-312]: Too long string\n
	//       ⦁ [-313]: Invalid date\n
	//       ⦁ [-317]: Invalid username\n
	//       ⦁ [-318]: Invalid password\n
	//       ⦁ [-410]: Invalid contract working hours\n
	//       ⦁ [-411]: Invalid contract vacation days"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-202]: No right to create user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '409':
	//     description: "__Conflict__\n\n
	//       ⦁ [-409]: User already exists"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Read API model from request
		var acu model.CreateUserData
		if err := readRequestBody(eCtx, &acu); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateCreateUser(&acu); err != nil {
			return err
		}

		// Convert to logic model
		user := mapper.FromCreateUserData(&acu)

		// Execute action
		if err := c.uServ.CreateUserData(getContext(eCtx), user); err != nil {
			return err
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		return writeResponse(eCtx, http.StatusOK, au)
	}
}

// GetUserHandler returns a handler for "GET /users/{id}".
func (c *UserController) GetUserHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-201]: No right to get user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get user ID from request
		userId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		user, err := c.uServ.GetUserDataByUserId(getContext(eCtx), userId)
		if err != nil {
			return err
		}

		// Check if a user was found
		if user == nil {
			err := e.NewError(e.LogicUserNotFound, fmt.Sprintf("Could not find user %d.", userId))
			log.Debug(err.StackTrace())
			return err
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		return writeResponse(eCtx, http.StatusOK, au)
	}
}

// UpdateUserHandler returns a handler for "PUT /users/{id}".
func (c *UserController) UpdateUserHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-309]: Negative number\n
	//       ⦁ [-310]: Negative or zero number\n
	//       ⦁ [-311]: Empty string\n
	//       ⦁ [-312]: Too long string\n
	//       ⦁ [-313]: Invalid date\n
	//       ⦁ [-317]: Invalid username\n
	//       ⦁ [-318]: Invalid password\n
	//       ⦁ [-410]: Invalid contract working hours\n
	//       ⦁ [-411]: Invalid contract vacation days"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-202]: No right to update user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '409':
	//     description: "__Conflict__\n\n
	//       ⦁ [-409]: User already exists"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Read API model from request
		var auu model.UpdateUserData
		if err := readRequestBody(eCtx, &auu); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateUpdateUser(&auu); err != nil {
			return err
		}

		// Convert to logic model
		user := mapper.FromUpdateUserData(id, &auu)

		// Execute action
		if err := c.uServ.UpdateUserData(getContext(eCtx), user); err != nil {
			return err
		}

		// Convert to API model and write response
		au := mapper.ToUserData(user)
		return writeResponse(eCtx, http.StatusOK, au)
	}
}

// DeleteUserHandler returns a handler for "DELETE /users/{id}".
func (c *UserController) DeleteUserHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-202]: No right to delete user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get user ID from request
		userId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		if err := c.uServ.DeleteUserById(getContext(eCtx), userId); err != nil {
			return err
		}

		// Write response
		return writeResponse(eCtx, http.StatusNoContent, nil)
	}
}

// UpdateUserPasswordHandler returns a handler for "PUT /users/{id}/password".
func (c *UserController) UpdateUserPasswordHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-318]: Invalid password"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-202]: No right to update user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		userId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Read API model from request
		var auupw model.UpdateUserPassword
		if err := readRequestBody(eCtx, &auupw); err != nil {
			return err
		}

		// Validate password
		if err := validator.ValidateUpdateUserPassword(&auupw); err != nil {
			return err
		}

		// Execute action
		if err := c.uServ.UpdateUserPassword(getContext(eCtx), userId, auupw.Password); err != nil {
			return err
		}

		// Write response
		return writeResponse(eCtx, http.StatusNoContent, nil)
	}
}

// GetUserRolesHandler returns a handler for "GET /users/{id}/roles".
func (c *UserController) GetUserRolesHandler() echo.HandlerFunc {
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
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-201]: No right to get user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get user ID from request
		userId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		userRoles, err := c.uServ.GetUserRoles(getContext(eCtx), userId)
		if err != nil {
			return err
		}

		// Convert to API model and write response
		aur := mapper.ToRoles(userRoles)
		return writeResponse(eCtx, http.StatusOK, aur)
	}
}

// UpdateUserRolesHandler returns a handler for "PUT /users/{id}/roles".
func (c *UserController) UpdateUserRolesHandler() echo.HandlerFunc {
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
	//     description: "__Bad Request__\n\n
	//       ⦁ [-303]: Invalid ID\n
	//       ⦁ [-312]: Too long string\n
	//       ⦁ [-315]: Empty array\n
	//       ⦁ [-318]: Invalid role"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-202]: No right to update user"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-407]: Role not found\n
	//       ⦁ [-408]: User not found"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '412':
	//     description: "__Precondition Failed__\n\n
	//       ⦁ [-103]: User not activated"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   default:
	//     "$ref": "#/responses/ErrorResponse"
	return func(eCtx echo.Context) error {
		// Get ID from request
		userId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Read API model from request
		var auurs model.UpdateUserRoles
		if err := readRequestBody(eCtx, &auurs); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateUpdateUserRoles(&auurs); err != nil {
			return err
		}

		// Convert to logic model
		userRoles := mapper.FromRoles(&auurs)

		// Execute action
		if err := c.uServ.SetUserRoles(getContext(eCtx), userId, userRoles); err != nil {
			return err
		}

		// Convert to API model and write response
		aurs := mapper.ToRoles(userRoles)
		return writeResponse(eCtx, http.StatusOK, aurs)
	}
}
