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

type handlerFunc func(eCtx echo.Context, ctx context.Context) error
type hxHandlerFunc func(eCtx echo.Context, ctx context.Context) error
type resourceHandlerFunc func(eCtx echo.Context, ctx context.Context) error

type baseController struct {
	uServ *service.UserService
	eServ *service.EntryService

	mapper *mapper.Mapper
}

func (c *baseController) handler(hf handlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		ctx := getContext(eCtx)
		return hf(eCtx, ctx)
	}
}

func (c *baseController) hxHandler(hf hxHandlerFunc) echo.HandlerFunc {
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

func (c *baseController) resourceHandler(hf resourceHandlerFunc) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return hf(eCtx, getContext(eCtx))
	}
}

func (c *baseController) getUser(ctx context.Context, userId int) (*model.User, error) {
	return c.uServ.GetUserById(ctx, userId)
}

func (c *baseController) getUserContract(ctx context.Context, userId int) (*model.Contract, error) {
	return c.uServ.GetUserContractByUserId(ctx, userId)
}

func (c *baseController) getUserInfoViewData(ctx context.Context) (*vm.UserInfo, error) {
	userId := getCurrentUserId(ctx)
	user, err := c.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	return c.mapper.CreateUserInfoViewModel(user), nil
}

func (c *baseController) getEntry(ctx context.Context, entryId int, userId int) (*model.Entry,
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

func (c *baseController) getEntryMasterData(ctx context.Context, entryTypeId int) ([]*model.EntryType,
	[]*model.EntryActivity, error) {
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

func (c *baseController) getEntryTypes(ctx context.Context) ([]*model.EntryType, error) {
	return c.eServ.GetEntryTypes(ctx)
}

func (c *baseController) getEntryActivities(ctx context.Context, entryTypeId int,
) ([]*model.EntryActivity, error) {
	if entryTypeId != model.EntryTypeIdWork {
		return []*model.EntryActivity{}, nil
	}
	return c.eServ.GetEntryActivities(ctx)
}

func (c *baseController) getEntryMasterDataMap(ctx context.Context) (map[int]*model.EntryType,
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

func (c *baseController) getEntryTypesMap(ctx context.Context) (map[int]*model.EntryType, error) {
	return c.eServ.GetEntryTypesMap(ctx)
}

func (c *baseController) getEntryActivitiesMap(ctx context.Context) (map[int]*model.EntryActivity,
	error) {
	return c.eServ.GetEntryActivitiesMap(ctx)
}

func (c *baseController) buildQueryString(filter model.EntryFilter) string {
	if filter == nil {
		return ""
	}

	switch f := filter.(type) {
	case *model.FieldEntryFilter:
		return c.buildAdvancedQueryString(f)
	case *model.TextEntryFilter:
		return c.buildBasicQueryString(f)
	default:
		return ""
	}
}

func (c *baseController) buildBasicQueryString(filter *model.TextEntryFilter) string {
	return fmt.Sprintf("txt:%s", c.formatQueryText(filter.Text))
}

func (c *baseController) buildAdvancedQueryString(filter *model.FieldEntryFilter) string {
	var qps []string
	// Add parameter/value for entry type
	if filter.ByType {
		qps = append(qps, fmt.Sprintf("typ:%d", filter.TypeId))
	}
	// Add parameter/value for entry start/end time
	if filter.ByTime {
		qps = append(qps, fmt.Sprintf("tim:%s", c.formatQueryDateRange(filter.StartTime,
			filter.EndTime)))
	}
	// Add parameter/value for entry activity
	if filter.ByActivity {
		qps = append(qps, fmt.Sprintf("act:%d", filter.ActivityId))
	}
	// Add parameter/value for entry description
	if filter.ByDescription {
		qps = append(qps, fmt.Sprintf("des:%s", c.formatQueryText(filter.Description)))
	}
	// Add parameter/value for entry labels
	if filter.ByLabel {
		qps = append(qps, fmt.Sprintf("lbl:%s", c.formatQueryLabels(filter.Labels)))
	}
	return strings.Join(qps[:], "|")
}

func (c *baseController) parseQueryString(userId int, isAdvanced bool, query string) (
	model.EntryFilter, error) {
	if !isAdvanced {
		return c.parseBasicQueryString(userId, query)
	}
	return c.parseAdvancedQueryString(userId, query)
}

func (c *baseController) parseBasicQueryString(userId int, query string) (*model.TextEntryFilter,
	error) {
	filter := model.NewTextEntryFilter()
	filter.ByUser = true
	filter.UserId = userId

	if query == "" {
		return filter, nil
	}

	pErr := c.parseQueryParts(query, func(p string, v string) error {
		var cErr error

		// Handle specific conversion
		switch p {
		// Convert value for text
		case "txt":
			filter.Text, cErr = c.parseQueryText(v)
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

func (c *baseController) parseAdvancedQueryString(userId int, query string) (*model.FieldEntryFilter,
	error) {
	filter := model.NewFieldEntryFilter()
	filter.ByUser = true
	filter.UserId = userId

	if query == "" {
		return filter, nil
	}

	pErr := c.parseQueryParts(query, func(p string, v string) error {
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
			filter.StartTime, filter.EndTime, cErr = c.parseQueryDateRange(v)
		// Convert value for entry activity
		case "act":
			filter.ByActivity = true
			filter.ActivityId, cErr = strconv.Atoi(v)
		// Convert value for entry description
		case "des":
			filter.ByDescription = true
			filter.Description, cErr = c.parseQueryText(v)
		// Convert value for entry labels
		case "lbl":
			filter.ByLabel = true
			filter.Labels, cErr = c.parseQueryLabels(v)
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

func (c *baseController) parseQueryParts(query string, partParserFunc func(string, string) error,
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

func (c *baseController) formatQueryDateRange(startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf("%s-%s", c.formatQueryDate(startDate), c.formatQueryDate(endDate))
}

func (c *baseController) parseQueryDateRange(dateRange string) (time.Time, time.Time, error) {
	se := strings.Split(dateRange, "-")
	if len(se) < 2 {
		return time.Time{}, time.Time{}, errors.New("invalid range")
	}
	startTime, err := c.parseQueryDate(se[0])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endTime, err := c.parseQueryDate(se[1])
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return startTime, endTime, nil
}

func (c *baseController) formatQueryDate(date time.Time) string {
	return date.Format(queryDateTimeFormat)
}

func (c *baseController) parseQueryDate(date string) (time.Time, error) {
	return time.ParseInLocation(queryDateTimeFormat, date, time.Local)
}

func (c *baseController) formatQueryLabels(labels []string) string {
	labelsStr := strings.Join(labels, ",")
	return util.EncodeBase64(labelsStr)
}

func (c *baseController) parseQueryLabels(labels string) ([]string, error) {
	labelsStr, err := util.DecodeBase64(labels)
	if err != nil {
		return nil, err
	}
	if labelsStr == "" {
		return []string{}, nil
	}
	return strings.Split(labelsStr, ","), nil
}

func (c *baseController) formatQueryText(text string) string {
	return util.EncodeBase64(text)
}

func (c *baseController) parseQueryText(text string) (string, error) {
	return util.DecodeBase64(text)
}

func (c *baseController) isFilterEmpty(filter model.EntryFilter) bool {
	if filter == nil {
		return true
	}

	switch f := filter.(type) {
	case *model.TextEntryFilter:
		return false
	case *model.FieldEntryFilter:
		return !f.ByType && !f.ByTime && !f.ByActivity && !f.ByDescription && !f.ByLabel
	default:
		return true
	}
}

func (c *baseController) getFilterDetailsViewData(ctx context.Context, filter model.EntryFilter) (
	vm.EntryFilterDetails, error) {
	if filter == nil {
		return nil, nil
	}

	switch f := filter.(type) {
	case *model.TextEntryFilter:
		return c.mapper.CreateBasicEntryFilterDetailsViewModel(f), nil
	case *model.FieldEntryFilter:
		entryTypes, entryActivities, err := c.getEntryMasterData(ctx, f.TypeId)
		if err != nil {
			return nil, err
		}
		return c.mapper.CreateAdvancedEntryFilterDetailsViewModel(f, entryTypes, entryActivities), nil
	default:
		return nil, nil
	}
}
