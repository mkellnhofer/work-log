package repo

import (
	"context"
	"database/sql"
	"fmt"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
	"kellnhofer.com/work-log/pkg/model"
)

type dbContract struct {
	initOvertimeHours float32
	initVacationDays  float32
	firstDay          string
}

type dbContractWorkingHours struct {
	firstDay   string
	dailyHours float32
}

type dbContractVacationDays struct {
	firstDay    string
	monthlyDays float32
}

// ContractRepo retrieves and stores contract related entities.
type ContractRepo struct {
	repo
}

// NewContractRepo creates a new contract repository.
func NewContractRepo(db *sql.DB) *ContractRepo {
	return &ContractRepo{repo{db}}
}

// --- Contract functions ---

// GetContractByUserId retrieves the contract information of a user by its ID.
func (r *ContractRepo) GetContractByUserId(ctx context.Context, userId int) (*model.Contract,
	*e.Error) {
	c, qErr := r.getContract(ctx, userId)
	if qErr != nil {
		return nil, qErr
	}

	cvd, qErr := r.getContractVacationDays(ctx, userId)
	if qErr != nil {
		return nil, qErr
	}
	c.VacationDays = cvd

	cwh, qErr := r.getContractWorkingHours(ctx, userId)
	if qErr != nil {
		return nil, qErr
	}
	c.WorkingHours = cwh

	return c, nil
}

// CreateContract creates the contract information of a user.
func (r *ContractRepo) CreateContract(ctx context.Context, userId int, contract *model.Contract,
) *e.Error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) *e.Error {
		if err := r.createContract(tx, userId, contract); err != nil {
			return err
		}
		if err := r.setContractVacationDays(tx, userId, contract.VacationDays); err != nil {
			return err
		}
		if err := r.setContractWorkingHours(tx, userId, contract.WorkingHours); err != nil {
			return err
		}
		return nil
	})
}

// UpdateContract updates the contract information of a user.
func (r *ContractRepo) UpdateContract(ctx context.Context, userId int, contract *model.Contract,
) *e.Error {
	return r.executeInTransaction(ctx, func(tx *sql.Tx) *e.Error {
		if err := r.updateContract(tx, userId, contract); err != nil {
			return err
		}
		if err := r.setContractVacationDays(tx, userId, contract.VacationDays); err != nil {
			return err
		}
		if err := r.setContractWorkingHours(tx, userId, contract.WorkingHours); err != nil {
			return err
		}
		return nil
	})
}

func (r *ContractRepo) getContract(ctx context.Context, userId int) (*model.Contract, *e.Error) {
	q := "SELECT init_overtime_hours, init_vacation_days, first_day FROM contract " +
		"WHERE user_id = ?"

	sr, qErr := r.queryRow(ctx, &scanContractHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read contract for user %d "+
			"from database.", userId), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}
	if sr == nil {
		return nil, nil
	}

	return sr.(*model.Contract), nil
}

func (r *ContractRepo) createContract(tx *sql.Tx, userId int, contract *model.Contract,
) *e.Error {
	c := toDbContract(contract)

	q := "INSERT INTO contract (user_id, init_overtime_hours, init_vacation_days, first_day) " +
		"VALUES (?, ?, ?, ?)"

	_, cErr := r.insertWithTx(tx, q, userId, c.initOvertimeHours, c.initVacationDays, c.firstDay)
	if cErr != nil {
		err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not create contract for user %d "+
			"in database.", userId), cErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

func (r *ContractRepo) updateContract(tx *sql.Tx, userId int, contract *model.Contract,
) *e.Error {
	c := toDbContract(contract)

	q := "UPDATE contract SET init_overtime_hours = ?, init_vacation_days = ?, first_day = ? " +
		"WHERE user_id = ?"

	uErr := r.execWithTx(tx, q, c.initOvertimeHours, c.initVacationDays, c.firstDay, userId)
	if uErr != nil {
		err := e.WrapError(e.SysDbUpdateFailed, fmt.Sprintf("Could not update contract for user %d "+
			"in database.", userId), uErr)
		log.Error(err.StackTrace())
		return err
	}

	return nil
}

func (r *ContractRepo) getContractVacationDays(ctx context.Context, userId int,
) ([]model.ContractVacationDays, *e.Error) {
	q := "SELECT first_day, monthly_days FROM contract_vacation_days WHERE user_id = ?"

	sr, qErr := r.query(ctx, &scanContractVacationDaysHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read contract for user %d "+
			"from database.", userId), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]model.ContractVacationDays), nil
}

func (r *ContractRepo) setContractVacationDays(tx *sql.Tx, userId int,
	vacationDays []model.ContractVacationDays) *e.Error {
	dErr := r.execWithTx(tx, "DELETE FROM contract_vacation_days WHERE user_id = ?", userId)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not update contract for user "+
			"%d in database.", userId), dErr)
		log.Error(err.StackTrace())
		return err
	}

	for _, vd := range vacationDays {
		cvd := toDbContractVacationDays(vd)

		cErr := r.execWithTx(tx, "INSERT INTO contract_vacation_days (user_id, first_day, "+
			"monthly_days) VALUES (?, ?, ?)", userId, cvd.firstDay, cvd.monthlyDays)
		if cErr != nil {
			err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not update contract for "+
				"user %d in database.", userId), cErr)
			log.Error(err.StackTrace())
			return err
		}
	}

	return nil
}

func (r *ContractRepo) getContractWorkingHours(ctx context.Context, userId int,
) ([]model.ContractWorkingHours, *e.Error) {
	q := "SELECT first_day, daily_hours FROM contract_working_hours WHERE user_id = ?"

	sr, qErr := r.query(ctx, &scanContractWorkingHoursHelper{}, q, userId)
	if qErr != nil {
		err := e.WrapError(e.SysDbQueryFailed, fmt.Sprintf("Could not read contract for user %d "+
			"from database.", userId), qErr)
		log.Error(err.StackTrace())
		return nil, err
	}

	return sr.([]model.ContractWorkingHours), nil
}

func (r *ContractRepo) setContractWorkingHours(tx *sql.Tx, userId int,
	workingHours []model.ContractWorkingHours) *e.Error {
	dErr := r.execWithTx(tx, "DELETE FROM contract_working_hours WHERE user_id = ?", userId)
	if dErr != nil {
		err := e.WrapError(e.SysDbDeleteFailed, fmt.Sprintf("Could not update contract for user "+
			"%d in database.", userId), dErr)
		log.Error(err.StackTrace())
		return err
	}

	for _, wh := range workingHours {
		cwh := toDbContractWorkingHours(wh)

		cErr := r.execWithTx(tx, "INSERT INTO contract_working_hours (user_id, first_day, "+
			"daily_hours) VALUES (?, ?, ?)", userId, cwh.firstDay, cwh.dailyHours)
		if cErr != nil {
			err := e.WrapError(e.SysDbInsertFailed, fmt.Sprintf("Could not update contract for "+
				"user %d in database.", userId), cErr)
			log.Error(err.StackTrace())
			return err
		}
	}

	return nil
}

// --- Helper functions ---

type scanContractHelper struct {
}

func (h *scanContractHelper) makeSlice() interface{} {
	return make([]*model.Contract, 0, 10)
}

func (h *scanContractHelper) scan(s scanner) (interface{}, error) {
	var dbC dbContract

	err := s.Scan(&dbC.initOvertimeHours, &dbC.initVacationDays, &dbC.firstDay)
	if err != nil {
		return nil, err
	}

	c := fromDbContract(&dbC)

	return c, nil
}

func (h *scanContractHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]*model.Contract), item.(*model.Contract))
}

func toDbContract(in *model.Contract) *dbContract {
	var out dbContract
	out.firstDay = *formatDate(&in.FirstDay)
	out.initOvertimeHours = in.InitOvertimeHours
	out.initVacationDays = in.InitVacationDays
	return &out
}

func fromDbContract(in *dbContract) *model.Contract {
	var out model.Contract
	out.FirstDay = *parseDate(&in.firstDay)
	out.InitOvertimeHours = in.initOvertimeHours
	out.InitVacationDays = in.initVacationDays
	return &out
}

type scanContractVacationDaysHelper struct {
}

func (h *scanContractVacationDaysHelper) makeSlice() interface{} {
	return make([]model.ContractVacationDays, 0, 10)
}

func (h *scanContractVacationDaysHelper) scan(s scanner) (interface{}, error) {
	var dbC dbContractVacationDays

	err := s.Scan(&dbC.firstDay, &dbC.monthlyDays)
	if err != nil {
		return nil, err
	}

	c := fromDbContractVacationDays(dbC)

	return c, nil
}

func (h *scanContractVacationDaysHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]model.ContractVacationDays), item.(model.ContractVacationDays))
}

func toDbContractVacationDays(in model.ContractVacationDays) dbContractVacationDays {
	var out dbContractVacationDays
	out.firstDay = *formatDate(&in.FirstDay)
	out.monthlyDays = in.Days
	return out
}

func fromDbContractVacationDays(in dbContractVacationDays) model.ContractVacationDays {
	var out model.ContractVacationDays
	out.FirstDay = *parseDate(&in.firstDay)
	out.Days = in.monthlyDays
	return out
}

type scanContractWorkingHoursHelper struct {
}

func (h *scanContractWorkingHoursHelper) makeSlice() interface{} {
	return make([]model.ContractWorkingHours, 0, 10)
}

func (h *scanContractWorkingHoursHelper) scan(s scanner) (interface{}, error) {
	var dbC dbContractWorkingHours

	err := s.Scan(&dbC.firstDay, &dbC.dailyHours)
	if err != nil {
		return nil, err
	}

	c := fromDbContractWorkingHours(dbC)

	return c, nil
}

func (h *scanContractWorkingHoursHelper) appendSlice(items interface{}, item interface{}) interface{} {
	return append(items.([]model.ContractWorkingHours), item.(model.ContractWorkingHours))
}

func toDbContractWorkingHours(in model.ContractWorkingHours) dbContractWorkingHours {
	var out dbContractWorkingHours
	out.firstDay = *formatDate(&in.FirstDay)
	out.dailyHours = in.Hours
	return out
}

func fromDbContractWorkingHours(in dbContractWorkingHours) model.ContractWorkingHours {
	var out model.ContractWorkingHours
	out.FirstDay = *parseDate(&in.firstDay)
	out.Hours = in.dailyHours
	return out
}
