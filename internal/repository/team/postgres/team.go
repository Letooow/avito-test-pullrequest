package postgres

import (
	"avito-test/internal/db"
	"avito-test/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type TeamRepository struct {
	db *db.Queries
}

func NewTeamRepository(db *db.Queries) *TeamRepository {
	return &TeamRepository{db: db}
}

func (t *TeamRepository) LinkUserToTeam(ctx context.Context, team *domain.Team, user *domain.User) error {
	err := t.db.SaveUserTeam(ctx, db.SaveUserTeamParams{Teamname: team.Name, Userid: user.ID})
	if err != nil {
		return fmt.Errorf("error saving user team: %w", err)
	}
	return nil
}

func (t *TeamRepository) SaveTeam(ctx context.Context, team *domain.Team) error {
	err := t.db.CreateTeam(ctx, team.Name)
	if err != nil {
		return fmt.Errorf("can't save new team: %w", err)
	}
	return nil
}

func (t *TeamRepository) GetTeamByName(ctx context.Context, name string) (*domain.Team, error) {
	team, err := t.db.GetTeamByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get team by name: %w", err)
	}
	members, err := t.db.GetUsersByTeamName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("can't get users by team name: %w", err)
	}
	result := make([]domain.User, len(members))
	for i, member := range members {
		result[i] = domain.User{ID: member.Userid, Username: member.Username, IsActive: member.Isactive}
	}
	return &domain.Team{Name: team, Members: result}, nil
}

func (t *TeamRepository) GetTeams(ctx context.Context) ([]domain.Team, error) {
	teams, err := t.db.GetTeams(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.Team{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get teams: %w", err)
	}
	result := make([]domain.Team, len(teams))
	for i, team := range teams {
		resTeam, err := t.GetTeamByName(ctx, team)
		if err != nil {
			return nil, fmt.Errorf("can't get team by name: %w", err)
		}
		result[i] = *resTeam
	}
	return result, nil
}
