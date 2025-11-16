package postgres

import (
	"avito-test/internal/domain"
	"avito-test/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRequestOwnerRepository_SaveRequestOwner(t *testing.T) {
	type args struct {
		ro *domain.RequestOwner
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ro: &domain.RequestOwner{
					UserID:    "user-1",
					RequestID: "pr-1",
					Role:      "REVIEWER",
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO users_pull_requests")).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				ro: &domain.RequestOwner{
					UserID:    "user-1",
					RequestID: "pr-1",
					Role:      "REVIEWER",
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO users_pull_requests")).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &RequestOwnerRepository{db: queries}

			err := repo.SaveRequestOwner(context.Background(), tt.args.ro)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SaveRequestOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestOwnerRepository_DeleteRequestOwner(t *testing.T) {
	type args struct {
		ro *domain.RequestOwner
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ro: &domain.RequestOwner{
					UserID:    "user-1",
					RequestID: "pr-1",
					Role:      "REVIEWER",
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(
					"DELETE FROM users_pull_requests WHERE pullrequestid = $1 AND userid = $2",
				)).
					WithArgs("pr-1", "user-1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				ro: &domain.RequestOwner{
					UserID:    "user-1",
					RequestID: "pr-1",
					Role:      "REVIEWER",
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(
					"DELETE FROM users_pull_requests WHERE pullrequestid = $1 AND userid = $2",
				)).
					WithArgs("pr-1", "user-1").
					WillReturnError(errors.New("delete failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &RequestOwnerRepository{db: queries}

			err := repo.DeleteRequestOwner(context.Background(), tt.args.ro)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteRequestOwner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestOwnerRepository_GetUsersByPullRequestID(t *testing.T) {
	type args struct {
		prID string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		want    []domain.RequestOwner
		wantErr bool
	}{
		{
			name: "no rows returns empty slice",
			args: args{prID: "pr-1"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT userid, role FROM users_pull_requests WHERE pullrequestid = $1",
				)).
					WithArgs("pr-1").
					WillReturnError(sql.ErrNoRows)
			},
			want:    []domain.RequestOwner{},
			wantErr: false,
		},
		{
			name: "multiple rows",
			args: args{prID: "pr-1"},
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"userid", "role"}).
					AddRow("user-1", "AUTHOR").
					AddRow("user-2", "REVIEWER")
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT userid, role FROM users_pull_requests WHERE pullrequestid = $1",
				)).
					WithArgs("pr-1").
					WillReturnRows(rows)
			},
			want: []domain.RequestOwner{
				{UserID: "user-1", RequestID: "pr-1", Role: "AUTHOR"},
				{UserID: "user-2", RequestID: "pr-1", Role: "REVIEWER"},
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{prID: "pr-1"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT userid, role FROM users_pull_requests WHERE pullrequestid = $1",
				)).
					WithArgs("pr-1").
					WillReturnError(errors.New("select failed"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &RequestOwnerRepository{db: queries}

			got, err := repo.GetUsersByPullRequestID(context.Background(), tt.args.prID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetUsersByPullRequestID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetUsersByPullRequestID() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
