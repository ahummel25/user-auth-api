package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" //Import sqlite3 for db flavor.

	"github.com/src/user-auth-api/utils"
)

const dbFile string = "file://mnt/storage/db/user_auth.db?cache=shared&_auth&_auth_user=admin&_auth_pass=admin&_auth_crypt=sha1"

// ConnectToDB will return a database connection
func ConnectToDB() (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	if db, err = sql.Open("sqlite3", dbFile); err != nil {
		return nil, fmt.Errorf("error opening Sqlite DB: %v", err)
	}

	return db, nil
}

// InitDB will initialize the DB and execute any DDL statements.
func InitDB() error {
	var (
		db   *sql.DB
		err  error
		file *os.File
	)

	if !utils.FileExists(dbFile) {
		log.Println("Creating database file: " + dbFile)
		if file, err = os.Create(dbFile); err != nil {
			return fmt.Errorf("error creating DB file %v", err)
		}
	}

	if err = file.Close(); err != nil {
		return fmt.Errorf("error closing DB file %v", err)
	}

	if db, err = ConnectToDB(); err != nil {
		return err
	}

	for _, ddlStmt := range ddlStmts {
		if _, err = db.Exec(ddlStmt); err != nil {
			return fmt.Errorf("error executing DDL statement%v", err)
		}
	}

	return nil
}
