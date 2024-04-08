package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/loc"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	"kellnhofer.com/work-log/web/pages"
)

type entryFormInput struct {
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

// GetCreateHandler returns a handler for "GET /create".
func (c *EntryController) GetCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /create.")
		return c.handleShowCreate(eCtx)
	}
}

// PostCreateHandler returns a handler for "POST /create".
func (c *EntryController) PostCreateHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /create.")
		return c.handleExecuteCreate(eCtx)
	}
}

// GetEditHandler returns a handler for "GET /edit/{id}".
func (c *EntryController) GetEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /edit/{id}.")
		return c.handleShowEdit(eCtx)
	}
}

// PostEditHandler returns a handler for "POST /edit/{id}".
func (c *EntryController) PostEditHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /edit/{id}.")
		return c.handleExecuteEdit(eCtx)
	}
}

// GetCopyHandler returns a handler for "GET /copy/{id}".
func (c *EntryController) GetCopyHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /copy/{id}.")
		return c.handleShowCopy(eCtx)
	}
}

// PostCopyHandler returns a handler for "POST /copy/{id}".
func (c *EntryController) PostCopyHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /copy/{id}.")
		return c.handleExecuteCopy(eCtx)
	}
}

// GetDeleteHandler returns a handler for "GET /delete/{id}".
func (c *EntryController) GetDeleteHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle GET /delete/{id}.")
		return c.handleShowDelete(eCtx)
	}
}

// PostDeleteHandler returns a handler for "POST /delete/{id}".
func (c *EntryController) PostDeleteHandler() echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		log.Verb("Handle POST /delete/{id}.")
		return c.handleExecuteDelete(eCtx)
	}
}

// --- Create handler functions ---

func (c *EntryController) handleShowCreate(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	entryTypeId := 0
	if len(entryTypes) > 0 {
		entryTypeId = entryTypes[0].Id
	}
	model := c.mapper.CreateInitialCreateViewModel(prevUrl, entryTypeId, time.Now(), 0, entryTypes,
		entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.CreateEntryPage(model))
}

func (c *EntryController) handleExecuteCreate(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(0, userId, input)
	if err != nil {
		return c.handleCreateError(eCtx, err, input)
	}

	// Create entry
	if err := c.eServ.CreateEntry(ctx, entry); err != nil {
		return c.handleCreateError(eCtx, err, input)
	}

	return c.handleCreateSuccess(eCtx)
}

func (c *EntryController) handleCreateSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleCreateError(eCtx echo.Context, err error, input *entryFormInput,
) error {
	// Get context
	ctx := getContext(eCtx)

	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.mapper.CreateCreateViewModel(prevUrl, em, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.CreateEntryPage(model))
}

// --- Edit handler functions ---

func (c *EntryController) handleShowEdit(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	model := c.mapper.CreateInitialEditViewModel(prevUrl, entry.Id, entry.TypeId, entry.StartTime,
		entry.EndTime, entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.EditEntryPage(model))
}

func (c *EntryController) handleExecuteEdit(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(entryId, userId, input)
	if err != nil {
		return c.handleEditError(eCtx, err, entryId, input)
	}

	// Update entry
	if err := c.eServ.UpdateEntry(ctx, entry); err != nil {
		return c.handleEditError(eCtx, err, entryId, input)
	}

	return c.handleEditSuccess(eCtx)
}

func (c *EntryController) handleEditSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleEditError(eCtx echo.Context, err error, id int,
	input *entryFormInput) error {
	// Get context
	ctx := getContext(eCtx)

	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.mapper.CreateEditViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.EditEntryPage(model))
}

// --- Copy handler functions ---

func (c *EntryController) handleShowCopy(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}
	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	model := c.mapper.CreateInitialCopyViewModel(prevUrl, entry.Id, entry.TypeId, entry.StartTime,
		entry.EndTime, entry.ActivityId, entry.Description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.CopyEntryPage(model))
}

func (c *EntryController) handleExecuteCopy(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get form inputs
	input := c.getEntryFormInput(eCtx)

	// Create model
	entry, err := c.createEntryModel(0, userId, input)
	if err != nil {
		return c.handleCopyError(eCtx, err, entryId, input)
	}

	// Create entry
	if err := c.eServ.CreateEntry(ctx, entry); err != nil {
		return c.handleCopyError(eCtx, err, entryId, input)
	}

	return c.handleCopySuccess(eCtx)
}

func (c *EntryController) handleCopySuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

func (c *EntryController) handleCopyError(eCtx echo.Context, err error, id int,
	input *entryFormInput) error {
	// Get context
	ctx := getContext(eCtx)

	// Get error message
	ec := getErrorCode(err)
	em := loc.GetErrorMessageString(ec)

	// Get entry master data
	entryTypes, entryActivities, err := c.getEntryMasterData(ctx)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	entryTypeId, _ := strconv.Atoi(input.typeId)
	entryActivityId, _ := strconv.Atoi(input.activityId)
	model := c.mapper.CreateCopyViewModel(prevUrl, em, id, entryTypeId, input.date, input.startTime,
		input.endTime, entryActivityId, input.description, entryTypes, entryActivities)

	// Render
	return web.Render(eCtx, http.StatusOK, pages.CopyEntryPage(model))
}

// --- Delete handler functions ---

func (c *EntryController) handleShowDelete(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Get entry
	entry, err := c.getEntry(ctx, entryId, userId)
	if err != nil {
		return err
	}

	// Create view model
	prevUrl := getPreviousUrl(eCtx)
	model := c.mapper.CreateDeleteViewModel(prevUrl, "", entry.Id)
	return web.Render(eCtx, http.StatusOK, pages.DeleteEntryPage(model))
}

func (c *EntryController) handleExecuteDelete(eCtx echo.Context) error {
	// Get context
	ctx := getContext(eCtx)

	// Get current user ID
	userId := getCurrentUserId(ctx)

	// Get ID
	entryId, err := getIdPathVar(eCtx)
	if err != nil {
		return err
	}

	// Delete entry
	if err := c.eServ.DeleteEntryByIdAndUserId(ctx, entryId, userId); err != nil {
		return err
	}

	return c.handleDeleteSuccess(eCtx)
}

func (c *EntryController) handleDeleteSuccess(eCtx echo.Context) error {
	prevUrl := getPreviousUrl(eCtx)
	return eCtx.Redirect(http.StatusFound, prevUrl)
}

// --- Form input retrieval functions ---

func (c *EntryController) getEntryFormInput(eCtx echo.Context) *entryFormInput {
	i := entryFormInput{}
	i.typeId = eCtx.FormValue("type")
	i.date = eCtx.FormValue("date")
	i.startTime = eCtx.FormValue("start-time")
	i.endTime = eCtx.FormValue("end-time")
	i.activityId = eCtx.FormValue("activity")
	i.description = eCtx.FormValue("description")
	return &i
}

// --- Model converter functions ---

func (c *EntryController) createEntryModel(id int, userId int, input *entryFormInput) (*model.Entry,
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
