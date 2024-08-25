package controller

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	vm "kellnhofer.com/work-log/web/model"
	"kellnhofer.com/work-log/web/view/hx"
)

// UserController handles requests for user endpoints.
type UserController struct {
	handlerHelper
	baseUserController
}

func NewUserController(uServ *service.UserService) *UserController {
	return &UserController{
		baseUserController: *newBaseUserController(uServ),
	}
}

// GetHxUserProfileModalHandler returns a handler for "GET /hx/user-profile-modal".
func (c *UserController) GetHxUserProfileModalHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		profileInfo, err := c.getUserProfileInfoViewData(ctx)
		if err != nil {
			return err
		}
		return web.RenderHx(eCtx, http.StatusOK, hx.UserProfileModal(profileInfo))
	})
}

// PostHxUserProfileModalCloseHandler returns a handler for "POST /hx/user-profile-modal/close".
func (c *UserController) PostHxUserProfileModalCloseHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		return eCtx.NoContent(http.StatusOK)
	})
}

func (c *UserController) getUserProfileInfoViewData(ctx context.Context) (*vm.UserProfileInfo, error) {
	userId := getCurrentUserId(ctx)
	user, err := c.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	userContract, err := c.getUserContract(ctx, userId)
	if err != nil {
		return nil, err
	}
	return c.uMapper.CreateUserProfileInfoViewModel(user, userContract), nil
}
