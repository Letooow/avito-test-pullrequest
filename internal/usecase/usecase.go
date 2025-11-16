package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
)

var (
	ErrInvalidTeamName                = errors.New("invalid team name")
	ErrTeamAlreadyExists              = errors.New("team already exists")
	ErrTeamNotFound                   = errors.New("team not found")
	ErrMemberNotFound                 = errors.New("member not found")
	ErrPullRequestNotFound            = errors.New("pull request not found")
	ErrAuthorNotFound                 = errors.New("author not found")
	ErrPullRequestAlreadyExists       = errors.New("pull request already exists")
	ErrAuthorIsInactive               = errors.New("author is inactive")
	ErrCannotFindActiveMembers        = errors.New("cannot find active members")
	ErrPullRequestIsMerged            = errors.New("pull request is merged")
	ErrPullRequestRepositoryNotFound  = errors.New("pull request repository is nil")
	ErrUserRepositoryNotFound         = errors.New("user repository is nil")
	ErrTeamRepositoryNotFound         = errors.New("team repository is nil")
	ErrRequestOwnerRepositoryNotFound = errors.New("request owner repository is nil")
)

//go:generate mockgen -source usecase.go -package usecase -destination usecase_mock.go
type UserRepository interface {
	// SaveUser - функция сохранения пользователя
	SaveUser(ctx context.Context, user *domain.User) error
	// GetUserByID - функция получения пользователя по его ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	// GetUsersByTeamName - функция получения пользователей по названию команды
	GetUsersByTeamName(ctx context.Context, teamName string) ([]domain.User, error)
	// UpdateUser - функция обновления пользователя
	UpdateUser(ctx context.Context, user *domain.User) error
	// GetTeamsByUserID - функция получения списка команд пользователя
	GetTeamsByUserID(ctx context.Context, userID string) ([]domain.Team, error)
}

type RequestOwnerRepository interface {
	// SaveRequestOwner - функция сохранения связи между реквестом и пользователем
	SaveRequestOwner(ctx context.Context, request *domain.RequestOwner) error
	// DeleteRequestOwner - функция удаления связи между реквестом и пользователем
	DeleteRequestOwner(ctx context.Context, requestOwner *domain.RequestOwner) error
	// GetRequestsByUserID - функция получения реквестов по id участника
	GetRequestsByUserID(ctx context.Context, userID string) ([]domain.RequestOwner, error)
	// GetUsersByPullRequestID - функция получения пользователей по id pr
	GetUsersByPullRequestID(ctx context.Context, pullRequestID string) ([]domain.RequestOwner, error)
}
type TeamRepository interface {
	// SaveTeam - функция сохранения команды
	SaveTeam(ctx context.Context, team *domain.Team) error
	// GetTeamByName - функция получения команды по ее ID
	GetTeamByName(ctx context.Context, name string) (*domain.Team, error)
	// GetTeams - функция получения всех команд
	GetTeams(ctx context.Context) ([]domain.Team, error)
	// LinkUserToTeam - функция привязки пользователя к команде
	LinkUserToTeam(ctx context.Context, team *domain.Team, user *domain.User) error
}

type PullRequestRepository interface {
	// SavePullRequest - функция сохранения пул реквеста
	SavePullRequest(ctx context.Context, pull *domain.PullRequest) error
	// GetPullRequestByID - функция получения пул реквеста по его ID
	GetPullRequestByID(ctx context.Context, id string) (*domain.PullRequest, error)
	// GetPullRequests - функция получения всех пул реквестов
	GetPullRequests(ctx context.Context) ([]domain.PullRequest, error)
	// UpdatePullRequest - функция обновления пул реквеста
	UpdatePullRequest(ctx context.Context, pull *domain.PullRequest) error
}
