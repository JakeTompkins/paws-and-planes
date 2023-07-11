package migrate

import (
	"database/sql"
	"log"
	"os"
	db "paws-n-planes/pkg/db"
	m "paws-n-planes/pkg/models/migration"
	"sort"
	"strings"
	"time"
)

type MigrationManager struct {
	db                 *sql.DB
	tx                 *sql.Tx
	migrationDirectory string
}

func New() (*MigrationManager, error) {
	newManager := &MigrationManager{}

	db, err := db.GetInstance()

	if err != nil {
		return nil, err
	}

	newManager.db = db
	newManager.migrationDirectory = "/migrations"

	return newManager, nil
}

func (mm *MigrationManager) ensureMigrationTableExists() {
	var numTables int
	sql := `select count(TABLE_NAME) from information_schema.TABLES where TABLE_NAME = 'migrations' `

	result := mm.db.QueryRow(sql)

	err := result.Scan(numTables)

	if err != nil {
		log.Println("Did not find migrations table, attempting to create...")

		sql = `create table migrations (
		name text,
		date_run date not null default now(),
		PRIMARY KEY (name)
		)`

		_, err = mm.db.Exec(sql)

		if err != nil {
			log.Fatalf("Could not create migration table, aborting: %v", err)
		}
	}

	if numTables > 1 {
		log.Fatal("Found more than one migration table, manuall fix schema and try again")
	}
}

func (mm *MigrationManager) getPastMigrations() []m.Migration {
	db := mm.db

	if mm.tx == nil {
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Could not begin transaction: %v", err)
		}

		mm.tx = tx
	}

	var pastMigrations []m.Migration

	rows, err := mm.tx.Query(`
    select * from migrations
	`)

	if err != nil {
		log.Fatalf("Could not fetch past migrations: %v", err)
	}

	for rows.Next() {
		migration := &m.Migration{}

		err = rows.Scan(migration.Name, migration.DateRun)

		if err != nil {
			log.Fatalf("Could not parse past migrations: %v", err)
		}

		pastMigrations = append(pastMigrations, *migration)
	}

	return pastMigrations
}

func (mm *MigrationManager) getMigrationFileNames() []string {
	var fileNames []string
	fileEntries, err := os.ReadDir(mm.migrationDirectory)

	if err != nil {
		log.Fatalf("Could not read from migrations directory: %v", err)
	}

	for _, f := range fileEntries {
		fileNames = append(fileNames, f.Name())
	}

	return fileNames
}

func (mm *MigrationManager) runNewMigrations(newMigrationNames []string) {
	for _, fileName := range newMigrationNames {
		fileContent, err := os.ReadFile(mm.migrationDirectory + "/" + fileName)

		if err != nil {
			log.Fatalf("Could not read migration:%s\n%v", fileName, err)
		}

		_, err = mm.tx.Exec(string(fileContent))

		if err != nil {
			mm.tx.Rollback()
			log.Fatalf("Could not run migration:%s\n%v", fileName, err)
		}

		_, err = mm.tx.Exec(`
     insert into migrations (name, date_run) values ($1, $2)
		`, fileName, time.Now().Format("MM-DD-YYYY"))
	}

	mm.tx.Commit()
}

func (mm *MigrationManager) Run() {
	mm.ensureMigrationTableExists()

	migrationFiles := mm.getMigrationFileNames()
	pastMigrations := mm.getPastMigrations()
	newMigrationFileNames := filterAndSortNewMigrationFiles(pastMigrations, migrationFiles)

	mm.runNewMigrations(newMigrationFileNames)
}

func filterAndSortNewMigrationFiles(pastMigrations []m.Migration, allFileNames []string) []string {
	var oldNames map[string]bool
	var newNames []string

	for _, migration := range pastMigrations {
		oldNames[migration.Name] = true
	}

	for _, name := range allFileNames {
		if _, ok := oldNames[name]; ok == false {
			newNames = append(newNames, name)
		}
	}

	sort.Strings(newNames)

	return newNames
}
