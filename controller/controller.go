package controller

import (
	"net/http"
	"strconv"

	e "kellnhofer.com/work-log/error"
	"kellnhofer.com/work-log/log"
)

var errorMessages = map[int]string{
	// Authentication errors
	e.AuthUnknown: "Ein unbekannter Authentifizierungsfehler trat auf.",

	// Validation erros
	e.ValUnknown: "Ein unbekannter Validierungsfehler trat auf.",

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

func getStringQueryParam(r *http.Request, n string) string {
	qvs := r.URL.Query()
	return qvs.Get(n)
}

func getIntQueryParam(r *http.Request, n string) (int, error) {
	qv := getStringQueryParam(r, n)
	if qv == "" {
		return 0, nil
	}
	return strconv.Atoi(qv)
}
