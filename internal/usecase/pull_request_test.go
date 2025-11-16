package usecase

import (
	"avito-test/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestPullRequest_CreateRepository_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	author := &domain.User{
		ID:       "author-1",
		Username: "author",
		IsActive: true,
	}
	pr := &domain.PullRequest{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: author.ID,
	}

	coworkers := []domain.User{
		{ID: "u1", Username: "u1", IsActive: true},
		{ID: "u2", Username: "u2", IsActive: true},
		{ID: "u3", Username: "u3", IsActive: false},          // должен быть отфильтрован как неактивный
		{ID: "author-1", Username: "author", IsActive: true}, // должен быть отфильтрован как автор
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, author.ID).
		Return(author, nil)

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, pr.ID).
		Return(nil, ErrPullRequestNotFound)

	mockPRRepo.EXPECT().
		SavePullRequest(ctx, pr).
		Return(nil)

	authorCall := mockReqOwnerRepo.EXPECT().
		SaveRequestOwner(ctx, &domain.RequestOwner{
			RequestID: pr.ID,
			UserID:    author.ID,
			Role:      domain.UserRoleAuthor,
		}).
		Return(nil)

	selectedReviewers := map[string]bool{}
	reviewerCall1 := mockReqOwnerRepo.EXPECT().
		SaveRequestOwner(ctx, gomock.AssignableToTypeOf(&domain.RequestOwner{})).
		DoAndReturn(func(_ context.Context, ro *domain.RequestOwner) error {
			if ro.Role != domain.UserRoleReviewer {
				t.Fatalf("expected reviewer role, got %v", ro.Role)
			}
			if ro.RequestID != pr.ID {
				t.Fatalf("expected RequestID %s, got %s", pr.ID, ro.RequestID)
			}
			if ro.UserID != "u1" && ro.UserID != "u2" {
				t.Fatalf("unexpected reviewer id: %s", ro.UserID)
			}
			if selectedReviewers[ro.UserID] {
				t.Fatalf("duplicate reviewer assignment: %s", ro.UserID)
			}
			selectedReviewers[ro.UserID] = true
			return nil
		})

	reviewerCall2 := mockReqOwnerRepo.EXPECT().
		SaveRequestOwner(ctx, gomock.AssignableToTypeOf(&domain.RequestOwner{})).
		DoAndReturn(func(_ context.Context, ro *domain.RequestOwner) error {
			if ro.Role != domain.UserRoleReviewer {
				t.Fatalf("expected reviewer role, got %v", ro.Role)
			}
			if ro.RequestID != pr.ID {
				t.Fatalf("expected RequestID %s, got %s", pr.ID, ro.RequestID)
			}
			if ro.UserID != "u1" && ro.UserID != "u2" {
				t.Fatalf("unexpected reviewer id: %s", ro.UserID)
			}
			if selectedReviewers[ro.UserID] {
				t.Fatalf("duplicate reviewer assignment: %s", ro.UserID)
			}
			selectedReviewers[ro.UserID] = true
			return nil
		})

	gomock.InOrder(authorCall, reviewerCall1, reviewerCall2)

	mockUserRepo.EXPECT().
		GetTeamsByUserID(ctx, author.ID).
		Return([]domain.Team{{Name: "team-1"}}, nil)

	mockUserRepo.EXPECT().
		GetUsersByTeamName(ctx, "team-1").
		Return(coworkers, nil)

	// Act
	got, err := usecase.CreatePullRequest(ctx, pr)

	// Assert
	if err != nil {
		t.Fatalf("CreatePullRequest() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("CreatePullRequest() returned nil PullRequest")
	}
	if got.Status != domain.RequestStatusOpen {
		t.Fatalf("expected status %v, got %v", domain.RequestStatusOpen, got.Status)
	}
	if len(got.AssignedReviewersID) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(got.AssignedReviewersID))
	}
	for _, id := range got.AssignedReviewersID {
		if id != "u1" && id != "u2" {
			t.Fatalf("unexpected reviewer id in AssignedReviewersID: %s", id)
		}
	}
}

func TestPullRequest_CreateRepository_RequestNil(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := NewPullRequest(nil, nil, nil, nil)
	ctx := context.Background()

	// Act
	got, err := usecase.CreatePullRequest(ctx, nil)

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrAuthorNotFound) {
		t.Fatalf("expected ErrAuthorNotFound, got %v", err)
	}
}

func TestPullRequest_CreateRepository_AuthorNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	pr := &domain.PullRequest{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: "author-1",
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, "author-1").
		Return(nil, ErrMemberNotFound)

	// Act
	got, err := usecase.CreatePullRequest(ctx, pr)

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrAuthorNotFound) {
		t.Fatalf("expected ErrAuthorNotFound, got %v", err)
	}
}

func TestPullRequest_CreateRepository_PullRequestAlreadyExists(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	author := &domain.User{
		ID:       "author-1",
		Username: "author",
		IsActive: true,
	}
	pr := &domain.PullRequest{
		ID:       "pr-1",
		Name:     "Test PR",
		AuthorID: author.ID,
	}

	mockUserRepo.EXPECT().
		GetUserByID(ctx, author.ID).
		Return(author, nil)

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, pr.ID).
		Return(&domain.PullRequest{}, nil)

	// Act
	got, err := usecase.CreatePullRequest(ctx, pr)

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrPullRequestAlreadyExists) {
		t.Fatalf("expected ErrPullRequestAlreadyExists, got %v", err)
	}
}

func TestPullRequest_UpdateRepository_MergedNoUpdate(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	existing := &domain.PullRequest{
		ID:     "pr-1",
		Status: domain.RequestStatusMerged,
	}
	req := &domain.PullRequest{
		ID:     "pr-1",
		Status: domain.RequestStatusOpen,
	}

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, req.ID).
		Return(existing, nil)

	// Act
	got, err := usecase.UpdatePullRequest(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("UpdatePullRequest() unexpected error: %v", err)
	}
	if got != existing {
		t.Fatalf("expected existing PR to be returned, got %#v", got)
	}
}

func TestPullRequest_UpdateRepository_PullRequestNotFound(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	req := &domain.PullRequest{ID: "pr-1"}

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, req.ID).
		Return(nil, ErrPullRequestNotFound)

	// Act
	got, err := usecase.UpdatePullRequest(ctx, req)

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrPullRequestNotFound) {
		t.Fatalf("expected ErrPullRequestNotFound, got %v", err)
	}
}

func TestPullRequest_UpdateRepository_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	old := &domain.PullRequest{
		ID:     "pr-1",
		Status: domain.RequestStatusOpen,
	}
	updated := &domain.PullRequest{
		ID:     "pr-1",
		Status: domain.RequestStatusMerged,
	}

	req := &domain.PullRequest{ID: "pr-1", Status: domain.RequestStatusMerged}

	gomock.InOrder(
		mockPRRepo.EXPECT().
			GetPullRequestByID(ctx, req.ID).
			Return(old, nil),
		mockPRRepo.EXPECT().
			UpdatePullRequest(ctx, req).
			Return(nil),
		mockPRRepo.EXPECT().
			GetPullRequestByID(ctx, req.ID).
			Return(updated, nil),
	)

	// Act
	got, err := usecase.UpdatePullRequest(ctx, req)

	// Assert
	if err != nil {
		t.Fatalf("UpdatePullRequest() unexpected error: %v", err)
	}
	if got.Status != domain.RequestStatusMerged {
		t.Fatalf("expected status %v, got %v", domain.RequestStatusMerged, got.Status)
	}
}

func TestPullRequest_ReassignRequest_MergedNoOp(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	usecase := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	pr := &domain.PullRequest{ID: "pr-1"}
	stored := &domain.PullRequest{
		ID:       "pr-1",
		AuthorID: "author-1",
		Status:   domain.RequestStatusMerged,
	}

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, pr.ID).
		Return(stored, nil)

	// Act
	got, _, err := usecase.ReassignRequest(ctx, pr.ID, "old-reviewer")

	// Assert
	if err != nil {
		t.Fatalf("ReassignRequest() unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil result for merged PR, got %#v", got)
	}
}

func TestPullRequest_ReassignRequest_NoCandidates(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	uc := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	stored := &domain.PullRequest{
		ID:       "pr-1",
		AuthorID: "author-1",
		Status:   domain.RequestStatusOpen,
	}
	author := &domain.User{
		ID:       "author-1",
		Username: "author",
		IsActive: true,
	}

	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, stored.ID).
		Return(stored, nil)

	mockUserRepo.EXPECT().
		GetUserByID(ctx, stored.AuthorID).
		Return(author, nil)

	mockReqOwnerRepo.EXPECT().
		DeleteRequestOwner(ctx, &domain.RequestOwner{
			RequestID: stored.ID,
			UserID:    "old-reviewer",
			Role:      domain.UserRoleReviewer,
		}).
		Return(nil)

	mockReqOwnerRepo.EXPECT().
		GetUsersByPullRequestID(ctx, stored.ID).
		Return([]domain.RequestOwner{
			{UserID: "used-1", RequestID: stored.ID, Role: domain.UserRoleReviewer},
		}, nil)

	// ВАЖНО: теперь usecase вызывает GetTeamsByUserID по pr.AuthorID,
	// а не по req.AuthorID, поэтому аргумент — stored.AuthorID
	mockUserRepo.EXPECT().
		GetTeamsByUserID(ctx, stored.AuthorID).
		Return([]domain.Team{{Name: "team-1"}}, nil)

	// Все кандидаты либо неактивные, либо уже назначены / автор
	coworkers := []domain.User{
		{ID: "used-1", Username: "used", IsActive: true},     // уже назначен
		{ID: "author-1", Username: "author", IsActive: true}, // автор
		{ID: "u3", Username: "u3", IsActive: false},          // неактивный
	}
	mockUserRepo.EXPECT().
		GetUsersByTeamName(ctx, "team-1").
		Return(coworkers, nil)

	// Act
	got, _, err := uc.ReassignRequest(ctx, stored.ID, "old-reviewer")

	// Assert
	if got != nil {
		t.Fatalf("expected nil result, got %#v", got)
	}
	if !errors.Is(err, ErrCannotFindActiveMembers) {
		t.Fatalf("expected ErrCannotFindActiveMembers, got %v", err)
	}
}

func TestPullRequest_ReassignRequest_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockPRRepo := NewMockPullRequestRepository(ctrl)
	mockTeamRepo := NewMockTeamRepository(ctrl)
	mockUserRepo := NewMockUserRepository(ctrl)
	mockReqOwnerRepo := NewMockRequestOwnerRepository(ctrl)

	uc := NewPullRequest(mockPRRepo, mockTeamRepo, mockUserRepo, mockReqOwnerRepo)

	requestID := "pr-1"

	stored := &domain.PullRequest{
		ID:       requestID,
		AuthorID: "author-1",
		Status:   domain.RequestStatusOpen,
	}
	author := &domain.User{
		ID:       "author-1",
		Username: "author",
		IsActive: true,
	}

	// Первый вызов GetPullRequestByID — в начале usecase
	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, requestID).
		Return(stored, nil)

	mockUserRepo.EXPECT().
		GetUserByID(ctx, stored.AuthorID).
		Return(author, nil)

	mockReqOwnerRepo.EXPECT().
		DeleteRequestOwner(ctx, &domain.RequestOwner{
			RequestID: stored.ID,
			UserID:    "old-reviewer",
			Role:      domain.UserRoleReviewer,
		}).
		Return(nil)

	// Уже назначенный ревьюер, которого нельзя назначать снова
	mockReqOwnerRepo.EXPECT().
		GetUsersByPullRequestID(ctx, stored.ID).
		Return([]domain.RequestOwner{
			{UserID: "used-1", RequestID: stored.ID, Role: domain.UserRoleReviewer},
		}, nil)

	// ВАЖНО: GetTeamsByUserID теперь вызывается с pr.AuthorID
	mockUserRepo.EXPECT().
		GetTeamsByUserID(ctx, stored.AuthorID).
		Return([]domain.Team{{Name: "team-1"}}, nil)

	coworkers := []domain.User{
		{ID: "used-1", Username: "used", IsActive: true},
		{ID: "free-1", Username: "free", IsActive: true},
	}
	mockUserRepo.EXPECT().
		GetUsersByTeamName(ctx, "team-1").
		Return(coworkers, nil)

	mockReqOwnerRepo.EXPECT().
		SaveRequestOwner(ctx, gomock.AssignableToTypeOf(&domain.RequestOwner{})).
		DoAndReturn(func(_ context.Context, ro *domain.RequestOwner) error {
			if ro.Role != domain.UserRoleReviewer {
				t.Fatalf("expected reviewer role, got %v", ro.Role)
			}
			if ro.RequestID != stored.ID {
				t.Fatalf("expected RequestID %s, got %s", stored.ID, ro.RequestID)
			}
			if ro.UserID == "used-1" || ro.UserID == "author-1" {
				t.Fatalf("assigned invalid reviewer id: %s", ro.UserID)
			}
			return nil
		})

	// Новый usecase в конце ещё раз читает PR из репозитория:
	// pr, err = p.pullRequestRepository.GetPullRequestByID(ctx, requestID)
	// поэтому нужно второе ожидание
	mockPRRepo.EXPECT().
		GetPullRequestByID(ctx, requestID).
		Return(stored, nil)

	// Act
	got, _, err := uc.ReassignRequest(ctx, requestID, "old-reviewer")

	// Assert
	if err != nil {
		t.Fatalf("ReassignRequest() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil PullRequest")
	}
	if got.ID != requestID {
		t.Fatalf("expected PullRequest ID %s, got %s", requestID, got.ID)
	}
}
