package validator

import (
	vm "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/pkg/model"
)

// --- Entry activity API model valdidation functions ---

// ValidateCreateEntryActivity validates information of a CreateEntryActivity API model.
func ValidateCreateEntryActivity(data *vm.CreateEntryActivity) error {
	return checkEntryActivityDescription(data.Description)
}

// ValidateUpdateEntryActivity validates information of a UpdateEntryActivity API model.
func ValidateUpdateEntryActivity(data *vm.UpdateEntryActivity) error {
	return checkEntryActivityDescription(data.Description)
}

// --- Basic entry activity validation functions ---

func checkEntryActivityDescription(desc string) error {
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
func ValidateCreateEntry(data *vm.CreateEntry) error {
	if err := checkEntryUserId(data.UserId); err != nil {
		return err
	}
	if err := checkEntryStartTime(data.StartTime); err != nil {
		return err
	}
	if err := checkEntryEndTime(data.EndTime); err != nil {
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
func ValidateUpdateEntry(data *vm.UpdateEntry) error {
	if err := checkEntryUserId(data.UserId); err != nil {
		return err
	}
	if err := checkEntryStartTime(data.StartTime); err != nil {
		return err
	}
	if err := checkEntryEndTime(data.EndTime); err != nil {
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

func checkEntryUserId(id int) error {
	return checkIdValid("userId", id)
}

func checkEntryStartTime(timestamp string) error {
	return checkTimestampValid("startTime", timestamp)
}

func checkEntryEndTime(timestamp string) error {
	return checkTimestampValid("startTime", timestamp)
}

func checkEntryTypeId(id int) error {
	return checkIdValid("typeId", id)
}

func checkEntryActivityId(id int) error {
	return checkIdValid("activityId", id)
}

func checkEntryDescription(desc string) error {
	if err := checkStringNotTooLong("description", desc, m.MaxLengthEntryDescription); err != nil {
		return err
	}
	return nil
}
