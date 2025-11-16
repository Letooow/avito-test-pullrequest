package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
)

type Team struct {
	teamRepository TeamRepository
	userRepository UserRepository
}

func NewTeam(teamRepository TeamRepository, userRepository UserRepository) Team {
	return Team{
		teamRepository: teamRepository,
		userRepository: userRepository,
	}
}

func (t *Team) CreateTeam(ctx context.Context, team *domain.Team, members []domain.User) (*domain.Team, error) {
	if team == nil || team.Name == "" {
		return nil, ErrInvalidTeamName
	}
	if len(members) == 0 {
		return nil, ErrMemberNotFound
	}
	if t.teamRepository == nil {
		return nil, ErrTeamRepositoryNotFound
	}
	if t.userRepository == nil {
		return nil, ErrUserRepositoryNotFound
	}

	existing, err := t.teamRepository.GetTeamByName(ctx, team.Name)
	if err != nil && !errors.Is(err, ErrTeamNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTeamAlreadyExists
	}

	if err := t.teamRepository.SaveTeam(ctx, team); err != nil {
		return nil, err
	}

	for _, member := range members {
		user, err := t.userRepository.GetUserByID(ctx, member.ID)
		if err != nil {
			if errors.Is(err, ErrMemberNotFound) {

				user = &domain.User{
					ID:       member.ID,
					Username: member.Username,
					IsActive: member.IsActive,
				}
				if err := t.userRepository.SaveUser(ctx, user); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			if user.IsActive != member.IsActive {
				user.IsActive = member.IsActive
				if err := t.userRepository.UpdateUser(ctx, user); err != nil {
					return nil, err
				}
			}
		}

		if err := t.teamRepository.LinkUserToTeam(ctx, team, user); err != nil {
			return nil, err
		}
	}

	team, err = t.teamRepository.GetTeamByName(ctx, team.Name)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	if t.teamRepository == nil {
		return nil, ErrInvalidTeamName
	} else if t.userRepository == nil {
		return nil, ErrMemberNotFound
	}
	if t.teamRepository == nil {
		return nil, ErrInvalidTeamName
	} else if t.userRepository == nil {
		return nil, ErrMemberNotFound
	}
	team, err := t.teamRepository.GetTeamByName(ctx, teamName)
	if err != nil {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
