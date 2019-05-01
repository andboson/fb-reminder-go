package reminders

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/andboson/fb-reminder-go/internal"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
)

type remindersSuite struct {
	db     *sql.DB
	dbAddr string
	m      *migrate.Migrate
	rm     *Manager

	internal.DockerSuite
}

func (s *remindersSuite) Test_Create_GetByID() {
	t := time.Now().Add(1 * time.Hour)
	to := t.Add(1 * time.Hour)
	rem := Reminder{
		Text:             "texttext",
		UserID:           "1223567",
		RemindAt:         t,
		RemindAtOriginal: to.String(),
	}

	err := s.rm.Create(rem)

	remStored, err := s.rm.GetByID(1)
	s.Require().NoError(err)
	s.Equal(rem.Text, remStored.Text)
}

func (s *remindersSuite) Test_GetToday() {
	t := time.Now().Add(1 * time.Hour)
	to := t.Add(1 * time.Hour)
	rem := Reminder{
		Text:             "texttext",
		UserID:           "1223567",
		RemindAt:         t,
		RemindAtOriginal: to.String(),
	}

	err := s.rm.Create(rem)

	rems, err := s.rm.GetTodayByUser("1223567")
	s.Require().NoError(err)
	s.Require().Equal(len(rems), 1)
	s.Equal(rem.Text, rems[0].Text)
}

func (s *remindersSuite) SetupSuite() {
	var err error
	s.Setup("rem_test")
	addr := s.SetupPSQL("classes")
	s.dbAddr = addr

	err = internal.Retry(func() error {
		var err error
		const dsn = "postgres://%s:%s@%s/%s?connect_timeout=%d&sslmode=disable"
		result := fmt.Sprintf(dsn, "admin", "admin", s.dbAddr, "classes", 2)

		s.db, err = sql.Open("postgres", result)
		if err != nil {
			return err
		}
		_, err = s.db.Exec("SELECT 1")
		return err
	}, 1*time.Minute)
	s.Require().NoError(err, `Could not connect to db docker: %s`, err)

	// миграции
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	s.NoError(err)
	s.NotNil(driver)

	s.m, err = migrate.NewWithDatabaseInstance("file://../migrations", "postgres", driver)
	s.Require().NoError(err)

	err = s.m.Up()
	s.NoError(err)

	s.rm = NewManager(s.db)
}

func (s *remindersSuite) TearDownSuite() {
	s.Down()
}

func TestRemindersTestSuite(t *testing.T) {
	suite.Run(t, new(remindersSuite))
}
