package service

import (
	"kellnhofer.com/work-log/db/tx"
)

type service struct {
	tm *tx.TransactionManager
}
