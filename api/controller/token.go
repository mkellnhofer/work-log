package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/api/mapper"
	"kellnhofer.com/work-log/api/model"
	"kellnhofer.com/work-log/api/validator"
	"kellnhofer.com/work-log/pkg/service"
)

// TokenController handles requests for token endpoints.
type TokenController struct {
	tServ *service.TokenService
}

// NewTokenController create a new token controller.
func NewTokenController(ts *service.TokenService) *TokenController {
	return &TokenController{ts}
}

// --- Parameters ---

// swagger:parameters createToken
type CreateTokenParameters struct {
	// in: body
	// required: true
	Body model.CreateToken
}

// swagger:parameters getToken deleteToken
type GetTokenParameters struct {
	// The ID of the token.
	//
	// in: path
	// required: true
	Id int `json:"id"`
}

// --- Responses ---

// The list of tokens.
// swagger:response GetTokensResponse
type GetTokensResponse struct {
	// in: body
	Body model.TokenList
}

// The token.
// swagger:response GetTokenResponse
type GetTokenResponse struct {
	// in: body
	Body model.Token
}

// The created token.
// swagger:response CreateTokenResponse
type CreateTokenResponse struct {
	// in: body
	Body model.Token
}

// --- Endpoints ---

// GetTokensHandler returns a handler for "GET /user/tokens".
func (c *TokenController) GetTokensHandler() echo.HandlerFunc {
	// swagger:operation GET /user/tokens user listTokens
	//
	// Lists all tokens of the current user.
	//
	// ---
	//
	// security:
	// - Basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetTokensResponse"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials\n
	//       ⦁ [-104]: Invalid token"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-105]: Bearer auth not allowed"
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
		tokens, err := c.tServ.GetCurrentUserTokens(getContext(eCtx))
		if err != nil {
			return err
		}

		// Convert to API model and write response
		at := mapper.ToTokens(tokens)
		return writeResponse(eCtx, http.StatusOK, at)
	}
}

// CreateTokenHandler returns a handler for "POST /user/tokens".
func (c *TokenController) CreateTokenHandler() echo.HandlerFunc {
	// swagger:operation POST /user/tokens user createToken
	//
	// Create a token.
	//
	// Creates a new API token for the current user. The full token string is only returned in the
	// response of this request.
	//
	// # Input Rules
	//
	// __Name:__
	//
	// ⦁ Minimum length: 1
	// ⦁ Maximum length: 30
	//
	// ---
	//
	// security:
	// - Basic: []
	//
	// consumes:
	// - application/json
	// produces:
	// - application/json
	//
	// responses:
	//   '201':
	//     "$ref": "#/responses/CreateTokenResponse"
	//   '400':
	//     description: "__Bad Request__\n\n
	//       ⦁ [-311]: Empty string\n
	//       ⦁ [-312]: Too long string"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials\n
	//       ⦁ [-104]: Invalid token"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-105]: Bearer auth not allowed"
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
		var act model.CreateToken
		if err := readRequestBody(eCtx, &act); err != nil {
			return err
		}

		// Validate model
		if err := validator.ValidateCreateToken(&act); err != nil {
			return err
		}

		// Execute action
		token, err := c.tServ.CreateCurrentUserToken(getContext(eCtx), act.Name)
		if err != nil {
			return err
		}

		// Convert to API model and write response
		at := mapper.ToTokenFull(token)
		return writeResponse(eCtx, http.StatusCreated, at)
	}
}

// GetTokenHandler returns a handler for "GET /user/tokens/{id}".
func (c *TokenController) GetTokenHandler() echo.HandlerFunc {
	// swagger:operation GET /user/tokens/{id} user getToken
	//
	// Get a token.
	//
	// ---
	//
	// security:
	// - Basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '200':
	//     "$ref": "#/responses/GetTokenResponse"
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials\n
	//       ⦁ [-104]: Invalid token"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-105]: Bearer auth not allowed"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-413]: Token not found"
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
		// Get token ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		token, err := c.tServ.GetCurrentUserTokenById(getContext(eCtx), id)
		if err != nil {
			return err
		}

		// Convert to API model and write response
		at := mapper.ToToken(token)
		return writeResponse(eCtx, http.StatusOK, at)
	}
}

// DeleteTokenHandler returns a handler for "DELETE /user/tokens/{id}".
func (c *TokenController) DeleteTokenHandler() echo.HandlerFunc {
	// swagger:operation DELETE /user/tokens/{id} user deleteToken
	//
	// Delete a token.
	//
	// ---
	//
	// security:
	// - Basic: []
	//
	// produces:
	// - application/json
	//
	// responses:
	//   '204':
	//     description: No content.
	//   '401':
	//     description: "__Unauthorized__\n\n
	//       ⦁ [-101]: Invalid authentication data\n
	//       ⦁ [-102]: Invalid credentials\n
	//       ⦁ [-104]: Invalid token"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '403':
	//     description: "__Forbidden__\n\n
	//       ⦁ [-105]: Bearer auth not allowed"
	//     schema:
	//       "$ref": "#/definitions/Error"
	//   '404':
	//     description: "__Not Found__\n\n
	//       ⦁ [-413]: Token not found"
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
		// Get token ID from request
		id, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		// Execute action
		if err := c.tServ.DeleteCurrentUserTokenById(getContext(eCtx), id); err != nil {
			return err
		}

		// Write response
		return eCtx.NoContent(http.StatusNoContent)
	}
}
