package postgres

import (
	"avito-test/internal/db"
	"avito-test/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RequestOwnerRepository struct {
	db *db.Queries
}

func NewRequestOwnerRepository(db *db.Queries) *RequestOwnerRepository {
	return &RequestOwnerRepository{db: db}
}

func (r *RequestOwnerRepository) SaveRequestOwner(ctx context.Context, requestOwner *domain.RequestOwner) error {
	err := r.db.AssignUserPullRequest(ctx, db.AssignUserPullRequestParams{Userid: requestOwner.UserID, Pullrequestid: requestOwner.RequestID, Role: string(requestOwner.Role)})
	if err != nil {
		return fmt.Errorf("can't save request owner: %w", err)
	}
	return nil
}

func (r *RequestOwnerRepository) DeleteRequestOwner(ctx context.Context, requestOwner *domain.RequestOwner) error {
	err := r.db.DeletePullRequestAssignOfUser(ctx, db.DeletePullRequestAssignOfUserParams{Userid: requestOwner.UserID, Pullrequestid: requestOwner.RequestID})
	if err != nil {
		return fmt.Errorf("can't delete request owner: %w", err)
	}
	return nil
}

func (r *RequestOwnerRepository) GetRequestsByUserID(ctx context.Context, userID string) ([]domain.RequestOwner, error) {
	user, err := r.db.GetUsersAssignedPullRequest(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.RequestOwner{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get requests by user id: %w", err)
	}
	requestOwners := make([]domain.RequestOwner, len(user))
	for i, u := range user {
		requestOwners[i] = domain.RequestOwner{UserID: userID, RequestID: u.Pullrequestid, Role: domain.Role(u.Role)}
	}
	return requestOwners, nil
}

func (r *RequestOwnerRepository) GetUsersByPullRequestID(ctx context.Context, pullRequestID string) ([]domain.RequestOwner, error) {
	users, err := r.db.GetListOfUsersByPullRequestID(ctx, pullRequestID)
	if errors.Is(err, sql.ErrNoRows) {
		return []domain.RequestOwner{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("can't get users by pull request id: %w", err)
	}
	requestOwners := make([]domain.RequestOwner, len(users))
	for i, u := range users {
		requestOwners[i] = domain.RequestOwner{UserID: u.Userid, RequestID: pullRequestID, Role: domain.Role(u.Role)}
	}
	return requestOwners, nil
}
