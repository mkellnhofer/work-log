package validator

import (
	vm "kellnhofer.com/work-log/api/model"
	m "kellnhofer.com/work-log/pkg/model"
)

// ValidateCreateToken validates information of a CreateToken API model.
func ValidateCreateToken(data *vm.CreateToken) error {
	if err := checkStringNotEmpty("name", data.Name); err != nil {
		return err
	}
	if err := checkStringNotTooLong("name", data.Name, m.MaxLengthTokenName); err != nil {
		return err
	}
	return nil
}
