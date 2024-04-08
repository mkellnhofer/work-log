package controller

import (
	"context"
	"fmt"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
	"kellnhofer.com/work-log/pkg/service"
)

type baseController struct {
	uServ *service.UserService
	eServ *service.EntryService
}

func (c *baseController) getUserIdAndUserContract(ctx context.Context) (int, *model.Contract,
	error) {
	userId := getCurrentUserId(ctx)
	userContract, err := c.uServ.GetUserContractByUserId(ctx, userId)
	if err != nil {
		return 0, nil, err
	}
	return userId, userContract, nil
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

func (c *baseController) getEntryMasterData(ctx context.Context) ([]*model.EntryType,
	[]*model.EntryActivity, error) {
	entryTypes, err := c.getEntryTypes(ctx)
	if err != nil {
		return nil, nil, err
	}
	entryActivities, err := c.getEntryActivities(ctx)
	if err != nil {
		return nil, nil, err
	}
	return entryTypes, entryActivities, nil
}

func (c *baseController) getEntryTypes(ctx context.Context) ([]*model.EntryType, error) {
	return c.eServ.GetEntryTypes(ctx)
}

func (c *baseController) getEntryActivities(ctx context.Context) ([]*model.EntryActivity, error) {
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
