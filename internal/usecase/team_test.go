package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestTeam_CreateTeam_TeamNil(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)

	usecase := NewTeam(mockTeamRepo, mockUserRepo)

	// Act
	err := usecase.CreateTeam(ctx, nil, []domain.User{})

	// Assert
	if !errors.Is(err, ErrInvalidTeamName) {
		t.Fatalf("expected ErrInvalidTeamName, got %v", err)
	}
}

func TestTeam_CreateTeam_TeamAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)

	usecase := NewTeam(mockTeamRepo, mockUserRepo)

	team := &domain.Team{Name: "team-1"}
	members := []domain.User{
		{ID: "user-1", Username: "u1", IsActive: true},
	}

	mockTeamRepo.EXPECT().
		GetTeamByName(ctx, "team-1").
		Return(&domain.Team{Name: "team-1"}, nil)

	// Act
	err := usecase.CreateTeam(ctx, team, members)

	// Assert
	if !errors.Is(err, ErrTeamAlreadyExists) {
		t.Fatalf("expected ErrTeamAlreadyExists, got %v", err)
	}
}

func TestTeam_CreateTeam_SuccessExistingUser(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)

	usecase := NewTeam(mockTeamRepo, mockUserRepo)

	team := &domain.Team{Name: "team-1"}
	member := domain.User{ID: "user-1", Username: "u1", IsActive: true}

	mockTeamRepo.EXPECT().
		GetTeamByName(ctx, team.Name).
		Return(nil, nil)

	existingUser := &domain.User{ID: "user-1", Username: "u1", IsActive: true}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, member.ID).
		Return(existingUser, nil)

	mockTeamRepo.EXPECT().
		LinkUserToTeam(ctx, team, existingUser).
		Return(nil)

	// Act
	err := usecase.CreateTeam(ctx, team, []domain.User{member})

	// Assert
	if err != nil {
		t.Fatalf("CreateTeam() unexpected error: %v", err)
	}
}

func TestTeam_GetTeam_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)

	usecase := NewTeam(mockTeamRepo, mockUserRepo)

	expected := &domain.Team{Name: "team-1"}

	mockTeamRepo.EXPECT().
		GetTeamByName(ctx, "team-1").
		Return(expected, nil)

	// Act
	got, err := usecase.GetTeam(ctx, "team-1")

	// Assert
	if err != nil {
		t.Fatalf("GetTeam() unexpected error: %v", err)
	}
	if got != expected {
		t.Fatalf("expected %#v, got %#v", expected, got)
	}
}

func TestTeam_GetTeam_ErrorWrappedAsTeamNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)

	usecase := NewTeam(mockTeamRepo, mockUserRepo)

	mockTeamRepo.EXPECT().
		GetTeamByName(ctx, "team-1").
		Return(nil, errors.New("db error"))

	// Act
	got, err := usecase.GetTeam(ctx, "team-1")

	// Assert
	if got != nil {
		t.Fatalf("expected nil team, got %#v", got)
	}
	if !errors.Is(err, ErrTeamNotFound) {
		t.Fatalf("expected ErrTeamNotFound, got %v", err)
	}
}
