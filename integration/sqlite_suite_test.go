// NOTE: currently have a warning when build tags present. Don't forget to address this later.

package integration

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SQLiteRepoSuite struct {
	suite.Suite
	db *sqlx.DB
}

const sqliteDSN = "sqlite3://db/test.db"

func (s *SQLiteRepoSuite) SetupSuite() {
	var err error
	s.db, err = sqlx.Connect("sqlite3", "./db/test.db")
	if err != nil {
		log.Fatalf("couldn't connect to sqlite database: %v", err)
	}
}

func (s *SQLiteRepoSuite) SetupTest() {
	m, err := migrate.New("file://db/migrations", sqliteDSN)
	assert.NoError(s.T(), err)

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

func (s *SQLiteRepoSuite) TearDownTest() {
	m, err := migrate.New("file://db/migrations", sqliteDSN)
	assert.NoError(s.T(), err)
	assert.NoError(s.T(), m.Down())
}
