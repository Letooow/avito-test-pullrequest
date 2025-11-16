package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestUser_SetActive_UserRepositoryNil(t *testing.T) {
	// Arrange
	ctx := context.Background()
	u := &User{
		userRepository: nil,
	}

	// Act
	_, err := u.SetActive(ctx, "user-1", true)

	// Assert
	if !errors.Is(err, ErrMemberNotFound) {
		t.Fatalf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestUser_SetActive_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepo := NewMockUserRepository(ctrl)
	u := &User{
		userRepository: mockUserRepo,
	}

	existing := &domain.User{
		ID:       "user-1",
		Username: "test",
		IsActive: false,
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, "user-1").
		Return(existing, nil)

	mockUserRepo.EXPECT().
		UpdateUser(ctx, gomock.AssignableToTypeOf(&domain.User{})).
		DoAndReturn(func(_ context.Context, user *domain.User) error {
			if !user.IsActive {
				t.Fatalf("expected user to be active in UpdateUser call")
			}
			if user.ID != "user-1" {
				t.Fatalf("expected ID user-1, got %s", user.ID)
			}
			return nil
		})

	// Act
	_, err := u.SetActive(ctx, "user-1", true)

	// Assert
	if err != nil {
		t.Fatalf("SetActive() unexpected error: %v", err)
	}
}

func TestUser_GetUserPullRequests_UserNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)
	mockPRRepo := NewMockPullRequestRepository(ctrl)

	u := &User{
		userRepository:         mockUserRepo,
		requestOwnerRepository: mockReqOwnerRepo,
		pullRequestRepository:  mockPRRepo,
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, "user-1").
		Return(nil, ErrMemberNotFound)

	// Act
	got, err := u.GetUserPullRequests(ctx, "user-1")

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrMemberNotFound) {
		t.Fatalf("expected ErrMemberNotFound, got %v", err)
	}
}

func TestUser_GetUserPullRequests_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)
	mockPRRepo := NewMockPullRequestRepository(ctrl)

	u := &User{
		userRepository:         mockUserRepo,
		requestOwnerRepository: mockReqOwnerRepo,
		pullRequestRepository:  mockPRRepo,
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, "user-1").
		Return(&domain.User{ID: "user-1", Username: "u1", IsActive: true}, nil)

	reqOwners := []domain.RequestOwner{
		{UserID: "user-1", RequestID: "pr-1", Role: domain.UserRoleReviewer},
		{UserID: "user-1", RequestID: "pr-2", Role: domain.UserRoleAuthor}, // должен быть проигнорирован
	}
	mockReqOwnerRepo.EXPECT().
		GetRequestsByUserID(ctx, "user-1").
		Return(reqOwners, nil)

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, "pr-1").
		Return(&domain.PullRequest{ID: "pr-1", Name: "PR 1"}, nil)

	// Act
	got, err := u.GetUserPullRequests(ctx, "user-1")

	// Assert
	if err != nil {
		t.Fatalf("GetUserPullRequests() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(got))
	}
	if got[0].ID != "pr-1" {
		t.Fatalf("expected PR id pr-1, got %s", got[0].ID)
	}
}

func TestUser_GetUserPullRequests_PullRequestErrorWrapped(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)
	mockPRRepo := NewMockPullRequestRepository(ctrl)

	u := &User{
		userRepository:         mockUserRepo,
		requestOwnerRepository: mockReqOwnerRepo,
		pullRequestRepository:  mockPRRepo,
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, "user-1").
		Return(&domain.User{ID: "user-1", Username: "u1", IsActive: true}, nil)

	mockReqOwnerRepo.EXPECT().
		GetRequestsByUserID(ctx, "user-1").
		Return([]domain.RequestOwner{
			{UserID: "user-1", RequestID: "pr-1", Role: domain.UserRoleReviewer},
		}, nil)

	underlyingErr := errors.New("db error")
	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, "pr-1").
		Return(nil, underlyingErr)

	// Act
	got, err := u.GetUserPullRequests(ctx, "user-1")

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrPullRequestNotFound) {
		t.Fatalf("expected error to contain ErrPullRequestNotFound, got %v", err)
	}
}
