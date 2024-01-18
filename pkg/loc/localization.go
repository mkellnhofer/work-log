package loc

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	e "kellnhofer.com/work-log/pkg/error"
	"kellnhofer.com/work-log/pkg/log"
)

// Localization stores messages for a specific language.
type Localization struct {
	Language string     `xml:"language,attr"`
	Messages []*Message `xml:"message"`
}

// Message stores properties of localizable string.
type Message struct {
	Key  string `xml:"key,attr"`
	Text string `xml:"text"`
}

// LngTag holds the language tag of the configured localization.
var LngTag language.Tag

// LoadLocalization loads a localization.
func LoadLocalization(lang string) {
	fileName := fmt.Sprintf("localization-%s.xml", lang)

	// Open file
	file, ofErr := os.Open("config/localizations/" + fileName)
	if ofErr != nil {
		log.Fatalf("Could not read localization file '%s': %s", fileName, ofErr)
	}
	defer file.Close()

	// Read file
	byteValue, rfErr := ioutil.ReadAll(file)
	if rfErr != nil {
		log.Fatalf("Could not read localization file '%s': %s", fileName, rfErr)
	}

	// Parse file
	var loc Localization
	pfErr := xml.Unmarshal(byteValue, &loc)
	if pfErr != nil {
		log.Fatalf("Could not read localization file '%s': %s", fileName, pfErr)
	}

	// Parse language tag
	var pltErr error
	LngTag, pltErr = language.Parse(loc.Language)
	if pltErr != nil {
		log.Fatalf("Invalid language tag '%s': %s", loc.Language, pltErr)
	}

	// Register messages
	for _, m := range loc.Messages {
		if m.Key == "" {
			log.Fatalf("Message missing key!")
		}
		rmErr := message.SetString(LngTag, m.Key, m.Text)
		if rmErr != nil {
			log.Fatalf("Invalid message for key '%s': %s", m.Key, rmErr)
		}
	}
}

// CreateString creates a localized string.
func CreateString(key string, args ...interface{}) string {
	printer := message.NewPrinter(LngTag)
	return printer.Sprintf(key, args...)
}

var errorMessageKeys = map[int]string{
	// Authentication errors
	e.AuthUnknown:            "errAuthUnknown",
	e.AuthCredentialsInvalid: "errAuthCredentialsInvalid",

	// Permission errors
	e.PermUnknown:             "errPermUnknown",
	e.PermGetUserData:         "errPermMissing",
	e.PermChangeUserData:      "errPermMissing",
	e.PermGetUserAccount:      "errPermMissing",
	e.PermChangeUserAccount:   "errPermMissing",
	e.PermGetEntryCharacts:    "errPermMissing",
	e.PermChangeEntryCharacts: "errPermMissing",
	e.PermGetAllEntries:       "errPermMissing",
	e.PermChangeAllEntries:    "errPermMissing",
	e.PermGetOwnEntries:       "errPermMissing",
	e.PermChangeOwnEntries:    "errPermMissing",

	// Validation erros
	e.ValUnknown:              "errValUnknown",
	e.ValPageNumberInvalid:    "errValPageNumberInvalid",
	e.ValIdInvalid:            "errValIdInvalid",
	e.ValDateInvalid:          "errValDateInvalid",
	e.ValStartDateInvalid:     "errValStartDateInvalid",
	e.ValEndDateInvalid:       "errValEndDateInvalid",
	e.ValStartTimeInvalid:     "errValStartTimeInvalid",
	e.ValEndTimeInvalid:       "errValEndTimeInvalid",
	e.ValBreakDurationInvalid: "errValBreakDurationInvalid",
	e.ValDescriptionTooLong:   "errValDescriptionTooLong",
	e.ValSearchInvalid:        "errValSearchInvalid",
	e.ValSearchQueryInvalid:   "errValSearchQueryInvalid",
	e.ValMonthInvalid:         "errValMonthInvalid",
	e.ValPasswordEmpty:        "errValPasswordEmpty",
	e.ValPasswordTooShort:     "errValPasswordTooShort",
	e.ValPasswordTooLong:      "errValPasswordTooLong",
	e.ValPasswordInvalid:      "errValPasswordInvalid",
	e.ValPasswordsNotMatching: "errValPasswordsNotMatching",

	// Logic errors
	e.LogicUnknown:                        "errLogicUnknown",
	e.LogicEntryNotFound:                  "errLogicEntryNotFound",
	e.LogicEntryTypeNotFound:              "errLogicEntryTypeNotFound",
	e.LogicEntryActivityNotFound:          "errLogicEntryActivityNotFound",
	e.LogicEntryTimeIntervalInvalid:       "errLogicEntryTimeIntervalInvalid",
	e.LogicEntryBreakDurationTooLong:      "errLogicEntryBreakDurationTooLong",
	e.LogicEntrySearchDateIntervalInvalid: "errLogicEntrySearchDateIntervalInvalid",

	// System errors
	e.SysUnknown:             "errSysUnknown",
	e.SysDbUnknown:           "errSysDbUnknown",
	e.SysDbConnectionFailed:  "errSysDbConnectionFailed",
	e.SysDbTransactionFailed: "errSysDbTransactionFailed",
	e.SysDbQueryFailed:       "errSysDbQueryFailed",
	e.SysDbInsertFailed:      "errSysDbInsertFailed",
	e.SysDbUpdateFailed:      "errSysDbUpdateFailed",
	e.SysDbDeleteFailed:      "errSysDbDeleteFailed",
}

// GetErrorMessageString returns a localized error message string.
func GetErrorMessageString(errorCode int) string {
	emk, ok := errorMessageKeys[errorCode]
	if !ok {
		log.Debugf("Unexpected error code %d. Using fallback error message.", errorCode)
		emk, _ = errorMessageKeys[e.SysUnknown]
	}
	printer := message.NewPrinter(LngTag)
	return printer.Sprintf(emk)
}
