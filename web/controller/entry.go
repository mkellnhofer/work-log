package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	"kellnhofer.com/work-log/web/view/hx"
)

type entryInput struct {
	typeId      string
	date        string
	startTime   string
	endTime     string
	activityId  string
	description string
}

// EntryController handles requests for entry endpoints.
type EntryController struct {
	baseController

	mapper *mapper.EntryMapper
}

// NewEntryController creates a new entry controller.
func NewEntryController(uServ *service.UserService, eServ *service.EntryService) *EntryController {
	return &EntryController{
		baseController: baseController{
			uServ: uServ,
			eServ: eServ,
		},
		mapper: mapper.NewEntryMapper(),
	}
}

// --- Endpoints ---

// GetActivitiesHandler returns a handler for "GET /entry-modal/activities".
func (c *EntryController) GetActivitiesHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			return c.handleGetActivities(eCtx, ctx)
		})
	}
}

// GetCreateHandler returns a handler for "GET /entry-modal/create".
func (c *EntryController) GetCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			return c.handleShowCreate(eCtx, ctx)
		})
	}
}

// PostCreateHandler returns a handler for "POST /entry-modal/create".
func (c *EntryController) PostCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			input := c.getEntryInput(eCtx)
			return c.handleExecuteCreate(eCtx, ctx, input)
		})
	}
}

// GetCopyHandler returns a handler for "GET /entry-modal/copy/{id}".
func (c *EntryController) GetCopyHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			id, err := getIdPathVar(eCtx)
			if err != nil {
				return err
			}
			return c.handleShowCopy(eCtx, ctx, id)
		})
	}
}

// GetEditHandler returns a handler for "GET /entry-modal/edit/{id}".
func (c *EntryController) GetEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			id, err := getIdPathVar(eCtx)
			if err != nil {
				return err
			}
			return c.handleShowEdit(eCtx, ctx, id)
		})
	}
}

// PostEditHandler returns a handler for "POST /entry-modal/edit/{id}".
func (c *EntryController) PostEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			id, err := getIdPathVar(eCtx)
			if err != nil {
				return err
			}
			input := c.getEntryInput(eCtx)
			return c.handleExecuteEdit(eCtx, ctx, id, input)
		})
	}
}

// GetDeleteHandler returns a handler for "GET /entry-modal/delete/{id}".
func (c *EntryController) GetDeleteHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			id, err := getIdPathVar(eCtx)
			if err != nil {
				return err
			}
			return c.handleShowDelete(eCtx, ctx, id)
		})
	}
}

// PostDeleteHandler returns a handler for "POST /entry-modal/delete/{id}".
func (c *EntryController) PostDeleteHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			id, err := getIdPathVar(eCtx)
			if err != nil {
				return err
			}
			return c.handleExecuteDelete(eCtx, ctx, id)
		})
	}
}

// PostCancelHandler returns a handler for "POST /entry-modal/cancel".
func (c *EntryController) PostCancelHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return c.wrapHandler(eCtx, func(ctx context.Context) error {
			return eCtx.NoContent(http.StatusOK)
		})
	}
}

func (c *EntryController) wrapHandler(eCtx echo.Context, hf func(context.Context) error) error {
	isHtmxReq := web.IsHtmxRequest(eCtx)
	if !isHtmxReq {
		err := e.NewError(e.ValUnknown, "Not a HTMX request.")
		log.Debug(err.StackTrace())
		return err
	}
	return hf(getContext(eCtx))
}

func (c *EntryController) getEntryInput(eCtx echo.Context) *entryInput {
	return &entryInput{
		typeId:      eCtx.FormValue("type"),
		date:        eCtx.FormValue("date"),
		startTime:   eCtx.FormValue("start-time"),
		endTime:     eCtx.FormValue("end-time"),
		activityId:  eCtx.FormValue("activity"),
		description: eCtx.FormValue("description"),
	}
}

// --- Handler functions ---

func (c *EntryController) handleGetActivities(eCtx echo.Context, ctx context.Context) error {
	entryTypeId, err := getTypeIdQueryParam(eCtx)
	if err != nil {
		return err
	}

	// Get entry master data
	entryActivities, err := c.getEntryActivities(ctx, entryTypeId)
	if err != nil {
		return err
	}

	// Create view model
	viewData := c.mapper.CreateEntryActivitiesViewModel(entryActivities)

	// Render
	return c.handleShowSuccess(eCtx, hx.EntryModalActivityOptions(viewData))
}

func (c *EntryController) handleShowCreate(eCtx echo.Context, ctx context.Context) error {
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, model.EntryTypeIdWork)
	if err != nil {
		return err
	}

	// Create view model
	entry := model.NewEntry()
	entry.TypeId = model.EntryTypeIdWork
	entry.StartTime = time.Now()
	entry.EndTime = time.Now()
	entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

	// Render
	return c.handleShowSuccess(eCtx, hx.EntryModalCreate(entryViewData))
}

func (c *EntryController) handleExecuteCreate(eCtx echo.Context, ctx context.Context,
	input *entryInput) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Create model
	entry, err := c.createEntryModel(0, userId, input)
	if err != nil {
		return c.handleExecuteError(eCtx, err)
	}

	// Create entry
	if err := c.eServ.CreateEntry(ctx, entry); err != nil {
		return c.handleExecuteError(eCtx, err)
	}

	// Return empty response
	return c.handleExecuteSuccess(eCtx)
}

func (c *EntryController) handleShowCopy(eCtx echo.Context, ctx context.Context, entryId int) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, entry.TypeId)
	if err != nil {
		return err
	}

	// Create view model
	entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

	// Render
	return c.handleShowSuccess(eCtx, hx.EntryModalCopy(entryViewData))
}

func (c *EntryController) handleShowEdit(eCtx echo.Context, ctx context.Context, entryId int) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx, entry.TypeId)
	if err != nil {
		return err
	}

	// Create view model
	entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

	// Render
	return c.handleShowSuccess(eCtx, hx.EntryModalEdit(entryViewData))
}

func (c *EntryController) handleExecuteEdit(eCtx echo.Context, ctx context.Context, entryId int,
	input *entryInput) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Create model
	entry, err := c.createEntryModel(entryId, userId, input)
	if err != nil {
		return c.handleExecuteError(eCtx, err)
	}

	// Update entry
	if err := c.eServ.UpdateEntry(ctx, entry); err != nil {
		return c.handleExecuteError(eCtx, err)
	}

	// Return empty response
	return c.handleExecuteSuccess(eCtx)
}

func (c *EntryController) handleShowDelete(eCtx echo.Context, ctx context.Context, entryId int,
) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}

	// Render
	return c.handleShowSuccess(eCtx, hx.EntryModalDelete(entry.Id))
}

func (c *EntryController) handleExecuteDelete(eCtx echo.Context, ctx context.Context,
	entryId int) error {
	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Delete entry
	if err := c.eServ.DeleteEntryByIdAndUserId(ctx, entryId, userId); err != nil {
		return err
	}

	// Return empty response
	return c.handleExecuteSuccess(eCtx)
}

func (c *EntryController) handleShowSuccess(eCtx echo.Context, t templ.Component) error {
	// Render
	return web.RenderHx(eCtx, http.StatusOK, t)
}

func (c *EntryController) handleExecuteSuccess(eCtx echo.Context) error {
	// Set HTMX triggers
	web.HtmxTrigger(eCtx, "wlChangedEntries")
	// Return empty response
	return eCtx.NoContent(http.StatusOK)
}

func (c *EntryController) handleExecuteError(eCtx echo.Context, err error) error {
	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)
	// Render
	web.HtmxRetarget(eCtx, "#wl-entry-modal-error")
	return web.RenderHx(eCtx, http.StatusOK, hx.EntryModalError(em))
}

// --- Model converter functions ---

func (c *EntryController) createEntryModel(id int, userId int, input *entryInput) (*model.Entry,
	error) {
	entry := model.NewEntry()
	entry.Id = id
	entry.UserId = userId

	var err error

	// Convert type ID
	entry.TypeId, err = parseId(input.typeId, false)
	if err != nil {
		return nil, err
	}

	// Convert start/end time
	if _, err := parseDateTime(input.date, "00:00", e.ValDateInvalid); err != nil {
		return nil, err
	}
	entry.StartTime, err = parseDateTime(input.date, input.startTime, e.ValStartTimeInvalid)
	if err != nil {
		return nil, err
	}
	entry.EndTime, err = parseDateTime(input.date, input.endTime, e.ValEndTimeInvalid)
	if err != nil {
		return nil, err
	}

	// Convert activity ID
	entry.ActivityId, err = parseId(input.activityId, true)
	if err != nil {
		return nil, err
	}

	// Validate description
	if err = validateStringLength(input.description, model.MaxLengthEntryDescription,
		e.ValDescriptionTooLong); err != nil {
		return nil, err
	}
	entry.Description = input.description

	return entry, nil
}
