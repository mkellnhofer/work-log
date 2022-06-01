package db

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"kellnhofer.com/work-log/db/tx"

	_ "github.com/go-sql-driver/mysql"

	"kellnhofer.com/work-log/config"
	"kellnhofer.com/work-log/db/repo"
	"kellnhofer.com/work-log/log"
)

const curDbVers = 2

// Db abstracts the database access and provides repositories execute CRUD operations.
type Db struct {
	config *config.Config

	db    *sql.DB
	txm   *tx.TransactionManager
	uRepo *repo.UserRepo
	cRepo *repo.ContractRepo
	sRepo *repo.SessionRepo
	eRepo *repo.EntryRepo
}

// NewDb creates a new Db for the supplied configuration.
func NewDb(config *config.Config) *Db {
	return &Db{config, nil, nil, nil, nil, nil, nil}
}

// --- Public functions ---

// OpenDb opens the underlying database.
func (db *Db) OpenDb() {
	con := db.config.DbUsername + ":" + db.config.DbPassword +
		"@tcp(" + db.config.DbHost + ":" + strconv.Itoa(db.config.DbPort) + ")" +
		"/" + db.config.DbScheme

	var err error

	db.db, err = sql.Open("mysql", con)
	if err != nil {
		log.Fatalf("Could not open database connection!\nError: %s", err)
	}

	err = db.db.Ping()
	if err != nil {
		log.Fatalf("Could not open database connection!\nError: %s", err)
	}
}

// UpdateDb updates the underlying database to the current version.
func (db *Db) UpdateDb() {
	dbVers := getDbVersion(db.db)
	if dbVers == 0 {
		log.Info("Creating database ...")
		createDb(db.db)
		updateDb(db.db, 1)
		log.Info("Successfully created database.")
	} else if dbVers < curDbVers {
		log.Info("Updating database ...")
		updateDb(db.db, dbVers)
		log.Info("Successfully updated database.")
	}
}

// ClearDb deletes all tables from the underlying database.
func (db *Db) ClearDb() {
	clearDb(db.db)
}

// CloseDb closes the underlying database.
func (db *Db) CloseDb() {
	if db.db == nil {
		return
	}

	err := db.db.Close()
	if err != nil {
		log.Fatalf("Could not close database connection!\nError: %s", err)
	}
}

// GetTransactionManager provides the TransactionManager.
func (db *Db) GetTransactionManager() *tx.TransactionManager {
	if db.txm == nil {
		db.txm = tx.NewTransactionManager(db.db)
	}

	return db.txm
}

// GetUserRepo provides the UserRepo.
func (db *Db) GetUserRepo() *repo.UserRepo {
	if db.uRepo == nil {
		db.uRepo = repo.NewUserRepo(db.db)
	}

	return db.uRepo
}

// GetContractRepo provides the ContractRepo.
func (db *Db) GetContractRepo() *repo.ContractRepo {
	if db.cRepo == nil {
		db.cRepo = repo.NewContractRepo(db.db)
	}

	return db.cRepo
}

// GetSessionRepo provides the SessionRepo.
func (db *Db) GetSessionRepo() *repo.SessionRepo {
	if db.sRepo == nil {
		db.sRepo = repo.NewSessionRepo(db.db)
	}

	return db.sRepo
}

// GetEntryRepo provides the EntryRepo.
func (db *Db) GetEntryRepo() *repo.EntryRepo {
	if db.eRepo == nil {
		db.eRepo = repo.NewEntryRepo(db.db)
	}

	return db.eRepo
}

// --- Private functions ---

func getDbVersion(db *sql.DB) int {
	var name string
	row := db.QueryRow("SHOW TABLES LIKE 'setting'")
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Fatalf("Could not query database version! (Error: %s)", err)
	}

	var version int
	row = db.QueryRow("SELECT setting_value FROM setting WHERE setting_key LIKE 'db_version'")
	err = row.Scan(&version)
	if err != nil {
		log.Fatalf("Could not query database version! (Error: %s)", err)
	}

	return version
}

func createDb(db *sql.DB) {
	dbStmts := readDbFile("db_create.sql")
	executeDbStmts(db, dbStmts)
}

func updateDb(db *sql.DB, dbVers int) {
	for i := dbVers + 1; i <= curDbVers; i++ {
		fileName := fmt.Sprintf("db_update_v%d.sql", i)
		dbStmts := readDbFile(fileName)
		log.Infof("Executing database update v%d ...", i)
		executeDbStmts(db, dbStmts)
		updateDbVersion(db, i)
	}
}

func clearDb(db *sql.DB) {
	dbStmts := readDbFile("db_clear.sql")
	executeDbStmts(db, dbStmts)
}

func readDbFile(name string) []string {
	// Open file
	file, err := os.Open("scripts/db/" + name)
	if err != nil {
		log.Fatalf("Could not open database update script %s! (Error: %s)", name, err)
	}
	defer file.Close()

	var stmts []string

	// Read statements from file
	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()

		// If line is empty: Skip
		if line == "" {
			continue
		}

		// Add line to current statement
		buffer.WriteString(line)

		// If line is end of statement: Add current statement to result
		if strings.HasSuffix(line, ";") {
			stmts = append(stmts, buffer.String())
			buffer.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not read database update script %s! (Error: %s)", name, err)
	}

	return stmts
}

func executeDbStmts(db *sql.DB, stmts []string) {
	for _, stmt := range stmts {
		executeDbStmt(db, stmt)
	}
}

func executeDbStmt(db *sql.DB, stmt string) {
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatalf("Could not execute database statement! (Error: %s)", err)
	}
}

func updateDbVersion(db *sql.DB, dbVers int) {
	_, err := db.Exec("UPDATE setting SET setting_value = ? WHERE setting_key = 'db_version'", dbVers)
	if err != nil {
		log.Fatalf("Could not update database version! (Error: %s)", err)
	}
}
