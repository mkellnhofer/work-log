package validator

import (
	"fmt"
	"strings"
	"time"

	"kellnhofer.com/work-log/constant"
	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

func checkIdValid(name string, id int) *e.Error {
	if id <= int(0) {
		err := e.NewError(e.ValIdInvalid, fmt.Sprintf("'%s' must be positive.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkIntNotNegative(name string, num int) *e.Error {
	if num < 0 {
		err := e.NewError(e.ValNumberNegative, fmt.Sprintf("'%s' must be zero or positive.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkFloatNotNegative(name string, num float32) *e.Error {
	if num < 0 {
		err := e.NewError(e.ValNumberNegative, fmt.Sprintf("'%s' must be zero or positive.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkIntNotNegativeOrZero(name string, num float32) *e.Error {
	if num <= 0 {
		err := e.NewError(e.ValNumberNegativeOrZero, fmt.Sprintf("'%s' must be positive.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkFloatNotNegativeOrZero(name string, num float32) *e.Error {
	if num <= 0 {
		err := e.NewError(e.ValNumberNegativeOrZero, fmt.Sprintf("'%s' must be positive.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkStringNotEmpty(name string, str string) *e.Error {
	s := strings.TrimSpace(str)
	if len(s) == 0 {
		err := e.NewError(e.ValStringEmpty, fmt.Sprintf("'%s' must not be empty.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkStringNotTooLong(name string, str string, length int) *e.Error {
	if len(str) >= length {
		err := e.NewError(e.ValStringTooLong, fmt.Sprintf("'%s' must not be longer than %d.", name,
			length))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkNotNil(name string, obj interface{}) *e.Error {
	if obj == nil {
		err := e.NewError(e.ValFieldNil, fmt.Sprintf("'%s' must not be null.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkArrayLengthNotZero(name string, length int) *e.Error {
	if length == 0 {
		err := e.NewError(e.ValArrayEmpty, fmt.Sprintf("'%s' must not be empty.", name))
		log.Debug(err.StackTrace())
		return err
	}
	return nil
}

func checkStringArrayNotEmpty(name string, strs []string) *e.Error {
	for _, str := range strs {
		s := strings.TrimSpace(str)
		if len(s) == 0 {
			err := e.NewError(e.ValStringEmpty, fmt.Sprintf("Elements of '%s' must not be empty.",
				name))
			log.Debug(err.StackTrace())
			return err
		}
	}
	return nil
}

func checkStringArrayNotTooLong(name string, strs []string, length int) *e.Error {
	for _, str := range strs {
		if len(str) >= length {
			err := e.NewError(e.ValStringTooLong, fmt.Sprintf("Elements of '%s' must not be longer "+
				"than %d.", name, length))
			log.Debug(err.StackTrace())
			return err
		}
	}
	return nil
}

func checkDateValid(name string, date string) *e.Error {
	_, pErr := time.Parse(constant.ApiDateFormat, date)
	if pErr != nil {
		err := e.WrapError(e.ValDateInvalid, fmt.Sprintf("'%s' must have format 'YYYY-MM-DD'.",
			name), pErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}

func checkTimestampValid(name string, timestamp string) *e.Error {
	_, pErr := time.Parse(constant.ApiTimestampFormat, timestamp)
	if pErr != nil {
		err := e.WrapError(e.ValTimestampInvalid, fmt.Sprintf("'%s' must have format "+
			"'YYYY-MM-DDTHH:mm:ss'.", name), pErr)
		log.Error(err.StackTrace())
		return err
	}
	return nil
}
