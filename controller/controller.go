package controller

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

var errorMessages = map[int]string{
	// Authentication errors
	e.AuthUnknown:            "Ein unbekannter Authentifizierungsfehler trat auf.",
	e.AuthInvalidCredentials: "Falscher Benutzername oder Passwort.",

	// Validation erros
	e.ValUnknown:                   "Ein unbekannter Validierungsfehler trat auf.",
	e.ValPageNumberInvalid:         "Seitennummer ungültig. (Seitennummer muss numerisch und positiv sein.)",
	e.ValIdInvalid:                 "ID ungültig. (ID muss numerisch und positiv sein.)",
	e.ValEntryDateInvalid:          "Datum ungültig!",
	e.ValEntryStartTimeInvalid:     "Startzeit ungültig!",
	e.ValEntryEndTimeInvalid:       "Endzeit ungültig!",
	e.ValEntryBreakDurationInvalid: "Pausendauer ungültig!",
	e.ValEntryDescriptionTooLong:   "Beschreibung darf nicht länger als 200 Zeichen sein!",

	// Logic errors
	e.LogicUnknown:                   "Ein unbekannter Logikfehler trat auf.",
	e.LogicEntryNotFound:             "Der Eintrag konnte nicht gefunden werden.",
	e.LogicEntryTypeNotFound:         "Der Eintragstyp konnte nicht gefunden werden.",
	e.LogicEntryActivityNotFound:     "Die Eintragstätigkeit konnte nicht gefunden werden.",
	e.LogicEntryTimeIntervalInvalid:  "Startzeit-Endzeit-Interval ungültig!",
	e.LogicEntryBreakDurationTooLong: "Pausendauer zu lang!",

	// System errors
	e.SysUnknown:             "Ein unbekannter Systemfehler trat auf.",
	e.SysDbUnknown:           "Ein unbekannter Datenbankfehler trat auf.",
	e.SysDbConnectionFailed:  "Die Verbindung zur Datenbank wurde unterbrochen.",
	e.SysDbTransactionFailed: "Die Datenbanktransaktion schlug fehl.",
	e.SysDbQueryFailed:       "Die Datenbankabfrage schlug fehl.",
	e.SysDbInsertFailed:      "Ein Datenbankeintrag konnte nicht erstellt werden.",
	e.SysDbUpdateFailed:      "Ein Datenbankeintrag konnte nicht geändert werden.",
	e.SysDbDeleteFailed:      "Ein Datenbankeintrag konnte nicht gelöscht werden.",
}

func getErrorMessage(errorCode int) string {
	em, ok := errorMessages[errorCode]
	if !ok {
		log.Debugf("Unexpected error code %d. Using fallback error message.", errorCode)
		return errorMessages[e.SysUnknown]
	}
	return em
}

func getStringPathVar(r *http.Request, n string) (string, bool) {
	vs := mux.Vars(r)
	v, ok := vs[n]
	return v, ok
}

func getIdPathVar(r *http.Request) int {
	v, ok := getStringPathVar(r, "id")
	if !ok {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Varible missing.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	id, pErr := strconv.Atoi(v)
	if pErr != nil {
		err := e.WrapError(e.ValIdInvalid, "Invalid ID. (Varible must be numeric.)", pErr)
		log.Debug(err.StackTrace())
		panic(err)
	}

	if id <= 0 {
		err := e.NewError(e.ValIdInvalid, "Invalid ID. (Varible must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return id
}

func getStringQueryParam(r *http.Request, n string) (string, bool) {
	qvs := r.URL.Query()
	qv := qvs.Get(n)
	if qv == "" {
		return qv, false
	}
	return qv, true
}

func getErrorCodeQueryParam(r *http.Request) *int {
	v, ok := getStringQueryParam(r, "error")
	if !ok {
		return nil
	}

	ec, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}

	return &ec
}

func getPageNumberQueryParam(r *http.Request) *int {
	v, ok := getStringQueryParam(r, "page")
	if !ok {
		return nil
	}

	page, err := strconv.Atoi(v)
	if err != nil {
		err := e.WrapError(e.ValPageNumberInvalid, "Invalid page number. (Varible must be numeric.)",
			err)
		log.Debug(err.StackTrace())
		panic(err)
	}

	if page <= 0 {
		err := e.NewError(e.ValPageNumberInvalid, "Invalid page number. (Varible must be positive.)")
		log.Debug(err.StackTrace())
		panic(err)
	}

	return &page
}
