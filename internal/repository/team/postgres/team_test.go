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

func TestTeamRepository_SaveTeam(t *testing.T) {
	type args struct {
		team *domain.Team
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
				team: &domain.Team{Name: "team-1"},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO teams (teamname)")).
					WithArgs("team-1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				team: &domain.Team{Name: "team-1"},
			},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectExec(regexp.QuoteMeta("INSERT INTO teams (teamname)")).
					WithArgs("team-1").
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

			repo := &TeamRepository{db: queries}

			err := repo.SaveTeam(context.Background(), tt.args.team)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SaveTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTeamRepository_GetTeamByName(t *testing.T) {
	type args struct {
		name string
	}

	tests := []struct {
		name    string
		args    args
		mock    func(sqlmock.Sqlmock)
		want    *domain.Team
		wantErr bool
	}{
		{
			name: "not found returns nil",
			args: args{name: "missing"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM teams WHERE teamname =")).
					WithArgs("missing").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{name: "team-1"},
			mock: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta("FROM teams WHERE teamname =")).
					WithArgs("team-1").
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

			repo := &TeamRepository{db: queries}

			got, err := repo.GetTeamByName(context.Background(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTeamByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetTeamByName() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTeamRepository_GetTeams(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(sqlmock.Sqlmock)
		want    []domain.Team
		wantErr bool
	}{
		{
			name: "empty",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"teamname"})
				m.ExpectQuery(regexp.QuoteMeta("SELECT teamname FROM teams")).
					WillReturnRows(rows)
			},
			want:    []domain.Team{},
			wantErr: false,
		},
		{
			name: "multiple teams",
			mock: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"teamname"}).
					AddRow("team-1").
					AddRow("team-2")
				m.ExpectQuery(regexp.QuoteMeta("SELECT teamname FROM teams")).
					WillReturnRows(rows)
			},
			want: []domain.Team{
				{Name: "team-1"},
				{Name: "team-2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queries, mock, cleanup := usecase.NewTestQueries(t)
			defer cleanup()

			tt.mock(mock)

			repo := &TeamRepository{db: queries}

			got, err := repo.GetTeams(context.Background())
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetTeams() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("GetTeams() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
