package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
	"kellnhofer.com/work-log/pkg/util"
	"kellnhofer.com/work-log/web"
	"kellnhofer.com/work-log/web/mapper"
	vm "kellnhofer.com/work-log/web/model"
)

const queryDateTimeFormat = "200601021504"

// --- Handler Helper ---

type handlerFunc func(eCtx echo.Context, ctx context.Context) error
type hxHandlerFunc func(eCtx echo.Context, ctx context.Context) error
type resourceHandlerFunc func(eCtx echo.Context, ctx context.Context) error

type handlerHelper struct {}

func (hh *handlerHelper) handler(hf handlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		ctx := getContext(eCtx)
		return hf(eCtx, ctx)
	}
}

func (hh *handlerHelper) hxHandler(hf hxHandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		isHtmxReq := web.IsHtmxRequest(eCtx)
		if !isHtmxReq {
			err := e.NewError(e.ValUnknown, "Not a HTMX request.")
			log.Debug(err.StackTrace())
			return err
		}
		ctx := getContext(eCtx)
		return hf(eCtx, ctx)
	}
}

func (hh *handlerHelper) resourceHandler(hf resourceHandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return hf(eCtx, getContext(eCtx))
	}
}

// --- Base User Controller ---

type baseUserController struct {
	uServ  *service.UserService
	uMapper *mapper.UserMapper
}

func newBaseUserController(uServ *service.UserService) *baseUserController {
	return &baseUserController{
		uServ: uServ,
		uMapper: mapper.NewUserMapper(),
	}
}

func (c *baseUserController) getUser(ctx context.Context, userId int) (*model.User, error) {
	return c.uServ.GetUserById(ctx, userId)
}

func (c *baseUserController) getUserContract(ctx context.Context, userId int) (*model.Contract,
	error) {
	return c.uServ.GetUserContractByUserId(ctx, userId)
}

func (c *baseUserController) getUserInfoViewData(ctx context.Context) (*vm.UserInfo, error) {
	userId := getCurrentUserId(ctx)
	user, err := c.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return c.uMapper.CreateUserInfoViewModel(user), nil
}

// --- Base Entry Controller ---

type baseEntryController struct {
	eServ *service.EntryService
	eMapper *mapper.EntryMapper
}

func newBaseEntryController(eServ *service.EntryService) *baseEntryController {
	return &baseEntryController{
		eServ: eServ,
		eMapper: mapper.NewEntryMapper(),
	}
}

func (c *baseEntryController) getEntry(ctx context.Context, entryId int, userId int) (*model.Entry,
	error) {
	entry, err := c.eServ.GetEntryByIdAndUserId(ctx, entryId, userId)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		err := e.NewError(e.LogicEntryNotFound, fmt.Sprintf("Could not find entry %d.", entryId))
		log.Debug(err.StackTrace())
		return nil, err
	}
	return entry, nil
}

func (c *baseEntryController) getEntryMasterData(ctx context.Context, entryTypeId int,
	) ([]*model.EntryType, []*model.EntryActivity, error) {
	entryTypes, err := c.getEntryTypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	entryActivities, err := c.getEntryActivities(ctx, entryTypeId)
	if err != nil {
		return nil, nil, err
	}
	return entryTypes, entryActivities, nil
}

func (c *baseEntryController) getEntryTypes(ctx context.Context) ([]*model.EntryType, error) {
	return c.eServ.GetEntryTypes(ctx)
}

func (c *baseEntryController) getEntryActivities(ctx context.Context, entryTypeId int,
) ([]*model.EntryActivity, error) {
	if entryTypeId != model.EntryTypeIdWork {
		return []*model.EntryActivity{}, nil
	}
	return c.eServ.GetEntryActivities(ctx)
}

func (c *baseEntryController) getEntryMasterDataMap(ctx context.Context) (map[int]*model.EntryType,
	map[int]*model.EntryActivity, error) {
	entryTypesMap, err := c.getEntryTypesMap(ctx)
	if err != nil {
		return nil, nil, err
	}
	entryActivitiesMap, err := c.getEntryActivitiesMap(ctx)
	if err != nil {
		return nil, nil, err
	}
	return entryTypesMap, entryActivitiesMap, nil
}

func (c *baseEntryController) getEntryTypesMap(ctx context.Context) (map[int]*model.EntryType,
	error) {
	return c.eServ.GetEntryTypesMap(ctx)
}

func (c *baseEntryController) getEntryActivitiesMap(ctx context.Context) (map[int]*model.EntryActivity,
	error) {
	return c.eServ.GetEntryActivitiesMap(ctx)
}

func (c *baseEntryController) getFilterDetailsViewData(ctx context.Context, filter model.EntryFilter) (
	vm.EntryFilterDetails, error) {
	if filter == nil {
		return nil, nil
	}

	switch f := filter.(type) {
	case *model.TextEntryFilter:
		return c.eMapper.CreateBasicEntryFilterDetailsViewModel(f), nil
	case *model.FieldEntryFilter:
		entryTypes, entryActivities, err := c.getEntryMasterData(ctx, f.TypeId)
		if err != nil {
			return nil, err
		}
		return c.eMapper.CreateAdvancedEntryFilterDetailsViewModel(f, entryTypes, entryActivities), nil
	default:
		return nil, nil
	}
}

// --- Entry Filter Helper ---

type entryFilterHelper struct {}

func (fh *entryFilterHelper) buildQueryString(filter model.EntryFilter) string {
	if filter == nil {
		return ""
	}

	switch f := filter.(type) {
	case *model.FieldEntryFilter:
		return fh.buildAdvancedQueryString(f)
	case *model.TextEntryFilter:
		return fh.buildBasicQueryString(f)
	default:
		return ""
	}
}

func (fh *entryFilterHelper) buildBasicQueryString(filter *model.TextEntryFilter) string {
	return fmt.Sprintf("txt:%s", fh.formatQueryText(filter.Text))
}

func (fh *entryFilterHelper) buildAdvancedQueryString(filter *model.FieldEntryFilter) string {
	var qps []string
	// Add parameter/value for entry type
	if filter.ByType {
		qps = append(qps, fmt.Sprintf("typ:%d", filter.TypeId))
	}
	// Add parameter/value for entry start/end time
	if filter.ByTime {
		qps = append(qps, fmt.Sprintf("tim:%s", fh.formatQueryDateRange(filter.StartTime,
			filter.EndTime)))
	}
	// Add parameter/value for entry activity
	if filter.ByActivity {
		qps = append(qps, fmt.Sprintf("act:%d", filter.ActivityId))
	}
	// Add parameter/value for entry project
	if filter.ByProject {
		qps = append(qps, fmt.Sprintf("prj:%s", fh.formatQueryText(filter.Project)))
	}
	// Add parameter/value for entry description
	if filter.ByDescription {
		qps = append(qps, fmt.Sprintf("des:%s", fh.formatQueryText(filter.Description)))
	}
	// Add parameter/value for entry labels
	if filter.ByLabel {
		qps = append(qps, fmt.Sprintf("lbl:%s", fh.formatQueryLabels(filter.Labels)))
	}
	return strings.Join(qps[:], "|")
}

func (fh *entryFilterHelper) parseQueryString(userId int, isAdvanced bool, query string) (
	model.EntryFilter, error) {
	if !isAdvanced {
		return fh.parseBasicQueryString(userId, query)
	}
	return fh.parseAdvancedQueryString(userId, query)
}

func (fh *entryFilterHelper) parseBasicQueryString(userId int, query string) (*model.TextEntryFilter,
	error) {
	filter := model.NewTextEntryFilter()
	filter.ByUser = true
	filter.UserId = userId

	if query == "" {
		return filter, nil
	}

	pErr := fh.parseQueryParts(query, func(p string, v string) error {
		var cErr error

		// Handle specific conversion
		switch p {
		// Convert value for text
		case "txt":
			filter.Text, cErr = fh.parseQueryText(v)
		// Unknown parameter
		default:
			cErr = e.NewError(e.ValQueryInvalid, fmt.Sprintf("Unknown query parameter '%s'.", p))
		}

		return cErr
	})
	if pErr != nil {
		return nil, pErr
	}

	return filter, nil
}

func (fh *entryFilterHelper) parseAdvancedQueryString(userId int, query string,
	) (*model.FieldEntryFilter, error) {
	filter := model.NewFieldEntryFilter()
	filter.ByUser = true
	filter.UserId = userId

	if query == "" {
		return filter, nil
	}

	pErr := fh.parseQueryParts(query, func(p string, v string) error {
		var cErr error

		// Handle specific conversion
		switch p {
		// Convert value for entry type
		case "typ":
			filter.ByType = true
			filter.TypeId, cErr = strconv.Atoi(v)
		// Convert values for entry start/end time
		case "tim":
			filter.ByTime = true
			filter.StartTime, filter.EndTime, cErr = fh.parseQueryDateRange(v)
		// Convert value for entry activity
		case "act":
			filter.ByActivity = true
			filter.ActivityId, cErr = strconv.Atoi(v)
		// Convert value for entry project
		case "prj":
			filter.ByProject = true
			filter.Project, cErr = fh.parseQueryText(v)
		// Convert value for entry description
		case "des":
			filter.ByDescription = true
			filter.Description, cErr = fh.parseQueryText(v)
		// Convert value for entry labels
		case "lbl":
			filter.ByLabel = true
			filter.Labels, cErr = fh.parseQueryLabels(v)
		// Unknown parameter
		default:
			cErr = e.NewError(e.ValQueryInvalid, fmt.Sprintf("Query parameter '%s' is unknown.", p))
		}

		return cErr
	})
	if pErr != nil {
		return nil, pErr
	}

	return filter, nil
}

func (fh *entryFilterHelper) parseQueryParts(query string, partParserFunc func(string, string) error,
) error {
	qps := strings.Split(query, "|")

	for _, qp := range qps {
		pv := strings.Split(qp, ":")
		// Check if query part is invalid
		if len(pv) < 2 {
			err := e.NewError(e.ValQueryInvalid, "Query part is invalid.")
			log.Debug(err.StackTrace())
			return err
		}

		p := pv[0]
		v := pv[1]
		// Parse parameter/value
		cErr := partParserFunc(p, v)

		// Check if a error occurred
		if cErr != nil {
			err := e.WrapError(e.ValQueryInvalid, fmt.Sprintf("Query parameter '%s' has invalid "+
				"value.", p), cErr)
			log.Debug(err.StackTrace())
			return err
		}
	}

	return nil
}

func (fh *entryFilterHelper) formatQueryDateRange(startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf("%s-%s", fh.formatQueryDate(startDate), fh.formatQueryDate(endDate))
}

func (fh *entryFilterHelper) parseQueryDateRange(dateRange string) (time.Time, time.Time, error) {
	se := strings.Split(dateRange, "-")
	if len(se) < 2 {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	startTime, err := fh.parseQueryDate(se[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err := fh.parseQueryDate(se[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startTime, endTime, nil
}

func (fh *entryFilterHelper) formatQueryDate(date time.Time) string {
	return date.Format(queryDateTimeFormat)
}

func (fh *entryFilterHelper) parseQueryDate(date string) (time.Time, error) {
	return time.ParseInLocation(queryDateTimeFormat, date, time.Local)
}

func (fh *entryFilterHelper) formatQueryLabels(labels []string) string {
	labelsStr := strings.Join(labels, ",")
	return util.EncodeBase64(labelsStr)
}

func (fh *entryFilterHelper) parseQueryLabels(labels string) ([]string, error) {
	labelsStr, err := util.DecodeBase64(labels)
	if err != nil {
		return nil, err
	}
	if labelsStr == "" {
		return []string{}, nil
	}
	return strings.Split(labelsStr, ","), nil
}

func (fh *entryFilterHelper) formatQueryText(text string) string {
	return util.EncodeBase64(text)
}

func (fh *entryFilterHelper) parseQueryText(text string) (string, error) {
	return util.DecodeBase64(text)
}

func (fh *entryFilterHelper) isFilterEmpty(filter model.EntryFilter) bool {
	if filter == nil {
		return true
	}

	switch f := filter.(type) {
	case *model.TextEntryFilter:
		return false
	case *model.FieldEntryFilter:
		return !f.ByType && !f.ByTime && !f.ByActivity && !f.ByProject && !f.ByDescription &&
			!f.ByLabel
	default:
		return true
	}
}
