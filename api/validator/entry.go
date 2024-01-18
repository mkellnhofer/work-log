package validator

import (
	vm "kellnhofer.com/work-log/api/model"
	e "kellnhofer.com/work-log/pkg/error"
	m "kellnhofer.com/work-log/pkg/model"
)

// --- Entry activity API model valdidation functions ---

// ValidateCreateEntryActivity validates information of a CreateEntryActivity API model.
func ValidateCreateEntryActivity(data *vm.CreateEntryActivity) *e.Error {
	return checkEntryActivityDescription(data.Description)
}

// ValidateUpdateEntryActivity validates information of a UpdateEntryActivity API model.
func ValidateUpdateEntryActivity(data *vm.UpdateEntryActivity) *e.Error {
	return checkEntryActivityDescription(data.Description)
}

// --- Basic entry activity validation functions ---

func checkEntryActivityDescription(desc string) *e.Error {
	if err := checkStringNotEmpty("description", desc); err != nil {
		return err
	}
	if err := checkStringNotTooLong("description", desc, m.MaxLengthEntryActivityDescription); err !=
		nil {
		return err
	}
	return nil
}

// --- Entry API model valdidation functions ---

// ValidateCreateEntry validates information of a CreateEntryA API model.
func ValidateCreateEntry(data *vm.CreateEntry) *e.Error {
	if err := checkEntryUserId(data.UserId); err != nil {
		return err
	}
	if err := checkEntryStartTime(data.StartTime); err != nil {
		return err
	}
	if err := checkEntryEndTime(data.EndTime); err != nil {
		return err
	}
	if err := checkEntryBreakDuration(data.BreakDuration); err != nil {
		return err
	}
	if err := checkEntryTypeId(data.TypeId); err != nil {
		return err
	}
	if err := checkEntryActivityId(data.ActivityId); err != nil {
		return err
	}
	return checkEntryDescription(data.Description)
}

// ValidateUpdateEntry validates information of a UpdateEntry API model.
func ValidateUpdateEntry(data *vm.UpdateEntry) *e.Error {
	if err := checkEntryUserId(data.UserId); err != nil {
		return err
	}
	if err := checkEntryStartTime(data.StartTime); err != nil {
		return err
	}
	if err := checkEntryEndTime(data.EndTime); err != nil {
		return err
	}
	if err := checkEntryBreakDuration(data.BreakDuration); err != nil {
		return err
	}
	if err := checkEntryTypeId(data.TypeId); err != nil {
		return err
	}
	if err := checkEntryActivityId(data.ActivityId); err != nil {
		return err
	}
	return checkEntryDescription(data.Description)
}

// --- Basic entry validation functions ---

func checkEntryUserId(id int) *e.Error {
	return checkIdValid("userId", id)
}

func checkEntryStartTime(timestamp string) *e.Error {
	return checkTimestampValid("startTime", timestamp)
}

func checkEntryEndTime(timestamp string) *e.Error {
	return checkTimestampValid("startTime", timestamp)
}

func checkEntryBreakDuration(num int) *e.Error {
	return checkIntNotNegative("breakDuration", num)
}

func checkEntryTypeId(id int) *e.Error {
	return checkIdValid("typeId", id)
}

func checkEntryActivityId(id int) *e.Error {
	return checkIdValid("activityId", id)
}

func checkEntryDescription(desc string) *e.Error {
	if err := checkStringNotTooLong("description", desc, m.MaxLengthEntryDescription); err != nil {
		return err
	}
	return nil
}
