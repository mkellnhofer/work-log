package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
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

// GetHxActivitiesHandler returns a handler for "GET /hx/entry-modal/activities".
func (c *EntryController) GetHxActivitiesHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		entryTypeId, err := getTypeIdQueryParam(eCtx)
		if err != nil {
			return err
		}

		entryActivities, err := c.getEntryActivities(ctx, entryTypeId)
		if err != nil {
			return err
		}

		viewData := c.mapper.CreateEntryActivitiesViewModel(entryActivities)

		return c.handleShowSuccess(eCtx, hx.EntryModalActivityOptions(viewData))
	})
}

// GetHxCreateHandler returns a handler for "GET /hx/entry-modal/create".
func (c *EntryController) GetHxCreateHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		entryTypes, entryActivities, err := c.getEntryMasterData(ctx, model.EntryTypeIdWork)
		if err != nil {
			return err
		}

		entry := model.NewEntry()
		entry.TypeId = model.EntryTypeIdWork
		entry.StartTime = time.Now()
		entry.EndTime = time.Now()
		entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

		return c.handleShowSuccess(eCtx, hx.EntryModalCreate(entryViewData))
	})
}

// PostHxCreateHandler returns a handler for "POST /hx/entry-modal/create".
func (c *EntryController) PostHxCreateHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		input := c.getEntryInput(eCtx)

		entry, err := c.createEntryModel(0, userId, input)
		if err != nil {
			return c.handleExecuteError(eCtx, err)
		}

		if err := c.eServ.CreateEntry(ctx, entry); err != nil {
			return c.handleExecuteError(eCtx, err)
		}

		return c.handleExecuteSuccess(eCtx)
	})
}

// GetHxCopyHandler returns a handler for "GET /hx/entry-modal/copy/{id}".
func (c *EntryController) GetHxCopyHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		entryId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		entry, err := c.getEntry(ctx, entryId, userId)
		if err != nil {
			return err
		}
		entryTypes, entryActivities, err := c.getEntryMasterData(ctx, entry.TypeId)
		if err != nil {
			return err
		}

		entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

		return c.handleShowSuccess(eCtx, hx.EntryModalCopy(entryViewData))
	})
}

// GetHxEditHandler returns a handler for "GET /hx/entry-modal/edit/{id}".
func (c *EntryController) GetHxEditHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		entryId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		entry, err := c.getEntry(ctx, entryId, userId)
		if err != nil {
			return err
		}
		entryTypes, entryActivities, err := c.getEntryMasterData(ctx, entry.TypeId)
		if err != nil {
			return err
		}

		entryViewData := c.mapper.CreateEntryDataViewModel(entry, entryTypes, entryActivities)

		return c.handleShowSuccess(eCtx, hx.EntryModalEdit(entryViewData))
	})
}

// PostHxEditHandler returns a handler for "POST /hx/entry-modal/edit/{id}".
func (c *EntryController) PostHxEditHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		entryId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}
		input := c.getEntryInput(eCtx)

		entry, err := c.createEntryModel(entryId, userId, input)
		if err != nil {
			return c.handleExecuteError(eCtx, err)
		}

		if err := c.eServ.UpdateEntry(ctx, entry); err != nil {
			return c.handleExecuteError(eCtx, err)
		}

		return c.handleExecuteSuccess(eCtx)
	})
}

// GetHxDeleteHandler returns a handler for "GET /hx/entry-modal/delete/{id}".
func (c *EntryController) GetHxDeleteHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		entryId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		entry, err := c.getEntry(ctx, entryId, userId)
		if err != nil {
			return err
		}

		return c.handleShowSuccess(eCtx, hx.EntryModalDelete(entry.Id))
	})
}

// PostHxDeleteHandler returns a handler for "POST /hx/entry-modal/delete/{id}".
func (c *EntryController) PostHxDeleteHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		userId := getCurrentUserId(ctx)
		entryId, err := getIdPathVar(eCtx)
		if err != nil {
			return err
		}

		if err := c.eServ.DeleteEntryByIdAndUserId(ctx, entryId, userId); err != nil {
			return err
		}

		return c.handleExecuteSuccess(eCtx)
	})
}

// PostHxCancelHandler returns a handler for "POST /hx/entry-modal/cancel".
func (c *EntryController) PostHxCancelHandler() echo.HandlerFunc {
	return c.hxHandler(func(eCtx echo.Context, ctx context.Context) error {
		return eCtx.NoContent(http.StatusOK)
	})
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
	web.HtmxRetarget(eCtx, "#wl-modal-error-container")
	return web.RenderHx(eCtx, http.StatusOK, hx.ModalError(em))
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
