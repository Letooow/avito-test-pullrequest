package postgres

import (
	"avito-test/internal/db"
	"avito-test/internal/domain"
	"avito-test/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	if u.db == nil {
		return errors.New("db is nil")
	}
	if user == nil {
		return errors.New("user is nil")
	}
	err := u.db.UpdateUser(ctx, db.UpdateUserParams{Userid: user.ID, Username: user.Username, Isactive: user.IsActive})
	if err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}
	return nil
}

func (u *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if u.db == nil {
		return errors.New("db is nil")
	}
	if user == nil {
		return errors.New("user is nil")
	}
	err := u.db.SaveUser(ctx, db.SaveUserParams{Userid: user.ID, Username: user.Username, Isactive: user.IsActive})
	if err != nil {
		return fmt.Errorf("can't save new team: %w", err)
	}
	return nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if u.db == nil {
		return nil, errors.New("db is nil")
	}
	user, err := u.db.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, usecase.ErrMemberNotFound
	} else if err != nil {
		return nil, fmt.Errorf("can't get user by id: %w", err)
	}
	return &domain.User{ID: user.Userid, Username: user.Username, IsActive: user.Isactive}, nil
}

func (u *UserRepository) GetUsersByTeamName(ctx context.Context, teamName string) ([]domain.User, error) {
	if u.db == nil {
		return nil, errors.New("db is nil")
	}
	gotUser, err := u.db.GetUsersByTeamName(ctx, teamName)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.User{}, usecase.ErrMemberNotFound
	} else if err != nil {
		return nil, fmt.Errorf("can't get gotUser by team name: %w", err)
	}
	result := make([]domain.User, len(gotUser))
	for i, user := range gotUser {
		result[i] = domain.User{ID: user.Userid, Username: user.Username, IsActive: user.Isactive}
	}
	return result, nil
}

func (u *UserRepository) GetTeamsByUserID(ctx context.Context, userID string) ([]domain.Team, error) {
	if u.db == nil {
		return nil, errors.New("db is nil")
	}
	team, err := u.db.GetUsersTeams(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Team{}, usecase.ErrTeamNotFound
	} else if err != nil {
		return nil, fmt.Errorf("can't get teams: %w", err)
	}
	result := make([]domain.Team, len(team))
	for i, t := range team {
		result[i] = domain.Team{Name: t}
	}
	return result, nil
}
