package postgres

import (
	"avito-test/internal/domain"
	"avito-test/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPullRequestRepository_SavePullRequest(t *testing.T) {
	type args struct {
		pr *domain.PullRequest
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
				pr: &domain.PullRequest{
					ID:     "pr-1",
					Name:   "Test PR",
					Status: domain.RequestStatusOpen,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO pull_requests")).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				pr: &domain.PullRequest{
					ID:     "pr-1",
					Name:   "Test PR",
					Status: domain.RequestStatusOpen,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO pull_requests")).
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

			repo := &PullRequestRepository{db: queries}

			err := repo.SavePullRequest(context.Background(), tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SavePullRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPullRequestRepository_UpdatePullRequest(t *testing.T) {
	type args struct {
		pr *domain.PullRequest
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
				pr: &domain.PullRequest{
					ID:     "pr-1",
					Status: domain.RequestStatusMerged,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("UPDATE pull_requests")).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				pr: &domain.PullRequest{
					ID:     "pr-1",
					Status: domain.RequestStatusMerged,
				},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("UPDATE pull_requests")).
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

			repo := &PullRequestRepository{db: queries}

			err := repo.UpdatePullRequest(context.Background(), tt.args.pr)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdatePullRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPullRequestRepository_GetPullRequestByID_Errors(t *testing.T) {
	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		wantNil bool
		wantErr bool
	}{
		{
			name: "not found returns nil",
			args: args{id: "missing"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM pull_requests WHERE pullrequestid =")).
					WithArgs("missing").
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{id: "pr-1"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM pull_requests WHERE pullrequestid =")).
					WithArgs("pr-1").
					WillReturnError(errors.New("select failed"))
			},
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &PullRequestRepository{db: queries}

			got, err := repo.GetPullRequestByID(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPullRequestByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && got != nil {
				t.Fatalf("GetPullRequestByID() expected nil, got %#v", got)
			}
		})
	}
}

func TestPullRequestRepository_GetPullRequests_Errors(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(sqlmock.Sqlmock)
		wantNil bool
		wantErr bool
	}{
		{
			name: "no rows returns empty slice",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM pull_requests")).
					WillReturnError(sql.ErrNoRows)
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "db error",
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM pull_requests")).
					WillReturnError(errors.New("select failed"))
			},
			wantNil: true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &PullRequestRepository{db: queries}

			got, err := repo.GetPullRequests(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPullRequests() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && got != nil {
				t.Fatalf("GetPullRequests() expected nil, got %#v", got)
			}
			if !tt.wantNil && got == nil {
				t.Fatalf("GetPullRequests() expected non-nil slice")
			}
		})
	}
}
