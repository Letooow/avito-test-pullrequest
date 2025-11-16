package usecase

import (
	"avito-test/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func NewTestQueries(t *testing.T) (*db.Queries, sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}

	cleanup := func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
		sqlDB.Close()
	}

	return db.New(sqlDB), mock, cleanup
}
