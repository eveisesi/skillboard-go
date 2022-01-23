package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eveisesi/skillz/migrations"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	commands = append(
		commands,
		&cli.Command{
			Name: "migrate",
			Subcommands: []*cli.Command{
				{
					Name:   "up",
					Action: migrateUpCommand,
				},
				{
					Name:   "down",
					Action: migrateDownCommand,
				},
				{
					Name:   "create",
					Action: migrateCreateCommand,
				},
				{
					Name:   "current",
					Action: migrateCurrentCommand,
				},
			},
		},
	)
}

var migrationDir = "migrations/"

var migrationFS = migrations.NewFS()

func migrateUpCommand(c *cli.Context) error {

	successes := 0
	strSteps := c.Args().Get(0)
	steps := 0
	if strSteps != "" {
		steps, _ = strconv.Atoi(strSteps)
	}

	files, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		logger.WithError(err).Fatal("failed to read migrations")
	}

	var fileNames = make([]string, 0, len(files))
	for _, file := range files {
		if file.Name() == "embed.go" {
			continue
		}
		if !strings.Contains(file.Name(), ".up.") {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	sort.Strings(fileNames)

	err = initializeMigrations()
	if err != nil {
		logger.WithError(err).Fatal("failed to initialize migrations table")
	}

	logger.Info("migrations table initialize successfully")

	printMostRecentMigration()

	for _, file := range fileNames {

		name := strings.TrimSuffix(file, ".up.sql")

		_, err := getMigration(name)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.WithError(err).Fatal("failed to check if migration has been executed")
		}

		if err == nil {
			continue
		}

		entry := logger.WithField("name", name)

		handle, err := migrationFS.Open(file)
		if err != nil {
			entry.WithError(err).Fatal("failed to open file for migration")
		}

		data, err := io.ReadAll(handle)
		if err != nil {
			entry.WithError(err).Fatal("failed to read file for migration")
		}

		if len(data) == 0 {
			entry.WithError(err).Fatal("empty migration file detected, halting execution")
		}

		query := string(data)

		_, err = dbConn.Exec(query)
		if err != nil {
			entry.WithError(err).Fatal("failed to execute migration")
		}

		err = createMigration(name)
		if err != nil {
			entry.WithError(err).Fatal("failed to log migration execute in migration table")
		}

		entry.Info("migration executed successfully")
		successes++
		if successes >= steps && steps > 0 {
			break
		}

	}

	return nil

}

func migrateDownCommand(c *cli.Context) error {

	var ctx = context.Background()

	strSteps := c.Args().First()
	steps := 0
	if strSteps != "" {
		steps, _ = strconv.Atoi(strSteps)
	}

	if steps == 0 {
		logger.Fatal("Destructive down command requires explict negating of steps. To run all down migration, please pass -1 else pass steps > 0")
	}

	successes := 0

	migrations, err := getMigrations()
	if err != nil {
		logger.WithError(err).Fatal("failed to fetch migration from database")
	}

	printMostRecentMigration()

	for i := len(migrations) - 1; i >= 0; i-- {
		migration := migrations[i]

		fileName := fmt.Sprintf("%s.down.sql", migration.Name)

		entry := logger.WithFields(logrus.Fields{
			"name":     migration.Name,
			"fileName": fileName,
		})
		file, err := migrationFS.Open(fileName)
		if err != nil {
			entry.WithError(err).Fatal("failed to open migration file")
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			entry.WithError(err).Fatal("failed to read migration file")
		}

		if len(data) == 0 {
			entry.WithError(err).Fatal("empty migration file detected, halting execution")
		}

		query := string(data)

		_, err = dbConn.ExecContext(ctx, query)
		if err != nil {
			entry.WithError(err).Fatal("failed to execute query")
		}

		err = deleteMigration(migration.Name)
		if err != nil {
			entry.WithError(err).Fatal("failed to remove migration from migrations table")
		}

		entry.Info("migration executed successfully")
		successes++
		if successes >= steps && steps > 0 {
			break
		}

	}

	return nil

}

const createTableStmt = `CREATE TABLE %s (
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
) COLLATE = 'utf8mb4_unicode_ci' ENGINE = InnoDB;`

func migrateCreateCommand(c *cli.Context) error {

	name := c.Args().First()
	if name == "" {
		return fmt.Errorf("name is required, received empty string")
	}

	now := time.Now()

	filename := fmt.Sprintf("%s%s_%s.%%s.sql", migrationDir, now.Format("20060102150405"), name)
	up := fmt.Sprintf(filename, "up")
	entry := logger.WithFields(logrus.Fields{
		"name": name,
	})
	upFile, err := os.Create(up)
	if err != nil {
		entry.WithField("up", up).WithError(err).Fatal("failed to create up file")
	}
	defer upFile.Close()
	_, _ = upFile.WriteString(
		fmt.Sprintf(
			createTableStmt,
			name,
		),
	)

	entry.WithField("up", up).Info("migration created successfully")
	down := fmt.Sprintf(filename, "down")
	downFile, err := os.Create(down)
	if err != nil {
		entry.WithField("down", down).WithError(err).Fatal("failed to create down file")
	}
	defer downFile.Close()
	_, _ = downFile.WriteString(fmt.Sprintf("DROP TABLE `%s`;", name))

	entry.WithField("down", down).Info("migration created successfully")

	return nil

}

func migrateCurrentCommand(c *cli.Context) error {
	printMostRecentMigration()
	return nil
}

func printMostRecentMigration() {

	migrations, err := getMigrations()
	if err != nil {
		logger.WithError(err).Fatal("failed to fetch migration from database")
	}

	if len(migrations) == 0 {
		logger.Info("no migrations bave been run")
		return
	}

	logger.WithField("current", migrations[len(migrations)-1].Name).Info("current migration")

}

type migration struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

const createMigrationsTableQuery = `
	CREATE TABLE migrations (
		id INT UNSIGNED NOT NULL AUTO_INCREMENT, 
		name VARCHAR(255) NOT NULL,                  
		created_at DATETIME NOT NULL,               
		PRIMARY KEY (id) USING BTREE,                
		UNIQUE INDEX migrations_name_unique_idx (name)   
	) COLLATE = 'utf8mb4_unicode_ci' ENGINE = INNODB;
`

const checkTableExistsQuery = `
	SELECT 
		COUNT(*)
	FROM information_schema.tables
	WHERE 
		table_schema = ? AND table_name = ?
	LIMIT 1;
`

func initializeMigrations() error {

	var count uint
	err := dbConn.Get(&count, checkTableExistsQuery, cfg.MySQL.DB, "migrations")
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = dbConn.Exec(createMigrationsTableQuery)
	return err

}

func getMigrations() ([]*migration, error) {

	query := `
		SELECT id,name,created_at FROM migrations
	`

	var migrations = make([]*migration, 0)
	err := dbConn.Select(&migrations, query)
	return migrations, err

}

func getMigration(name string) (*migration, error) {

	query := `
		SELECT id,name,created_at FROM migrations WHERE name = ?
	`

	var migration = new(migration)
	err := dbConn.Get(migration, query, name)
	return migration, err

}

func createMigration(name string) error {
	query := `
		INSERT INTO migrations (name, created_at)VALUES(?, NOW())
	`

	_, err := dbConn.Exec(query, name)
	return err
}

func deleteMigration(name string) error {
	query := `
		DELETE FROM migrations where name = ?
	`

	_, err := dbConn.Exec(query, name)
	return err
}
