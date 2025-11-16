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

func TestUserRepository_SaveUser(t *testing.T) {
	type args struct {
		user *domain.User
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
				user: &domain.User{
					ID:       "user-1",
					Username: "alice",
					IsActive: true,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				// ожидаем корректную запись всех полей
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO users (userid, username, isactive) VALUES ($1, $2, $3)")).
					WithArgs("user-1", "alice", true).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				user: &domain.User{
					ID:       "user-1",
					Username: "alice",
					IsActive: true,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO users (userid, username, isactive) VALUES ($1, $2, $3)")).
					WithArgs("user-1", "alice", true).
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

			repo := &UserRepository{db: queries}

			err := repo.SaveUser(context.Background(), tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SaveUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	type args struct {
		user *domain.User
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
				user: &domain.User{
					ID:       "user-1",
					Username: "alice",
					IsActive: true,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE users SET username = $1, isactive = $2 WHERE userid = $3",
				)).
					WithArgs("alice", true, "user-1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				user: &domain.User{
					ID:       "user-1",
					Username: "alice",
					IsActive: true,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta(
					"UPDATE users SET username = $1, isactive = $2 WHERE userid = $3",
				)).
					WithArgs("alice", true, "user-1").
					WillReturnError(errors.New("update failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &UserRepository{db: queries}

			err := repo.UpdateUser(context.Background(), tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		want    *domain.User
		wantErr bool
	}{
		{
			name: "found",
			args: args{id: "user-1"},
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"userid", "username", "isactive"}).
					AddRow("user-1", "alice", true)
				m.ExpectQuery(regexp.QuoteMeta(
					"FROM users WHERE userid = $1",
				)).
					WithArgs("user-1").
					WillReturnRows(rows)
			},
			want: &domain.User{
				ID:       "user-1",
				Username: "alice",

				IsActive: true,
			},
			wantErr: false,
		},
		{
			name: "not found returns nil",
			args: args{id: "missing"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"FROM users WHERE userid = $1",
				)).
					WithArgs("missing").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db error",
			args: args{id: "user-1"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"FROM users WHERE userid = $1",
				)).
					WithArgs("user-1").
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

			repo := &UserRepository{db: queries}

			got, err := repo.GetUserByID(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetUserByID() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
